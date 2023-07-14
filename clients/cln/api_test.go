package main

import (
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	clnGRPC "github.com/scaling-lightning/scaling-lightning/clients/cln/grpc"
	"github.com/scaling-lightning/scaling-lightning/clients/cln/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleWalletBalance(t *testing.T) {
	mockClient := mocks.NewNodeClient(t)

	// mock based test for the handleWalletBalance function
	mockClient.On("ListFunds", mock.Anything, mock.Anything).
		Return(&clnGRPC.ListfundsResponse{Outputs: []*clnGRPC.ListfundsOutputs{
			{AmountMsat: &clnGRPC.Amount{Msat: 20}},
			{AmountMsat: &clnGRPC.Amount{Msat: 1}},
		}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	handleWalletBalance(res, req, mockClient)

	bodyBytes, err := io.ReadAll(res.Result().Body)
	assert.Nil(t, err)
	assert.Contains(t, string(bodyBytes), "21")
}

func TestHandleNewAddress(t *testing.T) {
	mockClient := mocks.NewNodeClient(t)

	// mock based test for the handleGetNewAddress function
	address := "bcrt1qddzehdyj5e7w4sfya3h9qznnm80etc9gkpk0qd"
	mockClient.On("NewAddr", mock.Anything, mock.Anything).
		Return(&clnGRPC.NewaddrResponse{Bech32: &address}, nil)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	res := httptest.NewRecorder()

	handleNewAddress(res, req, mockClient)

	bodyBytes, err := io.ReadAll(res.Result().Body)
	assert.Nil(t, err)
	assert.Contains(t, string(bodyBytes), address)
}

func TestHandlePubKey(t *testing.T) {
	mockClient := mocks.NewNodeClient(t)

	// mock based test for the handleGetPubKey function
	pubKey := "02c3d4d2b6b4b8e2f5f9c6e3f0b1e8d5a2c3d4d2b6b4b8e2f5f9c6e3f0b1e8d5"
	pubKeyBinary, err := hex.DecodeString(pubKey)
	assert.Nil(t, err)

	mockClient.On("Getinfo", mock.Anything, mock.Anything).
		Return(&clnGRPC.GetinfoResponse{Id: pubKeyBinary}, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	handlePubKey(res, req, mockClient)

	bodyBytes, err := io.ReadAll(res.Result().Body)
	assert.Nil(t, err)
	assert.Contains(t, string(bodyBytes), pubKey)
}
	