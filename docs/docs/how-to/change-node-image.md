---
sidebar_position: 1
---

# Change node image

To change the image used for bitcoind, CLN or LND:

```yaml title="bitcoind"
values:
  - image:
      repository: ruimarinho/bitcoin-core
      tag: 24
```

```yaml title="cln"
values:
  - image:
      repository: elementsproject/lightningd
      tag: v23.05.1
```

```yaml title="lnd"
values:
  - image:
      repository: lightninglabs/lnd
      tag: v0.17.0-beta
```

Full example:

```yaml title="network.yaml"
repositories:
  - name: scalinglightning
    url: https://charts.scalinglightning.com
releases:
  - name: bitcoind
    namespace: sl
    chart: scalinglightning/bitcoind
    values:
      - image:
          repository: ruimarinho/bitcoin-core
          tag: 24
  - name: alice
    namespace: sl
    chart: scalinglightning/cln
    values:
      - image:
          repository: elementsproject/lightningd
          tag: v23.05.1
  - name: bob
    namespace: sl
    chart: scalinglightning/lnd
    values:
      - image:
          repository: lightninglabs/lnd
          tag: v0.17.0-beta
```

## Upgrade an existing node

To update an existing node to a new version you can call the CLI or Library command `create` again with the same `network.yaml` file but with the image value updated to the new value.

> **_NOTE:_** When changing the image repository or tag, a new pod will be created and unless volume was specified in the configuration all data for that node will be lost.

> **_NOTE:_** If changing bitcoind node to a new version without a volume all other lightning nodes will need to be destroyed and recreated as they will reference a blockchain that nolonger exists.
