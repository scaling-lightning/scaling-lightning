package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/scaling-lightning/scaling-lightning/pkg/types"
	"github.com/spf13/cobra"
)

var openchannelFromName string
var openchannelToName string
var openchannelLocalAmt uint64

var openchannelCmd = &cobra.Command{
	Use:   "openchannel",
	Short: "Open a channel between two nodes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		slnetwork, err := sl.DiscoverStartedNetwork("")
		if err != nil {
			fmt.Printf(
				"Problem with network discovery, is there a network running? Error: %v\n",
				err.Error(),
			)
			return
		}
		var openchannelFromNode sl.LightningNode
		var openchannelToNode sl.LightningNode
		for _, node := range slnetwork.LightningNodes {
			if node.GetName() == openchannelFromName {
				openchannelFromNode = node
				continue
			}
			if node.GetName() == openchannelToName {
				openchannelToNode = node
			}
		}
		allNames := []string{}
		for _, node := range slnetwork.LightningNodes {
			allNames = append(allNames, node.GetName())
		}
		if openchannelFromNode.Name == "" {
			fmt.Printf(
				"Can't find node with name %v, here are the lightnign nodes that are running: %v\n",
				openchannelFromName,
				allNames,
			)
		}
		if openchannelToNode.Name == "" {
			fmt.Printf(
				"Can't find node with name %v, here are the lightning nodes that are running: %v\n",
				openchannelToName,
				allNames,
			)
		}

		err = openchannelFromNode.OpenChannel(
			&openchannelToNode,
			types.NewAmountSats(openchannelLocalAmt),
		)
		if err != nil {
			fmt.Printf("Problem opening channel: %v\n", err.Error())
			return
		}

		fmt.Println("Open channel command received")
	},
}

func init() {
	rootCmd.AddCommand(openchannelCmd)

	openchannelCmd.Flags().
		StringVarP(&openchannelFromName, "from", "f", "", "Name of node to open channel from")
	openchannelCmd.MarkFlagRequired("from")

	openchannelCmd.Flags().
		StringVarP(&openchannelToName, "to", "t", "", "Name of node to open channel to")
	openchannelCmd.MarkFlagRequired("to")

	openchannelCmd.Flags().
		Uint64VarP(&openchannelLocalAmt, "amount", "a", 0, "Amount of satoshis to put into channel from the opening side")
	openchannelCmd.MarkFlagRequired("amount")

}
