# Using XNS names

Once the I2P router and XNS Resolver are running, applications can use `.xns` names through the normal system resolver.

Test the complete path:

```sh
curl http://example.xns/
```

No I2P proxy setting is needed in `curl`. The application receives a synthetic address from XNS Resolver and opens an ordinary TCP connection. The resolver carries it through the configured I2P SOCKS listener to the extended destination derived from the XNS owner key.

Subdomains follow the same route:

```text
http://example.xns/
http://indexer.example.xns/
```

The service must explicitly accept each Host value it intends to serve.

## What the visitor trusts

The configured indeXer reports the current XNS owner key. XNS Resolver validates the response and address derivation, but it does not scan Monero itself. Users may choose another indeXer or run one locally.

I2P locates the encrypted LeaseSet2 destination belonging to the derived public identity. The changing blinded publication keys do not change the stable owner key or extended address.

## Finding a failure

Test the path in order.

Query the indeXer:

```sh
./xns lookup \
  --indexer <xns-indexer-url> \
  example
```

Derive the expected address:

```sh
./xns-i2p address <owner-key>
```

On the publisher, compare this result with:

```sh
./xns-i2p inspect /path/to/service
```

Then open the i2pd console, select the service under **Local Destinations**, and compare it with **Encrypted B33 address**. Do not compare it with the ordinary 52-character address on the Server Tunnels page.

Check the I2P SOCKS listener on the visitor:

```sh
ss -ltn | grep ':4447'
```

Check local XNS resolution:

```sh
getent hosts example.xns
```

It should return an address from the resolver's virtual range.

Finally:

```sh
curl -v http://example.xns/
```

If name resolution succeeds but the connection fails, check the I2P SOCKS listener, encrypted LeaseSet2 publication and i2pd tunnel state. If the connection reaches the service but returns `404`, the reverse proxy probably does not accept `Host: example.xns`.

## DNS records

XNS Resolver owns address records for `.xns` names. An `A` query returns the synthetic address used by the resolver's virtual route. An `AAAA` query returns no address, so applications continue with IPv4 through the resolver.

Other record types are forwarded to the owner service over TCP port `53`. This lets the service publish records such as TXT while keeping routing under XNS Resolver:

```sh
dig TXT openalias.example.xns
```

The owner must run a DNS server on the same I2P destination if it wants these records to exist. If that server is missing or unreachable, non-address lookups fail instead of pretending the record does not exist.

On the publisher, verify the local backend separately:

```sh
curl -v -H 'Host: example.xns' http://127.0.0.1:<local-port>/
```

If nothing is listening on the configured `host` and `port`, the I2P destination may publish correctly while every incoming stream fails at the final local hop.

## Reading i2pd

The **I2P tunnels** page confirms that the server tunnel was loaded. The **Local Destinations** detail page confirms the Encrypted B33 address and shows LeaseSets and tunnel state.

Useful logs are:

```sh
sudo journalctl -u i2pd -n 100 --no-pager
sudo tail -n 100 /var/log/i2pd/i2pd.log
```

`Unknown section type` means the tunnel `type` is invalid. XNS uses `type = server`. An address mismatch usually means the intended `private.dat` was not loaded; with the standard i2pd service layout, use a data-directory-relative key path such as:

```ini
keys = xns/private.dat
```

The complete route is:

```text
XNS name -> Monero claim -> owner key -> extended I2P address
         -> encrypted LeaseSet2 -> local service
```

The ordinary destination hash may still be displayed and used by I2P software, but it is not part of the XNS derivation.
