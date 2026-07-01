package xnsrocks

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"flag"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

//go:embed templates static docs
var files embed.FS

func Main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	listen := flag.String("listen", "", "HTTP listen address")
	node := flag.String("node", "", "Monero daemon RPC URL")
	dataDir := flag.String("data-dir", "", "persistent state directory")
	flag.Parse()
	if *listen == "" {
		return errors.New("--listen is required")
	}
	if *node == "" {
		return errors.New("--node is required")
	}
	if *dataDir == "" {
		return errors.New("--data-dir is required")
	}

	templates, err := template.ParseFS(files, "templates/*.html")
	if err != nil {
		return err
	}
	static, err := fs.Sub(files, "static")
	if err != nil {
		return err
	}
	docs, err := loadDocs(files)
	if err != nil {
		return err
	}
	donations, err := startDonations(*node, filepath.Join(*dataDir, "donations"))
	if err != nil {
		return err
	}
	defer donations.Close()

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))
	mux.HandleFunc("GET /{$}", render(templates, "index.html"))
	mux.Handle("GET /docs", docs.handler(templates))
	mux.Handle("GET /docs/{page...}", docs.handler(templates))
	mux.Handle("GET /donate", donations.handler(templates))

	log.Printf("xns.rocks listening on http://%s", *listen)
	server := &http.Server{Addr: *listen, Handler: mux}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	errs := make(chan error, 1)
	go func() {
		errs <- server.ListenAndServe()
	}()
	select {
	case err := <-errs:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdown)
	}
}

func render(templates *template.Template, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := templates.ExecuteTemplate(w, name, nil); err != nil {
			log.Printf("render %s: %v", name, err)
		}
	}
}

type docPage struct {
	Title    string
	URL      string
	Category string
	Content  template.HTML
}

type docCategory struct {
	Title  string
	Slug   string
	Topics []*docPage
}

type docsSite struct {
	Main       *docPage
	Categories []*docCategory
	pages      map[string]*docPage
}

type docsView struct {
	Page       *docPage
	Categories []*docCategory
}

var orderPrefix = regexp.MustCompile(`^\d+-`)
var orderedName = regexp.MustCompile(`^(\d+)-`)

func loadDocs(source fs.FS) (*docsSite, error) {
	markdown := goldmark.New(
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithHardWraps()),
	)
	site := &docsSite{pages: make(map[string]*docPage)}

	mainPage, err := readDoc(markdown, source, "docs/index.md", "/docs", "")
	if err != nil {
		return nil, err
	}
	site.Main = mainPage
	site.pages[mainPage.URL] = mainPage

	entries, err := fs.ReadDir(source, "docs")
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool {
		return documentationNameLess(entries[i].Name(), entries[j].Name())
	})

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		categorySlug := slug(entry.Name())
		categoryURL := "/docs/" + categorySlug
		categoryIndex := path.Join("docs", entry.Name(), "index.md")
		categoryTitle, err := readDocTitle(source, categoryIndex)
		if err != nil {
			return nil, err
		}
		category := &docCategory{Title: categoryTitle, Slug: categorySlug}
		topics, err := fs.ReadDir(source, path.Join("docs", entry.Name()))
		if err != nil {
			return nil, err
		}
		sort.Slice(topics, func(i, j int) bool {
			return documentationNameLess(topics[i].Name(), topics[j].Name())
		})

		for _, topic := range topics {
			if topic.IsDir() || topic.Name() == "index.md" || path.Ext(topic.Name()) != ".md" {
				continue
			}
			topicSlug := slug(strings.TrimSuffix(topic.Name(), ".md"))
			topicURL := categoryURL + "/" + topicSlug
			topicPath := path.Join("docs", entry.Name(), topic.Name())
			page, err := readDoc(markdown, source, topicPath, topicURL, categorySlug)
			if err != nil {
				return nil, err
			}
			if _, exists := site.pages[topicURL]; exists {
				return nil, &docsError{Path: topicPath, Message: "duplicate documentation URL"}
			}
			site.pages[topicURL] = page
			category.Topics = append(category.Topics, page)
		}
		site.Categories = append(site.Categories, category)
	}

	return site, nil
}

func readDocTitle(source fs.FS, filename string) (string, error) {
	sourceText, err := fs.ReadFile(source, filename)
	if err != nil {
		return "", err
	}
	title := firstHeading(string(sourceText))
	if title == "" {
		return "", &docsError{Path: filename, Message: "missing H1 title"}
	}
	return title, nil
}

func readDoc(markdown goldmark.Markdown, source fs.FS, filename, url, category string) (*docPage, error) {
	sourceText, err := fs.ReadFile(source, filename)
	if err != nil {
		return nil, err
	}
	title := firstHeading(string(sourceText))
	if title == "" {
		return nil, &docsError{Path: filename, Message: "missing H1 title"}
	}

	var rendered bytes.Buffer
	if err := markdown.Convert(sourceText, &rendered); err != nil {
		return nil, err
	}
	return &docPage{
		Title:    title,
		URL:      url,
		Category: category,
		Content:  template.HTML(rendered.String()),
	}, nil
}

func firstHeading(markdown string) string {
	for _, line := range strings.Split(markdown, "\n") {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return ""
}

func slug(name string) string {
	return orderPrefix.ReplaceAllString(name, "")
}

func documentationNameLess(a, b string) bool {
	aOrder, aNumbered := documentationOrder(a)
	bOrder, bNumbered := documentationOrder(b)
	if aNumbered != bNumbered {
		return aNumbered
	}
	if aNumbered && aOrder != bOrder {
		return aOrder < bOrder
	}
	return a < b
}

func documentationOrder(name string) (int, bool) {
	match := orderedName.FindStringSubmatch(name)
	if match == nil {
		return 0, false
	}
	order, err := strconv.Atoi(match[1])
	return order, err == nil
}

func (site *docsSite) handler(templates *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := strings.TrimSuffix(r.URL.Path, "/")
		if url == "" {
			url = "/"
		}
		page, ok := site.pages[url]
		if !ok {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		view := docsView{Page: page, Categories: site.Categories}
		if err := templates.ExecuteTemplate(w, "docs.html", view); err != nil {
			log.Printf("render docs.html: %v", err)
		}
	})
}

type docsError struct {
	Path    string
	Message string
}

func (err *docsError) Error() string {
	return err.Path + ": " + err.Message
}
