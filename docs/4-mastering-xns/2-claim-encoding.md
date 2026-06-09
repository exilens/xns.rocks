# Claim encoding

An XNS action is recognized from an incoming Monero transfer to the protocol wallet. Both the amount and `tx_extra` must satisfy the protocol.

The application payload is exactly 64 bytes:

```text
name32 || owner_key32
```

It is carried as a Monero nonce field:

```text
02 40 name32 owner_key32
```

`02` is the `tx_extra` nonce tag. `40` is hexadecimal for 64, the one-byte length of the nonce data.

A normal XNS transaction also contains its transaction public key field, so its relevant `tx_extra` structure is:

```text
01 <32-byte transaction public key>
02 40 <32-byte name> <32-byte owner key>
```

Other Monero fields may be present when they are understood by the parser, but the transaction must contain exactly one valid XNS payload. A missing payload, malformed `tx_extra` or multiple payloads makes the XNS action invalid even though the Monero transaction itself remains valid.

## Name field

The name occupies 32 bytes. Its ASCII bytes begin at the first position and every unused byte is zero:

```text
alice 00 00 00 ... 00
```

The decoded name ends at the first zero byte. All remaining bytes must also be zero, preventing multiple binary encodings of the same visible name.

Validation is byte-based:

- length from 1 through 32
- lowercase `a-z`
- digits `0-9`
- hyphen `-`
- first and last bytes must be letters or digits

No case folding, Unicode normalization or dot hierarchy is performed. Invalid bytes make the whole payload invalid.

## Owner field

The owner field is a canonical 32-byte compressed Ed25519 point. It must decode onto the curve, must not be the identity, and must belong to the prime-order subgroup.

The subgroup check prevents small-order or mixed-order points from becoming owner identities. In particular, the order-two point used as the protocol wallet spend public key is not accepted as an owner key.

The protocol stores the compressed point exactly as 64 lowercase hexadecimal characters in indeXer responses. How an application uses the corresponding private key for authentication or resolution is above the base registry; XNS establishes the name-to-key ownership.

## Amount field

One year costs `10,000,000,000` Monero atomic units:

```text
0.01 XMR
```

The amount received by the protocol wallet must be at least this value and divisible by it without remainder. The number of years is:

```text
years = received_amount / 10,000,000,000
```

An amount below `0.01 XMR` or between whole-year multiples is ignored by XNS. The Monero payment is still irreversible and remains burned.

Each year adds `262800` blocks. Expiration is therefore derived entirely from integer amounts and block heights; wall-clock timestamps do not participate in consensus.
