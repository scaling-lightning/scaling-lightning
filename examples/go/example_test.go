package main

import (
	"testing"

	"github.com/rs/zerolog/log"
	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	assert := assert.New(t)
	err := sl.StartViaHelmfile("../helmfiles/helmfile.yaml")
	assert.NoError(err)
	// defer sl.StopViaHelmfile("../helmfiles/helmfile.yaml")

	if err = sl.Send("btcd", "cln1", 1_000_000); err != nil {
		log.Fatal().Err(err).Msg("Failed to send")
	}

	log.Fatal().Msg("Testing main")
}
