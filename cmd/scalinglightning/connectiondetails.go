package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

func init() {

	var connectionDetailsCmd = &cobra.Command{
		Use:   "connectiondetails",
		Short: "Output the connection details for a node or all nodes",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			processDebugFlag(cmd)
			_ = cmd.Flag("node").Value.String()
			_, err := cmd.Flags().GetBool("all")
			if err != nil {
				fmt.Printf("Problem getting all flag: %v\n", err.Error())
				return
			}

			slnetwork, err := sl.DiscoverStartedNetwork(kubeConfigPath, apiHost, apiPort)
			if err != nil {
				fmt.Printf(
					"Problem with network discovery, is there a network running? Error: %v\n",
					err.Error(),
				)
				return
			}
			connectionDetails, err := slnetwork.GetConnectionDetailsForAllNodes()
			if err != nil {
				fmt.Printf("Problem getting connection details: %v\n", err.Error())
				return
			}
			previousNodeName := ""
			for _, conDetails := range connectionDetails {
				if conDetails.NodeName != previousNodeName {
					if previousNodeName != "" {
						fmt.Println()
					}
					fmt.Println(conDetails.NodeName)
				}
				fmt.Println("  type: ", conDetails.Type)
				fmt.Printf("  host: %v\n", conDetails.Host)
				fmt.Printf("  port: %v\n\n", conDetails.Port)

				previousNodeName = conDetails.NodeName
			}
		},
	}

	rootCmd.AddCommand(connectionDetailsCmd)

	connectionDetailsCmd.Flags().
		StringP("node", "n", "", "The name of the node to get connection details for")

	connectionDetailsCmd.Flags().
		BoolP("all", "a", false, "Get connection details for all nodes")

}
