# Using XNS names

Once Tor and XNS Resolver are running, `.xns` names behave like local system names for TCP applications.

Test the complete path:

```sh
curl http://example.xns/
```

No application proxy setting is needed. `curl` resolves the name through the operating system, receives the synthetic address from XNS Resolver, and opens an ordinary TCP connection. The resolver carries that connection to the derived onion service through Tor.

The same path can be used by browsers and other TCP clients:

```text
http://example.xns/
http://indexer.example.xns/
```

An application does not need to understand XNS, the indeXer API or onion-address encoding. It only needs to use the system resolver and connect to the returned address.

For `indexer.example.xns`, XNS Resolver looks up `example` and preserves the full hostname for the application. The onion service must explicitly accept `indexer.example.xns`; otherwise a correctly configured reverse proxy should return `404`.

## What the visitor trusts

The configured indeXer reports the current XNS owner key. XNS Resolver validates the response and onion conversion, but it does not independently scan Monero. A dishonest indeXer can lie to its own client.

Users may choose another indeXer, compare several, or run one locally. The XNS claim remains on Monero regardless of which HTTP endpoint is used to read it.

Tor authenticates the onion destination derived from the owner key. If the indeXer returns the correct key, the connection cannot quietly terminate at a different onion identity with a similar name.

## Finding a failure

The route has distinct layers. Test them in order.

First, query the indeXer:

```sh
./xns lookup \
  --indexer <xns-indexer-url> \
  example
```

The result must be found and contain the expected owner key.

Next, derive the onion address:

```sh
./xns-onion address <owner-key>
```

Compare it with the publisher's onion hostname.

Check local name resolution:

```sh
getent hosts example.xns
```

It should return an address inside the resolver's configured virtual range.

Then test HTTP:

```sh
curl -v http://example.xns/
```

If DNS succeeds but the connection fails, check that Tor SOCKS is running, the onion service is published and the resolver is using the intended SOCKS port. If the connection succeeds but returns `404`, the service probably does not accept `Host: example.xns`. If it returns the wrong site, inspect the reverse proxy's host matching and default server.

Direct onion access is a useful comparison:

```sh
curl --socks5-hostname <tor-socks-address> \
  http://<onion-address>/
```

If the onion address works directly but `.xns` does not, the problem is before or at XNS Resolver. If neither works, the publisher's Tor or local service configuration is the more likely cause.

## The result

The visible name is short, but no registrar maps it to a conventional server address. Monero establishes which Ed25519 key owns the name. The key determines the onion service. Tor carries the connection.

The full route is therefore reproducible from public protocol state:

```text
XNS name -> Monero claim -> owner key -> onion address -> onion service
```

There is no DNS zone, account dashboard or hidden database between the name and its cryptographic destination.
