# Before you claim

An XNS claim joins two different keys that serve two different purposes. The Monero wallet pays for the claim, while the owner key becomes the identity behind the name. They do not have to be related.

The paying wallet is an ordinary Monero wallet with enough unlocked balance for the chosen duration and the transaction fee. The owner is a 32-byte compressed Ed25519 public key, written as 64 hexadecimal characters. XNS records only this public key. Its private key never enters the claim transaction and should never be given to an indeXer, a Monero node or the XNS claim tool.

Keep the owner private key as carefully as the identity itself. XNS has no recovery operation and no authority capable of replacing an owner key. Losing it does not immediately release the name, but it makes the active name useless to you and prevents renewal under the same identity. The name becomes claimable by somebody else only after its expiration height is reached.

## Requirements

Clone and build XNS:

```sh
git clone https://github.com/exilens/xns
cd xns
go build -o xns ./cmd/xns
```

The claim command also requires `monero-wallet-rpc` in `PATH`. Use the version distributed with Monero. The current transaction construction was developed and tested against Monero `v0.18.5.0`.

You will need:

- the path of the full Monero wallet file
- its password, if it has one
- an accessible Monero daemon RPC URL
- the XNS name
- the owner's raw Ed25519 public key
- the number of years to claim

The wallet file should not already be open in another Monero wallet process. Monero locks an open wallet, so XNS cannot use the same file concurrently with the GUI, CLI or another `monero-wallet-rpc`.

The supplied daemon is started as a trusted daemon by the current claim implementation. Use your own node or one whose operator you trust. A remote node cannot obtain the wallet spend key, but it observes wallet synchronization and transaction activity.

## Name rules

An XNS name is between 1 and 32 bytes. It may contain lowercase letters, numbers and hyphens:

```text
a-z
0-9
-
```

The first and last character must be a letter or number. Uppercase letters, underscores, dots, spaces and Unicode are invalid.

Valid examples:

```text
xns
alice
alice-1
1984
```

Invalid examples:

```text
Alice
alice_name
-alice
alice-
alice.xns
```

There are deliberately no dots in the base registry. XNS assigns names, not hierarchical domains. A future naming layer may define subdomains without changing the ownership of the base name.

## Duration and cost

One XNS year is `262800` Monero blocks and costs exactly `0.01 XMR`. The claim tool accepts duration as whole years and derives the burn amount itself:

```text
1 year  = 0.01 XMR
2 years = 0.02 XMR
5 years = 0.05 XMR
```

The transaction fee is separate from this amount. The protocol wallet cannot spend its outputs, so the claim payment is permanently burned.

Before broadcasting, resolve the name through an indeXer and inspect its owner and expiration. This cannot reserve the name for you; another valid claim may be mined first. It does prevent the simpler mistake of knowingly burning XMR on a name that is already active under another owner.

Continue with [Claiming a name](/docs/claiming-a-name/claiming-a-name) when the wallet, name and owner key are ready.
