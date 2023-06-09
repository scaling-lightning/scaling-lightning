package main

import (
	"strings"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
)

func prepareBitcoind(client *rpcclient.Client) error {
	walletInfo, err := client.GetWalletInfo()
	if err != nil && !strings.Contains(err.Error(), "No wallet is loaded") {
		return errors.Wrap(err, "Getting wallet info")
	}

	if walletInfo == nil || walletInfo.WalletName == "" {
		_, err := client.CreateWallet(walletName)
		if err != nil {
			return errors.Wrap(err, "Creating bitcoind wallet")
		}
	}

	walletInfo, err = client.GetWalletInfo()
	if err != nil && !strings.Contains(err.Error(), "No wallet is loaded") {
		return errors.Wrap(err, "Getting wallet info")
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
