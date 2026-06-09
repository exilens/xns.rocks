# Registry rules

The XNS registry is the result of replaying protocol-wallet transfers in canonical order. It is not a collection of independent transactions that can be judged without context.

For each transfer, the indeXer first derives the number of years from the received amount and decodes exactly one XNS payload. Invalid amounts or payloads produce ignored events and do not modify the registry.

For a valid action at block height `H`, the extension is:

```text
extension = years * 262800
```

The action then falls into one of three cases.

## New claim

If the name has never been claimed, or its existing expiration height is less than or equal to `H`, the transaction creates a new active entry:

```text
owner_key         = payload owner
expiration_height = H + extension
first_claim_height = H
last_update_height = H
source_txids       = [current txid]
```

An expired name does not preserve any privilege for its former owner. A new claim starts a new source transaction history.

## Renewal

If the name is active and the payload owner equals the current owner, the transaction renews it:

```text
expiration_height += extension
last_update_height  = H
source_txids       += current txid
```

The extension begins at the existing expiration height, not the renewal block. Early renewal never wastes remaining time.

## Conflict

If the name is active under another owner, the transaction is ignored. It does not shorten, extend or otherwise alter the active entry.

The payment is already a valid Monero burn and cannot be refunded. XNS validity is an interpretation applied after Monero accepts the transaction.

## Ordering

Blocks are processed from lower height to higher height. Transactions within one block need an order because two claims for the same free name may otherwise both appear to be first.

For every transaction in the block, XNS computes:

```text
Keccak-256(block_hash || txid)
```

The byte strings are sorted lexicographically. If two keys were ever equal, the transaction ID itself is the final tie-breaker.

An ordinary claimant cannot know the block hash while constructing the transaction, so the final same-block position cannot be selected in advance from the transaction ID alone. A block producer influences the block through proof of work, but cannot freely choose arbitrary hashes without doing the corresponding work.

Every indeXer must verify that the transaction belongs to the canonical block whose hash enters this calculation. Sorting by wallet return order, transaction ID alone or daemon presentation order would define a different protocol.

## Events and entries

An indeXer may retain an event for every incoming transfer, including ignored ones. Events explain how a replay reached its conclusion, but only active entries are returned by lookup.

An entry is active when:

```text
expiration_height > current_height
```

At equality, the name is expired. This boundary is also used while replaying a new claim at that height.
