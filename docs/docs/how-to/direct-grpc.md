---
sidebar_position: 3
---

# Direct (g)RPC connection to nodes

One of the prerequisits for scaling lightning is to have traefik installed on the cluster to enable ingress into the cluster from outside. The command specified in the [Getting Started](/docs/getting-started) references an values file where we specify 40 endpoints to use with nodes in the network. One of those endpoints is used for scaling lightning itself to talk to nodes from the cli or library. This leaves 39 endpoints which can be allocated to bitcoind or lightning nodes to gain direct external acccess to their (g)RPC services.

The available endpoints are named `endpoint1` to `endpoint39`

## Enable external access

The following configuration enables direct external access to bitcoind, cln and lnd:

```yaml
repositories:
  - name: scalinglightning
    url: https://charts.scalinglightning.com
releases:
  - name: bitcoind
    namespace: sl
    chart: scalinglightning/bitcoind
    values:
      - rpcEntryPoint: endpoint37
      - zmqPubBlockEntryPoint: endpoint38
      - zmqPubTxEntryPoint: endpoint39
  - name: alice
    namespace: sl
    chart: scalinglightning/cln
    values:
      - gRPCEntryPoint: endpoint1
  - name: bob
    namespace: sl
    chart: scalinglightning/lnd
    values:
      - gRPCEntryPoint: endpoint2
```

## Get connection details

Connection details can be queried using the CLI or library

```shell
./scaling-lightning connectiondetails --all
```

Example output:

```shell
bitcoind
  type:  rpc
  host: localhost
  port: 28137

  type:  zmq blocks
  host: localhost
  port: 28138

  type:  zmp txs
  host: localhost
  port: 28139

alice
  type:  grpc
  host: localhost
  port: 28101

bob
  type:  grpc
  host: localhost
  port: 28102
```

Or in code:

```go
connectionDetails, err := network.GetConnectionDetails("alice")
assert.NoError(err)
log.Info().Msgf("alice connection host: %v", connectionDetails[0].Host)
log.Info().Msgf("alice connection host: %d", connectionDetails[0].Port)
```
