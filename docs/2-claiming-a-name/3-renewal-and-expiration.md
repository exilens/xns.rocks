# Renewal and expiration

XNS has no separate renewal transaction. Renewal is an ordinary claim for an active name using exactly the same owner key.

Suppose `alice` was claimed at height `3,700,000` for one year. Its expiration height is:

```text
3,700,000 + 262,800 = 3,962,800
```

If the same owner makes a two-year claim before that height, the new time is added to the existing expiration:

```text
3,962,800 + (2 * 262,800) = 4,488,400
```

It is not added to the height of the renewal transaction. Renewing early therefore does not discard time already paid for.

Every successful renewal is appended to `source_txids`, allowing the active entry to be traced through its original claim and later extensions.

## Renew before expiration

The protocol treats a name as expired when its expiration height is less than or equal to the height of the transaction currently being applied. At that point the previous owner has no priority.

This means the owner should not wait for the final block. A renewal that remains in the transaction pool while the expiration height passes may be ordered as a fresh claim rather than a renewal. If another owner reaches the canonical chain first, the old owner's transaction will be ignored.

Renewal should be done with enough time for ordinary wallet synchronization, transaction relay, mining delays and possible reorgs.

## Claims by another owner

While a name is active, a valid claim carrying a different owner key does not alter the registry. The XMR is still burned because Monero cannot know that the XNS action is invalid. The indeXer records the transaction as ignored with the reason that the name is already active for a different owner.

XNS does not refund protocol mistakes. Before claiming or renewing, resolve the name and compare the complete owner key.

## After expiration

An expired entry is no longer returned as found by the indeXer API. The next valid claim becomes a new ownership period regardless of whether it uses the old owner key or a different one.

The new expiration is calculated from the new claim's block height. Its `first_claim_height` and `source_txids` begin again, because it is a new active ownership rather than a continuation of the expired one.

There is no grace period, recovery period or privileged renewal window. Such exceptions would require somebody to decide who deserves special treatment. XNS uses the block height and the owner key, and nothing else.
