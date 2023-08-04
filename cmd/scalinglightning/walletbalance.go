package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

var walletbalanceCmd = &cobra.Command{
	Use:   "walletbalance",
	Short: "Get the onchain wallet balance of a node",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		balanceNodeName := cmd.Flag("node").Value.String()
		slnetwork, err := sl.DiscoverStartedNetwork("")
		if err != nil {
			fmt.Printf(
				"Problem with network discovery, is there a network running? Error: %v\n",
				err.Error(),
			)
			return
		}
		allNodes := slnetwork.GetAllNodes()
		for _, node := range allNodes {
			if node.GetName() == balanceNodeName {
				walletBalance, err := node.GetWalletBalance()
				if err != nil {
					fmt.Printf("Problem getting wallet balance: %v\n", err.Error())
					return
				}
				fmt.Printf("%d sats\n", walletBalance.AsSats())
				return
			}
		}
		allNames := []string{}
		for _, node := range allNodes {
			allNames = append(allNames, node.GetName())
		}
		fmt.Printf(
			"Can't find node with name %v, here are the nodes that are running: %v\n",
			balanceNodeName,
			allNames,
		)
	},
}

func init() {
	rootCmd.AddCommand(walletbalanceCmd)

	walletbalanceCmd.Flags().
		StringP("node", "n", "", "The name of the node to get the wallet balance of")
	walletbalanceCmd.MarkFlagRequired("node")
}
