package network

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const exampleInitialState = `
- OpenChannels:
    - lnd1 lnd2 2_000_000 1_000_000 firstChannel
    - lnd1 lnd2 2_000_000 1_000_000 secondChannel
    - lnd1 lnd2 2_000_000 1_000_000 thirdChannel
- ConnectPeer:
    - lnd4 lnd5
- CloseChannels:
    - secondChannel
- PayInvoice:
    - lnd1 lnd2 2000
- PayOnChain:
    - bitcoind lnd2 2000
- OpenChannels:
    - lnd1 lnd2 2_000_000 1_000_000 firstChannel
`

func TestParseInitialStateFile(t *testing.T) {
	assert := assert.New(t)
	assert.Nil(nil)

	initialState, err := newInitialState([]byte(exampleInitialState))
	assert.Nil(err)

	log.Printf("%v", initialState)
}
