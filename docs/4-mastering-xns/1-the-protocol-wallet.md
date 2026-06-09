# The protocol wallet

XNS needs a Monero address that can receive ordinary transactions, reveal those transactions to an indeXer and never allow the received outputs to be spent. A normal wallet satisfies the first two requirements but not the third: somewhere, a valid private spend key still exists.

The XNS protocol wallet is constructed differently. Its view keypair is valid, but its spend public key is the compressed Edwards point `(0, -1)`:

```text
ecffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7f
```

This point has order two. Monero spend public keys are scalar multiples of the Ed25519 base point and belong to its large prime-order subgroup. `(0, -1)` is a valid curve point but lies outside that subgroup, so no Monero spend scalar can produce it.

This distinction is important. Random invalid bytes would not be enough. A sender must be able to decode the spend public key as a curve point in order to derive the one-time output key. The order-two point is accepted for that calculation while remaining impossible to match with a valid Monero private spend key.

The protocol wallet therefore receives normal Monero outputs that every node can validate, while nobody can construct the corresponding key images and spend them.

## The view key

The private view key is deliberately public:

```text
c2694f7b2ba66ada8548e31c4ef1616ddce71b47969a337c293f7a72d5804909
```

Its public key is:

```text
51c3b119bbe36d8f56e08d95acc7ed944765cff3387e2e212e3bac41fbd1483d
```

Publishing a private view key would be a privacy failure for an ordinary wallet. For XNS it is the mechanism that makes independent indexing possible. Anyone can create the same view wallet, discover incoming outputs, decode their amounts and reconstruct the registry.

The word "private" in Monero describes the kind of key, not a requirement that XNS keep it secret.

## Address derivation

Each network address is derived locally from:

```text
network prefix || spend public key || view public key
```

The first four bytes of Keccak-256 over that body form the checksum, and the result is encoded with Monero Base58.

The resulting addresses are canonical XNS protocol values. The implementation derives and verifies them from the view secret, spend point and network prefix rather than relying on copied address strings. A reviewer can repeat the same derivation independently.

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
