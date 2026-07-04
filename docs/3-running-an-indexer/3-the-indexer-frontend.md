# The indeXer frontend

The indeXer daemon deliberately exposes JSON rather than a built-in website. This keeps protocol indexing separate from presentation and allows one daemon to serve command-line clients, applications and human-facing resolvers without carrying their interface code.

The official indeXer frontend is a separate Go program that places a small web interface in front of any XNS indeXer. It contains no JavaScript. A lookup is an ordinary HTML form submission, the Go server requests `/lookup?name=` from the configured indeXer, and the result is rendered into HTML.

The frontend repository is separate from both the XNS implementation and xns.rocks:

```sh
git clone https://github.com/exilens/indexer-frontend
cd indexer-frontend
go build -o indexer .
```

Both runtime addresses must be specified:

```sh
./indexer \
  --indexer <xns-indexer-url> \
  --listen <listen-address>
```

`--indexer` is the base URL of the XNS indeXer API. A trailing slash is accepted and removed.

`--listen` is the address on which the web frontend serves users. It has no default and does not need to match the indeXer's listen address.

Open the address selected with `--listen`.

## Request flow

The browser submits:

```text
GET /?name=alice
```

The frontend validates the name using the XNS character and length rules, then requests:

```text
GET <indexer>/lookup?name=alice
```

The request is performed by the frontend server, not by browser-side JavaScript. This means the indeXer API can remain bound to a private interface while only the frontend is exposed publicly.

The upstream request has a 15-second timeout, and the response body is limited to one mebibyte before JSON decoding. If the indeXer returns a structured error such as `indexer is synchronizing`, the frontend presents that message. Connection failures, malformed JSON and unexpected HTTP responses are reduced to simple lookup errors rather than exposing internal response data.

## Displayed state

For an active name, the page shows:

- name
- owner key
- expiration height
- estimated remaining time
- source transaction IDs

The remaining time is calculated by the frontend from `remaining_blocks`, using Monero's two-minute target block time:

```text
30 blocks     = 1 hour
720 blocks    = 1 day
262800 blocks = 1 year
```

This value is an estimate for humans. The expiration height and remaining block count are the protocol values.

An unclaimed or expired name is shown as not claimed. Invalid input and indeXer errors appear on the same page, while unrelated paths return `404 Not Found`.

## Deployment

A common deployment keeps both programs on localhost:

```text
public HTTPS
    |
reverse proxy
    |
indeXer frontend on 127.0.0.1:8090
    |
XNS indeXer on 127.0.0.1:8087
    |
Monero node
```

The reverse proxy should provide TLS and forward requests to the frontend. The XNS indeXer does not need to be publicly reachable unless other remote clients also use its JSON API.

The frontend has no administration endpoint, database or claim capability. Restarting or replacing it cannot change XNS state. It is a disposable view over one configured indeXer.

That final point is also its trust boundary. The frontend shows whatever its configured indeXer reports. The recommended deployment is to run the frontend beside an indeXer you operate yourself.
