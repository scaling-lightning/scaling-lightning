package scalinglightning

import (
	"fmt"
	"log"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

func init() {

	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate bitcoin blocks",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			processDebugFlag(cmd)
			nodeName := cmd.Flag("node").Value.String()
			numBlocks, err := cmd.Flags().GetUint32("blocks")
			if err != nil {
				fmt.Println("Not a valid number of blocks")
				return
			}

			slnetwork, err := sl.DiscoverRunningNetwork(kubeConfigPath, apiHost, apiPort, namespace)
			if err != nil {
				fmt.Printf(
					"Problem with network discovery, is there a network running? Error: %v\n",
					err.Error(),
				)
				return
			}

			generateRes, err := slnetwork.Generate(nodeName, numBlocks)
			if err != nil {
				fmt.Printf("Problem sending funds: %v\n", err.Error())
				return
			}

			fmt.Println("Generated blocks:")
			for _, blockHash := range generateRes {
				fmt.Printf("%v", blockHash)
			}

		},
	}

	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().
		StringP("node", "n", "", "The name of the node to generate blocks on")
	err := generateCmd.MarkFlagRequired("node")
	if err != nil {
		log.Fatalf("Problem marking node flag as required: %v", err.Error())
	}

	generateCmd.Flags().
		Uint32P("blocks", "b", 50, "How many blocks to generate")

}
