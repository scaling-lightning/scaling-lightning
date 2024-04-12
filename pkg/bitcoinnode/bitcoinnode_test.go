package bitcoinnode

import (
	"testing"

	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/common"
	"github.com/scaling-lightning/scaling-lightning/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBitcoinNode(t *testing.T) {
	assert := assert.New(t)

	bitcoinNode := &BitcoinNode{
		Namespace: "sl",
	}

	mockGrpcClient := common.NewMockCommonClient(t)
	mockGrpcClient.On("Send", mock.Anything, mock.Anything).Return(&common.SendResponse{}, nil)
	response, err := bitcoinNode.SendToAddress(mockGrpcClient, "address", types.NewAmountSats(1000))
	assert.Nil(err)
	assert.NotNil(response)
}
