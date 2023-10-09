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
	ConnectPeer(fromNodeName string, toNodeName string) error
	OpenChannel(fromNodeName string, toNodeName string, localAmt uint64) (types.ChannelPoint, error)
}

type yamlCommands []map[string][]map[string]interface{}
type command struct {
	commandType string
	args map[string]interface{}
}

type initialState struct {
	commands []command
	network SLNetworkInterface
}

func NewInitialState(yamlBytes []byte) (*initialState, error) {
	initialState := initialState{}
	err := initialState.parseYAML(yamlBytes)
	if err != nil {
		return nil, errors.Wrap(err, "Parsing initial state yaml")
	}
	return &initialState, nil
}

func (is *initialState) parseYAML(yamlBytes []byte) error {

	yamlData := yamlCommands{}
	err := yaml.Unmarshal(yamlBytes, &yamlData)
	if err != nil {
		return errors.Wrap(err, "Unmarshalling yaml")
	}
	for _, commandType := range yamlData {
		for commandName, commands := range commandType {
			for _, c := range commands {
				is.commands = append(is.commands, command {
					commandType: commandName,
					args: c,
				})
			}
			break
		}
	}

	return nil
}

// function to run through commands and execute them on the network
func (is *initialState) ApplyToNetwork(network string) error {
	for _, command := range is.commands {
			switch command.commandType {
			case "OpenChannels":
				_, err := is.network.OpenChannel("lnd1", "lnd2", 200_000)
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
				return errors.Errorf("Unknown command %v", command.commandType)
			}
	}
	return nil
}
