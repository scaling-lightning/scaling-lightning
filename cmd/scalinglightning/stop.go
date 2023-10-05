package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

func init() {

	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop the network",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			processDebugFlag(cmd)
			node := cmd.Flag("node").Value.String()
			all, err := cmd.Flags().GetBool("all")
			if err != nil {
				fmt.Printf("Problem getting all flag: %v\n", err.Error())
				return
			}

			slnetwork, err := sl.DiscoverRunningNetwork(kubeConfigPath, apiHost, apiPort)
			if err != nil {
				fmt.Printf(
					"Problem with network discovery, is there a network running? Error: %v\n",
					err.Error(),
				)
				return
			}
			if all {
				fmt.Println("Stopping the entire network. Start again with the start command. Volume data will be preserved.")
				err := slnetwork.Stop()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println("Network stopped")
				return
			}

			fmt.Println("Stopping node. Start again with the start command. Volume data will be preserved.")
			err = slnetwork.StopNode(node)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println("Node stopped")
		},
	}

	rootCmd.AddCommand(stopCmd)

	stopCmd.Flags().
		StringP("node", "n", "", "The name of the node to stop")

	stopCmd.Flags().
		BoolP("all", "a", false, "Stop all nodes")

}
