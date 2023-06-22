package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/apierrors"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/lightning"
)

// Probably better mock against our own interface
//go:generate mockery --srcpkg=github.com/lightningnetwork/lnd/lnrpc --name=LightningClient

func registerHandlers(standardclient lightning.StandardClient, lndClient lnrpc.LightningClient) {
	standardclient.RegisterWalletBalanceHandler(func(w http.ResponseWriter, r *http.Request) {
		handleWalletBalance(w, r, lndClient)
	})
	standardclient.RegisterNewAddressHandler(func(w http.ResponseWriter, r *http.Request) {
		handleNewAddress(w, r, lndClient)
	})
}

func handleWalletBalance(w http.ResponseWriter, r *http.Request, lndClient lnrpc.LightningClient) {
	response, err := lndClient.WalletBalance(context.Background(), &lnrpc.WalletBalanceRequest{})
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem getting wallet balance")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Wallet balance is: %v", response.TotalBalance)))
}

type newAddressRes struct {
	Address string `json:"address"`
}

func handleNewAddress(w http.ResponseWriter, r *http.Request, lndClient lnrpc.LightningClient) {
	newAddress, err := lndClient.NewAddress(context.Background(), &lnrpc.NewAddressRequest{})
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem getting new address")
		return
	}
	response := newAddressRes{Address: newAddress.Address}
	responseJson, err := json.Marshal(response)
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem marshalling new address json")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}
