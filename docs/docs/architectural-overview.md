---
sidebar_position: 2
---

# Architectural overview

Diagram of SL

Definition for which nodes to run and their configuration options is all done via helm. Helm charts currently exist for bitcoind, lnd and cln and can be found by version at https://charts.scalinglightning.com/. Helmfile can be configured to use the SL helm repository so there is no need to add it to your local helm installtion.

All components of scaling are versioned with the same tag. So if using version x of the API or CLI, it is recommended to use version x of the helm charts which in turn will reference version x of the sidecar client containers.

![SL Architecture](./img/SLArchitecture.jpeg)
