There are a few tools, scripts and commands that we use to start investigating a problem while working on Fury. This document collects a few of those.

# Tools worth knowing about.
## Fury specific tools
* [Block explorer](https://explorer.fury.trading/) ([repo](https://github.com/elysiumstation/explorer)), which contains [an API for decoding blocks](https://github.com/elysiumstation/explorer#api).
* [`furystream`](https://github.com/elysiumstation/furytools)

## API Specific tools
### GraphQL
* [`GraphQURL`](https://github.com/hasura/graphqurl) is a curl like CLI tool that supports subscriptions
* [`GraphQL Playground`](https://github.com/prisma-labs/graphql-playground) is served by Fury nodes serving a GraphQL API
* Simplest of all: CURL can be used to post queries.

### GRPC
### REST
# Some quick things

## Decoding a public key
```bash
echo -n 'DrXug9ukpvwdMVAG1c2jOG+TCYVSvqCshL8dr6z+Kd8=' | base64 -d | hexdump -C | cut -b11-58 | tr -dc '[:alnum:]'
```

# Some hypothetical situations

<details>
  <summary><strong>Is [insert network] down?</strong></summary>

  The quickest check is [`stats.fury.trading`](https://stats.fury.trading) ([repo](https://github.com/elysiumstation/stats/)). You should see the network there, and most or all of the stats rows should have a green block, implying it's healthy.

  Stats is a really simply web view of the REST [statistics endpoint](https://docs.testnet.fury.xyz/api/rest/#operation/Statistics), so you could also use curl. Choose a node serving REST from this [`devops repo document`](https://github.com/elysiumstation/devops-infra/blob/master/doc/fury_environments.md) and then curl the statistics endpoint:
  ```bash
  curl https://n04.d.fury.xyz/statistics
  ```

  If this fails, totally it could be that the node itself is down, while the network is fine. If you get a 502 error, then the machine is up, the HTTPS proxy is working, but the Fury node is not running.

  If you want to skip Fury and see if Tendermint is healthy, you can try going straight to Tendermint's RPC port. Choose a node that exposes the Tendermint RPC from this [`devops` repo document](https://github.com/elysiumstation/devops-infra/blob/master/doc/fury_environments.md) and then fetch the status endpoint:
  ```bash
  curl https://n01.d.fury.xyz/tm/status
  ```

 If those two fail, you can try `SSHing` to the machine to see what's up. The [`devops repo`](https://github.com/elysiumstation/devops-infra/blob/master/doc/fury_environments.md) will list all of the nodes, and how you can connect to them to investigate further.
</details
