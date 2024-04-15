---
sidebar_position: 4
---

# Configuration options

Configuration is provided via a [helmfile](https://helmfile.readthedocs.io/en/latest/).

## Simplest possible helmfile

```yaml title="network.yaml"
repositories:
  - name: scalinglightning
    url: https://charts.scalinglightning.com
releases:
  - name: bitcoind
    namespace: sl
    chart: scalinglightning/bitcoind
  - name: cln
    namespace: sl
    chart: scalinglightning/cln
  - name: lnd
    namespace: sl
    chart: scalinglightning/lnd
```

## More complex annotated helmfile example

```yaml title="network.yaml"
repositories:
  - name: scalinglightning
    url: https://charts.scalinglightning.com
releases:
  - name: bitcoind # bitcoind node must be called bitcoind and there must be one of these nodes
    namespace: sl # Must be sl namespace
    chart: scalinglightning/bitcoind # Using public helm chart repo added above
    version: 7.7.7 # Specifies version of chart. Should be the same version as cli and library
    values:
      - clientImage:
          tag: 7.7.7 # Specifies version of sidecar container. Should be the same as chart version and cli/library version
      - volume: # Optional: if specified will create a volume and data will be persisted between restarts and upgrades
          size: "1Gi" # Size of volume in kubernetes notation. Here 1 Gibibyte (1,073,741,824 bytes) is specified.
      - autoGen: false # Optional, default value is true. Sets if auto mining is enabled which will mine a block every 10 seconds.
      - rpcEntryPoint: endpoint37 # Allocate endpoint37 to bitcoind's rpc interface. Allocating an endpoint gives access to outside the cluster.
      - zmqPubBlockEntryPoint: endpoint38 # Allocate endpoint38 to bitcoind's zmq block interface. Gives external access.
      - zmqPubTxEntryPoint: endpoint39 # Allotcate endpoint39 to bitcoind's zmq tx interface. Gives external access.
  - name: alice # Friendly name, can be anything. Used to refer to node with the cli and library
    namespace: sl # Namespace has to be sl
    chart: ../../charts/cln # Use a chart from the local filesystem (usually used when developing the scaling lightning project itself)
    values:
      - gRPCEntryPoint: endpoint1
      - image:
          tag: v23.05.1 # Version of CLN docker image to use
      - clientImage:
          repository: scalingln/cln-client # Use a specific image for sidecar container (usually used when developing the scaling lightning project itself)
          pullPolicy: IfNotPresent # K8s Pull Policy for sidecar image. IfNotPresent helps locate updated local images.
      - volume: # Optional: if specified will create a volume and data will be persisted between restarts and upgrades
          size: "1Gi" # Size of volume in kubernetes notation. Here 1 Gibibyte (1,073,741,824 bytes) is specified.
  - name: bob # Friendly name, can be anything. Used to refer to node with the cli and library
    namespace: sl # Namespace has to be sl
    chart: scalinglightning/lnd # LND chart from public helm chart repo specified above
    version: 7.7.7 # Version of chart. Match with version of cli / library
    values:
      - image:
          tag: v0.17.0-beta.rc3 # Version of LND image to use
      - clientImage:
          tag: 7.7.7 # Version of sidecar container to use
      - gRPCEntryPoint: endpoint2
      - volume: # Optional: if specified will create a volume and data will be persisted between restarts and upgrades
          size: "1Gi" # Size of volume in kubernetes notation. Here 1 Gibibyte (1,073,741,824 bytes) is specified.
```
