# Reorgs and consistency

A Monero reorganization replaces one or more recent blocks with another valid chain. Since XNS chronology comes from Monero, the registry must follow the replacement chain as though the orphaned blocks had never existed.

The indeXer does not try to patch individual names after a reorg. Each scan orders all known protocol transfers and replays the registry from the beginning. This is simple enough to audit and avoids retaining conclusions that were valid only under an earlier chain.

## Canonical transaction checks

The protocol wallet reports incoming transfers, but the indeXer does not accept that list blindly. For every transfer it fetches the block at the reported height and verifies that the transaction ID is present in that canonical block.

It then obtains the full transaction from the daemon and checks that its hash and block height match the wallet transfer and that it is no longer in the transaction pool. Only after these checks does it parse `tx_extra` and apply the XNS rules.

At the beginning of a scan, the indeXer records the wallet tip block hash. It requests the same hash again after replay. If the value changed while the scan was running, the result is discarded with `chain changed during scan`. A later scan starts again against one coherent chain.

## Durable block anchor

SQLite stores the hash of the block at its ten-confirmation durable height. On startup, the indeXer asks the selected node for the canonical hash at that height.

A matching hash proves that the boundary block behind the snapshot is still on the chain reported by the node. A mismatch causes the snapshot to be withheld until a fresh wallet refresh and replay produce state for the new chain.

This check works even if the indeXer was offline while the reorganization happened. It does not depend on observing the reorg in real time.

## Deep reorgs

Monero wallet caches handle ordinary reorganizations themselves. If `monero-wallet-rpc` reports that a reorg exceeds its maximum allowed depth, the indeXer treats the cache as disposable.

It stops the wallet RPC, deletes the protocol wallet cache, recreates the view wallet from the protocol restore height and scans again. SQLite is then replaced by the newly reconstructed durable state.

Transient errors such as a refused connection are not treated as deep reorgs. Destroying a healthy wallet cache would turn a temporary network failure into an unnecessary historical rescan.

## Agreement between indeXers

Two indeXers connected to nodes on the same canonical Monero chain and running the same XNS rules should produce the same registry. Their SQLite files, scan times and HTTP availability may differ, but those are local implementation details.

Near the tip, indeXers may briefly disagree because their nodes have different heights or have not yet converged after a reorg. The source transaction IDs and expiration heights make these differences observable. Once both nodes share the same chain and both indeXers complete a scan, their registry state should converge without communicating with each other.
