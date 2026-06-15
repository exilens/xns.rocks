# Running XNS Resolver

XNS Resolver makes `.xns` names available to ordinary Linux applications. In I2P mode it asks an indeXer for the owner key, derives the extended address, and opens application TCP streams through the router's SAM interface.

The route is:

```text
.xns hostname -> rightmost XNS name -> owner key -> extended I2P address -> SAM
```

## Requirements

The current implementation requires:

- Linux
- `systemd-resolved`
- an I2P router with a SAM TCP listener
- root privileges for TUN and route setup
- an XNS indeXer URL chosen by the user

Clone and build it:

```sh
git clone https://github.com/exilens/xns-resolver
cd xns-resolver
go build -o xns-resolver ./cmd/xns-resolver
```

Run it in I2P mode:

```sh
sudo ./xns-resolver \
  --network i2p \
  --indexer <xns-indexer-url> \
  --i2p-sam 127.0.0.1:7656
```

The network, indeXer and SAM address are all required. XNS Resolver does not silently choose any of them.

Port `7656` is the conventional SAM TCP port. i2pd normally enables SAM on `127.0.0.1:7656`. The Java I2P router uses the same conventional port but may require SAM to be enabled explicitly. Keep SAM on loopback unless there is a specific protected deployment that requires otherwise.

Check the listener:

```sh
ss -ltn | grep ':7656'
```

## What it changes

XNS Resolver creates a TUN interface and assigns a private virtual IPv4 range to it. It starts a DNS listener on the virtual gateway and asks `systemd-resolved` to send only `.xns` queries there.

For an active name, the resolver:

1. requests `/lookup?name=` from the configured indeXer
2. validates the returned name, owner key and source transactions
3. derives the extended Ed25519 I2P address
4. returns a synthetic IPv4 address to the application
5. maps TCP connections through a transient SAM STREAM session

Other DNS names remain on the system's normal resolver path.

Stop the process with `Ctrl+C` or `SIGTERM`. It closes the SAM session, reverts the `systemd-resolved` route, removes the virtual interface and discards all cached mappings.

## Local network options

The implementation defaults are:

```text
--tun xns0
--cidr 10.204.0.0/16
--gateway 10.204.0.1
--dns 10.204.0.1:53
--mtu 1500
```

They are not XNS protocol constants. Change them when they conflict with an existing route, interface or listener.

## Cache and confirmations

The cache exists only in memory. Finalized records are cached until their estimated expiration based on `remaining_blocks` and Monero's two-minute target. A record whose latest claim or renewal has fewer than ten confirmations is checked again after at most one minute.

Restarting the resolver clears the cache and obtains current registry state again.

## Subdomains

The rightmost label before `.xns` is looked up:

```text
example.xns           -> example
indexer.example.xns   -> example
api.v2.example.xns    -> example
```

The complete hostname remains visible to the application. XNS Resolver carries TCP streams; the service reached by a destination is determined by its I2P tunnel configuration.

Continue with [using XNS names](/docs/xns-and-i2p/using-xns-names) and testing every layer of the route.
