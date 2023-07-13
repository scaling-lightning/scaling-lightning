package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	clnGRPC "github.com/scaling-lightning/scaling-lightning/clients/cln/grpc"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/apierrors"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/lightning"
)

// Probably better mock against our own interface
//go:generate mockery --dir=grpc --name=NodeClient

func registerHandlers(standardclient lightning.StandardClient, clnClient clnGRPC.NodeClient) {
	standardclient.RegisterWalletBalanceHandler(func(w http.ResponseWriter, r *http.Request) {
		handleWalletBalance(w, r, clnClient)
	})
	standardclient.RegisterNewAddressHandler(func(w http.ResponseWriter, r *http.Request) {
		handleNewAddress(w, r, clnClient)
	})
}

type newAddressRes struct {
	Address string `json:"address"`
}

func handleNewAddress(w http.ResponseWriter, r *http.Request, clnClient clnGRPC.NodeClient) {
	newAddress, err := clnClient.NewAddr(context.Background(), &clnGRPC.NewaddrRequest{})	
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem getting new address")
		return
	}
	response := newAddressRes{Address: *newAddress.Bech32}
	responseJson, err := json.Marshal(response)
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem marshalling new address json")
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

func handleWalletBalance(w http.ResponseWriter, r *http.Request, clnClient clnGRPC.NodeClient) {
	response, err := clnClient.ListFunds(context.Background(), &clnGRPC.ListfundsRequest{})
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem getting wallet balance")
		return
	}
	
	total := 0
	for _, output := range response.Outputs {
		total += int(output.AmountMsat.Msat)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Wallet balance is: %v msats", total)))
}
