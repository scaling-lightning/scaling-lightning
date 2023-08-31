package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	clnGRPC "github.com/scaling-lightning/scaling-lightning/clients/cln/grpc"
	standardclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/lightning"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools/grpc_helpers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type appConfig struct {
	clientCertificate string
	clientKey         string
	caCert            string
	grpcPort          int
	grpcAddress       string
	apiPort           int
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

	var client clnGRPC.NodeClient

	err = tools.Retry(func() error {

		cert, err := os.ReadFile(appConfig.clientCertificate)
		if err != nil {
			log.Warn().Err(err).Msg("Problem reading client certificate")
			return err
		}

		certKey, err := os.ReadFile(appConfig.clientKey)
		if err != nil {
			log.Warn().Err(err).Msg("Problem reading client key")
			return err
		}

		ca, err := os.ReadFile(appConfig.caCert)
		if err != nil {
			log.Warn().Err(err).Msg("Problem reading certificate authority cert")
			return err
		}

		conn, err := grpcConnect(
			fmt.Sprintf("%s:%d", appConfig.grpcAddress, appConfig.grpcPort),
			cert,
			certKey,
			ca,
		)
		if err != nil {
			log.Warn().Err(err).Msg("Problem connecting to CLN's gRPC server")
			return err
		}
		client = clnGRPC.NewNodeClient(conn)
		info, err := client.Getinfo(context.Background(), &clnGRPC.GetinfoRequest{})
		if err != nil {
			log.Warn().Err(err).Msg("Problem getting info from CLN's gRPC server")
			return err
		}

		log.Info().Msg("CLN Info:")
		log.Info().Msg(info.String())

		_ = fmt.Sprintf("%s:%d", appConfig.grpcAddress, appConfig.grpcPort)
		return nil

	}, 15*time.Second, 5*time.Minute)
	if err != nil {
		log.Fatal().Err(err).Msg("Starting CLN Client")
	}

	// // start api
	// restServer := lightning.NewStandardClient()
	// registerHandlers(restServer, client)
	// err = restServer.Start(appConfig.apiPort)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Starting REST service")
	// }

	startGRPCServer(appConfig.apiPort, client)
}

type server struct {
	standardclient.UnimplementedLightningClientServer
	client clnGRPC.NodeClient
}

func startGRPCServer(port int, client clnGRPC.NodeClient) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.Wrapf(err, "Listening on port %d", port)
	}
	s := grpc.NewServer()
	standardclient.RegisterLightningClientServer(s, &server{client: client})
	log.Info().Msgf("Starting gRPC server on port %d", port)
	if err := s.Serve(lis); err != nil {
		return errors.Wrap(err, "Serving gRPC server")
	}
	return nil
}

func parseFlags(appConfig *appConfig) error {
	var help = flag.Bool("help", false, "Show help")

	flag.StringVar(
		&appConfig.clientCertificate,
		"clientcert",
		"",
		"File location for CLN's client certificate",
	)
	flag.StringVar(&appConfig.clientKey, "clientkey", "", "File location for CLN's client key")
	flag.StringVar(
		&appConfig.caCert,
		"cacert",
		"",
		"File location for CLN's certificate authority cert",
	)
	flag.IntVar(&appConfig.grpcPort, "grpcport", 8383, "Optional: CLN's gRPC port")
	flag.StringVar(&appConfig.grpcAddress, "grpcaddress", "", "CLN's gRPC address")
	flag.IntVar(&appConfig.apiPort, "apiport", 8181, "Optional: Port to run REST API on")

	flag.Parse()

	if *help {
		return helpRequested
	}

	return validateFlags(appConfig)
}

func validateFlags(appConfig *appConfig) error {
	if appConfig.clientCertificate == "" {
		return errors.New("Client certificate required. Please use the -clientcert flag")
	}
	if appConfig.clientKey == "" {
		return errors.New("Client key file path required. Please use the -clientkey flag")
	}
	if appConfig.caCert == "" {
		return errors.New(
			"Certificate authoritiy cert file path required. Please use the -cacert flag",
		)
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
