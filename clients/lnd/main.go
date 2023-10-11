package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/lightninglabs/lndclient"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/rs/zerolog/log"
	stdcommonclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/common"
	stdlightningclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/lightning"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools"
	"google.golang.org/grpc"
)

type appConfig struct {
	tlsFilePath      string
	macaroonFilePath string
	grpcPort         int
	grpcAddress      string
	apiPort          int
}

func main() {

	var helpRequested = errors.New("Help requested")

	appConfig := appConfig{}

	err := parseFlags(&appConfig, helpRequested)
	if err != nil {
		if errors.Is(err, helpRequested) {
			flag.Usage()
			return
		}
		log.Fatal().Err(err).Msg("Problem parsing flags")
	}

	var client lnrpc.LightningClient
	err = tools.Retry(func(cancel context.CancelFunc) error {
		grpc := fmt.Sprintf("%s:%d", appConfig.grpcAddress, appConfig.grpcPort)
		client, err = lndclient.NewBasicClient(
			grpc,
			appConfig.tlsFilePath,
			appConfig.macaroonFilePath,
			"regtest",
		)
		if err != nil {
			log.Warn().Err(err).Msg("Problem when connecting to LND's gRPC, perhaps it's not ready")
			return errors.Wrap(err, "New basic client fail")
		}
		return nil
	}, 15*time.Second, 5*time.Minute)
	if err != nil {
		log.Fatal().Err(err).Msg("Starting LND Client")
	}

	log.Info().Msg("Waiting for command")

	// start api
	err = startGRPCServer(appConfig.apiPort, client)
	if err != nil {
		log.Fatal().Err(err).Msg("Starting gRPC api server")
	}
}

type lightningServer struct {
	stdlightningclient.UnimplementedLightningServer
	client lnrpc.LightningClient
}
type commonServer struct {
	stdcommonclient.UnimplementedCommonServer
	client lnrpc.LightningClient
}

func startGRPCServer(port int, client lnrpc.LightningClient) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.Wrapf(err, "Listening on port %d", port)
	}
	s := grpc.NewServer()
	stdcommonclient.RegisterCommonServer(s, &commonServer{client: client})
	stdlightningclient.RegisterLightningServer(s, &lightningServer{client: client})

	log.Info().Msgf("Starting gRPC server on port %d", port)
	if err := s.Serve(lis); err != nil {
		return errors.Wrap(err, "Serving gRPC server")
	}
	return nil
}

func parseFlags(appConfig *appConfig, helpRequested error) error {
	var help = flag.Bool("help", false, "Show help")

	flag.StringVar(&appConfig.tlsFilePath, "tlsfilepath", "", "File location for LND's tls file")
	flag.StringVar(
		&appConfig.macaroonFilePath,
		"macaroonfilepath",
		"",
		"File location for LND's macaroon file",
	)
	flag.IntVar(&appConfig.grpcPort, "grpcport", 10009, "Optional: LND's gRPC port")
	flag.StringVar(&appConfig.grpcAddress, "grpcaddress", "", "LND's gRPC address")
	flag.IntVar(&appConfig.apiPort, "apiport", 8181, "Optional: Port to run gRPC API on")

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
