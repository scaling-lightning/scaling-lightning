package main

import (
	"context"
	"fmt"
	"net/http"

	clnGRPC "github.com/scaling-lightning/scaling-lightning/clients/cln/grpc"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/apierrors"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/lightning"
)

func registerHandlers(standardclient lightning.StandardClient, clnClient clnGRPC.NodeClient) {
	standardclient.RegisterWalletBalanceHandler(func(w http.ResponseWriter, r *http.Request) {
		handleWalletBalance(w, r, clnClient)
	})
}

func handleWalletBalance(w http.ResponseWriter, r *http.Request, clnClient clnGRPC.NodeClient) {
	response, err := clnClient.ListFunds(context.Background(), &clnGRPC.ListfundsRequest{})
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem getting wallet balance")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Wallet balance is: %v", response.Outputs[0].AmountMsat)))
}
