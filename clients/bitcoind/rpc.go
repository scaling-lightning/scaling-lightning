package main

import (
	"strings"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
)

//go:generate mockery --name rpcClient --exported
type rpcClient interface {
	CreateWallet(
		name string,
		opts ...rpcclient.CreateWalletOpt,
	) (*btcjson.CreateWalletResult, error)
	LoadWallet(name string) (*btcjson.LoadWalletResult, error)
	GenerateToAddress(
		numBlocks int64,
		address btcutil.Address,
		maxTries *int64,
	) ([]*chainhash.Hash, error)
	GetBalance(account string) (btcutil.Amount, error)
	GetNewAddress(account string) (btcutil.Address, error)
	GetWalletInfo() (*btcjson.GetWalletInfoResult, error)
	SendToAddress(address btcutil.Address, amount btcutil.Amount) (*chainhash.Hash, error)
}

func initialiseBitcoind(client rpcClient) error {
	walletInfo, err := client.GetWalletInfo()
	if err != nil && !strings.Contains(err.Error(), "No wallet is loaded") {
		return errors.Wrap(err, "Getting wallet info")
	}

	if walletInfo == nil || walletInfo.WalletName == "" {
		log.Info().Msg("No wallet loaded, trying to load scalinglightning wallet")
		loadWalletResult, err := client.LoadWallet(walletName)
		if err != nil {
			log.Info().Msg("Couldn't load scalinglightning wallet, trying to create it")
			if err != nil {
				log.Info().Msgf("Load wallet err was: %v", err.Error())
			}
			_, err := client.CreateWallet(walletName)
			if err != nil {
				return errors.Wrap(err, "Creating bitcoind wallet")
			}
		}
		if loadWalletResult != nil && loadWalletResult.Warning != "" {
			log.Info().Msgf("Load wallet warning was: %v", loadWalletResult.Warning)
		}
	}

	walletInfo, err = client.GetWalletInfo()
	if err != nil {
		return errors.Wrap(err, "Getting wallet info, wallet should exist by now")
	}

	log.Info().Msgf("Loaded wallet: %v", walletInfo.WalletName)

	address, err := client.GetNewAddress(walletInfo.WalletName)
	if err != nil {
		return errors.Wrap(err, "Getting new address")
	}

	log.Info().Msgf("New address created to receive mined coins: %v", address.String())

	// maxTries := int64(1000000)
	_, err = client.GenerateToAddress(1000, address, nil)
	if err != nil {
		return errors.Wrapf(err, "Generating blocks to address: %v", address.String())
	}

	log.Info().Msg("Bitcoind is setup and ready to serve clients")
	return nil
}
