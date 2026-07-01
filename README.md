# xns.rocks

XNS website and documentation.

```sh
go build -o xns-rocks ./cmd/xns.rocks
./xns-rocks \
  --listen <listen-address> \
  --node <mainnet-node-url> \
  --data-dir <data-directory>
```

The listen address must be specified explicitly.
