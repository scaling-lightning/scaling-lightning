package initialstate

import (
	"os"

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

type yamlCommands []map[string][]map[string]any
type command struct {
	commandType string
	args map[string]interface{}
}

type initialState struct {
	commands []command
	network SLNetworkInterface
}

func NewInitialStateFromBytes(yamlBytes []byte, network SLNetworkInterface) (*initialState, error) {
	initialState := initialState{}
	err := initialState.parseYAML(yamlBytes)
	if err != nil {
		return nil, errors.Wrap(err, "Parsing initial state yaml")
	}
	initialState.network = network
	return &initialState, nil
}

func NewInitialStateFromFile(yamlFilePath string, network SLNetworkInterface) (*initialState, error) {
	yamlBytes, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "Reading yaml file")
	}
	return NewInitialStateFromBytes(yamlBytes, network)
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
func (is *initialState) Apply() error {
	for _, command := range is.commands {
			switch command.commandType {
			case "SendOnChain":
				_, err := is.network.Send(
					command.args["from"].(string),
					command.args["to"].(string),
					uint64(command.args["amountSats"].(int)))
				if err != nil {
					return errors.Wrap(err, "Sending on chain")
				}
			case "ConnectPeer":
				err := is.network.ConnectPeer(
					command.args["from"].(string),
					command.args["to"].(string))
				if err != nil {
					return errors.Wrap(err, "Connecting peer")
				}
			case "OpenChannels":
				_, err := is.network.OpenChannel(
					command.args["from"].(string),
					command.args["to"].(string),
					uint64(command.args["localAmountSats"].(int)))
				if err != nil {
					return errors.Wrap(err, "Opening channels")
				}
			case "SendOverChannel":
				invoice, err := is.network.CreateInvoice(
					command.args["to"].(string),
					uint64(command.args["amountMSat"].(int))/1000)
				if err != nil {
					return errors.Wrap(err, "Creating invoice")
				}
				_, err = is.network.PayInvoice(
					command.args["from"].(string),
					invoice)
				if err != nil {
					return errors.Wrap(err, "Paying invoice")
				}
			default:
				return errors.Errorf("Unknown command %v", command.commandType)
			}
	}
	return nil
}
