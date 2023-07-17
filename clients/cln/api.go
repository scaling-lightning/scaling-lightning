package main

import (
	"context"
	"encoding/hex"
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
	standardclient.RegisterPubKeyHandler(func(w http.ResponseWriter, r *http.Request) {
		handlePubKey(w, r, clnClient)
	})
	standardclient.RegisterConnectPeerHandler(func(w http.ResponseWriter, r *http.Request) {
		handleConnectPeer(w, r, clnClient)
	})
	standardclient.RegisterOpenChannelHandler(func(w http.ResponseWriter, r *http.Request) {
		handleOpenChannel(w, r, clnClient)
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

type pubKeyRes struct {
	PubKey string `json:"pubkey"`
}

func handlePubKey(w http.ResponseWriter, r *http.Request, clnClient clnGRPC.NodeClient) {
	info, err := clnClient.Getinfo(context.Background(), &clnGRPC.GetinfoRequest{})
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem getting node info")
		return
	}
	response := pubKeyRes{PubKey: hex.EncodeToString(info.Id)}
	responseJson, err := json.Marshal(response)
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem marshalling pubkey json")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

type connectPeerReq struct {
	PubKey string `json:"pubKey"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}

func handleConnectPeer(w http.ResponseWriter, r *http.Request, clnClient clnGRPC.NodeClient) {
	var connectPeerReq connectPeerReq
	if err := json.NewDecoder(r.Body).Decode(&connectPeerReq); err != nil {
		apierrors.SendBadRequestFromErr(w, err, "Problem reading request")
		return
	}

	port := uint32(connectPeerReq.Port)

	_, err := clnClient.ConnectPeer(context.Background(),
		&clnGRPC.ConnectRequest{Id: connectPeerReq.PubKey, Host: &connectPeerReq.Host, Port: &port})
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem connecting to peer")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Connect peer request received"))
}

type openChannelReq struct {
	PubKey   string `json:"pubKey"`
	LocalAmt int64  `json:"localAmount"`
}

type openChannelRes struct {
	FundingTx   string `json:"fundingTx"`
	OutputIndex uint32 `json:"outputIndex"`
}

func handleOpenChannel(w http.ResponseWriter, r *http.Request, clnClient clnGRPC.NodeClient) {
	var openChannelReq openChannelReq
	if err := json.NewDecoder(r.Body).Decode(&openChannelReq); err != nil {
		apierrors.SendBadRequestFromErr(w, err, "Problem reading request")
		return
	}

	pubKeyHex, err := hex.DecodeString(openChannelReq.PubKey)
	if err != nil {
		apierrors.SendBadRequestFromErr(w, err, "Problem decoding pubKey to hex")
		return
	}

	amount := clnGRPC.AmountOrAll{
		Value: &clnGRPC.AmountOrAll_Amount{
			Amount: &clnGRPC.Amount{Msat: uint64(openChannelReq.LocalAmt) * 1000},
		},
	}

	chanPoint, err := clnClient.FundChannel(context.Background(),
		&clnGRPC.FundchannelRequest{Id: pubKeyHex, Amount: &amount})
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem opening channel")
		return
	}

	response := openChannelRes{
		FundingTx:   hex.EncodeToString(chanPoint.Txid),
		OutputIndex: chanPoint.Outnum,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		apierrors.SendServerErrorFromErr(w, err, "Problem marshalling funding tx and index json")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}
