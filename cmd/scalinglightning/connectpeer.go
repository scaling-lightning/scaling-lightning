package scalinglightning

import (
	"fmt"
	"log"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

func init() {

	var connectpeerCmd = &cobra.Command{
		Use:   "connectpeer",
		Short: "Connect peers",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			processDebugFlag(cmd)
			connectpeerFromName := cmd.Flag("from").Value.String()
			connectpeerToName := cmd.Flag("to").Value.String()
			slnetwork, err := sl.DiscoverStartedNetwork(kubeConfigPath, apiHost, apiPort)
			if err != nil {
				fmt.Printf(
					"Problem with network discovery, is there a network running? Error: %v\n",
					err.Error(),
				)
				return
			}
			err = slnetwork.ConnectPeer(connectpeerFromName, connectpeerToName)
			if err != nil {
				fmt.Printf("Problem connecting peer: %v\n", err.Error())
				return
			}
			fmt.Println("Connect peer command received")
		},
	}

	rootCmd.AddCommand(connectpeerCmd)

	connectpeerCmd.Flags().
		StringP("from", "f", "", "Name of the node to connect from")
	err := connectpeerCmd.MarkFlagRequired("from")
	if err != nil {
		log.Fatalf("Problem marking from flag as required: %v", err.Error())
	}

	connectpeerCmd.Flags().
		StringP("to", "t", "", "Name of the node to connect from")
	err = connectpeerCmd.MarkFlagRequired("from")
	if err != nil {
		log.Fatalf("Problem marking to flag as required: %v", err.Error())
	}
}
