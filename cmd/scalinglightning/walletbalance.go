package scalinglightning

import (
	"fmt"
	"log"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

func init() {
	var walletbalanceCmd = &cobra.Command{
		Use:   "walletbalance",
		Short: "Get the onchain wallet balance of a node",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			processDebugFlag(cmd)
			balanceNodeName := cmd.Flag("node").Value.String()
			slnetwork, err := sl.DiscoverStartedNetwork(kubeConfigPath, apiHost, apiPort)
			if err != nil {
				fmt.Printf(
					"Problem with network discovery, is there a network running? Error: %v\n",
					err.Error(),
				)
				return
			}

			walletBalance, err := slnetwork.GetWalletBalance(balanceNodeName)
			if err != nil {
				fmt.Printf("Problem getting wallet balance: %v\n", err.Error())
				return
			}
			fmt.Printf("%d sats\n", walletBalance.AsSats())
		},
	}

	rootCmd.AddCommand(walletbalanceCmd)

	walletbalanceCmd.Flags().
		StringP("node", "n", "", "The name of the node to get the wallet balance of")
	err := walletbalanceCmd.MarkFlagRequired("node")
	if err != nil {
		log.Fatalf("Problem marking node flag as required: %v", err.Error())
	}
}
