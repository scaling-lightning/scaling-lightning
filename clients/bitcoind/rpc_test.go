package main

import (
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/cockroachdb/errors"
	"github.com/scaling-lightning/scaling-lightning/clients/bitcoind/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitialiseBitcoind(t *testing.T) {
	assert := assert.New(t)

	mockClient := mocks.NewRpcClient(t)

	mockClient.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{}, nil)

	mockClient.On("LoadWallet", mock.AnythingOfType("string")).
		Return(&btcjson.LoadWalletResult{Name: "", Warning: ""}, errors.New("Wallet doesn't exist"))

	mockClient.On("CreateWallet", mock.AnythingOfType("string")).
		Return(&btcjson.CreateWalletResult{}, nil)

	newAddress, _ := btcutil.DecodeAddress(
		"bcrt1qddzehdyj5e7w4sfya3h9qznnm80etc9gkpk0qd",
		&chaincfg.Params{Name: "regtest"},
	)
	mockClient.On("GetNewAddress", mock.AnythingOfType("string")).Return(newAddress, nil)

	mockClient.On("GenerateToAddress", mock.Anything, mock.Anything, mock.AnythingOfType("*int64")).
		Return([]*chainhash.Hash{}, nil)

	err := initialiseBitcoind(mockClient)
	assert.Nil(err)
}
