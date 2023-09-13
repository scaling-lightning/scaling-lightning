package scalinglightning

import (
	"fmt"

	"github.com/scaling-lightning/scaling-lightning/cmd/build"
	"github.com/spf13/cobra"
)

func init() {

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Get the version of this binary",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			processDebugFlag(cmd)
			fmt.Println(build.ExtendedVersion())
		},
	}

	rootCmd.AddCommand(versionCmd)
}
