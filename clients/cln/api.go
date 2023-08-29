package main

import (
	"context"

	"github.com/cockroachdb/errors"
	clnGRPC "github.com/scaling-lightning/scaling-lightning/clients/cln/grpc"
	stdcommonclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/common"
	stdlightningclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/lightning"
	"github.com/scaling-lightning/scaling-lightning/pkg/types"
)

// Probably better mock against our own interface
//go:generate mockery --dir=grpc --name=NodeClient

func (s *commonServer) WalletBalance(
	ctx context.Context,
	in *stdcommonclient.WalletBalanceRequest,
) (*stdcommonclient.WalletBalanceResponse, error) {
	response, err := s.client.ListFunds(context.Background(), &clnGRPC.ListfundsRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "Problem listing funds")
	}

	var total uint64
	for _, output := range response.Outputs {
		total += output.AmountMsat.Msat / 1000
	}

	return &stdcommonclient.WalletBalanceResponse{Balance: total}, nil
}

func (s *commonServer) NewAddress(
	ctx context.Context,
	in *stdcommonclient.NewAddressRequest,
) (*stdcommonclient.NewAddressResponse, error) {
	response, err := s.client.NewAddr(context.Background(), &clnGRPC.NewaddrRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "Problem getting new address")
	}
	return &stdcommonclient.NewAddressResponse{Address: *response.Bech32}, nil
}

func (s *lightningServer) PubKey(
	ctx context.Context,
	in *stdlightningclient.PubKeyRequest,
) (*stdlightningclient.PubKeyResponse, error) {
	response, err := s.client.Getinfo(context.Background(), &clnGRPC.GetinfoRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "Problem getting info")
	}
	return &stdlightningclient.PubKeyResponse{PubKey: response.Id}, nil
}

func (s *lightningServer) ConnectPeer(
	ctx context.Context,
	req *stdlightningclient.ConnectPeerRequest,
) (*stdlightningclient.ConnectPeerResponse, error) {

	pubkey := types.NewPubKeyFromByte(req.PubKey)
	_, err := s.client.ConnectPeer(context.Background(),
		&clnGRPC.ConnectRequest{Id: pubkey.AsHexString(), Host: &req.Host, Port: &req.Port})
	if err != nil {
		return nil, errors.Wrap(err, "Problem connecting to peer")
	}
	return &stdlightningclient.ConnectPeerResponse{}, nil
}

func (s *lightningServer) OpenChannel(
	ctx context.Context,
	req *stdlightningclient.OpenChannelRequest,
) (*stdlightningclient.OpenChannelResponse, error) {
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
	return &stdlightningclient.OpenChannelResponse{
		FundingTxId:          chanPoint.Txid,
		FundingTxOutputIndex: chanPoint.Outnum,
	}, nil
}

func (s *lightningServer) CreateInvoice(
	ctx context.Context,
	req *stdlightningclient.CreateInvoiceRequest,
) (*stdlightningclient.CreateInvoiceResponse, error) {
	amount := clnGRPC.AmountOrAny{
		Value: &clnGRPC.AmountOrAny_Amount{
			Amount: &clnGRPC.Amount{Msat: uint64(req.AmtSats) * 1000},
		},
	}
	invoice, err := s.client.Invoice(context.Background(),
		&clnGRPC.InvoiceRequest{AmountMsat: &amount})
	if err != nil {
		return nil, errors.Wrap(err, "Creating invoice via cln's gRPC")
	}
	return &stdlightningclient.CreateInvoiceResponse{Invoice: invoice.Bolt11}, nil
}
