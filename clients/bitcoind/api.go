package main

import (
	"context"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cockroachdb/errors"
	stdbitcoinclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/bitcoin"
	stdcommonclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/common"
)

func (s *commonServer) WalletBalance(ctx context.Context,
	in *stdcommonclient.WalletBalanceRequest,
) (*stdcommonclient.WalletBalanceResponse, error) {
	walletBalance, err := s.client.GetBalance("*")
	if err != nil {
		return nil, errors.Wrap(err, "Getting wallet balance from Bitcoin RPC")
	}
	return &stdcommonclient.WalletBalanceResponse{BalanceSats: uint64(walletBalance)}, nil
}

func (s *commonServer) NewAddress(
	ctx context.Context,
	in *stdcommonclient.NewAddressRequest,
) (*stdcommonclient.NewAddressResponse, error) {
	newAddress, err := s.client.GetNewAddress(walletName)
	if err != nil {
		return nil, errors.Wrap(err, "Getting new address from Bitcoin RPC")
	}
	return &stdcommonclient.NewAddressResponse{Address: newAddress.String()}, nil
}

func (s *commonServer) Send(
	ctx context.Context,
	req *stdcommonclient.SendRequest,
) (*stdcommonclient.SendResponse, error) {

	// TODO: pass in real network
	newAddress, err := btcutil.DecodeAddress(
		req.Address,
		&chaincfg.Params{Name: "regtest"},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Decoding address")
	}

	txid, err := s.client.SendToAddress(
		newAddress,
		btcutil.Amount(req.Amount),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Sending to address from Bitcoin RPC")
	}
	return &stdcommonclient.SendResponse{TxId: txid.String()}, nil
}

func (s *bitcoinServer) GenerateToAddress(
	ctx context.Context,
	req *stdbitcoinclient.GenerateToAddressRequest,
) (*stdbitcoinclient.GenerateToAddressResponse, error) {
	// TODO: pass in real network
	address, err := btcutil.DecodeAddress(
		req.Address,
		&chaincfg.Params{Name: "regtest"},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Decoding address")
	}
	genHashes, err := s.client.GenerateToAddress(
		int64(req.NumOfBlocks),
		address,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Generating to address from Bitcoin RPC")
	}

	hashes := []string{}
	for _, hash := range genHashes {
		hashes = append(hashes, hash.String()+"\n")
	}
	return &stdbitcoinclient.GenerateToAddressResponse{Hashes: hashes}, nil
}
