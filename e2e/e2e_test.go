package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/scaling-lightning/scaling-lightning/pkg/initialstate"
	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools"
)

// will need a longish (few mins) timeout
func Test_E2E(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// use locally built helm files, so must have them built first!
	network, err := sl.NewSLNetwork("./e2e.yaml", "", sl.Regtest, sl.DefaultNamespace)
	if err != nil {
		t.Fatalf("Problem creating network: %v", err)
	}

	err = network.CreateAndStart()
	if err != nil {
		t.Fatalf("Problem starting network: %v", err)
	}

	defer func() {
		err = network.Destroy()
		if err != nil {
			t.Errorf("Problem destroying network: %v", err)
		}
	}()

	// check that the network is running
	waitNodeRunning(t, network, "cln1")
	waitNodeRunning(t, network, "cln2")
	waitNodeRunning(t, network, "bitcoind")

	cln1, err := network.GetLightningNode("cln1")
	if err != nil {
		t.Fatalf("Problem getting cln1 node: %v", err)
	}

	cln2, err := network.GetLightningNode("cln2")
	if err != nil {
		t.Fatalf("Problem getting cln2 node: %v", err)
	}

	bitcoind, err := network.GetBitcoinNode("bitcoind")
	if err != nil {
		t.Fatalf("Problem getting bitcoind: %v", err)
	}

	newAddr, err := network.GetNewAddress(bitcoind.GetName())
	if err != nil {
		t.Fatalf("Problem getting new address: %v", err)
	}

	log.Info().Msgf("got new bitcoind address: %s", newAddr)

	cln1Pubkey, err := network.GetPubKey(cln1.GetName())
	if err != nil {
		t.Fatalf("Problem getting cln1 pubkey: %v", err)
	}

	log.Info().Msgf("got cln1 pubkey: %s", cln1Pubkey.AsHexString())

	cln2Pubkey, err := network.GetPubKey(cln2.GetName())
	if err != nil {
		t.Fatalf("Problem getting cln2 pubkey: %v", err)
	}

	log.Info().Msgf("got cln2 pubkey: %s", cln2Pubkey.AsHexString())

	// TODO add LND nodes also

	const initialStateYAML = `
- SendOnChain:
    - { from: bitcoind, to: cln1, amountSats: 1_000_000 }
- ConnectPeer:
    - { from: cln1, to: cln2 }
- OpenChannels:
    - { from: cln1, to: cln2, localAmountSats: 200_000 }
    - { from: cln1, to: cln2, localAmountSats: 300_000 }
- SendOverChannel:
    - { from: cln1, to: cln2, amountMSat: 2_000_000 }
`
	initialState, err := initialstate.NewInitialStateFromBytes([]byte(initialStateYAML), &network)
	if err != nil {
		t.Fatalf("Problem creating initial state: %v", err)
	}

	err = initialState.Apply()
	if err != nil {
		t.Fatalf("Problem applying initial state: %v", err)
	}

	balance, err := network.GetWalletBalance(cln1.GetName())
	if err != nil {
		t.Fatalf("Problem get wallet balance: %v", err)
	}

	log.Info().Msgf("cln1 balance: %d", balance.AsSats())

	connectionDetails, err := network.GetConnectionDetails(cln2.GetName())
	if err != nil {
		t.Fatalf("Problem getting connection details: %v", err)
	}

	log.Info().Msgf("cln2 connection host: %v", connectionDetails[0].Host)
	log.Info().Msgf("cln2 connection host: %d", connectionDetails[0].Port)

	connectionFiles, err := cln2.GetConnectionFiles(network.Network.String(), "")
	if err != nil {
		t.Fatalf("Problem getting connection files: %v", err)
	}

	log.Info().Msgf("cln2 client cert size : %v", len(connectionFiles.CLN.ClientCert))
}

func waitNodeRunning(t *testing.T, network sl.SLNetwork, nodeName string) {
	err := tools.Retry(func(cancel context.CancelFunc) error {
		ok, err := network.IsNodeRunning(nodeName)
		if err != nil {
			return errors.Wrap(err, "Checking if node is running")
		}

		if !ok {
			return errors.Errorf("Node %s is not running", nodeName)
		}

		return nil
	}, time.Second, time.Minute)
	if err != nil {
		t.Fatalf("Node did not start: %v", err)
	}
}
