package main

import (
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools"
	"github.com/scaling-lightning/scaling-lightning/pkg/types"
	"github.com/stretchr/testify/assert"
)

// will need a longish (few mins) timeout
func TestMain(t *testing.T) {
	assert := assert.New(t)
	network := sl.NewSLNetwork("../helmfiles/2cln2lnd.yaml", "")
	err := network.Start()
	assert.NoError(err)

	bitcoind, err := network.GetBitcoinNode("bitcoind")
	assert.NoError(err)

	cln1, err := network.GetLightningNode("cln1")
	assert.NoError(err)

	cln2, err := network.GetLightningNode("cln2")
	assert.NoError(err)

	assert.NoError(err)
	defer network.Stop()

	// this one will take a little while as the network is starting up
	err = tools.Retry(func() error {
		return bitcoind.Send(cln1, types.NewAmountSats(1_000_000))
	}, time.Second*15, time.Minute*2)
	assert.NoError(err)

	err = tools.Retry(func() error {
		return cln1.ConnectPeer(cln2)
	}, time.Second*15, time.Minute*3)
	assert.NoError(err)

	err = tools.Retry(func() error {
		return cln1.OpenChannel(cln2, types.NewAmountSats(40_001))
	}, time.Second*15, time.Minute*3)
	assert.NoError(err)

	err = tools.Retry(func() error {
		balance, err := cln1.GetWalletBalance()
		if err != nil {
			return err
		}
		log.Info().Msgf("cln1 balance: %d", balance.AsSats())
		return nil
	}, time.Second*15, time.Minute*2)
	assert.NoError(err)
}
