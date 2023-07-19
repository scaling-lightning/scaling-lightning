package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

var helmfile string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the network",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting the network")
		err := sl.StartViaHelmfile(helmfile)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().
		StringVarP(&helmfile, "helmfile", "f", "", "Location of helmfile.yaml (required)")
	startCmd.MarkFlagRequired("helmfile")
}
