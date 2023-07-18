package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/apierrors"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/bitcoin"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/types"
)

func registerHandlers(standardclient bitcoin.StandardClient, rpcClient rpcClient) {
	standardclient.RegisterWalletBalanceHandler(func(w http.ResponseWriter, r *http.Request) {
		handleWalletBalance(w, r, rpcClient)
	})
	standardclient.RegisterSendToAddressHandler(func(w http.ResponseWriter, r *http.Request) {
		handleSendToAddress(w, r, rpcClient)
	})
	standardclient.RegisterGenerateToAddressHandler(func(w http.ResponseWriter, r *http.Request) {
		handleGenerateToAddress(w, r, rpcClient)
	})
	standardclient.RegisterNewAddressHandler(func(w http.ResponseWriter, r *http.Request) {
		handleNewAddress(w, r, rpcClient)
	})
}

func handleWalletBalance(w http.ResponseWriter, r *http.Request, rpcClient rpcClient) {
	response, err := rpcClient.GetBalance("*")
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem getting wallet balance")
		return
	}
	w.Write([]byte(fmt.Sprintf("Wallet balance is: %v", response.String())))
}

func handleSendToAddress(w http.ResponseWriter, r *http.Request, rpcClient rpcClient) {
	var sendToAddressReq types.SendToAddressReq
	if err := json.NewDecoder(r.Body).Decode(&sendToAddressReq); err != nil {
		apierrors.SendBadRequestFromErr(w, err, "Problem reading request")
		return
	}

	// TODO: pass in real network
	newAddress, err := btcutil.DecodeAddress(
		sendToAddressReq.Address,
		&chaincfg.Params{Name: "regtest"},
	)
	if err != nil {
		apierrors.SendBadRequestFromErr(w, err, "Unable to decode address")
		return
	}
	response, err := rpcClient.SendToAddress(newAddress, btcutil.Amount(sendToAddressReq.AmtSats))
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem sending to address")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Payment sent. Hash: %v", response.String())))
}

func handleGenerateToAddress(w http.ResponseWriter, r *http.Request, rpcClient rpcClient) {
	var generateToAddressReq types.GenerateToAddressReq
	if err := json.NewDecoder(r.Body).Decode(&generateToAddressReq); err != nil {
		apierrors.SendBadRequestFromErr(w, err, "Problem reading request")
		return
	}

	// TODO: pass in real network
	address, err := btcutil.DecodeAddress(
		generateToAddressReq.Address,
		&chaincfg.Params{Name: "regtest"},
	)
	if err != nil {
		apierrors.SendBadRequestFromErr(w, err, "Unable to decode address")
		return
	}
	response, err := rpcClient.GenerateToAddress(
		int64(generateToAddressReq.NumOfBlocks),
		address,
		nil,
	)
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem generating to address")
		return
	}
	w.WriteHeader(http.StatusOK)
	hashes := []string{}
	for _, hash := range response {
		hashes = append(hashes, hash.String()+"\n")
	}
	w.Write([]byte(fmt.Sprintf("Generated. Hashes: %v", hashes)))
}

func handleNewAddress(w http.ResponseWriter, r *http.Request, rpcClient rpcClient) {
	newAddress, err := rpcClient.GetNewAddress(walletName)
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem generating to address")
		return
	}
	response := types.NewAddressRes{Address: newAddress.String()}
	responseJson, err := json.Marshal(response)
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem marshalling new address json")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}
