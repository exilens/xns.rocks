# IndeXer API

The indeXer exposes read-only JSON. It does not create claims, change owners or administer the registry. Claims happen on Monero; the HTTP daemon only serves the registry state it derives from Monero.

The simplest lookup endpoint remains:

```text
GET /lookup?name=<name>
```

This endpoint is kept for command-line clients and resolvers that only need to resolve one name.

Example:

```sh
curl 'http://127.0.0.1:8087/lookup?name=alice'
```

An active name returns `200 OK`:

```json
{
  "found": true,
  "name": "alice",
  "owner_key": "5866666666666666666666666666666666666666666666666666666666666666",
  "expiration_height": 4000000,
  "remaining_blocks": 200000,
  "finalized": true,
  "source_txids": [
    "<initial claim>",
    "<renewal>"
  ]
}
```

An unclaimed or expired name also returns `200 OK`:

```json
{
  "found": false,
  "name": "alice"
}
```

The distinction between an HTTP failure and `"found": false` matters. Not finding a name is an ordinary registry result.

## Names

The name detail endpoint is linkable and returns the current entry for one name:

```text
GET /names/<name>
```

Example:

```sh
curl 'http://127.0.0.1:8087/names/alice'
```

An existing current entry returns:

```json
{
  "name": "alice",
  "owner_key": "5866666666666666666666666666666666666666666666666666666666666666",
  "expiration_height": 4000000,
  "first_claim_height": 3737200,
  "last_update_height": 3737200,
  "remaining_blocks": 200000,
  "active": true,
  "finalized": true,
  "source_txids": [
    "<initial claim>"
  ]
}
```

`owner_key` is the raw compressed Ed25519 public key recorded by the claim.

`expiration_height` is the first block height at which the current ownership is no longer active.

`remaining_blocks` is calculated against the wallet height used for the current indeXer state. With Monero's target block time of two minutes, applications may estimate wall-clock time as `remaining_blocks * 2 minutes`, but the protocol itself uses heights rather than clocks.

`active` is true when `expiration_height` is greater than the current indeXer height.

`finalized` is `true` when the latest transaction contributing to the current entry has at least ten confirmations. A recent claim or renewal returns `false` until that transaction reaches the indeXer's durable confirmation boundary.

`source_txids` contains the initial claim and every same-owner renewal belonging to the current ownership period.

List current name entries with:

```text
GET /names
```

Supported filters:

```text
q=<text>
owner_key=<64-hex-ed25519-public-key>
limit=<n>
offset=<n>
```

`q` is a substring search over current names. It follows the same character rules as names: lowercase `a-z`, digits and `-`.

`owner_key` is an exact owner-key match.

Results are ordered by name.

## Events

Events describe how the registry reached its current state. They include successful claims, renewals and ignored protocol-wallet transfers.

Fetch one event by transaction ID:

```text
GET /events/<txid>
```

List events with:

```text
GET /events
```

Supported filters:

```text
name=<name>
owner_key=<64-hex-ed25519-public-key>
action=claimed|renewed|ignored
reason=<reason>
height_min=<height>
height_max=<height>
limit=<n>
offset=<n>
```

Results are ordered from newest to oldest by block height, then by transaction ID.

Ignored events use machine-readable reasons:

```text
INVALID_AMOUNT
INVALID_NAME
INVALID_OWNER
INVALID_PAYLOAD
ACTIVE_FOR_DIFFERENT_OWNER
```

`INVALID_AMOUNT` means the payment amount does not encode a whole number of XNS years.

`INVALID_NAME` means the XNS payload contains a name that fails the protocol name rules.

`INVALID_OWNER` means the owner key in the payload is not a valid owner public key.

`INVALID_PAYLOAD` covers missing payloads, malformed payloads and other payload-level failures that are not specifically name or owner failures.

`ACTIVE_FOR_DIFFERENT_OWNER` means the payload is valid, but the name is still active under another owner.

## Pagination

List endpoints use the same pagination parameters:

```text
limit=<n>
offset=<n>
```

The default limit is `50`. There is no protocol maximum. Operators who need stricter limits can place another API or reverse proxy in front of the core indeXer.

List responses use:

```json
{
  "limit": 50,
  "offset": 0,
  "count": 2,
  "items": []
}
```

`count` is the number of items returned in that response, not the total number of matching rows.

## Errors

Malformed names, owner keys, transaction IDs, filters and pagination values return `400 Bad Request`:

```json
{
  "error": "unsupported query parameter"
}
```

Unknown query parameters also return `400 Bad Request`. The API rejects them so mistakes do not silently change the meaning of a request.

Missing `/names/<name>` and `/events/<txid>` resources return `404 Not Found`.

While the indeXer is starting or rebuilding state, it returns `503 Service Unavailable`:

```json
{
  "error": "indexer is synchronizing"
}
```

Other paths return `404 Not Found`.

## Command-line lookup

The XNS client wraps `/lookup`:

```sh
./xns lookup \
  --indexer http://127.0.0.1:8087 \
  alice
```

It validates the name locally, sends the HTTP request and prints the indeXer's JSON response unchanged. `--indexer` is required for every lookup.

The recommended setup is to run your own indeXer. It is not a heavy service, and relying on a permanent public endpoint recreates the kind of authority XNS is meant to avoid. Public indeXers are useful for testing and convenience, not as the sovereign way to use the system.

The API response is useful evidence, not a cryptographic proof by itself. The transaction IDs and events allow a verifier to inspect the underlying Monero claims and reproduce the registry rules independently.
