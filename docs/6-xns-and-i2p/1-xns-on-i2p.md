# XNS on I2P

The base XNS protocol assigns a name to an Ed25519 public key. I2P can turn the same key into a reachable service, but the relationship is less obvious than it is with Tor because I2P commonly displays more than one address for one destination.

An ordinary `.b32.i2p` address is a hash of the complete I2P Destination. It includes more than the Ed25519 signing key, so it cannot be reconstructed from the XNS owner key alone.

Encrypted LeaseSet2 destinations also have an extended base32 address. i2pd calls it the **Encrypted B33 address**. This address carries the unblinded signing public key and its signing types, protected by a checksum. When the signing key is Ed25519, the address can be derived directly from the XNS owner key.

That is the address XNS uses.

For a name such as `example`, the path is:

```text
example.xns
    |
XNS indeXer
    |
Ed25519 owner public key
    |
extended .b32.i2p address
    |
encrypted LeaseSet2 destination
```

The service remains an I2P destination. XNS does not store an I2P address, LeaseSet or router record. It stores the Ed25519 public key from which the extended address is derived.

## Two I2P addresses

After a service is loaded, the i2pd **Server Tunnels** page shows an address such as:

```text
66xmjbxlwnjpekspa6mctq63minmnub7oe2gfagew6cwcfiqgswq.b32.i2p
```

This is the ordinary 52-character destination hash. It identifies the same local destination, but it is not the address derived by XNS.

Open that destination under **Local Destinations** and i2pd also shows:

```text
Encrypted B33 address:
gnaj6ifbnm3yo6pg63gyy7ljjyrfo6q2x4bym6u3romqvc2xedv6wui5.b32.i2p
```

The label before `.b32.i2p` is 56 characters. This is the XNS-compatible address. Both addresses can refer to the same service, but only the extended address contains enough information to recover the XNS owner key.

## Publishing and visiting

The publisher:

1. creates an Ed25519 identity
2. generates an i2pd destination from it
3. publishes the destination as encrypted LeaseSet2
4. claims an XNS name with the Ed25519 public key
5. configures the local service to accept the `.xns` hostname

The visitor:

1. runs an I2P router with a SOCKS5 listener
2. runs XNS Resolver in I2P mode with a chosen indeXer
3. opens `http://example.xns` in an ordinary application

XNS Resolver derives the extended address, asks the I2P router's SOCKS listener to open a stream to it, and leaves the application protocol unchanged.

## Subdomains

The rightmost label before `.xns` is the claimed XNS name:

```text
example.xns           -> XNS name: example
indexer.example.xns   -> XNS name: example
api.v2.example.xns    -> XNS name: example
```

Every hostname reaches the I2P destination derived from `example`. The labels to the left remain in the HTTP `Host` header, so the service decides which subdomains exist through its ordinary virtual-host configuration.

The next article introduces [xns-i2p](/docs/xns-and-i2p/the-xns-i2p-toolkit), the toolkit that creates and verifies these destinations.
