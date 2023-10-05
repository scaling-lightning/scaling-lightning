package scalinglightning

import (
	"fmt"
	"log"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

func init() {
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create and start the network",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			processDebugFlag(cmd)
			helmfile := cmd.Flag("helmfile").Value.String()
			fmt.Println("Creating and starting the network")
			slnetwork := sl.NewSLNetwork(helmfile, kubeConfigPath, sl.Regtest)
			err := slnetwork.CreateAndStart()
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}

	rootCmd.AddCommand(createCmd)

	createCmd.Flags().
		StringP("helmfile", "f", "", "Location of helmfile.yaml (required)")
	err := createCmd.MarkFlagRequired("helmfile")
	if err != nil {
		log.Fatalf("Problem marking helmfile flag as required: %v", err.Error())
	}
}
