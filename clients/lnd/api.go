package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/rs/zerolog/log"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient"
)

// Probably better mock against our own interface
//go:generate mockery --srcpkg=github.com/lightningnetwork/lnd/lnrpc --name=LightningClient

func registerHandlers(standardclient standardclient.StandardClient, lndClient lnrpc.LightningClient) {
	standardclient.HandleWalletBalance(func(w http.ResponseWriter, r *http.Request) {
		handleWalletBalance(w, r, lndClient)
	})
}

func handleWalletBalance(w http.ResponseWriter, r *http.Request, lndClient lnrpc.LightningClient) {
	response, err := lndClient.WalletBalance(context.Background(), &lnrpc.WalletBalanceRequest{})
	if err != nil {
		log.Error().Err(err).Msg("Problem getting wallet balance")
		return
	}
	w.Write([]byte(fmt.Sprintf("Wallet balance is: %v", response.TotalBalance)))
}
