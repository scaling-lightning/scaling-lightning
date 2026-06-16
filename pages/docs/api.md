---
sidebar_position: 6
---

# Golang API

As an alternative to the CLI, a golang library is available with all the same functions. In the future we aim to release libraries for other programming languages. Please get in touch to request which language we should work on next.

## Installation

    go get github.com/scaling-lightning/scaling-lightning

## Example Golang test

Before running this test, ensure the prerequisits are installed on the local system and kubernetes cluster. Please see the [Getting Started](/docs/getting-started) guide.

```go
package main

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/scaling-lightning/scaling-lightning/pkg/initialstate"
	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/stretchr/testify/assert"
)

// will need a longish (few mins) timeout
func TestMainExample(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	assert := assert.New(t)
	network := sl.NewSLNetwork("../helmfiles/public.yaml", "", sl.Regtest)
	err := network.CreateAndStart()
	if err != nil {
		log.Fatal().Err(err).Msg("Problem starting network")
	}

	const initialStateYAML = `
- SendOnChain:
    - { from: bitcoind, to: cln1, amountSats: 1_000_000 }
- ConnectPeer:
    - { from: cln1, to: cln2 }
- OpenChannels:
    - { from: cln1, to: cln2, localAmountSats: 200_000 }
    - { from: cln1, to: cln2, localAmountSats: 300_000 }
- SendOverChannel:
    - { from: cln1, to: cln2, amountMSat: 2_000_000 }
`
	initialState, err := initialstate.NewInitialStateFromBytes([]byte(initialStateYAML), &network)
	assert.NoError(err)
	err = initialState.Apply()
	assert.NoError(err)

	cln2, err := network.GetLightningNode("cln2")
	assert.NoError(err)

	assert.NoError(err)
	defer func() {
		err = network.Destroy()
		assert.NoError(err)
	}()

	balance, err := network.GetWalletBalance("cln1")
	assert.NoError(err)
	log.Info().Msgf("cln1 balance: %d", balance.AsSats())

	connectionDetails, err := network.GetConnectionDetails("cln2")
	assert.NoError(err)
	log.Info().Msgf("cln2 connection host: %v", connectionDetails[0].Host)
	log.Info().Msgf("cln2 connection host: %d", connectionDetails[0].Port)

	connectionFiles, err := cln2.GetConnectionFiles(network.Network.String(), "")
	assert.NoError(err)
	log.Info().Msgf("cln2 client cert size : %v", len(connectionFiles.CLN.ClientCert))
}
```
