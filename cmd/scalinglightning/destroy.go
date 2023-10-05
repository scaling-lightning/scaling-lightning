package scalinglightning

import (
	"fmt"
	"log"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

func init() {

	var destroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "destroy the network",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			processDebugFlag(cmd)
			destroyHelmfile := cmd.Flag("helmfile").Value.String()
			fmt.Println("Destroying the network")
			slnetwork := sl.NewSLNetwork(destroyHelmfile, kubeConfigPath, sl.Regtest)
			err := slnetwork.Destroy()
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}

	rootCmd.AddCommand(destroyCmd)

	destroyCmd.Flags().
		StringP("helmfile", "f", "", "Location of helmfile.yaml (required)")
	err := destroyCmd.MarkFlagRequired("helmfile")
	if err != nil {
		log.Fatalf("Problem marking helmfile flag as required: %v", err.Error())
	}
}
