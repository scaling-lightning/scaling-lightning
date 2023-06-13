package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/lightninglabs/lndclient"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/rs/zerolog/log"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools"
)

const walletName = "scalinglightning"

type appConfig struct {
	tlsFilePath      string
	macaroonFilePath string
	grpcPort         int
	grpcAddress      string
}

var helpRequested = errors.New("Help requested")

func main() {

	appConfig := appConfig{}

	err := parseFlags(&appConfig)
	if err != nil {
		if errors.Is(err, helpRequested) {
			flag.Usage()
			return
		}
		log.Fatal().Err(err).Msg("Problem parsing flags")
	}

	var client lnrpc.LightningClient
	tools.Retry(func() error {

		grpc := fmt.Sprintf("%s:%d", appConfig.grpcAddress, appConfig.grpcPort)
		client, err = lndclient.NewBasicClient(grpc, appConfig.tlsFilePath, appConfig.macaroonFilePath, "regtest")
		if err != nil {
			log.Warn().Err(err).Msg("Problem when connecting to LND's gRPC, perhaps it's not ready")
			return errors.Wrap(err, "New basic client fail")
		}
		return nil

	}, 15*time.Second, 5*time.Minute)

	response, err := client.WalletBalance(context.Background(), &lnrpc.WalletBalanceRequest{})
	if err != nil {
		log.Error().Err(err).Msg("Problem getting wallet balance")
	}
	log.Info().Msgf("Wallet balance is: %v", response.TotalBalance)

	log.Info().Msg("Waiting for command")

	for {
	}
}

func parseFlags(appConfig *appConfig) error {
	var help = flag.Bool("help", false, "Show help")

	flag.StringVar(&appConfig.tlsFilePath, "tlsfilepath", "", "File location for LND's tls file")
	flag.StringVar(&appConfig.macaroonFilePath, "macaroonfilepath", "", "File location for LND's macaroon file")
	flag.IntVar(&appConfig.grpcPort, "grpcport", 10009, "Optional: LND's gRPC port")
	flag.StringVar(&appConfig.grpcAddress, "grpcaddress", "", "LND's gRPC address")

	flag.Parse()

	if *help {
		return helpRequested
	}

	return validateFlags(appConfig)
}

func validateFlags(appConfig *appConfig) error {
	if appConfig.tlsFilePath == "" {
		return errors.New("TLS file path required. Please use the -tlsfilepath flag")
	}
	if appConfig.macaroonFilePath == "" {
		return errors.New("Macaroon file path required. Please use the -macaroon flag")
	}
	if appConfig.grpcAddress == "" {
		return errors.New("gRPC address required. Please use the -grpcaddress flag")
	}
	return nil
}
