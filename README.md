![ScalingLN Twitter Banner](https://github.com/ohenrik/scaling-lightning/assets/647617/8511c586-7549-4e2b-ad6d-bf87419a624c)

# Scaling Lightning - A Testing Toolkit for the Lightning Network

This initiative aims to build a testing toolkit for the Lightning Network protocol, its implementations, and 
applications that depend on the Lightning Network.

The goal is to collaborate as an industry to help scale the Lightning Network and the applications that depend on it.

## Getting started

> ***The project is still far from complete***. 

### Requirements:

* Kubernetes. 

    > If you are developing locally you can use Docker Desktop and enable 
    Kubernetes in the dashboard.

* You also need Helm 3 and Helmfile. If you are on a Mac you can install them with Homebrew
    
        brew install helm helmfile

* Helm Diff:

      helm plugin install https://github.com/databus23/helm-diff

* Traefik:

      helm repo add traefik https://traefik.github.io/charts
      helm repo update
      helm install traefik traefik/traefik -n traefik --create-namespace

### Starting a Network

To spin up an example network with 2 cln nodes and 4 lnd nodes, run:

    go run . start -f examples/helmfiles/helmfile.yaml

To destroy the network run: 

    go run . stop -f examples/helmfiles/helmfile.yaml
    
### Example CLI Commands

    # list nodes on the network (names were taken from the helmfile)
    go run . list
    
    # get wallet balance of node named bitcoind
    go run . walletbalance -n bitcoind
    
    # get wallet balance of node named lnd2
    go run . walletbalance -n lnd2

    # send on-chain 1 million satoshis from bitcoind to cln1
    go run . send -f bitcoind -t cln1 -a 1000000

    # get the pubkey of a node named lnd1
    go run . pubkey -n lnd1

    # peer lnd1 and cln1 from lnd1
    go run . connectpeer -f lnd1 -t cln1

    # open channel between cln1 and lnd1 with a local balance on cln1 of 70k satoshis
    go run . openchannel -f cln1 -t lnd1 -a 70000

    # have bitcoind generate some blocks and pay itself the block reward
    go run . generate -n bitcoind

### Run the above from code

See [examples/go/example_test.go](examples/go/example_test.go) 

#### Your own configuration

This project is still in its infancy, so we don't have a lot of configuration options yet. 
But you can create your own `helmfile.yaml` in the `helmfiles` folder with the following content:

```yaml
# examples/helmfiles/my-helmfile.yaml
# All files except helmfile.yaml are ignored by git, so you can add them safely.
releases:
    - name: bitcoind
      chart: ../charts/bitcoind # these should point to the charts you want to use.
      ## Set additional values or override existing ones
      # values:
      #   - ./anywhere/bitcoind/values.yaml
    - name: lnd1
      chart: ../charts/lnd # these should point to the charts you want to use.
      ## Set additional values or override existing ones
      ## original values: charts/bitcoind/values.yaml
      values:
          - gRPCNodePort: 30009
          - lndHostPath: {{env "PWD"}}/volumes/lnd1
    - name: lnd2
      chart: ../charts/lnd # these should point to the charts you want to use.
      ## Set additional values or override existing ones
      ## original values: charts/lnd/values.yaml
      values:
          - gRPCNodePort: 30010
          - lndHostPath: {{env "PWD"}}/volumes/lnd2
    - name: lnd3
      chart: ../charts/lnd # these should point to the charts you want to use.
      ## Set additional values or override existing ones
      ## original values: charts/lnd/values.yaml
      values:
          - gRPCNodePort: 30011
          - lndHostPath: {{env "PWD"}}/volumes/lnd3
    - name: lnd4
      chart: ../charts/lnd # these should point to the charts you want to use.
      ## Set additional values or override existing ones
      ## original values: charts/lnd/values.yaml
      values:
          - gRPCNodePort: 30012
          - lndHostPath: {{env "PWD"}}/volumes/lnd4
```

## Why is this important?

Currently, there are unknowns and untested assumptions about how the Lightning Network and its applications will react 
to shocks in transaction volume, channels, nodes, gossip messages, etc.

Having a set of tools and a signet Lightning Network will help:

* Developers test their applications.
* Researchers verify their assumptions.
* Operators test their infrastructure.
* Novices learn how the Lightning Network and various applications work in a somewhat realistic environment without 
  risking real coins.

## How will it work?

We are still in the early stages of planning, but the first tool we are building will be a tool to quickly generate one 
or more Lightning Nodes. These nodes can connect either to a public signet Lightning Network or a private Regtest 
Lightning Network for any combination of LN implementations (CLN, LND, LDK, Acinq etc.).

Other tools, made specifically for testing isolated parts of the protocol, are also relevant. These can help developers
and researchers test their assumptions in an isolated environment. An example of this is 
[The Million Channels Project](https://github.com/rustyrussell/million-channels-project-data) developed by Rusty Russell 
to test gossip.

### Yes, we know about Polar.

While Polar is an excellent project it is fundamentally different from what we want to achieve with the Scaling Lightning initiative. Polar is a desktop application meant to manually build a network through a drag and drop interface, which is great for novices and simple projects. However, it's not something that is suitable for a developer or testing environment at a startup or as a researcher. We have discussed and considered using parts of Polar or altering Polar, but it's not the right direction to go for Scaling Lightning.

## Project milestones

This is an outline of the project's milestones. We will further detail these milestones using the features of GitHub's milestones and project management tools:

* [ ] Create a basic kubernetes setup for running a Lightning Network.
* [ ] Improve configurability and ease of use.
* [ ] Create a client for communication with Lightning Nodes.
  * [ ] LND
  * [ ] CLN
  * [ ] LDK
  * [ ] Eclair
* [ ] Create a client for communication with Bitcoin Nodes.
  * [ ] Bitcoind
  * [ ] btcd
* [ ] Create a library for programmatically interacting with the clients.
* [ ] Create a cli version of the library.
* [ ] Create a tool for generating or simulating network activity.
* [x] Create a website to host the documentation, using Docusaurus.
* [ ] Thoroughly document the project and provide instructions for use.
* [ ] Document and provide links to supporting development resources.

## How can you help?

* We need developers to contribute to the toolkit.
* We invite operators to run nodes on the Signet Lightning Network.
* We encourage researchers to help design general tools relevant to them.
* Donate to the project to help fund development and maintain the signet Lightning Network.

## How can you reach us?

If you have any questions or want to join the project you can reach us here:

* Telgram: https://t.me/+AytRsS0QKH5mMzM8
* Twitter: https://twitter.com/ScalingLN
