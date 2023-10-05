package scalinglightning

import (
	"fmt"
	"log"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

func init() {

	var pubkeyCmd = &cobra.Command{
		Use:   "pubkey",
		Short: "Get the pubkey of a node",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			processDebugFlag(cmd)
			pubkeyNodeName := cmd.Flag("node").Value.String()
			slnetwork, err := sl.DiscoverRunningNetwork(kubeConfigPath, apiHost, apiPort)
			if err != nil {
				fmt.Printf(
					"Problem with network discovery, is there a network running? Error: %v\n",
					err.Error(),
				)
				return
			}
			pubkey, err := slnetwork.GetPubKey(pubkeyNodeName)
			if err != nil {
				fmt.Printf("Problem getting pubkey: %v\n", err.Error())
				return
			}
			fmt.Println(pubkey.AsHexString())
		},
	}

	rootCmd.AddCommand(pubkeyCmd)

	pubkeyCmd.Flags().
		StringP("node", "n", "", "The name of the node to get the pubkey of")
	err := pubkeyCmd.MarkFlagRequired("node")
	if err != nil {
		log.Fatalf("Problem marking node flag as required: %v", err.Error())
	}
}
