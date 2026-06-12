# Ed25519 keys and onion addresses

The owner field of an XNS claim is a raw compressed Ed25519 public key. A Tor v3 onion service also has a long-term Ed25519 identity. Using that identity as the XNS owner joins the name and service without introducing another key or record format.

There are three forms to keep distinct:

```text
Ed25519 seed
Ed25519 public key
Tor expanded secret
```

The seed is the original 32-byte private value from which the keypair is derived. OpenSSL stores it inside an Ed25519 PKCS#8 private-key file. The public key is the 32-byte compressed point written into XNS. Tor stores the private side differently because onion-service key blinding requires the 64-byte expanded Ed25519 secret.

The expansion is one-way:

```text
seed32 -> SHA-512 -> expanded_secret64
```

The first half is clamped for use as the private scalar, and the second half is retained for Ed25519 signing. Tor writes this expanded value into `hs_ed25519_secret_key`.

The original seed cannot be recovered from the expanded secret. This does not prevent an existing Tor service from being used with XNS: its public key can still be derived exactly. It only means that a Tor-generated service does not give the operator the original portable seed afterward.

## The preferred starting point

For a new XNS identity, begin with the ordinary Ed25519 seed:

```sh
openssl genpkey -algorithm Ed25519 -out service.pem
```

This leaves one simple source identity from which the public key, XNS owner key, onion address and Tor service files can all be derived.

Check the key type:

```sh
openssl pkey -in service.pem -text -noout
```

It must be an Ed25519 private key. `xns-onion service` accepts unencrypted PKCS#8 PEM files generated in this form and rejects other key types.

The PEM and Tor secret file both control the onion identity. Keep them private. The PEM is the more general root material, while the Tor file is what the running onion service needs.

## Existing onion services

An existing onion service does not need to be replaced. Its `hostname`, `hs_ed25519_public_key` and `hs_ed25519_secret_key` already describe one identity.

Read and verify it with:

```sh
./xns-onion inspect /path/to/hidden_service
```

The printed `owner_key` is suitable for an XNS claim. The inability to reconstruct the original seed does not change ownership of the onion service or the validity of its public key.

Do not create a new PEM and combine its files with an existing service directory. Tor requires the public and secret files to belong to the same keypair, and `xns-onion inspect` rejects a mixed directory.

## From key to address

For version 3 onion services, Tor computes:

```text
B = public_key32 || checksum2 || 03
address = lowercase_base32(B) || ".onion"
```

The checksum is:

```text
SHA3-256(".onion checksum" || public_key32 || 03)[0:2]
```

This is why an XNS owner key can become an onion address without a network request, and why an onion address can yield its owner key without consulting an indeXer.

The next step is [creating the onion service](/docs/xns-and-tor/creating-an-onion-service).
