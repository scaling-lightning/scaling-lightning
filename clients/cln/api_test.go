package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	clnGRPC "github.com/scaling-lightning/scaling-lightning/clients/cln/grpc"
	"github.com/scaling-lightning/scaling-lightning/clients/cln/mocks"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleWalletBalance(t *testing.T) {
	mockClient := mocks.NewNodeClient(t)

	mockClient.On("ListFunds", mock.Anything, mock.Anything).
		Return(&clnGRPC.ListfundsResponse{
			Outputs: []*clnGRPC.ListfundsOutputs{
				{AmountMsat: &clnGRPC.Amount{Msat: 20}},
				{AmountMsat: &clnGRPC.Amount{Msat: 1}}}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	handleWalletBalance(res, req, mockClient)

	bodyBytes, err := io.ReadAll(res.Result().Body)
	assert.Nil(t, err)
	assert.Contains(t, string(bodyBytes), "21")
}

func TestHandleNewAddress(t *testing.T) {
	mockClient := mocks.NewNodeClient(t)

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

func TestHandleConnectPeer(t *testing.T) {
	assert := assert.New(t)
	mockClient := mocks.NewNodeClient(t)

	mockClient.On("ConnectPeer", mock.Anything, mock.Anything).
		Return(&clnGRPC.ConnectResponse{}, nil)

	pubKey := "037c70cddec9b27c92af73a6b04cf09672fb29b18eca86890d835779979ff61c40"
	host := "lnd1.myfancysats.com"
	port := 9745

	connectPeerReq := types.ConnectPeerReq{PubKey: pubKey, Host: host, Port: port}
	connectPeerBytes, err := json.Marshal(connectPeerReq)

	assert.Nil(err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(connectPeerBytes))
	res := httptest.NewRecorder()

	handleConnectPeer(res, req, mockClient)

	bodyBytes, err := io.ReadAll(res.Result().Body)
	assert.Nil(err)
	assert.Contains(string(bodyBytes), "request received")
}

func TestHandleOpenChannel(t *testing.T) {
	assert := assert.New(t)
	mockClient := mocks.NewNodeClient(t)

	outIndex := 615
	txIdString := "71c73940758ac6ebe34a8f228e28300e"
	TxId, err := hex.DecodeString(txIdString) // real one would be larger
	assert.Nil(err)

	mockClient.On("FundChannel", mock.Anything, mock.Anything).
		Return(&clnGRPC.FundchannelResponse{Txid: TxId, Outnum: uint32(outIndex)}, nil)

	pubKey := "037c70cddec9b27c92af73a6b04cf09672fb29b18eca86890d835779979ff61c40"
	amount := 1000000

	openChannelReq := types.OpenChannelReq{PubKey: pubKey, LocalAmtSats: uint64(amount)}
	openChannelBytes, err := json.Marshal(openChannelReq)

	assert.Nil(err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(openChannelBytes))
	res := httptest.NewRecorder()

	handleOpenChannel(res, req, mockClient)

	bodyBytes, err := io.ReadAll(res.Result().Body)
	assert.Nil(err)
	assert.Contains(string(bodyBytes), "615")
	assert.Contains(string(bodyBytes), "615")
}
