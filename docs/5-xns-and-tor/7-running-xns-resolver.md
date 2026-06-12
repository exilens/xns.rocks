# Running XNS Resolver

XNS Resolver makes `.xns` names available to ordinary Linux applications. It asks an indeXer for the owner key, derives the corresponding Tor v3 address, and carries the application's TCP connection through Tor.

It is a specialized fork of [TSR](https://github.com/aeeravsar/TSR). The general proxy-routing and configuration machinery was removed; XNS Resolver handles one route:

```text
.xns hostname -> rightmost XNS name -> owner key -> Tor onion service
```

## Requirements

The current implementation requires:

- Linux
- `systemd-resolved`
- Tor with a SOCKS5 listener
- root privileges for TUN and route setup
- an XNS indeXer URL chosen by the user

Clone and build it:

```sh
git clone https://github.com/exilens/xns-resolver
cd xns-resolver
go build -o xns-resolver ./cmd/xns-resolver
```

Run it with an indeXer:

```sh
sudo ./xns-resolver \
  --indexer <xns-indexer-url>
```

The indeXer is required. XNS Resolver does not silently choose one.

By default it expects Tor SOCKS5 on:

```text
socks5://127.0.0.1:9050
```

Use another listener explicitly:

```sh
sudo ./xns-resolver \
  --indexer <xns-indexer-url> \
  --tor-proxy socks5://127.0.0.1:9150
```

The remaining options control the local virtual network:

```text
--tun xns0
--cidr 10.204.0.0/16
--gateway 10.204.0.1
--dns 10.204.0.1:53
--mtu 1500
```

These values are implementation defaults, not XNS protocol constants. Change them when they conflict with the host's existing routes or interface names.

## What it changes

XNS Resolver creates a TUN interface and assigns a private virtual IPv4 range to it. It starts a DNS listener on the virtual gateway and asks `systemd-resolved` to send only `.xns` queries to that listener.

For an active name, the resolver:

1. requests `/lookup?name=` from the configured indeXer
2. validates the returned name, owner key and source transactions
3. converts the owner key into a Tor v3 onion address
4. returns a synthetic IPv4 address to the application
5. maps TCP connections to that address through Tor SOCKS5

Other DNS names remain on the system's normal resolver path.

Stop the process with `Ctrl+C` or `SIGTERM`. It reverts the `systemd-resolved` route, removes the virtual interface and discards every cached mapping.

## Cache and confirmations

The cache exists only in memory. Restarting XNS Resolver begins with no name mappings.

When an indeXer reports `"finalized": true`, the mapping is cached until the name's estimated expiration based on `remaining_blocks` and Monero's two-minute target. When the latest claim or renewal has fewer than ten confirmations, the mapping is checked again after at most one minute.

The resolver does not persist an indeXer response as authority. Every fresh process must obtain the current registry state again.

## Current limits

The rightmost label before `.xns` is looked up through the indeXer:

```text
example.xns           -> example
indexer.example.xns   -> example
api.v2.example.xns    -> example
```

The complete hostname is preserved for the application. All subdomains of one base name therefore reach the same onion identity, while the service's virtual-host configuration decides which of them are actually served.

The virtual network carries TCP. Tor does not provide general UDP transport, so UDP applications cannot be carried to onion services through the resolver.

Continue with [using XNS names](/docs/xns-and-tor/using-xns-names) and checking each layer of the path.
