package main

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/bitcoin"
)

func registerHandlers(standardclient bitcoin.StandardClient, rpcClient rpcClient) {
	standardclient.RegisterWalletBalanceHandler(func(w http.ResponseWriter, r *http.Request) {
		handleWalletBalance(w, r, rpcClient)
	})
}

func handleWalletBalance(w http.ResponseWriter, r *http.Request, rpcClient rpcClient) {
	response, err := rpcClient.GetBalance(walletName)
	if err != nil {
		log.Error().Err(err).Msg("Problem getting wallet balance")
		return
	}
	w.Write([]byte(fmt.Sprintf("Wallet balance is: %v", response.String())))
}
