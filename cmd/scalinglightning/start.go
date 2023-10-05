package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

func init() {

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start a stopped network",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			processDebugFlag(cmd)
			node := cmd.Flag("node").Value.String()
			all, err := cmd.Flags().GetBool("all")
			if err != nil {
				fmt.Printf("Problem getting all flag: %v\n", err.Error())
				return
			}

			if !all && node == "" {
				fmt.Println("Must specify a node or use the --all flag")
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
				err := slnetwork.Start()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println("Network started")
				return
			}

			err = slnetwork.StartNode(node)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println("Node started")
		},
	}

	rootCmd.AddCommand(startCmd)

	startCmd.Flags().
		StringP("node", "n", "", "The name of the node to start")

	startCmd.Flags().
		BoolP("all", "a", false, "Start all nodes")

}
