# Network constants

The XNS protocol wallet uses the same key material on every supported Monero network:

```text
view secret:
935830fc11250e25153160951c1ba9152e5fee00763890314b67532ae6385607

view public key:
6451616df9f393cba2315d252e04dcaf8f60fdeed95a74c77a4e6585e90f0823

spend derivation input:
XNS

spend public key derived with Monero hash-to-EC:
a0cd652c2b2b9ee7664079650b6a738c7ec3551034c7950a4bc2f1ca02adc999
```

The spend public key is generated from public protocol data, not from a spend secret. There is no known Monero private spend scalar for the hash-derived point.

## Mainnet

```text
address prefix: 18
restore height: 3690551
address:
47iYR5A7CRbfhs8LQaHSMxQVynmKbTbbi2itFFo4Wg6tSf6Qhzsxe3Lb4W3FbyKCFMWN9rT5EUjzTaNBgS6NCeNj4y35XNS
```

Select mainnet explicitly with `--mainnet`.

## Stagenet

```text
address prefix: 24
restore height: 2135563
address:
57vaVv54r2hfhs8LQaHSMxQVynmKbTbbi2itFFo4Wg6tSf6Qhzsxe3Lb4W3FbyKCFMWN9rT5EUjzTaNBgS6NCeNj536CQMd
```

Use `--stagenet` with both the claim tool and indeXer. A stagenet wallet, node and protocol address must be used together.

## Testnet

```text
address prefix: 53
restore height: 3018160
address:
9yG5uKpNUnhfhs8LQaHSMxQVynmKbTbbi2itFFo4Wg6tSf6Qhzsxe3Lb4W3FbyKCFMWN9rT5EUjzTaNBgS6NCeNj4ysvF9S
```

The Go protocol package records the testnet constants, but the current `xns` command-line interface exposes mainnet and stagenet only.

## Deriving the address

For a selected network:

1. multiply the Ed25519 base point by the canonical view secret to obtain the view public key
2. derive the spend public key as Monero hash-to-EC over `XNS`
3. encode the network prefix as a varint
4. append the 32-byte spend public key
5. append the 32-byte view public key
6. append the first four bytes of Keccak-256 over the preceding body
7. encode the complete value with Monero Base58

These are the canonical XNS protocol addresses. Implementations should derive and verify them from the protocol key material and network prefixes rather than trusting copied address strings.

Node URLs are intentionally absent. A Monero daemon is chosen by the user or operator and has no place among protocol constants.
