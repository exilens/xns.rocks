# Creating an onion service

For a new service, generate an ordinary Ed25519 private key first:

```sh
openssl genpkey -algorithm Ed25519 -out service.pem
```

Then create the files Tor expects:

```sh
./xns-onion service service.pem hidden_service
```

The command prints the XNS owner key and onion address and creates:

```text
hidden_service/
├── hostname
├── hs_ed25519_public_key
└── hs_ed25519_secret_key
```

Choose a permanent private location for this directory and make it owned by the account that runs Tor. Tor refuses unsafe onion-service directory permissions. The directory should remain `0700`, and the key files should remain `0600`.

Add the service to `torrc`:

```text
HiddenServiceDir /path/to/hidden_service
HiddenServicePort 80 127.0.0.1:8033
```

`HiddenServicePort` maps the public onion-service port to a local listener. In this example, a request to port 80 of the onion address is delivered to `127.0.0.1:8033`.

Restart or reload Tor, then verify the installed identity:

```sh
./xns-onion inspect /path/to/hidden_service
```

The output must match the values printed when the directory was created. It should also agree with:

```sh
cat /path/to/hidden_service/hostname
```

## Existing onion service

If the service already exists, do not generate another key. Inspect its current directory:

```sh
./xns-onion inspect /path/to/existing_hidden_service
```

This recovers the public owner key and verifies that the Tor files belong together. The existing onion address remains unchanged.

This path is provided for services that already have users, links or reputation attached to their onion identity. For a new XNS-first service, starting from the OpenSSL seed is preferable because the seed remains available as the original portable Ed25519 private key.

## What to preserve

Back up the PEM and the onion-service directory privately before claiming the name. The PEM can regenerate the Tor files. The Tor expanded secret can continue operating the onion service, but it cannot be converted back into the original seed.

Losing all private material does not remove the XNS claim. It leaves the active name pointing to an identity that nobody can operate, and XNS has no owner replacement operation.

Once the service identity is settled, continue with [serving XNS hostnames](/docs/xns-and-tor/serving-xns-hostnames).
