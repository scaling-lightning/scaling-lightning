package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/scaling-lightning/scaling-lightning/clients/lnd/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleWalletBalance(t *testing.T) {
	mockClient := mocks.NewLightningClient(t)

	mockClient.On("WalletBalance", mock.Anything, mock.Anything).Return(&lnrpc.WalletBalanceResponse{TotalBalance: 21}, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	handleWalletBalance(res, req, mockClient)

	bodyBytes, err := io.ReadAll(res.Result().Body)
	assert.Nil(t, err)
	assert.Contains(t, string(bodyBytes), "21")
}

func TestHandleNewAddress(t *testing.T) {
	mockClient := mocks.NewLightningClient(t)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	res := httptest.NewRecorder()

	addressStr := "bcrt1qddzehdyj5e7w4sfya3h9qznnm80etc9gkpk0qd"
	mockClient.On("NewAddress", mock.Anything, mock.Anything).Return(&lnrpc.NewAddressResponse{Address: addressStr}, nil)

	handleNewAddress(res, req, mockClient)

	bodyBytes, err := io.ReadAll(res.Result().Body)
	assert.Nil(t, err)
	assert.Contains(t, string(bodyBytes), addressStr)
}

func TestHandlePubKey(t *testing.T) {
	mockClient := mocks.NewLightningClient(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	pubKey := "037c70cddec9b27c92af73a6b04cf09672fb29b18eca86890d835779979ff61c40"
	mockClient.On("GetInfo", mock.Anything, mock.Anything).Return(&lnrpc.GetInfoResponse{IdentityPubkey: pubKey}, nil)

	handlePubKey(res, req, mockClient)

	bodyBytes, err := io.ReadAll(res.Result().Body)
	assert.Nil(t, err)
	assert.Contains(t, string(bodyBytes), pubKey)
}

func TestHandleConnectPeer(t *testing.T) {
	mockClient := mocks.NewLightningClient(t)
	assert := assert.New(t)

	pubKey := "037c70cddec9b27c92af73a6b04cf09672fb29b18eca86890d835779979ff61c40"
	host := "lnd1.myfancysats.com"
	port := 9745

	connectPeerReq := connectPeerReq{PubKey: pubKey, Host: host, Port: port}
	connectPeerBytes, err := json.Marshal(connectPeerReq)
	assert.Nil(err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(connectPeerBytes))
	res := httptest.NewRecorder()

	mockClient.On("ConnectPeer", mock.Anything, mock.Anything).Return(&lnrpc.ConnectPeerResponse{}, nil)

	handleConnectPeer(res, req, mockClient)

	bodyBytes, err := io.ReadAll(res.Result().Body)
	assert.Nil(err)
	assert.Contains(string(bodyBytes), "request received")
}
