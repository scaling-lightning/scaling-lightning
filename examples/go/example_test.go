package main

import (
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools"
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

	cln2, err := network.GetLightningNode("cln2")
	assert.NoError(err)

	assert.NoError(err)
	defer func() {
		err = network.Destroy()
		assert.NoError(err)
	}()

	// this one will take a little while as the network is starting up
	err = tools.Retry(func() error {
		_, err := network.Send("bitcoind", "cln1", 1_000_000)
		return errors.Wrap(err, "Sending million sats to cln1")
	}, time.Second*15, time.Minute*2)
	assert.NoError(err)

	err = tools.Retry(func() error {
		return errors.Wrap(network.ConnectPeer("cln1", "cln2"), "Connecting cln1 to cln2")
	}, time.Second*15, time.Minute*3)
	assert.NoError(err)

	err = tools.Retry(func() error {
		_, err := network.OpenChannel("cln1", "cln2", 40_001)
		return errors.Wrap(err, "Opening channel from cln1 to cln2")
	}, time.Second*15, time.Minute*3)
	assert.NoError(err)

	err = tools.Retry(func() error {
		balance, err := network.GetWalletBalance("cln1")
		if err != nil {
			return errors.Wrap(err, "Getting cln1 balance")
		}
		log.Info().Msgf("cln1 balance: %d", balance.AsSats())
		return nil
	}, time.Second*15, time.Minute*2)
	assert.NoError(err)

	err = tools.Retry(func() error {
		connectionDetails, err := network.GetConnectionDetails("cln2")
		if err != nil {
			return errors.Wrap(err, "Getting cln2 connection details")
		}
		log.Info().Msgf("cln2 connection host: %v", connectionDetails[0].Host)
		log.Info().Msgf("cln2 connection host: %d", connectionDetails[0].Port)
		return nil
	}, time.Second*15, time.Minute*2)

	assert.NoError(err)
	err = tools.Retry(func() error {
		connectionFiles, err := cln2.GetConnectionFiles(network.Network.String(), "")
		if err != nil {
			return errors.Wrap(err, "Getting cln2 connection files")
		}
		log.Info().Msgf("cln2 client cert size : %v", len(connectionFiles.CLN.ClientCert))
		return nil
	}, time.Second*15, time.Minute*2)
	assert.NoError(err)
}
