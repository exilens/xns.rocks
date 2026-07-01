package xnsrocks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	donationAddress       = "43cR4XamK1VA8MN6QtKxeR2qSMgoAwt5i5hSLuQ8YGvfNT5mEU76kmkJmF3yxmBcXRVjawj3eH7yEVMpgbnNsTKu5waCDZW"
	donationViewKey       = "ad86a4fa9430a8c183c465a7454e880ae09c98b79dc4c161d9fb3e9e0088d805"
	donationRestoreHeight = uint64(3692700)
)

type donation struct {
	TxID          string
	Amount        string
	Height        uint64
	Confirmations uint64
	Time          string
}

type donationView struct {
	Address   string
	ViewKey   string
	Donations []donation
	Ready     bool
	Error     string
}

type donationService struct {
	rpc  walletRPC
	proc *walletProcess

	mu        sync.RWMutex
	donations []donation
	ready     bool
	err       string
	stop      chan struct{}
	done      chan struct{}
}

type walletRPC struct {
	url    string
	client *http.Client
}

type walletProcess struct {
	cmd      *exec.Cmd
	log      *os.File
	done     chan struct{}
	stopOnce sync.Once
}

func startDonations(node, dataDir string) (*donationService, error) {
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		return nil, err
	}
	walletDir := filepath.Join(dataDir, "wallet")
	if err := os.MkdirAll(walletDir, 0o700); err != nil {
		return nil, err
	}
	walletFile := filepath.Join(walletDir, "xns_donations")
	port, err := freeLocalPort()
	if err != nil {
		return nil, err
	}
	proc, err := startDonationWalletRPC(node, port, dataDir, walletFile, walletDir)
	if err != nil {
		return nil, err
	}
	rpc := walletRPC{
		url:    fmt.Sprintf("http://127.0.0.1:%d/json_rpc", port),
		client: &http.Client{Timeout: 10 * time.Minute},
	}
	if err := waitDonationWallet(rpc, proc); err != nil {
		proc.stop()
		return nil, err
	}
	if _, err := os.Stat(walletFile + ".keys"); os.IsNotExist(err) {
		err = rpc.call("generate_from_keys", map[string]any{
			"restore_height": donationRestoreHeight,
			"filename":       "xns_donations",
			"address":        donationAddress,
			"viewkey":        donationViewKey,
			"password":       "",
			"language":       "English",
		}, nil)
		if err != nil {
			proc.stop()
			return nil, err
		}
	}
	var address struct {
		Address string `json:"address"`
	}
	if err := rpc.call("get_address", map[string]any{"account_index": 0}, &address); err != nil {
		proc.stop()
		return nil, err
	}
	if address.Address != donationAddress {
		proc.stop()
		return nil, errors.New("donation wallet address does not match website constant")
	}

	service := &donationService{
		rpc:  rpc,
		proc: proc,
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	go service.loop()
	return service, nil
}

func startDonationWalletRPC(node string, port int, dataDir, walletFile, walletDir string) (*walletProcess, error) {
	logPath := filepath.Join(dataDir, "wallet-rpc.log")
	args := []string{
		"--rpc-bind-ip", "127.0.0.1",
		"--rpc-bind-port", strconv.Itoa(port),
		"--disable-rpc-login",
		"--non-interactive",
		"--no-initial-sync",
		"--log-file", logPath,
		"--log-level", "0",
		"--daemon-address", node,
		"--trusted-daemon",
	}
	if _, err := os.Stat(walletFile + ".keys"); err == nil {
		args = append(args, "--wallet-file", walletFile, "--password", "")
	} else {
		args = append(args, "--wallet-dir", walletDir)
	}
	out, err := os.OpenFile(logPath+".stdout", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("monero-wallet-rpc", args...)
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Start(); err != nil {
		out.Close()
		return nil, err
	}
	proc := &walletProcess{cmd: cmd, log: out, done: make(chan struct{})}
	go func() {
		_ = cmd.Wait()
		close(proc.done)
	}()
	return proc, nil
}

func waitDonationWallet(rpc walletRPC, proc *walletProcess) error {
	deadline := time.Now().Add(30 * time.Second)
	var last error
	for time.Now().Before(deadline) {
		select {
		case <-proc.done:
			return errors.New("donation wallet-rpc exited before becoming ready")
		default:
		}
		if err := rpc.call("get_version", map[string]any{}, nil); err == nil {
			return nil
		} else {
			last = err
		}
		time.Sleep(250 * time.Millisecond)
	}
	return fmt.Errorf("donation wallet-rpc did not become ready: %w", last)
}

func (service *donationService) loop() {
	defer close(service.done)
	service.refresh()
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-service.stop:
			return
		case <-ticker.C:
			service.refresh()
		}
	}
}

