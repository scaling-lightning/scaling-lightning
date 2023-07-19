package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

var balanceNodeName string

var walletbalanceCmd = &cobra.Command{
	Use:   "walletbalance",
	Short: "Get the onchain wallet balance of a node",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		balance, err := sl.GetWalletBalanceSats(balanceNodeName)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(balance)
	},
}

func init() {
	rootCmd.AddCommand(walletbalanceCmd)

	walletbalanceCmd.Flags().
		StringVarP(&balanceNodeName, "node", "n", "", "The name of the node to get the wallet balance of")
	walletbalanceCmd.MarkFlagRequired("node")
}
