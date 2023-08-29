package main

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/lightningnetwork/lnd/lnrpc"
	stdcommonclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/common"
	stdlightningclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/lightning"
	basictypes "github.com/scaling-lightning/scaling-lightning/pkg/types"
)

// Probably better mock against our own interface
//go:generate mockery --srcpkg=github.com/lightningnetwork/lnd/lnrpc --name=LightningClient

func (s *commonServer) WalletBalance(
	ctx context.Context,
	in *stdcommonclient.WalletBalanceRequest,
) (*stdcommonclient.WalletBalanceResponse, error) {
	balance, err := s.client.WalletBalance(context.Background(), &lnrpc.WalletBalanceRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "Getting wallet balance from LND's gRPC")
	}

	return &stdcommonclient.WalletBalanceResponse{Balance: uint64(balance.TotalBalance)}, nil
}

func (s *commonServer) NewAddress(
	ctx context.Context,
	in *stdcommonclient.NewAddressRequest,
) (*stdcommonclient.NewAddressResponse, error) {
	response, err := s.client.NewAddress(context.Background(), &lnrpc.NewAddressRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "Getting new address from LND's gRPC")
	}
	return &stdcommonclient.NewAddressResponse{Address: *&response.Address}, nil
}

func (s *lightningServer) PubKey(
	ctx context.Context,
	in *stdlightningclient.PubKeyRequest,
) (*stdlightningclient.PubKeyResponse, error) {
	response, err := s.client.GetInfo(context.Background(), &lnrpc.GetInfoRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "Getting node info from LND's gRPC")
	}
	pubkey, err := basictypes.NewPubKeyFromHexString(response.IdentityPubkey)
	if err != nil {
		return nil, errors.Wrap(err, "Problem converting pubkey to hex")
	}
	return &stdlightningclient.PubKeyResponse{PubKey: pubkey.AsBytes()}, nil
}

func (s *lightningServer) ConnectPeer(
	ctx context.Context,
	req *stdlightningclient.ConnectPeerRequest,
) (*stdlightningclient.ConnectPeerResponse, error) {

	peerAddress := fmt.Sprintf("%v:%v", req.Host, req.Port)
	_, err := s.client.ConnectPeer(
		context.Background(),
		&lnrpc.ConnectPeerRequest{
			Addr: &lnrpc.LightningAddress{
				Pubkey: basictypes.NewPubKeyFromByte(req.PubKey).AsHexString(),
				Host:   peerAddress,
			},
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Problem connecting to peer")
	}
	return &stdlightningclient.ConnectPeerResponse{}, nil
}

func (s *lightningServer) OpenChannel(
	ctx context.Context,
	req *stdlightningclient.OpenChannelRequest,
) (*stdlightningclient.OpenChannelResponse, error) {
	chanPoint, err := s.client.OpenChannelSync(
		context.Background(),
		&lnrpc.OpenChannelRequest{
			NodePubkey:         req.PubKey,
			LocalFundingAmount: int64(req.LocalAmtSats),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Problem opening channel")
	}

	return &stdlightningclient.OpenChannelResponse{
		FundingTxId:          chanPoint.GetFundingTxidBytes(),
		FundingTxOutputIndex: chanPoint.OutputIndex,
	}, nil
}

func (s *lightningServer) CreateInvoice(
	ctx context.Context,
	req *stdlightningclient.CreateInvoiceRequest,
) (*stdlightningclient.CreateInvoiceResponse, error) {
	response, err := s.client.AddInvoice(
		context.Background(),
		&lnrpc.Invoice{
			Value: int64(req.AmtSats),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Adding invoice via LND's gRPC")
	}

	return &stdlightningclient.CreateInvoiceResponse{
		Invoice: response.PaymentRequest,
	}, nil
}

func (s *lightningServer) PayInvoice(
	ctx context.Context,
	req *stdlightningclient.PayInvoiceRequest,
) (*stdlightningclient.PayInvoiceResponse, error) {
	response, err := s.client.SendPaymentSync(
		context.Background(),
		&lnrpc.SendRequest{
			PaymentRequest: req.Invoice,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Paying invoice via LND's gRPC")
	}

	return &stdlightningclient.PayInvoiceResponse{
		PaymentPreimage: response.PaymentPreimage,
	}, nil
}
