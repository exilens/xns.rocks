# Lookup API

The indeXer intentionally exposes one operation:

```text
GET /lookup?name=<name>
```

There is no endpoint for creating claims, changing owners, listing privileged state or administering the registry. Claims happen on Monero; the HTTP daemon only resolves names.

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

`owner_key` is the raw compressed Ed25519 public key recorded by the claim.

`expiration_height` is the first block height at which the current ownership is no longer active.

`remaining_blocks` is calculated against the wallet height used for the current indeXer state. With Monero's target block time of two minutes, applications may estimate wall-clock time as `remaining_blocks * 2 minutes`, but the protocol itself uses heights rather than clocks.

`finalized` is `true` when the latest transaction contributing to the current entry has at least ten confirmations. A recent claim or renewal returns `false` until that transaction reaches the indeXer's durable confirmation boundary.

`source_txids` contains the initial claim and every same-owner renewal belonging to the current active ownership period.

## Not found

A valid but unclaimed or expired name also returns `200 OK`:

```json
{
  "found": false,
  "name": "alice"
}
```

The distinction between an HTTP failure and `"found": false` matters. Not finding a name is an ordinary registry result.

Malformed names return `400 Bad Request`:

```json
{
  "error": "name may only contain lowercase a-z, 0-9, and -"
}
```

While the indeXer is starting or rebuilding state, it returns `503 Service Unavailable`:

```json
{
  "error": "indexer is synchronizing"
}
```

Other paths return `404 Not Found`.

## Command-line lookup

The XNS client wraps the same endpoint:

```sh
./xns lookup \
  --indexer http://127.0.0.1:8087 \
  alice
```

It validates the name locally, sends the HTTP request and prints the indeXer's JSON response unchanged. `--indexer` is required for every lookup.

An application should allow users to choose their indeXer. Depending forever on one public server would recreate the exact kind of authority XNS is designed to avoid. For higher assurance, compare independent indeXers or operate one locally.

The API response is useful evidence, not a cryptographic proof by itself. The transaction IDs allow a verifier to inspect the underlying Monero claims and reproduce the registry rules independently.
