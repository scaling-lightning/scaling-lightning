package main

import (
	"testing"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	assert := assert.New(t)
	err := sl.StartViaHelmfile("../helmfiles/helmfile.yaml")
	assert.NoError(err)
	defer sl.StopViaHelmfile("../helmfiles/helmfile.yaml")

	err = sl.Send("bitcoind", "cln1", 1_000_000)
	assert.NoError(err)

	// log.Fatal().Msg("Testing main")
}
