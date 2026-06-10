package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/rs/zerolog/log"
	stdcommonclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/common"
	stdlightningclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/lightning"
	basictypes "github.com/scaling-lightning/scaling-lightning/pkg/types"
)

func (s *commonServer) WalletBalance(
	ctx context.Context,
	in *stdcommonclient.WalletBalanceRequest,
) (*stdcommonclient.WalletBalanceResponse, error) {
	balance, err := s.client.WalletBalance(context.Background(), &lnrpc.WalletBalanceRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "Getting wallet balance from LND's gRPC")
	}

	return &stdcommonclient.WalletBalanceResponse{BalanceSats: uint64(balance.TotalBalance)}, nil //nolint:gosec
}

func (s *commonServer) NewAddress(
	ctx context.Context,
	in *stdcommonclient.NewAddressRequest,
) (*stdcommonclient.NewAddressResponse, error) {
	response, err := s.client.NewAddress(context.Background(), &lnrpc.NewAddressRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "Getting new address from LND's gRPC")
	}
	return &stdcommonclient.NewAddressResponse{Address: response.Address}, nil
}

func (s *lightningServer) PubKey(
	ctx context.Context,
	in *stdlightningclient.PubKeyRequest,
) (*stdlightningclient.PubKeyResponse, error) {
	response, err := s.lightningClient.GetInfo(context.Background(), &lnrpc.GetInfoRequest{})
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
	_, err := s.lightningClient.ConnectPeer(
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
	chanPoint, err := s.lightningClient.OpenChannelSync(
		context.Background(),
		&lnrpc.OpenChannelRequest{
			NodePubkey:         req.PubKey,
			LocalFundingAmount: int64(req.LocalAmtSats), //nolint:gosec
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
	response, err := s.lightningClient.AddInvoice(
		context.Background(),
		&lnrpc.Invoice{
			Value: int64(req.AmtSats), //nolint:gosec
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
	stream, err := s.routerClient.SendPaymentV2(
		context.Background(),
		&routerrpc.SendPaymentRequest{
			PaymentRequest: req.Invoice,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Paying invoice via LND's gRPC")
	}

	for {
		select {
		case <-ctx.Done():
			return nil, errors.Wrap(ctx.Err(), "Context cancelled while waiting for payment response")
		default:
		}

		resp, err := stream.Recv()
		if err != nil {
			return nil, errors.Wrap(err, "Receiving response from LND's gRPC stream")
		}

		if resp == nil {
			log.Warn().Msg("Received nil response from LND's SendPaymentV2 gRPC stream, continuing to wait for response")
			continue
		}

		if resp.Status == lnrpc.Payment_SUCCEEDED {
			preimageBytes, err := hex.DecodeString(resp.PaymentPreimage)
			if err != nil {
				return nil, errors.Wrap(err, "Decoding payment preimage from hex string after payment succeeded")
			}

			return &stdlightningclient.PayInvoiceResponse{
				PaymentPreimage: preimageBytes,
			}, nil
		} else if resp.Status == lnrpc.Payment_FAILED {
			return nil, errors.Newf("Payment failed with failure reason: %v", resp.FailureReason)
		}
	}
}

func (s *lightningServer) ChannelBalance(
	ctx context.Context,
	req *stdlightningclient.ChannelBalanceRequest,
) (*stdlightningclient.ChannelBalanceResponse, error) {
	response, err := s.lightningClient.ChannelBalance(context.Background(), &lnrpc.ChannelBalanceRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "Getting channel balance via LND's gRPC")
	}
	return &stdlightningclient.ChannelBalanceResponse{
		BalanceSats: uint64(response.LocalBalance.Sat),
	}, nil
}
