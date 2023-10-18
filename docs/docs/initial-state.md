---
sidebar_position: 6
---

# Specify initial state

To specify an initial state of the network a yaml file (or yaml string) can be passsed to the golang library.

## File format

```yaml title="init.yaml"
- SendOnChain:
    - { from: bitcoind, to: alice, amountSats: 2_000_000 }
- ConnectPeer:
    - { from: alice, to: bob }
- OpenChannels:
    - { from: alice, to: bob, localAmountSats: 200_000 }
    - { from: alice, to: bob, localAmountSats: 300_000 }
- SendOverChannel:
    - { from: alice, to: bob, amountMSat: 2_000_000 }
- SendOnChain:
    - { from: bitcoind, to: alice, amountSats: 5_000_000 }
```

The file is read in top down proceedural fashion and command can be repeated (notice in this example SendOnChain is added twice).

## Use in code

```go
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/scaling-lightning/scaling-lightning/pkg/initialstate"
	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
)

// will need a longish (few mins) timeout
func TestMain2(t *testing.T) {
	network := sl.NewSLNetwork("network.yaml", "", sl.Regtest)
	err := network.CreateAndStart()
	assert.NoError(t, err)

	initialState, _ := initialstate.NewInitialStateFromFile("init.yaml", &network)
	err = initialState.Apply()
	assert.NoError(t, err)
}
```

> **_NOTE:_** This file format is currently only available with the go library and not the CLI.
