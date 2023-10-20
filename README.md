![ScalingLN Twitter Banner](https://github.com/ohenrik/scaling-lightning/assets/647617/8511c586-7549-4e2b-ad6d-bf87419a624c)

# Scaling Lightning - A Testing Toolkit for the Lightning Network

This initiative aims to build a testing toolkit for the Lightning Network protocol, its implementations, and
applications that depend on the Lightning Network.

[Project use cases](https://scalinglightning.com/docs/project-use-cases)

The goal is to collaborate as an industry to help scale the Lightning Network and the applications that depend on it.

[Current Roadmap](https://scalinglightning.com/blog/2023/08/26/roadmap)

[Full Documentation](https://scalinglightning.com/docs)

## Getting started

> **_The project is still far from complete_**.

The following is a quick start guide to get something running. Please refer to the [documentation](https://scalinglightning.com/docs) for more info.

### Prerequisites

- Kubernetes.

  - If you are developing locally you can use Docker Desktop and enable
    Kubernetes in the dashboard.
  - Alternatively minikube works as an alternative to Docker Desktop. Please use `minikube tunnel` to enable traefik to get an "external" ip which the library and cli requires to communicate in to the sidecar clients.
  - SL has also been tested on Digital Ocean's hosted K8s cluster
  - Please let us know if you have run SL on a different cluster distribution such as Kind, K3s K0s or any other cloud provider

- Helm 3 and Helmfile.

  Mac OS

      brew install helm helmfile

  Windows

      scoop install helm helmfile

  For Linux check your distros package manager but you may need to download the binaries for helm and helmfile.

- Helm Diff:

  helm plugin install https://github.com/databus23/helm-diff

  > **_NOTE:_** On Windows the plugin install does not complete correctly and you need to download the binary manually from https://github.com/databus23/helm-diff/releases . Unzip the diff.exe file and put it in the _helm/plugins/helm-diff/bin_ folder (the _bin_ folder has to be created). You can find the folder by running _"helm env HELM_DATA_HOME"_

- Traefik:

  helm repo add traefik https://traefik.github.io/charts
  helm repo update
  helm install traefik traefik/traefik -n sl-traefik --create-namespace -f https://raw.githubusercontent.com/scaling-lightning/scaling-lightning/main/charts/traefik-values.yml

### Installation

Download binary for your system from [Releases](https://github.com/scaling-lightning/scaling-lightning/releases)

    # untar to get binary
    tar -xzf scaling-lightning-[version]-[os]-[architecture].tar.gz

    # Mac OS only - mark file as safe so it will run
    xattr -dr com.apple.quarantine scaling-lightning

    # run - should print CLI help
    ./scaling-lightning

### Starting a Network

To spin up an example network with 2 cln nodes and 4 lnd nodes, run:

    # Download example helmfile which defines the nodes you want in your network.
    wget https://raw.githubusercontent.com/scaling-lightning/scaling-lightning/main/examples/helmfiles/public.yaml

    # Create and start the network. Scaling lightning will use your currently defined default k8s cluster
    # as specified in kubectl kubectl config get-contexts
    ./scaling-lightning create -f public.yaml

To destroy the network run:

    ./scaling-lightning destroy

### Example CLI Commands

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

### Run the above from code instead of CLI

See [examples/go/example_test.go](examples/go/example_test.go). This test takes around 3 minutes to pass on an M1 Macbook Pro so you may need to adjust your test runner's default timeout.

Example go test command with extra timeout:

    go test -run ^TestMain$ github.com/scaling-lightning/scaling-lightning/examples/go -count=1 -v -timeout=15m

### Helpful Kubernetes commands

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

### Your own configuration

This project is still in its infancy, so we don't have a lot of configuration options yet. Please take a look in the [examples](/examples/helmfiles) directory for examples of different networks.

## Why is this important?

Currently, there are unknowns and untested assumptions about how the Lightning Network and its applications will react
to shocks in transaction volume, channels, nodes, gossip messages, etc.

Having a set of tools and a signet Lightning Network will help:

- Developers test their applications.
- Researchers verify their assumptions.
- Operators test their infrastructure.
- Novices learn how the Lightning Network and various applications work in a somewhat realistic environment without
  risking real coins.

## How will it work?

We are still in the early stages of planning, but the first tool we are building will be a tool to quickly generate one
or more Lightning Nodes. These nodes can connect either to a public signet Lightning Network or a private Regtest
Lightning Network for any combination of LN implementations (CLN, LND, LDK, Acinq etc.).

Other tools, made specifically for testing isolated parts of the protocol, are also relevant. These can help developers
and researchers test their assumptions in an isolated environment. An example of this is
[The Million Channels Project](https://github.com/rustyrussell/million-channels-project-data) developed by Rusty Russell
to test gossip.

### How is this different from Polar?

While Polar is an excellent project it is fundamentally different from what we want to achieve with the Scaling Lightning initiative. Polar is a desktop application meant to manually build a network through a drag and drop interface, which is great for novices and simple projects. However, it's not something that is suitable for a developer or testing environment at a startup or as a researcher.

## Project milestones

This is an outline of the project's milestones. We will further detail these milestones using the features of GitHub's milestones and project management tools:

- [x] Create a basic kubernetes setup for running a Lightning Network
- [x] Define file format to describe initial state of the network
- [ ] Create a helm chart and sidecar client for each Lightning Network implementation
  - [x] LND
  - [x] CLN
  - [ ] LDK
  - [ ] Eclair
- [ ] Create a helm chart and sidecar client for communication with Bitcoin Nodes
  - [x] Bitcoind
  - [ ] btcd
- [ ] Create a library for programmatically interacting with the clients
  - [x] Go
  - [ ] Rust
  - [ ] Python
  - [ ] JavaScript
  - [ ] JVM
- [x] Create a cli version of the library
- [ ] Create or use a tool for generating or simulating network activity
  - [ ] Internal tool
  - [ ] SimLN
- [ ] Facilitate running a public signet node or network
- [ ] Facitlitate interoperability and regression testing of the main Lightning Network implementations
- [ ] Facilitate testing research questions such as new routing algos or channel jamming mitigations
- [x] Thoroughly document the project and provide instructions for use

## How can you help?

- Please give us feedback on the project and your lightning testing use cases
- Directly contributing code, issues and feature requests
- We encourage researchers to help design general tools relevant to them
- Donate to the project to help fund development and maintain the signet Lightning Network

## How can you reach us?

If you have any questions or want to join the project you can reach us here:

- Telgram: https://t.me/+AytRsS0QKH5mMzM8
- Twitter: https://twitter.com/ScalingLN
