---
title: Combining SimLN with Scaling Lightning
authors: [max]
tags: [simln, activity-generation]
draft: false
---

SimLN is a new tool to generate payment activity between lightning nodes. In this post I give a walkthough of running it against a Scaling Lightning cluster.

## Scaling Lightning setup

To start off let's create a local cluster with 2 nodes: Alice(CLN) and Bob(LND).

```yaml title="network.yaml"
repositories:
  - name: scalinglightning
    url: https://charts.scalinglightning.com
releases:
  - name: bitcoind
    namespace: sl
    chart: scalinglightning/bitcoind
  - name: alice
    namespace: sl
    chart: scalinglightning/cln
    values:
      - gRPCEntryPoint: endpoint1
  - name: bob
    namespace: sl
    chart: scalinglightning/cln
    values:
      - gRPCEntryPoint: endpoint2
```

Start the network. The scaling-lightning binary can be downloaded from [https://github.com/scaling-lightning/scaling-lightning/releases](https://github.com/scaling-lightning/scaling-lightning/releases).

```shell
scaling-lightning start -f network.yaml
```

{eer up Alice and Carol and open a channel and send some initial funds that can be sent back and forth.

```shell
scaling-lightning send -f bitcoind -t alice -a 1000000
scaling-lightning connectpeer -f alice -t bob
scaling-lightning openchannel -f alice -t bob -a 100000
scaling-lightning createinvoice -n bob -a 50000
scaling-lightning payinvoice -n alice -i <bolt11 invoice>
```

Scaling Lightning cluster should be running with Alice and Bob having a channel between them.

## SimLN setup

[SimLN](https://github.com/bitcoin-dev-project/sim-ln) also requires a config file.

```shell
git clone https://github.com/bitcoin-dev-project/sim-ln
cd sim-ln
```

Then create config file for SimLN to give it the identity of nodes in the Scaling Lightning cluster.

```json title="config.json"
{
  "nodes": [
    {
      "LND": {
        "id": "0248efcfe94e3c451f4995b471ef0707163f279d4681af23727279c9c696004b42",
        "address": "https://localhost:28102",
        "macaroon": "/Users/max/source/sim-ln/auth/bob/admin.macaroon",
        "cert": "/Users/max/source/sim-ln/auth/bob/tls.cert"
      }
    },
    {
      "CLN": {
        "id": "0221b76f4cce7ab42538127022bac6245e999804ab5901813d9337d5cadd6428df",
        "address": "https://localhost:28101",
        "ca_cert": "/Users/max/source/sim-ln/auth/alice/ca.pem",
        "client_cert": "/Users/max/source/sim-ln/auth/alice/client.pem",
        "client_key": "/Users/max/source/sim-ln/auth/alice/client-key.pem"
      }
    }
  ],
  "activity": [
    {
      "source": "0248efcfe94e3c451f4995b471ef0707163f279d4681af23727279c9c696004b42",
      "destination": "0221b76f4cce7ab42538127022bac6245e999804ab5901813d9337d5cadd6428de",
      "interval_secs": 1,
      "amount_msat": 2000
    },
    {
      "source": "0221b76f4cce7ab42538127022bac6245e999804ab5901813d9337d5cadd6428de",
      "destination": "0248efcfe94e3c451f4995b471ef0707163f279d4681af23727279c9c696004b42",
      "interval_secs": 1,
      "amount_msat": 2000
    }
  ]
}
```

To create your version of this config file you need three things: pubkey, address+port of GRPC api and auth files. Scaling Lightning has three commands for that.

```shell
scaling-lightning writeauthfiles -o ~/source/sim-ln/auth --all
scaling-lightning pubkey -n alice
scaling-lightning pubkey -n bob
scaling-lightning connectiondetails --all
```

Finally to run SimLN

```shell
cargo install --path sim-cli/
sim-cli --log-level debug --config config.json
```

We will be following the development of SimLN closely. Scaling Lightning will need it's own activity generator and perhaps that could be SimLN?