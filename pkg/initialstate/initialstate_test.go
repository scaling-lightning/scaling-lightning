package initialstate

import (
	"testing"

	"github.com/scaling-lightning/scaling-lightning/pkg/types"
	"github.com/stretchr/testify/assert"
)

// TODO: Implement close channels
const exampleInitialState = `
- SendOnChain:
    - { from: bitcoind, to: alice, amountSats: 2_000_000 }
- ConnectPeer:
    - { from: alice, to: bob }
- OpenChannels:
    - { from: alice, to: bob, localAmountSats: 200_000 }
    - { from: alice, to: bob, localAmountSats: 300_000 }
- CloseChannels:
    - { from: alice, havingPeer: bob, havingCapacity: 300_000 }
- OpenChannels:
    - { from: alice, to: bob, localAmountSats: 250_000 }
- SendOverChannel:
    - { from: alice, to: bob, amountMSat: 2_000_000 }
`
func TestParseInitialStateBytes(t *testing.T) {
	assert := assert.New(t)

	initialState, err := NewInitialStateFromBytes([]byte(exampleInitialState), nil)
	assert.Nil(err)

	assert.Equal(7, len(initialState.commands))

	assert.Equal("SendOnChain", initialState.commands[0].commandType)
	assert.Equal("ConnectPeer", initialState.commands[1].commandType)
	assert.Equal("OpenChannels", initialState.commands[2].commandType)
	assert.Equal("OpenChannels", initialState.commands[3].commandType)

	assert.Equal("bitcoind", initialState.commands[0].args["from"])
	assert.Equal(300_000, initialState.commands[3].args["localAmountSats"])
}

func TestSendOnChain(t *testing.T) {
	assert := assert.New(t)

	const initYAML = `
- SendOnChain:
    - { from: bitcoind, to: alice, amount: 2_000_000 }
`
	mockNetwork := NewMockSLNetworkInterface(t)
	mockNetwork.On("Send", "bitcoind", "alice", uint64(2_000_000)).Return("txid", nil)
	initialState, err := NewInitialStateFromBytes([]byte(initYAML), mockNetwork)
	assert.Nil(err)

	err = initialState.Apply()
	assert.Nil(err)
}

func TestConnectPeer(t *testing.T) {
	assert := assert.New(t)

	initYAML := `
- ConnectPeer:
    - { from: alice, to: bob }
`
	mockNetwork := NewMockSLNetworkInterface(t)
	mockNetwork.On("ConnectPeer", "alice", "bob").Return(nil)
	initialState, err := NewInitialStateFromBytes([]byte(initYAML), mockNetwork)
	assert.Nil(err)

	err = initialState.Apply()
	assert.Nil(err)
}

func TestOpenChannel(t *testing.T) {
	assert := assert.New(t)

	const initYAML = `
- OpenChannels:
    - { from: alice, to: bob, localAmountSats: 200_000 }
`
	mockNetwork := NewMockSLNetworkInterface(t)
	mockNetwork.On("OpenChannel", "alice", "bob", uint64(200_000)).Return(types.ChannelPoint{FundingTx: types.Transaction{}, OutputIndex: 21}, nil)
	initialState, err := NewInitialStateFromBytes([]byte(initYAML), mockNetwork)
	assert.Nil(err)

	err = initialState.Apply()
	assert.Nil(err)
}

func TestSendOnChannel(t *testing.T) {
	assert := assert.New(t)

	const initYAML = `
- SendOverChannel:
    - { from: alice, to: bob, amountMSat: 2_000_000 }
`
	mockNetwork := NewMockSLNetworkInterface(t)
	mockNetwork.On("CreateInvoice", "bob", uint64(2_000)).Return("bolt11inv", nil)
	mockNetwork.On("PayInvoice", "alice", "bolt11inv").Return("preimage", nil)
	initialState, err := NewInitialStateFromBytes([]byte(initYAML), mockNetwork)
	assert.Nil(err)

	err = initialState.Apply()
	assert.Nil(err)
}
