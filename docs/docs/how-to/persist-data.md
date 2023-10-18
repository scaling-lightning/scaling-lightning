---
sidebar_position: 2
---

# Persist node data

By default nodes in scaling lightning are ephemereal and will lost their data between network or node restarts and upgrades.

To persist data specify volume in the helmfile configuration:

```yaml title="network.yaml"
repositories:
  - name: scalinglightning
    url: https://charts.scalinglightning.com
releases:
  - name: bitcoind
    namespace: sl
    chart: scalinglightning/bitcoind
    values:
      - volume:
          size: "1Gi" # Size of volume in kubernetes notation. Here 1 Gibibyte (1,073,741,824 bytes) is specified.
  - name: alice
    namespace: sl
    chart: scalinglightning/cln
    values:
      - volume:
          size: "1Gi" # Size of volume in kubernetes notation. Here 1 Gibibyte (1,073,741,824 bytes) is specified.
  - name: bob
    namespace: sl
    chart: scalinglightning/lnd
    values:
      - volume:
          size: "1Gi" # Size of volume in kubernetes notation. Here 1 Gibibyte (1,073,741,824 bytes) is specified.
```
