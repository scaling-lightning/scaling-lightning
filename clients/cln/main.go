package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	clnGRPC "github.com/scaling-lightning/scaling-lightning/clients/cln/grpc"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/lightning"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools/grpc_helpers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type appConfig struct {
	tlsFilePath      string
	macaroonFilePath string
	grpcPort         int
	grpcAddress      string
	apiPort          int
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

	err = tools.Retry(func() error {

		_ = fmt.Sprintf("%s:%d", appConfig.grpcAddress, appConfig.grpcPort)
		return nil

	}, 15*time.Second, 5*time.Minute)
	if err != nil {
		log.Fatal().Err(err).Msg("Starting CLN Client")
	}

	cert, err := os.ReadFile("./client.pem")
	if err != nil {
		log.Fatal().Err(err).Msg("Problem reading client certificate")
	}

	certKey, err := os.ReadFile("./client-key.pem")
	if err != nil {
		log.Fatal().Err(err).Msg("Problem reading client key")
	}

	ca, err := os.ReadFile("./ca.pem")
	if err != nil {
		log.Fatal().Err(err).Msg("Problem reading certificate authority cert")
	}

	conn, err := grpcConnect("localhost:32274", cert, certKey, ca)
	if err != nil {
		log.Fatal().Err(err).Msg("Problem connecting to CLN's gRPC server")
	}

	client := clnGRPC.NewNodeClient(conn)
	info, err := client.Getinfo(context.Background(), &clnGRPC.GetinfoRequest{})

	log.Info().Msg("CLN Info:")

	log.Info().Msg(info.String())

	log.Info().Msg("Waiting for command")
	// start api
	restServer := lightning.NewStandardClient()
	registerHandlers(restServer, client)
	err = restServer.Start(appConfig.apiPort)
	if err != nil {
		log.Fatal().Err(err).Msg("Starting REST service")
	}

}

func parseFlags(appConfig *appConfig) error {
	var help = flag.Bool("help", false, "Show help")

	flag.StringVar(&appConfig.tlsFilePath, "tlsfilepath", "", "File location for CLN's tls file")
	flag.StringVar(&appConfig.macaroonFilePath, "macaroonfilepath", "", "File location for CLN's macaroon file")
	flag.IntVar(&appConfig.grpcPort, "grpcport", 10009, "Optional: CLN's gRPC port")
	flag.StringVar(&appConfig.grpcAddress, "grpcaddress", "", "CLN's gRPC address")
	flag.IntVar(&appConfig.apiPort, "apiport", 8181, "Optional: Port to run REST API on")

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

func grpcConnect(host string,
	certificate []byte,
	key []byte,
	caCertificate []byte) (*grpc.ClientConn, error) {

	clientCrt, err := tls.X509KeyPair(certificate, key)
	if err != nil {
		return nil, errors.New("CLN credentials: failed to create X509 KeyPair")
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCertificate)

	serverName := "localhost"
	if strings.Contains(host, "cln") {
		serverName = "cln"
	}

	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		ClientAuth:   tls.RequestClientCert,
		Certificates: []tls.Certificate{clientCrt},
		RootCAs:      certPool,
		ServerName:   serverName,
	}

	loggerOpts := grpc_helpers.GetLoggingOptions()

	logger := zerolog.New(os.Stderr)

	opts := []grpc.DialOption{
		grpc.WithReturnConnectionError(),
		grpc.FailOnNonTempDialError(true),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpc_helpers.RecvMsgSize)),
		grpc.WithChainUnaryInterceptor(
			timeout.UnaryClientInterceptor(grpc_helpers.UnaryTimeout),
			logging.UnaryClientInterceptor(grpc_helpers.InterceptorLogger(logger), loggerOpts...),
		),
		grpc.WithChainStreamInterceptor(
			logging.StreamClientInterceptor(grpc_helpers.InterceptorLogger(logger), loggerOpts...),
		),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, opts...)
	if err != nil {
		return nil, fmt.Errorf("cannot dial to CLN %v", err)
	}

	return conn, nil
}
