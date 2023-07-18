package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/scaling-lightning/scaling-lightning/clients/bitcoind/mocks"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleWalletBalance(t *testing.T) {
	mockClient := mocks.NewRpcClient(t)
	assert := assert.New(t)

	mockClient.On("GetBalance", mock.AnythingOfType("string")).Return(btcutil.Amount(615), nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	handleWalletBalance(res, req, mockClient)

	bodyBytes, err := io.ReadAll(res.Result().Body)
	assert.Nil(err)
	assert.Contains(string(bodyBytes), "615")
}

func TestHandleSendToAddress(t *testing.T) {
	mockClient := mocks.NewRpcClient(t)
	assert := assert.New(t)

	addressStr := "bcrt1qddzehdyj5e7w4sfya3h9qznnm80etc9gkpk0qd"
	amount := 615
	newAddress, _ := btcutil.DecodeAddress(addressStr, &chaincfg.Params{Name: "regtest"})

	hash, err := chainhash.NewHashFromStr("0")
	assert.Nil(err)
	mockClient.On("SendToAddress", newAddress, btcutil.Amount(amount)).Return(hash, nil)

	sendReq := types.SendToAddressReq{Address: addressStr, AmtSats: uint64(amount)}
	sendReqBytes, err := json.Marshal(sendReq)
	assert.Nil(err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(sendReqBytes))
	res := httptest.NewRecorder()

	handleSendToAddress(res, req, mockClient)

	bodyBytes, err := io.ReadAll(res.Result().Body)
	assert.Nil(err)
	assert.Contains(strings.ToLower(string(bodyBytes)), "payment sent")
}

// GenerateToAddress(numBlocks int64, address btcutil.Address, maxTries *int64) ([]*chainhash.Hash, error)
func TestHandleGenerateToAddress(t *testing.T) {
	mockClient := mocks.NewRpcClient(t)
	assert := assert.New(t)

	hash, err := chainhash.NewHashFromStr("0")
	assert.Nil(err)
	mockClient.On("GenerateToAddress", mock.AnythingOfType("int64"), mock.Anything, mock.AnythingOfType("*int64")).
		Return([]*chainhash.Hash{hash}, nil)

	addressStr := "bcrt1qddzehdyj5e7w4sfya3h9qznnm80etc9gkpk0qd"
	genReq := types.GenerateToAddressReq{Address: addressStr, NumOfBlocks: uint64(20)}
	genReqBytes, err := json.Marshal(genReq)
	assert.Nil(err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(genReqBytes))
	res := httptest.NewRecorder()

	handleGenerateToAddress(res, req, mockClient)

	bodyBytes, err := io.ReadAll(res.Result().Body)
	assert.Nil(err)
	assert.Contains(string(bodyBytes), hash.String())
}

func TestHandleNewAddress(t *testing.T) {
	mockClient := mocks.NewRpcClient(t)
	assert := assert.New(t)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	res := httptest.NewRecorder()

	addressStr := "bcrt1qddzehdyj5e7w4sfya3h9qznnm80etc9gkpk0qd"
	address, _ := btcutil.DecodeAddress(addressStr, &chaincfg.Params{Name: "regtest"})

	mockClient.On("GetNewAddress", mock.AnythingOfType("string")).Return(address, nil)

	handleNewAddress(res, req, mockClient)

	bodyBytes, err := io.ReadAll(res.Result().Body)
	assert.Nil(err)
	assert.Contains(string(bodyBytes), addressStr)
}
