# Running an indeXer

An XNS indeXer is an HTTP daemon that reads the protocol wallet and answers name lookups. It does not need a private spend key, accept claim submissions or communicate with other indeXers. Its only upstream dependency is a Monero node.

Clone and build XNS, and make sure `monero-wallet-rpc` is available in `PATH`:

```sh
git clone https://github.com/exilens/xns
cd xns
go build -o xns ./cmd/xns
```

Run a mainnet indeXer:

```sh
./xns indexer \
  --mainnet \
  --node <monero-node-url> \
  --listen <listen-address> \
  --data-dir <data-directory>
```

Run one on stagenet:

```sh
./xns indexer \
  --stagenet \
  --node <stagenet-node-url> \
  --listen <listen-address> \
  --data-dir <data-directory>
```

Exactly one of `--mainnet` or `--stagenet` is required. The user always supplies the node; daemon addresses are not protocol constants.

`--listen` is the public lookup server, not the managed wallet RPC. Expose it beyond localhost only when the surrounding network and reverse-proxy configuration are intentional.

`--data-dir` chooses where the wallet cache, SQLite database and logs live. It is required so an operator cannot accidentally create persistent protocol state in an assumed location.

Use a different directory for every network. If a data directory contains a database for another protocol address, the indeXer removes that stale database and builds the correct state.

## Managed protocol wallet

The indeXer starts its own `monero-wallet-rpc` on a free localhost port. The port is chosen at runtime and remains private to the process.

On first use, it derives the protocol address from the protocol view secret, invalid spend public key and selected network prefix, then creates a view-only wallet at:

```text
<data-dir>/wallet/xns_protocol_view
```

Later runs reuse the wallet cache after verifying that its address still matches the protocol address. A mismatched cache is deleted and generated again.

The wallet begins at the protocol restore height rather than genesis. These heights are protocol constants chosen when the current protocol wallet was created; no valid XNS claim can exist before them.

## Startup

The HTTP listener starts immediately, but lookups return `503 Service Unavailable` until the first refresh and replay finish. This prevents old disk state from being served as though it were current.

The first startup may take time because the protocol wallet must scan from its restore height. Later startups reuse the Monero wallet cache, although the indeXer still reloads and replays all protocol transfers before becoming ready.

Progress is written to standard error:

```text
scan: starting
scan: refreshing protocol wallet
scan: loading wallet transfers
scan: found 12 incoming transfers
scan: ordering transfers
scan: replaying transfers
scan: done in 4s, visible_names=3 durable_names=3
```

The scan repeats every 30 seconds. If one scan is still running when the next interval arrives, the new scan is skipped instead of running concurrently.

Stop the process with `Ctrl+C` or `SIGTERM`. The HTTP server shuts down, current scan work is allowed to finish, the wallet RPC is stopped, and SQLite is closed.

## Node requirements

The node must provide ordinary daemon RPC methods including block lookup, transaction decoding and chain height. The indeXer passes it to `monero-wallet-rpc` as a trusted daemon and also queries it directly to verify block hashes and canonical transaction membership.

For a public service, operating a local synchronized Monero node gives the indeXer the clearest trust and availability boundary. A remote node can delay, omit or misrepresent chain data to that indeXer, although it cannot create a claim that honest Monero nodes do not accept.

Continue with the [lookup API](/docs/running-an-indexer/lookup-api) to place applications in front of the daemon.
