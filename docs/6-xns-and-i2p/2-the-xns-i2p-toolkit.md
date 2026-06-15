# The xns-i2p toolkit

`xns-i2p` handles the boundary between ordinary Ed25519 keys, XNS owner keys and native i2pd destination files. It does not contact an indeXer, create a claim, run i2pd or configure a web server.

Clone and build it:

```sh
git clone https://github.com/exilens/xns-i2p
cd xns-i2p
go build -o xns-i2p ./cmd/xns-i2p
```

The toolkit has four commands:

```text
xns-i2p address OWNER_KEY
xns-i2p owner I2P_ADDRESS
xns-i2p service PRIVATE_KEY.pem DIRECTORY
xns-i2p inspect DIRECTORY
```

## Public conversions

`address` accepts the 64-character hexadecimal owner key used by XNS:

```sh
./xns-i2p address \
  20a16b378779e6f6cd8c7d694e22577a1abf03867a9b8b990a8b5720ebeb511d
```

It prints the extended address:

```text
gngj6ifbnm3yo6pg63gyy7ljjyrfo6q2x4bym6u3romqvc2xedv6wui5.b32.i2p
```

`owner` performs the reverse operation:

```sh
./xns-i2p owner \
  gngj6ifbnm3yo6pg63gyy7ljjyrfo6q2x4bym6u3romqvc2xedv6wui5.b32.i2p
```

It validates the base32 encoding, checksum, signing types and Ed25519 point before printing:

```text
20a16b378779e6f6cd8c7d694e22577a1abf03867a9b8b990a8b5720ebeb511d
```

These commands accept the 56-character extended address used by encrypted LeaseSet2. They do not accept the ordinary 52-character destination hash because that hash does not contain the owner key.

## Creating i2pd service files

`service` reads an unencrypted OpenSSL PKCS#8 Ed25519 private key:

```sh
./xns-i2p service service.pem service-directory
```

It creates:

```text
service-directory/
├── hostname
└── private.dat
```

`private.dat` is a native i2pd destination key file. The directory is written with mode `0700`, the files with mode `0600`, and existing files are never overwritten.

The command prints:

```text
owner_key: <XNS owner public key>
i2p_address: <extended I2P address>
```

## Inspecting a generated service

`inspect` verifies the complete destination:

```sh
./xns-i2p inspect /path/to/service-directory
```

It reads the Ed25519 seed from `private.dat`, derives the signing public key, X25519 encryption key, identity padding and extended address again, and rejects the file if any component differs. It also checks `hostname` when present.

`xns-i2p` never prints private key material. Its output contains only the public owner key and extended address.

Continue with [Ed25519 keys and I2P addresses](/docs/xns-and-i2p/ed25519-keys-and-i2p-addresses) before installing a service.
