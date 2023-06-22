package main

import (
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
