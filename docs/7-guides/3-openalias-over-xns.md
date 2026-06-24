# OpenAlias over XNS

OpenAlias lets a Monero wallet turn a name into a Monero address. The user may type an email-like alias such as `donate@example.xns`, and the wallet looks up the DNS name `donate.example.xns`. The alias is stored as a DNS `TXT` record, so XNS does not need a special OpenAlias field. Owner DNS can already publish that record below an owned `.xns` name.

The route is:

```text
donate@example.xns
    |
donate.example.xns
    |
XNS Resolver
    |
owner DNS over I2P port 53
    |
TXT record
    |
Monero address
```

The XNS claim establishes who controls `example`. The owner DNS server publishes what `example.xns` says about payment addresses.

## The TXT record

Add an OpenAlias record to the owner DNS zone:

```dns
@       IN TXT  "oa1:xmr recipient_address=<monero-address>; recipient_name=<name>;"
```

For a donation alias or another subdomain, place the record under that label:

```dns
donate  IN TXT  "oa1:xmr recipient_address=<monero-address>; recipient_name=<name>;"
```

Then the visible names are:

```text
example.xns
donate@example.xns
```

In the second case, the wallet converts the alias to `donate.example.xns` before it asks DNS for the TXT record.

The Monero address is not written into the XNS claim. It is ordinary owner-published DNS data. Changing the TXT record changes the address returned for the alias, but it does not change XNS ownership.

## Testing

After reloading the DNS server, query the record:

```sh
dig example.xns TXT +short
dig donate.example.xns TXT +short
```

or query XNS Resolver directly when using the default virtual network:

```sh
dig @10.204.0.1 example.xns TXT +short
```

The answer should contain the `oa1:xmr` record.

Monero wallets that use the system resolver can resolve `.xns` OpenAlias names while XNS Resolver is running. The wallet is still reading an OpenAlias TXT record; XNS only changes how that DNS name is reached.

## DNSSEC warnings

Many Monero wallets warn when an OpenAlias record is not protected by DNSSEC. That warning is correct for ordinary DNS names.

For `.xns`, DNSSEC is not the authority model. XNS Resolver does not ask the public DNS hierarchy who owns the name. XNS ownership comes from Monero: the indeXer scans the protocol wallet, reads the claim transactions, applies the registry rules, and reports the current owner key. XNS Resolver uses that owner key, derives the I2P destination from it, and forwards the TXT query to the owner DNS server at that destination.

The OpenAlias record is therefore not protected by a DNSSEC chain from the DNS root. It is reached through XNS ownership and the I2P destination identity. A wallet may still show a DNSSEC warning because it only knows the ordinary OpenAlias convention. When the user intentionally resolves through XNS Resolver, the absence of DNSSEC is expected.

## What is proven

XNS proves which owner key controls the name. I2P reaches the destination derived from that key. Owner DNS publishes the OpenAlias TXT record. The wallet reads the same OpenAlias format it already understands, but the name reaches its DNS record through XNS rather than the ordinary DNS root.
