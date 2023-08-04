package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/scaling-lightning/scaling-lightning/pkg/types"
	"github.com/spf13/cobra"
)

var connectpeerFromName string
var connectpeerToName string

var connectpeerCmd = &cobra.Command{
	Use:   "connectpeer",
	Short: "Connect peers",
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
		var connectpeerFromNode sl.LightningNode
		var connectpeerToNode sl.LightningNode
		for _, node := range slnetwork.LightningNodes {
			if node.GetName() == connectpeerFromName {
				connectpeerFromNode = node
				continue
			}
			if node.GetName() == connectpeerToName {
				connectpeerToNode = node
			}
		}
		if connectpeerFromNode.Name == "" {
			fmt.Printf(
				"Can't find node with name %v, here are the lightnign nodes that are running: %v\n",
				openchannelFromName,
				slnetwork.LightningNodes,
			)
		}
		if connectpeerToNode.Name == "" {
			fmt.Printf(
				"Can't find node with name %v, here are the lightning nodes that are running: %v\n",
				openchannelToName,
				slnetwork.LightningNodes,
			)
		}

		err = connectpeerFromNode.OpenChannel(
			&connectpeerToNode,
			types.NewAmountSats(openchannelLocalAmt),
		)
		if err != nil {
			fmt.Printf("Problem sending funds: %v\n", err.Error())
			return
		}

		fmt.Println("Connect peer command received")
	},
}

func init() {
	rootCmd.AddCommand(connectpeerCmd)

	connectpeerCmd.Flags().
		StringVarP(&connectpeerFromName, "from", "f", "", "Name of the node to connect from")
	connectpeerCmd.MarkFlagRequired("from")

	connectpeerCmd.Flags().
		StringVarP(&connectpeerToName, "to", "t", "", "Name of the node to connect from")
	connectpeerCmd.MarkFlagRequired("from")
}
