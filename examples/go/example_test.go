package main

import (
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools"
	"github.com/scaling-lightning/scaling-lightning/pkg/types"
	"github.com/stretchr/testify/assert"
)

// will need a longish (few mins) timeout
func TestMain(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	assert := assert.New(t)
	network := sl.NewSLNetwork("../helmfiles/public.yaml", "", sl.Regtest)
	err := network.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("Problem starting network")
	}

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
		_, err := bitcoind.Send(cln1, types.NewAmountSats(1_000_000))
		return err
	}, time.Second*15, time.Minute*2)
	assert.NoError(err)

	err = tools.Retry(func() error {
		return cln1.ConnectPeer(cln2)
	}, time.Second*15, time.Minute*3)
	assert.NoError(err)

	err = tools.Retry(func() error {
		_, err := cln1.OpenChannel(cln2, types.NewAmountSats(40_001))
		return err
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

	err = tools.Retry(func() error {
		connectionDetails, err := cln2.GetConnectionDetails()
		if err != nil {
			return err
		}
		log.Info().Msgf("cln2 connection host: %v", connectionDetails[0].Host)
		log.Info().Msgf("cln2 connection host: %d", connectionDetails[0].Port)
		return nil
	}, time.Second*15, time.Minute*2)

	assert.NoError(err)
	err = tools.Retry(func() error {
		connectionFiles, err := cln2.GetConnectionFiles()
		if err != nil {
			return err
		}
		log.Info().Msgf("cln2 client cert size : %v", len(connectionFiles.CLN.ClientCert))
		return nil
	}, time.Second*15, time.Minute*2)
	assert.NoError(err)
}
