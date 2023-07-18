package main

import (
	"testing"

	"github.com/rs/zerolog/log"
	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	assert := assert.New(t)
	err := sl.Start()
	assert.NoError(err)

	defer sl.Stop()

	log.Error().Msg("Testing main")
}
