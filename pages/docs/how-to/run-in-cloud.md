---
sidebar_position: 4
---

# Running in the cloud

An advantage of running on kubernetes is that running in the cloud is as easy as running locally. Create a cluster in the cloud and download the cluster configuration file. If you make this your default configuration then nothing else special needs to be done. Scaling lightning cli and library will use the default cluster.

To see what would be the default config (and context) on your machine run:

```shell
kubectl config get-contexts
```

To specify a specific kubernetes config file in cli:

```shell
./scaling-lightning --kubeconfig ~/cloud-k8s.yml create -f network.yaml
./scaling-lightning --kubeconfig ~/cloud-k8s.yml walletbalance -n bitcoind
```

or in code:

```go
network := sl.NewSLNetwork("network.yaml", "cloud-k8s.yml", sl.Regtest)
err := network.CreateAndStart()
if err != nil {
    log.Fatal().Err(err).Msg("Problem starting network")
}
```

