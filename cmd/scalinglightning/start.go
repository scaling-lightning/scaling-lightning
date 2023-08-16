package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the network",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		processDebugFlag(cmd)
		helmfile := cmd.Flag("helmfile").Value.String()
		fmt.Println("Starting the network")
		slnetwork := sl.NewSLNetwork(helmfile, kubeConfigPath, sl.Regtest)
		err := slnetwork.Start()
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().
		StringP("helmfile", "f", "", "Location of helmfile.yaml (required)")
	startCmd.MarkFlagRequired("helmfile")
}
