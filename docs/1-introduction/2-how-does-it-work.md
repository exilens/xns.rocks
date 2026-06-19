# How does it work?

XNS does not have its own blockchain, network, token, miners, validators or consensus.
It uses Monero as its only source of truth.

At the center of XNS is a protocol constant Monero wallet.
It has a valid view keypair, so every payment sent to it can be found and read.
But its public spend key is not generated from a spend secret.
It is derived by applying Monero's hash-to-EC map to the fixed string `XNS`.
There is no known private spend scalar behind it.
Nothing sent to this wallet can ever be spent.

In other words, it is a view-only burn wallet.

## Claiming a name

An XNS claim is a normal Monero transaction to the protocol wallet.
The amount defines how long the name will be owned, while the transaction's `tx_extra` contains the name and the owner's public key.

The XNS payload is:

```text
02 40 name32 owner_key32
```

`02` is the Monero `tx_extra` nonce tag and `40` is the hexadecimal length of the following 64 bytes.

`name32` is the name, encoded directly as lowercase ASCII and filled with zeroes until it is exactly 32 bytes.
Names can be between 1 and 32 bytes, can only contain `a-z`, `0-9` and `-`, and must start and end with a letter or number.
There are no dots.
Subdomains are a separate concern and do not belong in the base name registry.

`owner_key32` is the raw 32-byte compressed Ed25519 public key of the owner.
It must be a valid non-identity point in the Ed25519 prime-order subgroup.

There are no hashes or hidden off-chain claim data.
The name and owner key are written directly into the Monero transaction.
Anyone with the protocol view key can reconstruct XNS from Monero alone.

## Time is paid in Monero

One year costs `0.01 XMR`.
In protocol terms, one year is exactly `262800` Monero blocks.

The payment must be at least `0.01 XMR` and must be a whole multiple of it.
A payment of `0.01 XMR` claims the name for one year, `0.05 XMR` claims it for five years, and so on.

The XMR is not paid to a company, foundation or maintainer.
It is burned in the protocol wallet and can never be recovered.

## There are only claims

XNS has no update, transfer or sale operation.
There are only claims.

If a name is free or expired, the first valid claim assigns it to the owner key in that transaction.
If the same owner claims the active name again, the purchased time is added to its current expiration height.
This is how renewal works without risking loss of the name.
If a different owner claims an active name, the transaction is ignored by the registry.
Once the name expires, it becomes free and can be claimed again.

This is also what makes XNS unsquattable in practice.
An XNS name cannot be safely sold because there is no ownership transfer.
Selling the name would mean selling a private key, and the seller can always keep a copy of it.
There is no way for the buyer to know that the key was destroyed.
Without a trustworthy resale market, there is no practical reason to hoard names for resale.

## The indeXer

Monero holds the truth, but it does not know what XNS means.
An XNS indeXer is a small deterministic reader that scans the protocol wallet and interprets its transactions according to the XNS rules.

It begins at the protocol restore height, finds every incoming transaction, reads its XNS payload, and replays all valid and invalid actions from past to future.
Chronology matters.
A transaction that looks valid alone may be invalid when an earlier transaction has already claimed the same name.

Transactions are ordered by block height.
When multiple claims appear in the same block, they are ordered by:

```text
Keccak-256(block_hash || txid)
```

The block hash is not known when transactions are created, so an actor cannot know or directly choose the final same-block order in advance.
Every correct indeXer sees the same block, applies the same ordering and reaches the same registry.

The indeXer does not own names, approve claims or create truth.
It only reads Monero.
Anyone can run one, destroy its database, scan again and arrive at the same result.

## Confirmations and reorgs

A claim becomes visible after it is mined.
One confirmation is enough for an indeXer to serve it, but the newest state is still considered temporary.

Only transactions with at least 10 confirmations are written into the persistent registry.
Anything newer stays in memory and is reconstructed from Monero after a restart.

The indeXer stores the block hash behind its durable state and verifies the canonical chain while scanning.
If a reorganization changes recent blocks, it simply replays the transactions in their new canonical order.
If a deep reorganization invalidates the wallet cache itself, the protocol view wallet is rebuilt from the protocol restore height and the registry is reconstructed again.

There is no special recovery authority and no hidden state to restore.
As long as Monero exists, XNS can be derived from it.
