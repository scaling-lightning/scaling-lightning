package scalinglightning

import (
	"fmt"
	"log"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

func init() {

	var channelbalanceCmd = &cobra.Command{
		Use:   "channelbalance",
		Short: "Get the onchain wallet balance of a node",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			processDebugFlag(cmd)
			nodeName := cmd.Flag("node").Value.String()
			slnetwork, err := sl.DiscoverRunningNetwork(kubeConfigPath, apiHost, apiPort, namespace)
			if err != nil {
				fmt.Printf(
					"Problem with network discovery, is there a network running? Error: %v\n",
					err.Error(),
				)
				return
			}
			channelBalance, err := slnetwork.ChannelBalance(nodeName)
			if err != nil {
				fmt.Printf("Problem getting wallet balance: %v\n", err.Error())
				return
			}
			fmt.Printf("%d sats\n", channelBalance.AsSats())
		},
	}

	rootCmd.AddCommand(channelbalanceCmd)

	channelbalanceCmd.Flags().
		StringP("node", "n", "", "The name of the node to get the wallet balance of")
	err := channelbalanceCmd.MarkFlagRequired("node")
	if err != nil {
		log.Fatalf("Problem marking node flag as required: %v", err.Error())
	}
}
