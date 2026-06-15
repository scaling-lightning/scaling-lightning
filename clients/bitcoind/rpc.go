package main

import (
	"context"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
)

type RpcClient interface {
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
	GetBlockCount() (int64, error)
}

const generateUntilBlockheight = int64(1000)
const maxGenerateAtOnce = int64(100)

func initialiseBitcoind(client RpcClient) error {
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

	err = mineUntilTarget(client, walletInfo)
	if err != nil {
		return errors.Wrap(err, "Mining blocks until target block height")
	}

	log.Info().Msg("Bitcoind is setup and ready to serve clients")

	return nil
}

func mineUntilTarget(client RpcClient, walletInfo *btcjson.GetWalletInfoResult) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	var address btcutil.Address

	for {
		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "Timed out")
		default:
		}

		currentHeight, err := client.GetBlockCount()
		if err != nil {
			return errors.Wrap(err, "Getting mining info")
		}

		blocksToMine := generateUntilBlockheight - currentHeight

		if blocksToMine > maxGenerateAtOnce {
			blocksToMine = maxGenerateAtOnce
		}

		if blocksToMine <= 0 {
			log.Info().Msgf("Current block height %v is already above %v, no need to mine new blocks",
				currentHeight, generateUntilBlockheight)
			return nil
		}

		if address.String() == "" {
			address, err = client.GetNewAddress(walletInfo.WalletName)
			if err != nil {
				return errors.Wrap(err, "Getting new address")
			}

			log.Info().Msgf("New address created to receive mined coins: %v", address.String())
		}

		log.Info().Msgf("Mining %v blocks. Current height: %v. Target height %v.",
			blocksToMine, currentHeight, generateUntilBlockheight)

		_, err = client.GenerateToAddress(blocksToMine, address, nil)
		if err != nil {
			return errors.Wrap(err, "Generating blocks")
		}
	}
}
