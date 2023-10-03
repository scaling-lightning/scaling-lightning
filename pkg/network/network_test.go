package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNetwork(t *testing.T) {
	assert := assert.New(t)
	mockLightningNode := NewMockLightningNodeInterface(t)
	mockBitcoinNode := NewMockBitcoinNodeInterface(t)

	mockLightningNode.On("GetName").Return("alice")
	mockLightningNode.On("SendToAddress", mock.Anything, mock.Anything, mock.Anything).Return("txid", nil)
	mockBitcoinNode.On("GetName").Return("bitcoind")
	mockBitcoinNode.On("GetNewAddress", mock.Anything).Return("address", nil)
	mockBitcoinNode.On("Generate", mock.Anything, mock.Anything, mock.Anything).Return([]string{"hash"}, nil)

	assert.Equal(mockLightningNode, mockLightningNode)

	slNetwork := &SLNetwork{}
	slNetwork.LightningNodes = []LightningNodeInterface{mockLightningNode}
	slNetwork.BitcoinNodes = []BitcoinNodeInterface{mockBitcoinNode}

	response, err := slNetwork.Send("alice", "bitcoind", 1000)
	assert.Nil(err)
	assert.Equal("txid", response)
}
