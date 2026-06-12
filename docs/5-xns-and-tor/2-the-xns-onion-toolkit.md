# The xns-onion toolkit

`xns-onion` handles the boundary between ordinary Ed25519 keys, XNS owner keys and Tor v3 onion-service files. It does not contact an indeXer, run Tor, create an XNS claim or manage a web server. Its work is local and deterministic.

Clone and build it:

```sh
git clone https://github.com/exilens/xns-onion
cd xns-onion
go build -o xns-onion ./cmd/xns-onion
```

The toolkit has four commands:

```text
xns-onion address OWNER_KEY
xns-onion owner ONION_ADDRESS
xns-onion service PRIVATE_KEY.pem DIRECTORY
xns-onion inspect DIRECTORY
```

## Public conversions

`address` accepts the 64-character hexadecimal owner key used by XNS and prints its Tor v3 onion address:

```sh
./xns-onion address \
  20a16b378779e6f6cd8c7d694e22577a1abf03867a9b8b990a8b5720ebeb511d
```

The result is:

```text
ecqwwn4hphtpntmmpvuu4isxpinl6a4gpknyxgikrnlsb27lkeosjuqd.onion
```

`owner` performs the reverse operation:

```sh
./xns-onion owner \
  ecqwwn4hphtpntmmpvuu4isxpinl6a4gpknyxgikrnlsb27lkeosjuqd.onion
```

It verifies the base32 encoding, version, checksum and Ed25519 point before printing:

```text
20a16b378779e6f6cd8c7d694e22577a1abf03867a9b8b990a8b5720ebeb511d
```

These operations use public information only. They are useful for checking a claim, inspecting an onion address or confirming that two visible identities are the same key.

## Creating Tor service files

`service` reads an unencrypted OpenSSL PKCS#8 Ed25519 private key:

```sh
./xns-onion service service.pem service-directory
```

It creates:

```text
service-directory/
├── hostname
├── hs_ed25519_public_key
└── hs_ed25519_secret_key
```

The directory is written with mode `0700` and the files with mode `0600`. Existing key files are never overwritten. Before writing anything, the toolkit derives the public key, onion address and Tor expanded secret from the same seed and verifies that they agree.

The command prints:

```text
owner_key: <XNS owner public key>
onion_address: <Tor v3 onion address>
```

## Inspecting an existing service

`inspect` reads a Tor onion-service directory:

```sh
./xns-onion inspect /path/to/hidden_service
```

It derives the public key from `hs_ed25519_secret_key`, compares it with `hs_ed25519_public_key`, and checks `hostname` when that file exists. A mismatched or malformed directory is rejected.

On success it prints the owner key and onion address. This is the command to use when an onion service already exists and now needs an XNS name.

`xns-onion` never prints the private key. The only private operation is converting an explicitly supplied PEM into Tor's native files.

Continue with [Ed25519 keys and onion addresses](/docs/xns-and-tor/ed25519-keys-and-onion-addresses) before creating a service identity.
