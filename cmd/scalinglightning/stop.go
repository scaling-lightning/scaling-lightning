package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

var stopHelmfile string

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the network",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stop the network")
		err := sl.StopViaHelmfile(stopHelmfile)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	stopCmd.Flags().
		StringVarP(&stopHelmfile, "helmfile", "f", "", "Location of helmfile.yaml (required)")
	stopCmd.MarkFlagRequired("helmfile")
}
