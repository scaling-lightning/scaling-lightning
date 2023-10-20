---
sidebar_position: 2
---

# Getting started

## Prerequisites:

* Kubernetes.

    * If you are developing locally you can use Docker Desktop and enable
    Kubernetes in the dashboard.
    * Alternatively minikube works as an alternative to Docker Desktop. Please use `minikube tunnel` to enable traefik to get an "external" ip which the library and cli requires to communicate in to the sidecar clients.
    * SL has also been tested on Digital Ocean's hosted K8s cluster
    * Please let us know if you have run SL on a different cluster distribution such as Kind, K3s K0s or any other cloud provider

* Helm 3 and Helmfile.

  Mac OS

      brew install helm helmfile

  Windows

      scoop install helm helmfile

  For Linux check your distros package manager but you may need to download the binaries for helm and helmfile.

* Helm Diff:

    helm plugin install https://github.com/databus23/helm-diff

	> **_NOTE:_** On Windows the plugin install does not complete correctly and you need to download the binary manually from https://github.com/databus23/helm-diff/releases . Unzip the diff.exe file and put it in the _helm/plugins/helm-diff/bin_ folder (the _bin_ folder has to be created). You can find the folder by running _"helm env HELM_DATA_HOME"_

* Traefik:

    helm repo add traefik https://traefik.github.io/charts
    helm repo update
    helm install traefik traefik/traefik -n sl-traefik --create-namespace -f https://raw.githubusercontent.com/scaling-lightning/scaling-lightning/main/charts/traefik-values.yml

## Installation

Download binary for your system from [Releases](https://github.com/scaling-lightning/scaling-lightning/releases)

    # untar to get binary
    tar -xzf scaling-lightning-[version]-[os]-[architecture].tar.gz

    # Mac OS only - mark file as safe so it will run
    xattr -dr com.apple.quarantine scaling-lightning

    # run - should print CLI help
    ./scaling-lightning

## Starting a Network

To spin up an example network with 2 cln nodes and 4 lnd nodes, run:

    # Download example helmfile which defines the nodes you want in your network.
    wget https://raw.githubusercontent.com/scaling-lightning/scaling-lightning/main/examples/helmfiles/public.yaml

    # Create and start the network. Scaling lightning will use your currently defined default k8s cluster
    # as specified in kubectl kubectl config get-contexts
    ./scaling-lightning create -f public.yaml

To destroy the network run:

    ./scaling-lightning destroy

## Example CLI Commands

    # list nodes on the network (names were taken from the helmfile)
    ./scaling-lightning list

    # get wallet balance of node named bitcoind
    ./scaling-lightning walletbalance -n bitcoind

    # get wallet balance of node named lnd2
    ./scaling-lightning walletbalance -n lnd2

    # send on-chain 1 million satoshis from bitcoind to cln1
    ./scaling-lightning send -f bitcoind -t cln1 -a 1000000

    # get the pubkey of a node named lnd1
    ./scaling-lightning pubkey -n lnd1

    # peer lnd1 and cln1 from lnd1
    ./scaling-lightning connectpeer -f lnd1 -t cln1

    # open channel between cln1 and lnd1 with a local balance on cln1 of 70k satoshis
    ./scaling-lightning openchannel -f cln1 -t lnd1 -a 70000

    # have bitcoind generate some blocks and pay itself the block reward
    ./scaling-lightning generate -n bitcoind

## Run the above from code instead of CLI

See [examples/go/example_test.go](https://github.com/scaling-lightning/scaling-lightning/blob/main/examples/go/example_test.go). This test takes around 3 minutes to pass on an M1 Macbook Pro so you may need to adjust your test runner's default timeout.

Example go test command with extra timeout:

    go test -run ^TestMain$ github.com/scaling-lightning/scaling-lightning/examples/go -count=1 -v -timeout=15m

## Helpful Kubernetes commands

    # list pods
    kubectl -n sl get pods

    # describe cln1 pod in more detail
    kubectl -n sl describe pod cln1-0

    # view logs of lnd1 node
    kubectl -n sl logs -f lnd1-0

    # view logs of a crashed bitcoind pod
    kubectl -n sl logs -previous bitcoind-0

    # view logs of lnd1's scaling lightning sidecar client (it handles our api requests and forwards them to the node)
    kubectl -n sl logs -f -c lnd-client lnd1-0

    # same for cln and bitcoind
    kubectl -n sl logs -f -c cln-client cln1-0
    kubectl -n sl logs -f -c bitcoind-client bitcoind-0

    # get shell into lnd1
    kubectl -n sl exec -it lnd1-0 -- bash

    # view loadbalancer public ip from traefik
    kubectl -n sl-traefik get services

    # destroy all scaling lightning nodes
    kubectl delete namespace sl

    # uninstall traefik
    kubectl delete namespace sl-traefik

    # uninstall traefik alternative
    helm uninstall traefik -n sl-traefik


Note that the above commands assume you are using the default kubeconfig and context. You would need to add `--kubeconfig path/to/file.yml` or `--context mycluster` to all of the above commands if you wanted to look at a different cluster.
