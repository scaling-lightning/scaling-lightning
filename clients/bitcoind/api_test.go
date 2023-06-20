package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/scaling-lightning/scaling-lightning/clients/bitcoind/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleWalletBalance(t *testing.T) {
	mockClient := mocks.NewRpcClient(t)

	mockClient.On("GetBalance", mock.AnythingOfType("string")).Return(btcutil.Amount(615), nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	handleWalletBalance(res, req, mockClient)

	bodyBytes, err := io.ReadAll(res.Result().Body)
	assert.Nil(t, err)
	assert.Contains(t, string(bodyBytes), "615")
}
