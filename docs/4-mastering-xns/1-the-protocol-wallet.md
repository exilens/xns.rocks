# The protocol wallet

XNS needs a Monero address that can receive ordinary transactions, reveal those transactions to an indeXer and never allow the received outputs to be spent. A normal wallet satisfies the first two requirements but not the third: somewhere, a valid private spend key still exists.

The XNS protocol wallet is constructed differently. Its view keypair is valid, but its spend public key is not generated from a spend secret.

The spend public key is derived from the fixed string:

```text
XNS
```

using Monero's hash-to-EC map:

```text
cn_fast_hash(input)
ge_fromfe_frombytes_vartime(hash)
ge_mul8(point)
```

The resulting spend public key is:

```text
a0cd652c2b2b9ee7664079650b6a738c7ec3551034c7950a4bc2f1ca02adc999
```

This distinction is important. Random invalid bytes would not be enough. A sender must be able to use the spend public key in Monero's ordinary output-key derivation. The hash-derived point is a valid prime-order curve point, so ordinary wallets can send to it.

There is still no known spend scalar behind it. Anyone can reproduce the public point from `XNS`, but doing so does not produce a private spend key. Spending an output would require the discrete logarithm of the hash-derived point.

The protocol wallet therefore receives normal Monero outputs that every node can validate, while nobody can construct the corresponding key images and spend them.

## The view key

The private view key is deliberately public:

```text
935830fc11250e25153160951c1ba9152e5fee00763890314b67532ae6385607
```

Its public key is:

```text
6451616df9f393cba2315d252e04dcaf8f60fdeed95a74c77a4e6585e90f0823
```

Publishing a private view key would be a privacy failure for an ordinary wallet. For XNS it is the mechanism that makes independent indexing possible. Anyone can create the same view wallet, discover incoming outputs, decode their amounts and reconstruct the registry.

The word "private" in Monero describes the kind of key, not a requirement that XNS keep it secret.

The public view key does not make the wallet spendable. It lets a scanner derive the per-output shared secret. Spending still requires the unknown spend scalar for the hash-derived spend public key.

## Address derivation

Each network address is derived locally from:

```text
network prefix || spend public key || view public key
```

The first four bytes of Keccak-256 over that body form the checksum, and the result is encoded with Monero Base58.

The resulting addresses are canonical XNS protocol values. The implementation derives and verifies them from the view secret, spend-key derivation input and network prefix rather than relying on copied address strings. A reviewer can repeat the same derivation independently.

Since the key material is shared, mainnet, stagenet and testnet differ only by their standard Monero address prefix.

## Restore heights

The protocol wallet was created at known network heights. An indeXer begins there instead of scanning years of blocks in which the current protocol wallet did not exist:

```text
mainnet:  3690551
stagenet: 2135563
testnet:  3018160
```

These heights are part of the protocol constants because choosing a later value could hide valid claims. They are not command-line optimization hints.

The wallet cache built from them is still disposable. Given the constants and a Monero node, the complete protocol wallet can always be generated and scanned again.

See [Network constants](/docs/mastering-xns/network-constants) for the resulting addresses.
