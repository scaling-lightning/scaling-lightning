package scalinglightning

import (
	"fmt"

	"github.com/cockroachdb/errors"
	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

func init() {

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List the nodes in the network",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			processDebugFlag(cmd)
			slnetwork, err := sl.DiscoverRunningNetwork(kubeConfigPath, apiHost, apiPort, namespace)
			if err != nil {
				fmt.Printf(
					"Problem with network discovery, is there a network running? Error: %v\n",
					err.Error(),
				)
				return
			}
			makeRunningMessage := func(nodeName string) (string, error) {
				running, err := slnetwork.IsNodeRunning(nodeName)
				if err != nil {
					fmt.Printf("Problem checking if node is running: %v\n", err.Error())
					return "", errors.Wrap(err, "Checking if node is running")
				}
				runningMessage := ""
				if running {
					runningMessage = ""
				} else {
					runningMessage = " (stopped)"
				}
				return runningMessage, nil
			}
			fmt.Printf("Bitcoin nodes:\n\n")
			for _, node := range slnetwork.BitcoinNodes {
				runningMessage, err := makeRunningMessage(node.GetName())
				if err != nil {
					fmt.Printf("%v\n", err.Error())
					return
				}
				fmt.Printf("	%v%v\n", node.GetName(), runningMessage)
			}
			fmt.Printf("\nLightning nodes:\n\n")
			for _, node := range slnetwork.LightningNodes {
				runningMessage, err := makeRunningMessage(node.GetName())
				if err != nil {
					fmt.Printf("%v\n", err.Error())
					return
				}
				fmt.Printf("	%v%v\n", node.GetName(), runningMessage)
			}
		},
	}

	rootCmd.AddCommand(listCmd)
}
