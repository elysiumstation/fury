# Chain Replay

This is a guide to replay a chain using the backup of an existing chain (e.g. Testnet)

## How it works

A Tendermint Core and Fury Core node store their configuration and data to disk by default at `$HOME/.cometbft` and `$HOME/.fury`. When you start new instances of those nodes using a copy of these directories as their home, Tendermint re-submits (replays) historical blocks/transactions from the genesis height to Fury Core.

## Prerequisites

- [Google Cloud SDK ][gcloud]
- Fury Core Node
- Fury Wallet
- [Tendermint][tendermint]

## Chain backups

Note you need to first authenticate `gcloud`.

You can find backups for the Fury networks stored in Google Cloud Storage, e.g. For Testnet Node 01

```
$ gsutil ls gs://fury-chaindata-n01-testnet/chain_stores
```

## Steps

- Copy backups locally to `<path>`

- Overwrite Fury node wallet with your own development [node wallet][wallet]. 

```
$ cp -rp ~/.fury/node_wallets_dev <path>/.fury
$ cp ~/.fury/nodewalletstore <path>/.fury
```

- Update Fury node configuration

```
$ sed -i 's/\/home\/fury/<path>' <path>/.fury/config.toml
```

- Start Fury and Tendermint using backups

```
$ fury node --root-path=<path>/.fury --stores-enabled=false
$ tendermint node --home=<path>/.tendermint
```


## Tips

The Fury nodes adheres to the Tendermint ABCI contract, therefore breakpoints in the following methods are useful:

```
blockchain/abci/abci.go#BeginBlock
```

## Alternatives

Instead of a backup, which effectively replays the full chain from genesis, you can also use a snapshot of the chain at a given height to bootstrap the Tendermint node. Which only replays blocks/transactions from the given height. This however requires extra tooling.

## References

- https://github.com/tendermint/tendermint/blob/master/docs/introduction/quick-start.md
- https://docs.tendermint.com/master/spec/abci/apps.html
- https://github.com/tendermint/spec/blob/master/spec/abci/README.md
- https://docs.tendermint.com/master/spec/abci/apps.html#state-sync

[wallet]: https://github.com/elysiumstation/fury#configuration
[gcloud]: https://cloud.google.com/sdk/docs/install
[tendermint]: https://github.com/tendermint/tendermint/blob/master/docs/introduction/install.md