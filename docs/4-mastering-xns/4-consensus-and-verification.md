# Consensus and verification

XNS has no consensus network of its own. It inherits Monero's canonical chain and adds deterministic interpretation rules. Agreement requires both parts: two implementations reading different Monero chains may disagree temporarily, and two implementations applying different XNS rules will disagree even on the same chain.

The protocol boundary is therefore straightforward:

1. Monero decides which transactions exist, their amounts, their block heights and the canonical block hashes.
2. XNS decides which incoming protocol-wallet transfers contain valid claims and how those claims change the registry.

An indeXer is correct only when it verifies the first part and applies the second part exactly.

## Independent reconstruction

A verifier needs:

- a Monero node
- the protocol view secret
- the invalid protocol spend public key
- the network prefix and restore height
- the XNS payload and registry rules

From the key material, the verifier derives the protocol address rather than trusting a copied string. From the restore height, it scans all incoming outputs. For each transfer it fetches the canonical block and full transaction, decodes the amount and `tx_extra`, orders the actions, and replays the registry.

SQLite, public indeXers and the official XNS program are conveniences. None of them contain information that cannot be recovered from Monero and the protocol constants.

## What an indeXer can do

A dishonest indeXer can return false JSON, omit a name, report an old height or refuse service. This is no different from a dishonest DNS resolver or block explorer lying to its own client.

It cannot alter the Monero transaction, create a valid owner signature, make another honest indeXer reach the same false result or erase the source transaction from nodes following the canonical chain.

Clients choose how much verification they need. A casual lookup may use one public indeXer. An application protecting valuable identity can compare several. A fully independent application can run an indeXer beside its own Monero node.

The `source_txids` response field exists to connect convenient resolution back to the transactions that justify it.

## Privacy

Monero hides ordinary transaction relationships from the public, but the XNS protocol view key is public by design. Anyone can identify outputs received by the protocol wallet, recover their amounts and read the explicit name and owner key in `tx_extra`.

An XNS claim is therefore public protocol data. Monero still protects unrelated wallet history and does not place the sender address in the transaction, but claimants should not mistake XNS for a private name registry.

The node used by the claim tool also sees wallet synchronization and transaction submission. Operating a local node reduces that exposure.

## Protocol changes

The current rules describe one fixed protocol. Changing the wallet keys, restore heights, payload encoding, amount ratio, ordering function or registry transitions would change the state derived from Monero.

Such changes cannot be treated as an indeXer preference. They would define another protocol unless every relevant historical rule and activation boundary were specified unambiguously.

This is why the XNS core is deliberately small. The fewer rules that enter canonical interpretation, the easier it is for independent implementations to remain identical and for a reviewer to reproduce the result.
