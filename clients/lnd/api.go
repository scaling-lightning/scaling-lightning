package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/apierrors"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/lightning"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/types"
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
	standardclient.RegisterPubKeyHandler(func(w http.ResponseWriter, r *http.Request) {
		handlePubKey(w, r, lndClient)
	})
	standardclient.RegisterConnectPeerHandler(func(w http.ResponseWriter, r *http.Request) {
		handleConnectPeer(w, r, lndClient)
	})
	standardclient.RegisterOpenChannelHandler(func(w http.ResponseWriter, r *http.Request) {
		handleOpenChannel(w, r, lndClient)
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

func handleNewAddress(w http.ResponseWriter, r *http.Request, lndClient lnrpc.LightningClient) {
	newAddress, err := lndClient.NewAddress(context.Background(), &lnrpc.NewAddressRequest{})
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem getting new address")
		return
	}
	response := types.NewAddressRes{Address: newAddress.Address}
	responseJson, err := json.Marshal(response)
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem marshalling new address json")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

func handlePubKey(w http.ResponseWriter, r *http.Request, lndClient lnrpc.LightningClient) {
	pubKey, err := lndClient.GetInfo(context.Background(), &lnrpc.GetInfoRequest{})
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem getting node info")
		return
	}
	response := types.PubKeyRes{PubKey: pubKey.IdentityPubkey}
	responseJson, err := json.Marshal(response)
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem marshalling pubkey json")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

func handleConnectPeer(w http.ResponseWriter, r *http.Request, lndClient lnrpc.LightningClient) {
	var connectPeerReq types.ConnectPeerReq
	if err := json.NewDecoder(r.Body).Decode(&connectPeerReq); err != nil {
		apierrors.SendBadRequestFromErr(w, err, "Problem reading request")
		return
	}

	peerAddress := fmt.Sprintf("%v:%v", connectPeerReq.Host, connectPeerReq.Port)
	_, err := lndClient.ConnectPeer(
		context.Background(),
		&lnrpc.ConnectPeerRequest{
			Addr: &lnrpc.LightningAddress{Pubkey: connectPeerReq.PubKey, Host: peerAddress},
		},
	)
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem connecting to peer")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Connect peer request received"))
}

func handleOpenChannel(w http.ResponseWriter, r *http.Request, lndClient lnrpc.LightningClient) {
	var openChannelReq types.OpenChannelReq
	if err := json.NewDecoder(r.Body).Decode(&openChannelReq); err != nil {
		apierrors.SendBadRequestFromErr(w, err, "Problem reading request")
		return
	}

	pubKeyHex, err := hex.DecodeString(openChannelReq.PubKey)
	if err != nil {
		apierrors.SendBadRequestFromErr(w, err, "Problem decoding pubKey to hex")
		return
	}

	chanPoint, err := lndClient.OpenChannelSync(
		context.Background(),
		&lnrpc.OpenChannelRequest{
			NodePubkey:         pubKeyHex,
			LocalFundingAmount: openChannelReq.LocalAmtSats,
		},
	)
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem opening channel")
		return
	}

	response := types.OpenChannelRes{
		FundingTx:   chanPoint.GetFundingTxidStr(),
		OutputIndex: chanPoint.OutputIndex,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem marshalling funding tx and index json")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}
