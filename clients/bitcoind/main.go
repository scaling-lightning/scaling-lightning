package main

import (
	"flag"
	"fmt"
	"github.com/btcsuite/btcd/rpcclient"
	"log"
	"os"
)

type appConfig struct {
	rpcCookieFile string
	rpcHost       string
	rpcPort       int
	chain         string
}

func main() {

	appConfig := appConfig{}
	parseFlags(&appConfig)

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
		log.Fatal(err)
	}
	defer client.Shutdown()

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)
}

func parseFlags(appConfig *appConfig) error {
	var help = flag.Bool("help", false, "Show help")

	flag.StringVar(&appConfig.rpcCookieFile, "rpccookiefile", "", "File location for Bitcoind's .cookie file")
	flag.StringVar(&appConfig.rpcHost, "rpchost", "", "Bitcoind's RPC host")
	flag.IntVar(&appConfig.rpcPort, "rpcport", 0, "Optional: Bitcoind's RPC port. Will use defaults specified by -chain if not set. ")
	flag.StringVar(&appConfig.chain, "chain", "regtest", "Current chain. Valid options: regtest, signet.")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	return nil
}
