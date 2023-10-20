---
sidebar_position: 5
---

# CLI

CLI is one option for creating and interacting with the network (the other being code via a library).

## Installation

Download binary for your system from [Releases](https://github.com/scaling-lightning/scaling-lightning/releases)

    # untar to get binary
    tar -xzf scaling-lightning-[version]-[os]-[architecture].tar.gz

    # Mac OS only - mark file as safe so it will run
    xattr -dr com.apple.quarantine scaling-lightning

    # run - should print CLI help
    ./scaling-lightning

## Use

Before using the CLI, ensure the dependencies have been installed on your local system and cluster by following the prerequisits section of the [getting started guide](/docs/getting-started)

To see a full list of commands run:

    ./scaling-lightning help

Output of help:

```shell
A CLI for interacting with the scaling-lightning network

Usage:
sl [command]

Available Commands:
channelbalance    Get the onchain wallet balance of a node
completion        Generate the autocompletion script for the specified shell
connectiondetails Output the connection details for a node or all nodes
connectpeer       Connect peers
create            Create and start the network
createinvoice     Create lightning invoice
destroy           destroy the network
generate          Generate bitcoin blocks
help              Help about any command
list              List the nodes in the network
openchannel       Open a channel between two nodes
payinvoice        pay lightning invoice
pubkey            Get the pubkey of a node
send              Send on chain funds betwen nodes
start             Start a stopped network
stop              Stop the network
version           Get the version of this binary
walletbalance     Get the onchain wallet balance of a node
writeauthfiles    Output the auth files for a node or all nodes

Flags:
-d, --debug               Enable debug logging
-h, --help                help for sl
-H, --host string         Host of the scaling-lightning API
-k, --kubeconfig string   Location of Kubernetes config file (default "/Users/maxedwards/.kube/config")
-p, --port uint16         Port of the scaling-lightning API

Use "sl [command] --help" for more information about a command.
```

For specific flags for each command run:

    ./scaling-lightning send help

Output:

```shell
Send on chain funds betwen nodes

Usage:
  scaling-lightning send [flags]

Flags:
  -a, --amount uint   Amount of satoshis to send
  -f, --from string   Name of node to send from
  -h, --help          help for send
  -t, --to string     Name of node to send to

Global Flags:
  -d, --debug               Enable debug logging
  -H, --host string         Host of the scaling-lightning API
  -k, --kubeconfig string   Location of Kubernetes config file (default "/Users/maxedwards/.kube/config")
  -p, --port uint16         Port of the scaling-lightning API
```

Making the command to move onchain funds from bitcoind to alice be:

    ./scaling-lightning send -f bitcoind -t alice -a 1000000
