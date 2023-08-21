package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

var writeAuthFilesCmd = &cobra.Command{
	Use:   "writeauthfiles",
	Short: "Output the auth files for a node",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		processDebugFlag(cmd)
		nodeName := cmd.Flag("node").Value.String()
		authFilesDir := cmd.Flag("dir").Value.String()

		slnetwork, err := sl.DiscoverStartedNetwork(kubeConfigPath)
		if err != nil {
			fmt.Printf(
				"Problem with network discovery, is there a network running? Error: %v\n",
				err.Error(),
			)
			return
		}
		for _, node := range slnetwork.LightningNodes {
			if node.GetName() == nodeName {
				err := node.WriteAuthFilesToDirectory(authFilesDir)
				if err != nil {
					fmt.Printf("Problem writing auth files: %v\n", err.Error())
					return
				}
				fmt.Println("Files written")
				return
			}
		}

		allNames := []string{}
		for _, node := range slnetwork.LightningNodes {
			allNames = append(allNames, node.GetName())
		}
		fmt.Printf(
			"Can't find node with name %v, here are the lightning nodes that are running: %v\n",
			nodeName,
			allNames,
		)
	},
}

func init() {
	rootCmd.AddCommand(writeAuthFilesCmd)

	writeAuthFilesCmd.Flags().
		StringP("node", "n", "", "The name of the node to download the auth files for")
	writeAuthFilesCmd.MarkFlagRequired("node")

	writeAuthFilesCmd.Flags().
		StringP("dir", "o", "", "The directory to write the auth files to")
	writeAuthFilesCmd.MarkFlagRequired("dir")

}