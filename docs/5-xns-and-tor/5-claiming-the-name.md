# Claiming the name

Before claiming, obtain the public key from the service directory:

```sh
./xns-onion inspect /path/to/hidden_service
```

The output is:

```text
owner_key: <64-character Ed25519 public key>
onion_address: <56-character Tor v3 address>.onion
```

The `owner_key` is the value passed to XNS. Do not pass the onion address, PEM file, Tor secret file or a hash of any of them.

Claim the name with the normal XNS client:

```sh
./xns claim \
  --mainnet \
  --wallet-file <wallet-file> \
  --wallet-password <wallet-password> \
  --name example \
  --owner <owner-key> \
  --node <monero-node-url> \
  --years 1
```

The Monero wallet pays for the claim. It does not need to contain or know the onion private key. The owner key is public and is written directly into the claim payload.

After the transaction is mined, resolve it through an indeXer:

```sh
./xns lookup \
  --indexer <xns-indexer-url> \
  example
```

Compare the returned `owner_key` with:

```sh
./xns-onion inspect /path/to/hidden_service
```

They must match exactly. You can also derive the expected onion address from the lookup result:

```sh
./xns-onion address <owner-key>
```

This three-way comparison proves that the active XNS entry, Tor key files and onion hostname describe the same Ed25519 identity.

## Check before burning

An XNS claim cannot change the owner of an active name. Confirm every value before broadcasting:

- the name follows the XNS name rules
- the onion service starts with the intended key files
- `xns-onion inspect` accepts the directory
- the printed onion address is the service address you intend to publish
- the owner key passed to `xns claim` is copied exactly
- the name is not already active under another owner

If the wrong valid key is claimed, there is no update operation to replace it. The name remains attached to that key until expiration.

Renewal uses the same owner key. Claiming the active name again with the same service identity adds the purchased years to its current expiration.

The name is now correct at the registry layer. The next article explains [serving the .xns hostname](/docs/xns-and-tor/serving-xns-hostnames).
