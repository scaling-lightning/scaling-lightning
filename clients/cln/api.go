package main

import (
	"context"

	"github.com/cockroachdb/errors"
	clnGRPC "github.com/scaling-lightning/scaling-lightning/clients/cln/grpc"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/lightning"
	"github.com/scaling-lightning/scaling-lightning/pkg/types"
)

// Probably better mock against our own interface
//go:generate mockery --dir=grpc --name=NodeClient

func (s *lightningServer) WalletBalance(
	ctx context.Context,
	in *lightning.WalletBalanceRequest,
) (*lightning.WalletBalanceResponse, error) {
	response, err := s.client.ListFunds(context.Background(), &clnGRPC.ListfundsRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "Problem listing funds")
	}

	var total uint64
	for _, output := range response.Outputs {
		total += output.AmountMsat.Msat
	}

	return &lightning.WalletBalanceResponse{Balance: total}, nil
}

func (s *lightningServer) NewAddress(
	ctx context.Context,
	in *lightning.NewAddressRequest,
) (*lightning.NewAddressResponse, error) {
	response, err := s.client.NewAddr(context.Background(), &clnGRPC.NewaddrRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "Problem getting new address")
	}
	return &lightning.NewAddressResponse{Address: *response.Bech32}, nil
}

func (s *commonServer) PubKey(
	ctx context.Context,
	in *lightning.PubKeyRequest,
) (*lightning.PubKeyResponse, error) {
	response, err := s.client.Getinfo(context.Background(), &clnGRPC.GetinfoRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "Problem getting info")
	}
	return &lightning.PubKeyResponse{PubKey: response.Id}, nil
}

func (s *lightningServer) ConnectPeer(
	ctx context.Context,
	req *lightning.ConnectPeerRequest,
) (*lightning.ConnectPeerResponse, error) {

	pubkey := types.NewPubKeyFromByte(req.PubKey)
	_, err := s.client.ConnectPeer(context.Background(),
		&clnGRPC.ConnectRequest{Id: pubkey.AsHexString(), Host: &req.Host, Port: &req.Port})
	if err != nil {
		return nil, errors.Wrap(err, "Problem connecting to peer")
	}
	return &lightning.ConnectPeerResponse{}, nil
}

func (s *lightningServer) OpenChannel(
	ctx context.Context,
	req *lightning.OpenChannelRequest,
) (*lightning.OpenChannelResponse, error) {
	amount := clnGRPC.AmountOrAll{
		Value: &clnGRPC.AmountOrAll_Amount{
			Amount: &clnGRPC.Amount{Msat: uint64(req.LocalAmtSats) * 1000},
		},
	}
	chanPoint, err := s.client.FundChannel(context.Background(),
		&clnGRPC.FundchannelRequest{Id: req.PubKey, Amount: &amount})
	if err != nil {
		return nil, errors.Wrap(err, "Problem opening channel")
	}
	return &lightning.OpenChannelResponse{
		FundingTxId:          chanPoint.Txid,
		FundingTxOutputIndex: chanPoint.Outnum,
	}, nil
}
