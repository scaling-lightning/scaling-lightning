package initialstate

import (
	"github.com/cockroachdb/errors"
	"github.com/scaling-lightning/scaling-lightning/pkg/types"
	"gopkg.in/yaml.v3"
)

type SLNetworkInterface interface {
	Send(fromNodeName string, toNodeName string, amountSats uint64) (string, error)
	CreateInvoice(nodeName string, amountSats uint64) (string, error)
	PayInvoice(nodeName string, invoice string) (string, error)
	ChannelBalance(nodeName string) (types.Amount, error)
	ConnectPeer(nodeName string, pubkey types.PubKey) error
	OpenChannel(nodeName string, pubkey types.PubKey, localAmt types.Amount) (types.ChannelPoint, error)
}

type initialStateYAML []initialStateCommand
type initialStateCommand map[string][]string

type initialState struct {
	commands initialStateYAML
	network SLNetworkInterface
}

func newInitialState(yamlBytes []byte) (*initialState, error) {
	initialState := initialState{}
	err := initialState.parseYAML(yamlBytes)
	if err != nil {
		return nil, errors.Wrap(err, "Parsing initial state yaml")
	}
	return &initialState, nil
}

func (is *initialState) parseYAML(yamlBytes []byte) error {

	err := yaml.Unmarshal(yamlBytes, &is.commands)
	if err != nil {
		return errors.Wrap(err, "Unmarshalling yaml")
	}
	return nil
}

// function to run through commands and execute them on the network
func (is *initialState) ApplyToNetwork(network string) error {
	for _, command := range is.commands {
		for commandName, args := range command {
			switch commandName {
			case "OpenChannels":
				err := network.OpenChannels(args)
				if err != nil {
					return errors.Wrap(err, "Opening channels")
				}
			// case "CloseChannels":
			// 	err := network.CloseChannels(args)
			// 	if err != nil {
			// 		return errors.Wrap(err, "Closing channels")
			// 	}
			// case "ConnectPeer":
			// 	err := network.ConnectPeer(args)
			// 	if err != nil {
			// 		return errors.Wrap(err, "Connecting peer")
			// 	}
			// case "PayInvoice":
			// 	err := network.PayInvoice(args)
			// 	if err != nil {
			// 		return errors.Wrap(err, "Paying invoice")
			// 	}
			// case "PayOnChain":
			// 	err := network.PayOnChain(args)
			// 	if err != nil {
			// 		return errors.Wrap(err, "Paying on chain")
			// 	}
			default:
				return errors.Errorf("Unknown command %v %v", commandName, args)
			}
		}
	}
	return nil
}
