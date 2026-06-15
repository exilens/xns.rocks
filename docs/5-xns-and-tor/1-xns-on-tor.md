# XNS on Tor

The base XNS protocol assigns a name to an Ed25519 public key. It does not decide whether that key identifies a website, a signing identity, a messaging account or something else. This is deliberate: the registry establishes ownership, while applications decide what ownership means in their own field.

Tor v3 onion services are the most direct way to turn that ownership into a reachable service. Their addresses are built from Ed25519 public keys, the same kind of key stored by an XNS claim. Nothing has to translate the name into an IP address controlled by a hosting company. The owner key itself determines the onion destination.

For a name such as `example`, the path is:

```text
example.xns
    |
XNS indeXer
    |
Ed25519 owner public key
    |
Tor v3 onion address
    |
onion service
```

The conversion between the owner key and onion address is deterministic. A Tor v3 address contains:

```text
public_key32 || checksum2 || version1
```

These 35 bytes are encoded as lowercase base32 and followed by `.onion`. The version is `3`, and the checksum is derived from the public key and version. An onion address therefore carries its public identity inside itself; the public key can be recovered and the checksum verified without contacting Tor.

This gives XNS a natural web layer. The publisher claims a name with the public key of an onion service. A visitor resolves the name through an indeXer, derives the onion address, and reaches it through Tor. The indeXer never needs to store onion addresses, and the XNS protocol does not need a Tor-specific record type.

## Two sides of the connection

Publishing and visiting use different tools.

The publisher:

1. creates an Ed25519 identity
2. derives a Tor onion service from it
3. claims an XNS name with its public key
4. configures the service to accept the `.xns` hostname

The visitor:

1. runs Tor
2. runs XNS Resolver with a chosen indeXer
3. opens `http://example.xns` in an ordinary application

The website itself remains an onion service. XNS gives it a human name, and XNS Resolver makes that name usable by applications which know nothing about XNS.

## Subdomains

The rightmost label before `.xns` is always the claimed XNS name:

```text
xns.xns             -> XNS name: xns
indexer.xns.xns     -> XNS name: xns
api.v2.xns.xns      -> XNS name: xns
```

All three names reach the onion service derived from the `xns` owner key. The labels to the left are preserved for the application. A web server may serve `indexer.xns.xns`, reject `www.xns.xns`, and route `api.v2.xns.xns` elsewhere using its ordinary virtual-host configuration.

XNS does not store separate subdomain records. Ownership of the base name determines the onion identity, while the service behind that identity decides which subdomains exist.

## Why there is no onion field

An XNS claim does not contain:

```text
example -> something.onion
```

It contains:

```text
example -> Ed25519 public key
```

The onion address is derived from that key. Storing both would be redundant and would create a second value that could disagree with the owner identity. By keeping the raw public key as the only protocol value, XNS remains useful outside Tor while Tor resolution remains exact.

The next article introduces [xns-onion](/docs/xns-and-tor/the-xns-onion-toolkit), the small toolkit used to move between these forms.
