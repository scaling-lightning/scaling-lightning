package main

import (
	"flag"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
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
