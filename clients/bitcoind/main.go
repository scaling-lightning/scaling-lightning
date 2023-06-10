package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools"
)

const walletName = "scalinglightning"

type appConfig struct {
	rpcCookieFile string
	rpcHost       string
	rpcPort       int
	chain         string
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

	host := fmt.Sprintf("%s:%d", appConfig.rpcHost, appConfig.rpcPort)
	connCfg := &rpcclient.ConnConfig{
		Host:         host,
		CookiePath:   appConfig.rpcCookieFile,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}

	// Notification parameter is nil since notifications are not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Creating new rpc client")
	}
	defer client.Shutdown()

	log.Info().Msg("Attempting to initialise bitcoind")

	err = tools.Retry(func() error {
		err := prepareBitcoind(client)
		if err != nil {
			log.Warn().Err(err).Msg("Problem when initialising bitcoind, perhaps it's not ready")
			return errors.Wrap(err, "Preparing bitcoind")
		}
		return nil
	}, 10*time.Second, 5*time.Minute)
	if err != nil {
		log.Fatal().Err(err).Msg("Preparing bitcoind")
	}

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Warn().Err(err).Send()
	}
	log.Info().Msgf("Block count: %d", blockCount)

	log.Info().Msg("Waiting for command")

	for {
	}
}

func parseFlags(appConfig *appConfig) error {
	var help = flag.Bool("help", false, "Show help")

	flag.StringVar(&appConfig.rpcCookieFile, "rpccookiefile", "", "File location for Bitcoind's .cookie file")
	flag.StringVar(&appConfig.rpcHost, "rpchost", "", "Bitcoind's RPC host")
	flag.IntVar(&appConfig.rpcPort, "rpcport", 0, "Optional: Bitcoind's RPC port. Will use defaults specified by -chain if not set. ")
	flag.StringVar(&appConfig.chain, "chain", "regtest", "Current chain. Valid options: regtest, signet.")

	flag.Parse()

	if *help {
		return helpRequested
	}

	return validateFlags(appConfig)
}

func validateFlags(appConfig *appConfig) error {
	if appConfig.rpcCookieFile == "" {
		return errors.New("RPC Cookie File location required, please use the -rpccookiefile flag")
	}
	if appConfig.rpcHost == "" {
		return errors.New("RPC Host required, please use the -rpchost flag")
	}
	if appConfig.rpcPort == 0 {
		switch appConfig.chain {
		case "regtest":
			appConfig.rpcPort = 18443
		case "signet":
			appConfig.rpcPort = 38332
		}
	}
	return nil
}
