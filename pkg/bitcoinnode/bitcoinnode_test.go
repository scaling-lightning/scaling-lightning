package bitcoinnode

import (
	"testing"

	"github.com/scaling-lightning/scaling-lightning/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestBitcoinNode(t *testing.T) {
	assert := assert.New(t)

	bitcoinNode := &BitcoinNode{}
	response, err := bitcoinNode.SendToAddress(nil, "address", types.NewAmountSats(1000))
	assert.Nil(err)
	assert.NotNil(response)
}
