# Serving XNS hostnames

XNS Resolver changes the destination of the TCP connection. It does not rewrite the protocol inside it.

When a visitor opens:

```text
http://example.xns/
```

the resolver sends the TCP stream through I2P to the encrypted LeaseSet2 destination derived from the owner key. The HTTP request still contains:

```http
Host: example.xns
```

This is why the i2pd tunnel uses:

```ini
type = server
```

The raw server tunnel carries the request to the local listener without replacing the Host value. The web server must accept the `.xns` hostname.

Subdomains work the same way. `indexer.example.xns` resolves the base XNS name `example`, reaches the same I2P destination, and arrives with:

```http
Host: indexer.example.xns
```

The service decides whether that subdomain exists.

## Caddy

Suppose i2pd forwards to Caddy on `127.0.0.1:8033`, while the application listens on `127.0.0.1:8080`.

The i2pd tunnel is:

```ini
[xns]
type = server
host = 127.0.0.1
port = 8033
keys = xns/private.dat
signaturetype = 11
i2cp.leaseSetType = 5
```

Caddy can then route the XNS hostnames:

```caddyfile
:8033 {
    bind 127.0.0.1

    @service host example.xns indexer.example.xns

    handle @service {
        reverse_proxy 127.0.0.1:8080
    }

    handle {
        respond 404
    }
}
```

The site address is `:8033`, while `bind 127.0.0.1` restricts the listener to loopback. Writing `http://127.0.0.1:8033` would make `127.0.0.1` part of Caddy's host matching and prevent `Host: example.xns` from reaching the intended handler.

## nginx

The equivalent nginx configuration is:

```nginx
server {
    listen 127.0.0.1:8033;
    server_name example.xns indexer.example.xns;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

server {
    listen 127.0.0.1:8033 default_server;
    server_name _;
    return 404;
}
```

The default server prevents an unrelated Host value from reaching the application.

The current XNS Resolver carries TCP streams. I2P also has datagram protocols, but they are outside this resolver's route.

With the publisher side complete, visitors can continue with [running XNS Resolver](/docs/xns-and-i2p/running-xns-resolver).
