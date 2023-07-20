package main

import (
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools"
	"github.com/stretchr/testify/assert"
)

// will need a longish (few mins) timeout
func TestMain(t *testing.T) {
	assert := assert.New(t)

	err := sl.StartViaHelmfile("../helmfiles/2cln2lnd.yaml")
	assert.NoError(err)
	defer sl.StopViaHelmfile("../helmfiles/2cln2lnd.yaml")

	// this one will take a little while as the network is starting up
	err = tools.Retry(func() error {
		return sl.Send("bitcoind", "cln1", 1_000_000)
	}, time.Second*15, time.Minute*2)
	assert.NoError(err)

	err = tools.Retry(func() error {
		return sl.ConnectPeer("cln1", "cln2")
	}, time.Second*15, time.Minute*3)
	assert.NoError(err)

	err = tools.Retry(func() error {
		return sl.OpenChannel("cln1", "cln2", 40_001)
	}, time.Second*15, time.Minute*3)
	assert.NoError(err)

	err = tools.Retry(func() error {
		balance, err := sl.GetWalletBalanceSats("cln1")
		if err != nil {
			return err
		}
		log.Info().Msgf("cln1 balance: %v", balance)
		return nil
	}, time.Second*15, time.Minute*2)
	assert.NoError(err)
}
