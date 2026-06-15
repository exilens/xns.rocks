# Serving XNS hostnames

XNS Resolver changes where a TCP connection goes. It does not rewrite the application protocol carried inside that connection.

When a visitor opens:

```text
http://example.xns/
```

the resolver sends the TCP connection through Tor to the onion address derived from the XNS owner key. The HTTP request still contains:

```http
Host: example.xns
```

The onion service must therefore accept the XNS hostname. If the same service should also work through its direct onion address, accept both names.

Subdomains follow the same route. `indexer.example.xns` resolves the XNS name `example`, reaches the same onion service, and arrives with:

```http
Host: indexer.example.xns
```

Resolution does not mean the subdomain is automatically served. The reverse proxy or application decides that by accepting or rejecting the Host value.

## Caddy

Suppose Tor forwards onion port 80 to Caddy on `127.0.0.1:8033`, while the application listens on `127.0.0.1:8080`.

The Tor configuration is:

```text
HiddenServiceDir /path/to/hidden_service
HiddenServicePort 80 127.0.0.1:8033
```

A Caddy configuration for one XNS name is:

```caddyfile
:8033 {
    bind 127.0.0.1

    @service host example.xns indexer.example.xns exampleonionaddress.onion

    handle @service {
        reverse_proxy 127.0.0.1:8080
    }

    handle {
        respond 404
    }
}
```

Replace `exampleonionaddress.onion` with the address printed by `xns-onion inspect`.

The site address is `:8033`, while `bind 127.0.0.1` restricts the listener to loopback. Writing `http://127.0.0.1:8033` instead would make `127.0.0.1` part of Caddy's HTTP host matching, causing requests with `Host: example.xns` to miss the intended route.

The final handler rejects unrelated Host values instead of serving the application for every name that reaches the port.

## nginx

The equivalent nginx server is:

```nginx
server {
    listen 127.0.0.1:8033;
    server_name example.xns indexer.example.xns exampleonionaddress.onion;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

nginx rejects unmatched hosts only if another default server does not accept them. On a shared listener, define an explicit default:

```nginx
server {
    listen 127.0.0.1:8033 default_server;
    server_name _;
    return 404;
}
```

## HTTP and HTTPS

An onion service already provides authenticated end-to-end encryption between the Tor client and onion service. Plain HTTP over the onion connection does not travel as cleartext across the Internet.

HTTPS may still be used for application policy or browser presentation, but it is not required to obtain encrypted onion transport. It also introduces certificate validation for `.xns`, which ordinary public certificate authorities do not provide.

Keep the Tor target and application listeners on loopback unless there is a separate reason to expose them. The public entrance should be the onion service, not an accidentally reachable local reverse-proxy port.

With the publisher side complete, visitors can install and run [XNS Resolver](/docs/xns-and-tor/running-xns-resolver).
