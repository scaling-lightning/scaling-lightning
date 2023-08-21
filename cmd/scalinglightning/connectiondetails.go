package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

var connectionDetailsCmd = &cobra.Command{
	Use:   "connectiondetails",
	Short: "Output the connection details for a node or all nodes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		processDebugFlag(cmd)
		nodeName := cmd.Flag("node").Value.String()
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Printf("Problem getting all flag: %v\n", err.Error())
			return
		}

		slnetwork, err := sl.DiscoverStartedNetwork(kubeConfigPath)
		if err != nil {
			fmt.Printf(
				"Problem with network discovery, is there a network running? Error: %v\n",
				err.Error(),
			)
			return
		}
		foundANode := false
		for _, node := range slnetwork.LightningNodes {
			if node.GetName() == nodeName || all {
				connectionDetails := node.GetConnectionDetails()
				// if err != nil {
				// 	fmt.Printf("Problem getting connection details: %v\n", err.Error())
				// 	return
				// }
				foundANode = true
				fmt.Println(node.GetName())
				fmt.Printf("  host: %v\n", connectionDetails.Host)
				fmt.Printf("  port: %v\n\n", connectionDetails.Port)
			}
		}
		if foundANode {
			return
		}

		allNames := []string{}
		for _, node := range slnetwork.LightningNodes {
			allNames = append(allNames, node.GetName())
		}
		fmt.Printf(
			"Can't find node(s), here are the lightning nodes that are running: %v\n",
			allNames,
		)
	},
}

func init() {
	rootCmd.AddCommand(connectionDetailsCmd)

	connectionDetailsCmd.Flags().
		StringP("node", "n", "", "The name of the node to get connection details for")

	connectionDetailsCmd.Flags().
		BoolP("all", "a", false, "Get connection details for all nodes")

}
