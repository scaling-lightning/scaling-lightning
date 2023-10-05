package scalinglightning

import (
	"fmt"

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
			fmt.Println("Destroying the network")
			slnetwork := sl.NewSLNetwork("", kubeConfigPath, sl.Regtest)
			err := slnetwork.Destroy()
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}

	rootCmd.AddCommand(destroyCmd)
}
