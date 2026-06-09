# Inside the claim tool

Monero wallet RPC can create and sign ordinary transfers, but it does not expose a general option for inserting arbitrary bytes into `tx_extra`. XNS works around this without maintaining a patched Monero build.

The full wallet remains the source of truth throughout the process. It is refreshed first, then asked for its primary address, private view key, current wallet state and signed key images. The spend key never leaves the full wallet.

From this information, XNS creates a temporary watch-only copy of the sender wallet under a private system temporary directory, normally `/tmp` on Unix. After the full wallet has refreshed, XNS stores its current wallet cache and copies that cache into the temporary watch wallet. It then imports signed key images and verifies that the watch wallet has the same address, height, balance and unlocked balance as the full wallet.

The copied cache gives the temporary wallet the outputs, block history and transaction state that the full wallet has already discovered. It does not contain the private spend key; that remains only in the full wallet's `.keys` file. The watch wallet therefore does not repeat the full historical scan and still cannot sign or spend anything.

The watch wallet creates an unsigned transfer to the protocol address with relay disabled. XNS decrypts the unsigned transaction set using the sender view key, locates the dummy encrypted-payment-ID nonce created by the wallet, and replaces it with:

```text
02 40 name32 owner_key32
```

The modified transaction set is authenticated again and returned to the full wallet. Only then does the full wallet sign it. The signed transaction is submitted immediately, both wallet RPC processes are stopped, and the temporary directory is removed.

This division is necessary because a full wallet signs normal transfers directly, while a watch-only wallet returns their construction data as an unsigned transaction set. XNS needs that intermediate representation in order to change `tx_extra`.

## Local processes

The command starts two `monero-wallet-rpc` processes bound to `127.0.0.1`. Each receives a free local port selected by the operating system. The ports are not protocol constants and are never exposed as part of the user interface.

The full wallet RPC opens the wallet supplied by `--wallet-file`. The watch wallet RPC uses only the temporary wallet directory. Both are stopped on success or failure, and a process that does not exit after an interrupt is killed after ten seconds.

Because the full wallet file is opened by `monero-wallet-rpc`, it cannot already be locked by another wallet process.

## Fees

The wallet chooses the fee before XNS expands `tx_extra`. The final transaction is therefore slightly larger than the transaction for which the wallet originally calculated the fee.

The claim tool uses Monero wallet priority `2`, the normal priority, rather than automatic priority. This provides reasonable fee-rate headroom for the additional bytes without choosing an elevated fee level. The fee is not recalculated after patching, because doing so would require reconstructing and rebalancing the transaction rather than making the narrow `tx_extra` replacement.

The resulting transaction has been tested on stagenet with stock Monero wallet software. This mechanism still depends on the serialized unsigned transaction-set format expected by the implementation, so compatibility should be retested when changing the supported Monero version.

## Common failures

`wallet-rpc exited early` usually means the wallet path is wrong, the wallet is already open, the password is wrong, or the selected wallet network does not match `--stagenet`.

`could not find expected default tx_extra` means the unsigned transaction did not contain the exact dummy nonce that the patcher understands. This is generally a Monero version or transaction-construction compatibility problem.

`temporary watch wallet state does not match full wallet` means the copied cache and imported key images did not reproduce the full wallet's height and balance. XNS refuses to construct a transaction from mismatched state.

An RPC or broadcast error does not always prove that no transaction was submitted. If failure occurs during submission, inspect the paying wallet and search the returned or logged transaction hash before repeating the claim. Sending the same claim twice with the same owner renews it, but it also burns another `0.01 XMR` per year.
