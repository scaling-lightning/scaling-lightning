package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the nodes in the network",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		processDebugFlag(cmd)
		slnetwork, err := sl.DiscoverStartedNetwork(kubeConfigPath)
		if err != nil {
			fmt.Printf(
				"Problem with network discovery, is there a network running? Error: %v\n",
				err.Error(),
			)
			return
		}
		fmt.Printf("Bitcoin nodes:\n\n")
		for _, node := range slnetwork.BitcoinNodes {
			fmt.Printf("	%v\n", node.GetName())
		}
		fmt.Printf("\nLightning nodes:\n\n")
		for _, node := range slnetwork.LightningNodes {
			fmt.Printf("	%v\n", node.GetName())
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