func (service *donationService) refresh() {
	if err := service.rpc.call("refresh", map[string]any{}, nil); err != nil {
		service.setError(err)
		return
	}
	var transfers struct {
		In []struct {
			TxID          string `json:"txid"`
			Amount        uint64 `json:"amount"`
			Height        uint64 `json:"height"`
			Confirmations uint64 `json:"confirmations"`
			Timestamp     uint64 `json:"timestamp"`
		} `json:"in"`
	}
	if err := service.rpc.call("get_transfers", map[string]any{
		"in":            true,
		"out":           false,
		"pending":       false,
		"failed":        false,
		"pool":          false,
		"account_index": 0,
	}, &transfers); err != nil {
		service.setError(err)
		return
	}
	donations := make([]donation, 0, len(transfers.In))
	for _, transfer := range transfers.In {
		donations = append(donations, donation{
			TxID:          transfer.TxID,
			Amount:        formatXMR(transfer.Amount),
			Height:        transfer.Height,
			Confirmations: transfer.Confirmations,
			Time:          time.Unix(int64(transfer.Timestamp), 0).UTC().Format("2006-01-02 15:04 UTC"),
		})
	}
	sort.Slice(donations, func(i, j int) bool {
		if donations[i].Height != donations[j].Height {
			return donations[i].Height > donations[j].Height
		}
		return donations[i].TxID > donations[j].TxID
	})
	service.mu.Lock()
	service.donations = donations
	service.ready = true
	service.err = ""
	service.mu.Unlock()
}

func (service *donationService) setError(err error) {
	service.mu.Lock()
	service.err = err.Error()
	service.mu.Unlock()
}

func (service *donationService) handler(templates *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		service.mu.RLock()
		view := donationView{
			Address:   donationAddress,
			ViewKey:   donationViewKey,
			Donations: append([]donation(nil), service.donations...),
			Ready:     service.ready,
			Error:     service.err,
		}
		service.mu.RUnlock()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := templates.ExecuteTemplate(w, "donate.html", view); err != nil {
			http.Error(w, "render donation page", http.StatusInternalServerError)
		}
	})
}

func (service *donationService) Close() {
	close(service.stop)
	service.proc.stop()
	<-service.done
}

func (proc *walletProcess) stop() {
	proc.stopOnce.Do(func() {
		select {
		case <-proc.done:
		default:
			_ = proc.cmd.Process.Signal(os.Interrupt)
		}
		select {
		case <-proc.done:
		case <-time.After(10 * time.Second):
			_ = proc.cmd.Process.Kill()
			<-proc.done
		}
		_ = proc.log.Close()
	})
}

func (rpc walletRPC) call(method string, params, result any) error {
	raw, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      "0",
		"method":  method,
		"params":  params,
	})
	if err != nil {
		return err
	}
	resp, err := rpc.client.Post(rpc.url, "application/json", bytes.NewReader(raw))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("wallet rpc returned HTTP %d", resp.StatusCode)
	}
	var out struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return err
	}
	if out.Error != nil {
		return fmt.Errorf("wallet rpc %d: %s", out.Error.Code, out.Error.Message)
	}
	if result == nil {
		return nil
	}
	return json.Unmarshal(out.Result, result)
}

func freeLocalPort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func formatXMR(atomic uint64) string {
	whole := atomic / 1_000_000_000_000
	fraction := atomic % 1_000_000_000_000
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%d.%012d", whole, fraction), "0"), ".") + " XMR"
}
