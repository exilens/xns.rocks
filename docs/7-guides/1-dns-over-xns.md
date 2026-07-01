# DNS over XNS

XNS Resolver already gives `.xns` names to ordinary applications. For address lookups it does not ask the owner service at all. It returns a synthetic address from its virtual route and carries the resulting TCP connection through I2P to the destination derived from the XNS owner key.

Other DNS records are different. A `TXT`, `MX` or similar query does not decide where the `.xns` name routes. XNS Resolver forwards those queries to the owner service over TCP port `53`, on the same I2P destination that belongs to the XNS name.

The split is:

```text
A and AAAA records        -> answered by XNS Resolver
TXT, MX and other records -> forwarded to the owner DNS server
```

This lets the owner publish ordinary DNS data without giving DNS authority over XNS ownership. The XNS name still resolves from Monero, the owner key and the derived I2P destination. The DNS server only speaks for records below that already-owned name.

## Requirements

This guide assumes the I2P path from the previous section is already working:

- an active XNS name
- an i2pd destination generated from the owner key
- XNS Resolver running in I2P mode
- CoreDNS or another authoritative DNS server

The examples use `example.xns`. Replace it with the name you own.

## The I2P tunnel

The owner DNS server must be reachable on port `53` of the XNS I2P destination. With i2pd, add another server tunnel that uses the same destination key as the rest of the service:

```ini
[example_dns]
type = server
host = 127.0.0.1
inport = 53
port = 9053
keys = xns/private.dat
signaturetype = 11
i2cp.leaseSetType = 5
```

`inport = 53` is the port seen by visitors through I2P. `port = 9053` is the local DNS server on the publisher's machine. The `keys` value must point to the same destination identity that was claimed in XNS.

Restart i2pd after changing the tunnel configuration:

```sh
sudo systemctl restart i2pd
```

Then check that the local DNS target is listening:

```sh
ss -ltnup | grep ':9053'
```

## CoreDNS

CoreDNS is a small authoritative server that can serve a zone file directly. Create a directory for the DNS configuration and write a `Corefile`:

```text
example.xns:9053 {
    bind 127.0.0.1
    file example.xns.zone
    log
    errors
}
```

The listener is bound to loopback because the public entrance is the I2P tunnel, not the host network.

Create `example.xns.zone` beside the `Corefile`:

```dns
$ORIGIN example.xns.
$TTL 60

@       IN SOA  ns.example.xns. hostmaster.example.xns. (
            1       ; serial
            60      ; refresh
            60      ; retry
            3600    ; expire
            60 )    ; minimum

@       IN NS   ns.example.xns.
ns      IN A    127.0.0.1

@       IN TXT  "xns=example"
@       IN TXT  "oa1:xmr recipient_address=<monero-address>; recipient_name=<name>;"
```

The `ns` address exists so ordinary DNS tooling sees a complete zone. It is not what makes `example.xns` reachable through XNS.

Do not expect an `A` record in this zone to control where `example.xns` routes. XNS Resolver owns address records for `.xns` names so it can return the synthetic address used by its virtual network. The owner DNS server is mainly for the other record types.

Start CoreDNS from the directory containing the `Corefile`:

```sh
coredns -conf Corefile
```

## Local tests

First test the authoritative server without XNS Resolver or I2P:

```sh
dig @127.0.0.1 -p 9053 example.xns TXT +short
```

It should return the TXT records from the zone file. If this fails, fix CoreDNS before looking at i2pd or XNS Resolver.

Next test the same DNS server through the XNS address. This exercises the resolver's address mapping, the I2P stream and the i2pd tunnel:

```sh
dig @example.xns -p 53 example.xns TXT +short
```

Finally test the normal resolver path:

```sh
dig example.xns TXT +short
```

or query the XNS Resolver listener directly when using the default virtual network:

```sh
dig @10.204.0.1 example.xns TXT +short
```

An address lookup should still return a synthetic resolver address:

```sh
getent hosts example.xns
```

The returned address belongs to XNS Resolver's virtual range. It is not published by the owner DNS zone.

With owner DNS working, the name can publish OpenAlias, mail discovery and other records without adding new XNS protocol fields. Continue with the later guides for those specific uses.
