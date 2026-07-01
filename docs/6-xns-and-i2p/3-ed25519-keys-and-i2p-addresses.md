# Ed25519 keys and I2P addresses

An I2P Destination contains a signing key, an encryption key and identity data. The XNS owner field contains only a raw compressed Ed25519 public key. A deterministic construction is therefore needed to create the rest of the destination from one portable secret.

For a new service, `xns-i2p` begins with the 32-byte Ed25519 seed stored in an OpenSSL PKCS#8 file. The seed is expanded with SHA-512 and clamped into the RedDSA private scalar used by i2pd signing type `11`.

The destination also needs an X25519 encryption key. `xns-i2p` derives it separately from the same signing scalar:

```text
expanded       = SHA-512(ed25519_seed)
signing_scalar = clamp(expanded[0:32])
encryption_key = HMAC-SHA256(signing_scalar, "XNS" || 0x00)
identity_pad   = HMAC-SHA256(signing_scalar, "XNS" || 0x01)
```

The selector bytes keep the encryption key and padding distinct. The signing scalar is not converted or reused as the X25519 secret.

The result is one deterministic native i2pd file:

```text
Ed25519 seed
    |
RedDSA signing scalar
    |
    +-> Ed25519 public key
    +-> X25519 encryption keypair
    +-> deterministic identity padding
    |
private.dat
```

The same PEM recreates the same `private.dat` byte for byte.

## The extended address

For the XNS-compatible Ed25519 destination, the extended address encodes:

```text
flags1 || signing_type1 || blinded_type1 || public_key32
```

The flags are zero, the unblinded signing type is `11` for RedDSA-Ed25519, and the blinded signing type is also `11`. A CRC32 checksum of the public key is mixed into the first three bytes before the 35-byte result is encoded as lowercase base32.

This produces a 56-character label followed by `.b32.i2p`.

The address identifies the stable unblinded public key. Encrypted LeaseSet2 uses blinded keys that change with time when publishing and locating the service. The public extended address remains stable while the network lookup follows the current blinded destination.

## Ordinary base32 is different

The ordinary address shown beside a server tunnel is:

```text
base32(SHA-256(full_destination)) || ".b32.i2p"
```

Its label is 52 characters. It hashes the encryption key, signing key, certificate and identity padding together. The Ed25519 owner key cannot be recovered from it.

The extended address has a 56-character label and is displayed by i2pd as **Encrypted B33 address**. This is the address printed by `xns-i2p`, derived by XNS Resolver and reversible back to the owner key.

## Generate the root key

Create a new Ed25519 identity:

```sh
openssl genpkey -algorithm Ed25519 -out service.pem
```

Check its type:

```sh
openssl pkey -in service.pem -text -noout
```

Keep the PEM private. It is the portable root from which the complete i2pd destination can be regenerated.

The next step is [creating the I2P service](/docs/xns-and-i2p/creating-an-i2p-service).
