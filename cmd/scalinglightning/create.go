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

			// Namespace flag should not be used with create, since namespace is read from the helmfile
			if rootCmd.Flag("namespace").Changed {
				fmt.Println("Cannot create. Do not use namespace flag with create. Instead specify the namespace in the helmfile.")
				return
			}

			fmt.Println("Creating and starting the network")
			slnetwork, err := sl.NewSLNetworkWithoutNamespace(helmfile, kubeConfigPath, sl.Regtest)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			err = slnetwork.CreateAndStart()
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
