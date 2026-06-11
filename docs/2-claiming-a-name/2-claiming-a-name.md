# Claiming a name

The `xns claim` command performs the complete claim. It opens the wallet, synchronizes it, builds the transaction, inserts the XNS payload, signs it with the full wallet and broadcasts it. There is no dry-run mode and no unsigned result left for another wallet to finish.

```sh
./xns claim \
  --mainnet \
  --wallet-file <wallet-file> \
  --wallet-password <wallet-password> \
  --name alice \
  --owner 5866666666666666666666666666666666666666666666666666666666666666 \
  --node <monero-node-url> \
  --years 1
```

For stagenet, add `--stagenet` and use a stagenet node and wallet:

```sh
./xns claim \
  --stagenet \
  --wallet-file <stagenet-wallet-file> \
  --wallet-password <wallet-password> \
  --name alice \
  --owner 5866666666666666666666666666666666666666666666666666666666666666 \
  --node <stagenet-node-url> \
  --years 1
```

The flags have no hidden meaning:

- `--wallet-file` is the Monero wallet file used to pay
- `--wallet-password` opens that wallet
- `--name` is the XNS name
- `--owner` is the 64-character hexadecimal Ed25519 public key
- `--node` is the Monero daemon RPC URL chosen by the user
- `--years` is the number of whole XNS years
- exactly one of `--mainnet` or `--stagenet` is required

Every option is required. For a wallet with no password, specify the empty value explicitly:

```text
--wallet-password ''
```
- `--stagenet` changes the wallet network and protocol address

The command refuses a duration of zero, malformed names and invalid owner points before constructing the transaction.

## Output

After a successful broadcast, the command prints JSON:

```json
{
  "network": "mainnet",
  "name": "alice",
  "owner": "5866666666666666666666666666666666666666666666666666666666666666",
  "years": 1,
  "amount_atomic": 10000000000,
  "tx_hash_list": [
    "<transaction id>"
  ]
}
```

The transaction ID is the useful receipt. Once the transaction is mined, an indeXer can find the payment and read the name and owner key directly from `tx_extra`.

Broadcasting is not the same as claiming. XNS chronology is determined by mined blocks, not by the time a wallet submitted a transaction. A claim becomes effective only when it appears in the canonical Monero chain and the indeXer replays it.

## Checking the result

Resolve the name through an indeXer:

```sh
./xns lookup \
  --indexer https://your-indexer.example \
  alice
```

The indeXer URL is always explicit. XNS does not choose a resolver on the user's behalf.

A mined and active name returns its owner key, expiration height, remaining blocks and every transaction that established or renewed the current ownership:

```json
{
  "found": true,
  "name": "alice",
  "owner_key": "5866666666666666666666666666666666666666666666666666666666666666",
  "expiration_height": 4000000,
  "remaining_blocks": 200000,
  "finalized": true,
  "source_txids": [
    "<claim transaction id>"
  ]
}
```

An unclaimed or expired name still returns a successful response, but with `"found": false`.

The indeXer serves a claim after it has been mined. Its persistent state waits for ten confirmations, so a very recent claim may disappear during a Monero reorganization and reappear according to the replacement chain.
