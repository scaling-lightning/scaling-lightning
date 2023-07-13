package main

import (
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
