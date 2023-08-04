package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/scaling-lightning/scaling-lightning/pkg/types"
	"github.com/spf13/cobra"
)

var sendFromName string
var sendToName string
var sendAmount uint64

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send on chain funds betwen nodes",
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
		var sendFromNode sl.Node
		var sendToNode sl.Node
		allNodes := slnetwork.GetAllNodes()
		for _, node := range allNodes {
			if node.GetName() == sendFromName {
				sendFromNode = node
				continue
			}
			if node.GetName() == sendToName {
				sendToNode = node
			}
		}
		allNames := []string{}
		for _, node := range allNodes {
			allNames = append(allNames, node.GetName())
		}
		if sendFromNode == nil {
			fmt.Printf(
				"Can't find node with name %v, here are the nodes that are running: %v\n",
				sendFromName,
				allNames,
			)
		}
		if sendToNode == nil {
			fmt.Printf(
				"Can't find node with name %v, here are the nodes that are running: %v\n",
				sendToName,
				allNames,
			)
		}

		err = sendFromNode.Send(sendToNode, types.NewAmountSats(sendAmount))
		if err != nil {
			fmt.Printf("Problem sending funds: %v\n", err.Error())
			return
		}

		fmt.Println("Sent funds")
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringVarP(&sendFromName, "from", "f", "", "Name of node to send from")
	sendCmd.MarkFlagRequired("from")

	sendCmd.Flags().StringVarP(&sendToName, "to", "t", "", "Name of node to send to")
	sendCmd.MarkFlagRequired("to")

	sendCmd.Flags().Uint64VarP(&sendAmount, "amount", "a", 0, "Amount of satoshis to send")
	sendCmd.MarkFlagRequired("amount")

}
