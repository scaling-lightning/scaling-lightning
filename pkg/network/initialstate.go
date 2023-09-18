package network

import (
	"github.com/cockroachdb/errors"
	"gopkg.in/yaml.v3"
)

type initialStateYAML []initialStateCommand
type initialStateCommand map[string][]string

type initialState struct {
	commands initialStateYAML
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
func (is *initialState) ApplyToNetwork(network *SLNetwork) error {
	for _, command := range is.commands {
		for commandName, args := range command {
			switch commandName {
			// case "OpenChannels":
			// 	err := network.OpenChannels(args)
			// 	if err != nil {
			// 		return errors.Wrap(err, "Opening channels")
			// 	}
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
