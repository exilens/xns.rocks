# Creating an I2P service

Generate an Ed25519 private key:

```sh
openssl genpkey -algorithm Ed25519 -out service.pem
```

Create the i2pd destination:

```sh
./xns-i2p service service.pem service
```

The command creates:

```text
service/
├── hostname
└── private.dat
```

It also prints the XNS owner key and extended I2P address.

## Install the key for i2pd

Place the private destination under the i2pd data directory:

```sh
sudo install -d -o i2pd -g i2pd -m 700 /var/lib/i2pd/xns
sudo install -o i2pd -g i2pd -m 600 \
  service/private.dat /var/lib/i2pd/xns/private.dat
```

Create `/etc/i2pd/tunnels.d/xns.conf`:

```ini
[xns]
type = server
host = 127.0.0.1
port = 8080
keys = xns/private.dat
signaturetype = 11
i2cp.leaseSetType = 5
```

The `keys` value is relative to the i2pd data directory. With the standard service layout, `keys = xns/private.dat` refers to `/var/lib/i2pd/xns/private.dat`.

Do not write the absolute path there. i2pd passes the value through its data-directory path handling, so an absolute-looking value may be prefixed with `/var/lib/i2pd` and fail to load. If i2pd cannot open the requested file, it may create or use another destination instead of the XNS identity you intended.

`signaturetype = 11` matches the RedDSA destination written by `xns-i2p`.

`type = server` creates a raw TCP tunnel. This is appropriate for XNS because i2pd passes the HTTP request through unchanged, including:

```http
Host: example.xns
```

The `host` and `port` fields are the local application or reverse-proxy listener. In this example, i2pd forwards incoming I2P streams to `127.0.0.1:8080`.

## Encrypted LeaseSet2 is required

The line:

```ini
i2cp.leaseSetType = 5
```

tells i2pd to publish an encrypted LeaseSet2. Without it, the service may have an ordinary I2P destination, but XNS Resolver cannot locate it through the extended address derived from the owner key.

Type `5` does not enable per-client authentication and does not require visitors to know another secret. It selects the encrypted LeaseSet2 publication mechanism and its blinded destination lookup.

Restart i2pd:

```sh
sudo systemctl restart i2pd
```

## Verify the loaded identity

Open the i2pd console, normally:

```text
http://127.0.0.1:7070/
```

The **I2P tunnels** page should contain an `xns` entry under **Server Tunnels**. The address shown there is the ordinary 52-character destination hash.

Open **Local Destinations**, select that destination, and find **Encrypted B33 address**. It must match:

```sh
./xns-i2p inspect service
```

If the addresses differ, i2pd did not load the intended `private.dat`. Check the `keys` path, file ownership and logs:

```sh
sudo journalctl -u i2pd -n 100 --no-pager
sudo tail -n 100 /var/log/i2pd/i2pd.log
```

Also confirm that the local target is listening:

```sh
ss -ltn | grep ':8080'
```

Once the destination and extended address agree, continue with [serving XNS hostnames](/docs/xns-and-i2p/serving-xns-hostnames).
