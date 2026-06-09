# Why XNS?

The Internet was supposed to make every man a publisher, a merchant and a sovereign operator of his own services. In many ways it did. You can write your own software, run it on your own machine, encrypt every connection and speak to the entire world without owning a printing press or asking a broadcaster for airtime. Yet, after doing all of that yourself, the name through which people know and reach you is still rented from institutions.

It is a strange weakness to accept. You may own the server and every byte on it, but your identity remains attached to an entry that someone else can suspend, redirect or take away. Moving to another hosting company or hosting from your own home changes nothing. As long as your name exists by permission, the most important entrance to your work does not belong to you.

A name is not merely a shortcut to an address. People remember it, return to it and attach reputation to it. It tells them that the person they found today is the same person they trusted yesterday. Whoever controls that continuity holds something much more valuable than a technical record. They hold a part of your identity.

Cryptography has already shown that identity does not require an institution. If a message is signed by your key, nobody needs to ask a registrar who wrote it. Tor and I2P go further and make cryptographic identities reachable. The authority comes from possession of the private key, not from an account in somebody's database.

The problem is that cryptographic identities are made for machines. A public key may be pure and undeniable, but it is not something people can comfortably remember, speak or type. Traditional names are practical but unfree; cryptographic identities are free but impractical. What is missing is a simple name that belongs to a key without first belonging to an authority.

That is the problem XNS solves.

## Why another name system?

There have been attempts to escape DNS before, but many of them began by inventing a new blockchain and a new currency for the single purpose of maintaining names. Even when the engineering is honest, this creates a much larger problem than the one it was meant to solve. The new chain must gain miners, users, liquidity, security and enough economic weight to survive. Early participants receive an advantage over everyone who comes later, speculation becomes inseparable from usage, and a simple registry arrives with an entire artificial economy attached to it.

ENS took the opposite path and built its name system on Ethereum, a network that already had users and a reason to exist beyond names. It succeeded in its own field because it did not ask the world to adopt a blockchain merely to use a name. Still, ENS is made for the Ethereum world. It was never meant to be a general replacement for DNS or a name system for men who want an identity outside the permission of states, corporations and platforms.

XNS takes the useful lesson and removes everything unnecessary. Monero already exists with the security, economic weight and culture of private, permissionless exchange needed for irreversible public facts. XNS does not introduce another chain, token, premine, treasury or governance system on top of it. It only defines how a particular kind of Monero transaction is to be read.

Monero is therefore not a settlement layer underneath some separate XNS network. It is the history, the clock and the source of truth, while XNS is only the logic required to find names inside it.

## Ownership without a registrar

An XNS name belongs to an Ed25519 public key. Whoever holds the corresponding private key can prove ownership, and nobody who lacks it can pretend to be the owner. The relationship does not depend on an account, an administrator or a company continuing to recognize it.

Names are claimed for a period of time and renewed by claiming them again with the same owner key. If a name expires, it becomes free. While it remains active, a claim made by another key has no effect. These rules are deliberately small because every extra operation is another opportunity to introduce authority, ambiguity or abuse.

Most importantly, XNS has no transfer operation. This may sound limiting until one considers what selling a cryptographic identity actually means. The seller would have to give the buyer a private key, yet nothing can prove that he did not keep a copy. The buyer could never know whether the identity was exclusively his, and the seller could return at any time. Since names cannot be transferred safely, collecting them for resale does not produce a trustworthy market. XNS does not fight squatting with committees, trademarks or auctions; it removes the reason to squat in the first place.

## Monero remembers

An indeXer reads the XNS protocol wallet, follows its transactions from the beginning and applies the protocol rules in chronological order. It does not register names or decide who owns them. Its database is merely a cache of facts that already exist on Monero.

This distinction is what makes the system independent of any particular indeXer. A dishonest one can lie to its own users, just as any server can return false information, but it cannot change the underlying claim or prevent somebody else from reading it correctly. If every public indeXer disappeared, new ones could reconstruct the same registry from Monero. Nothing essential is held in their databases.

XNS is uncensorable because a valid claim cannot be removed from Monero by an indeXer. It is unstoppable because there is no central XNS service whose disappearance would destroy the registry. It is unseizable because an active name cannot be assigned to another owner without the rules rejecting it. And it is unsquattable because a name cannot become a reliable commodity for resale.

The exile should not have to rent his identity from the same institutions he left behind. He should not have to trade one authority for a newly invented blockchain and its early rulers either. He needs a key that is his, a name for that key, and a history that nobody can rewrite for him.

Monero provides the history. XNS gives the key a name.

See [How does it work?](/docs/introduction/how-does-it-work) for the protocol itself.
