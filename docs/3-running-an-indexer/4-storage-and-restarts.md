# Storage and restarts

The indeXer keeps two versions of state: a current in-memory view and a durable SQLite snapshot.

The in-memory view includes the current chain tip above the durable SQLite boundary. This is the state served by the API, so a claim becomes visible after its first confirmation.

The SQLite snapshot stops ten blocks behind the wallet. It contains only transactions at or below that durable height. The newest blocks remain in memory because short Monero reorganizations are expected to happen near the chain tip.

This arrangement gives new claims fast visibility without pretending that a freshly mined block is permanent.

## Files

The data directory contains:

```text
xns.db
xns.db-wal
xns.db-shm
wallet/
  xns_protocol_view
  xns_protocol_view.keys
wallet-rpc.log
wallet-rpc.log.stdout
```

The `wallet` directory is a Monero view-wallet cache. It holds the protocol view key and scanned wallet data, but no valid spend key.

`xns.db` stores the durable registry, claim events, source transaction IDs, protocol address, durable height, derived-state version and the hash of the block at that height.

The SQLite database uses write-ahead logging, full synchronization and foreign-key enforcement. Full rebuilds replace the previous snapshot inside one SQL transaction. Incremental scans promote newly durable events inside one SQL transaction. A crash before commit leaves the earlier complete state intact rather than a half-written registry.

## Restart behavior

At startup, the indeXer loads SQLite but does not immediately mark itself ready. It checks that the stored protocol address belongs to the selected network and asks the Monero node for the block hash at the durable height.

If the hash matches, the snapshot is a valid starting cache. If the hash changed, the stored state is not served. If the node cannot be queried, the stored state is also withheld rather than being presented without verification.

The indeXer then refreshes the protocol wallet. If the stored durable state is valid, it replays only transfers above the durable height and serves the resulting view. If the stored state is missing, stale or from an older derived-state version, it reconstructs the durable registry by replaying protocol wallet transfers. Only a successful scan makes the API ready.

Transactions with fewer than ten confirmations are intentionally absent from SQLite. If the process stops, they are recovered from the Monero wallet and canonical chain during the next startup.

## Backups

The database is reproducible and therefore not irreplaceable. Deleting `xns.db` causes the indeXer to create a new one from protocol wallet transfers on its next successful scan.

The wallet cache is also reproducible from the protocol address, public protocol constants and restore height. Deleting it makes startup slower because `monero-wallet-rpc` must scan that range again.

Backing up the data directory can reduce recovery time, but it is not a backup of XNS itself. Monero remains the authoritative record.

Do not copy a live SQLite database by taking only `xns.db` while WAL mode is active. Stop the indeXer first, or use a SQLite-aware backup method that includes committed WAL state.
