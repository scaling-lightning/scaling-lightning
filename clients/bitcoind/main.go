package main

import (
	"flag"
	"fmt"
	"github.com/btcsuite/btcd/rpcclient"
	"log"
	"os"
)

func main() {
	var help = flag.Bool("help", false, "Show help")

	rpcCookieFile := flag.String("rpccookiefile", "", "File location for Bitcoind's .cookie file")
	rpcHost := flag.String("rpchost", "", "Bitcoind's RPC host")
	rpcPort := flag.String("rpcport", "", "Optional: Bitcoind's RPC port. Will use defaults specified by -chain if not set. ")
	chain := flag.String("chain", "regtest", "Current chain. Valid options: regtest, signet.")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	fmt.Printf("Chain is: %v\n", *chain)
	fmt.Printf("RPCCookieFile is: %v\n", *rpcCookieFile)
	fmt.Printf("RPCHost is: %v\n", *rpcHost)
	fmt.Printf("RPCPort is: %v\n", *rpcPort)

	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:18443",
		User:         "foo",
		Pass:         "pass",
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
