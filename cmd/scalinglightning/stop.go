package scalinglightning

import (
	"fmt"
	"log"

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
			stopHelmfile := cmd.Flag("helmfile").Value.String()
			fmt.Println("Stopping the network")
			slnetwork := sl.NewSLNetwork(stopHelmfile, kubeConfigPath, sl.Regtest)
			err := slnetwork.Stop()
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}

	rootCmd.AddCommand(stopCmd)

	stopCmd.Flags().
		StringP("helmfile", "f", "", "Location of helmfile.yaml (required)")
	err := stopCmd.MarkFlagRequired("helmfile")
	if err != nil {
		log.Fatalf("Problem marking helmfile flag as required: %v", err.Error())
	}
}
