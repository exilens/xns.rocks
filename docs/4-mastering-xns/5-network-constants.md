# Network constants

The XNS protocol wallet uses the same key material on every supported Monero network:

```text
view secret:
c2694f7b2ba66ada8548e31c4ef1616ddce71b47969a337c293f7a72d5804909

view public key:
51c3b119bbe36d8f56e08d95acc7ed944765cff3387e2e212e3bac41fbd1483d

spend public key:
ecffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7f
```

The spend public key is the compressed order-two Edwards point `(0, -1)`. There is no corresponding valid Monero private spend key.

## Mainnet

```text
address prefix: 18
restore height: 3690551
address:
4Ac28oSg9AnjpXCZedGfVQjpXCZedGfVQjpXCZedGfVQNJA993cu61iQya98tknArxRoVUsqZ6XPs6YtpHAY6nhD7zWk1QC
```

Select mainnet explicitly with `--mainnet`.

## Stagenet

```text
address prefix: 24
restore height: 2135563
address:
5Ap4DeMdnmtjpXCZedGfVQjpXCZedGfVQjpXCZedGfVQNJA993cu61iQya98tknArxRoVUsqZ6XPs6YtpHAY6nhD7wEWn9y
```

Use `--stagenet` with both the claim tool and indeXer. A stagenet wallet, node and protocol address must be used together.

## Testnet

```text
address prefix: 53
restore height: 3018160
address:
A29Zd46wRXtjpXCZedGfVQjpXCZedGfVQjpXCZedGfVQNJA993cu61iQya98tknArxRoVUsqZ6XPs6YtpHAY6nhD7yP5JTv
```

The Go protocol package records and tests the testnet constants, but the current `xns` command-line interface exposes mainnet and stagenet only.

## Deriving the address

For a selected network:

1. multiply the Ed25519 base point by the canonical view secret to obtain the view public key
2. encode the network prefix as a varint
3. append the 32-byte spend public key
4. append the 32-byte view public key
5. append the first four bytes of Keccak-256 over the preceding body
6. encode the complete value with Monero Base58

These are the canonical XNS protocol addresses. Implementations should derive and verify them from the protocol key material and network prefixes rather than trusting copied address strings.

Node URLs are intentionally absent. A Monero daemon is chosen by the user or operator and has no place among protocol constants.
