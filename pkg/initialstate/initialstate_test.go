package initialstate

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const exampleInitialState = `
- SendOnChain:
    - { from: bitcoind, to: lnd1, amountSats: 2_000_000 }
- OpenChannels:
    - { from: lnd1, to: lnd2, localAmountSats: 200_000 }
    - { from: lnd1, to: lnd2, localAmountSats: 300_000 }
- ConnectPeer:
    - { from: lnd4, to: lnd5 }
- CloseChannels:
    - { from: lnd1, havingPeer: lnd2, havingCapacity: 300_000 }
- OpenChannels:
    - { from: lnd1, to: lnd2, localAmountSats: 250_000 }
- SendOverChannel:
    - { from: lnd1, to: lnd2, amountMSat: 2_000_000 }
`
func TestParseInitialStateFile(t *testing.T) {
	assert := assert.New(t)
	assert.Nil(nil)

	initialState, err := NewInitialState([]byte(exampleInitialState))
	assert.Nil(err)

	assert.Equal(7, len(initialState.commands))

	log.Printf("%v", initialState)
}
