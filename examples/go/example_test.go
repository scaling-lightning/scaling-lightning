package main

import (
	"testing"
	"time"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools"
	"github.com/stretchr/testify/assert"
)

// will need a longish (few mins) timeout
func TestMain(t *testing.T) {
	assert := assert.New(t)

	err := sl.StartViaHelmfile("../helmfiles/helmfile.yaml")
	assert.NoError(err)
	// defer sl.StopViaHelmfile("../helmfiles/helmfile.yaml")

	// this one will take a little while as the network is starting up
	err = tools.Retry(func() error {
		return sl.Send("bitcoind", "cln1", 1_000_000)
	}, time.Second*15, time.Minute*2)
	assert.NoError(err)

	err = tools.Retry(func() error {
		return sl.ConnectPeer("lnd1", "cln1")
	}, time.Second*15, time.Minute*2)
	assert.NoError(err)

	err = tools.Retry(func() error {
		return sl.OpenChannel("cln1", "lnd1", 40_001)
	}, time.Second*15, time.Minute*2)
	assert.NoError(err)

	// log.Fatal().Msg("Testing main")
}
