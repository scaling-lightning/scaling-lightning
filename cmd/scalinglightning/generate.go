package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate bitcoin blocks",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		nodeName := cmd.Flag("node").Value.String()
		numOfBlocks, err := cmd.Flags().GetUint32("blocks")
		if err != nil {
			fmt.Println("Not a valid number of blocks")
			return
		}

		slnetwork, err := sl.DiscoverStartedNetwork("")
		if err != nil {
			fmt.Printf(
				"Problem with network discovery, is there a network running? Error: %v\n",
				err.Error(),
			)
			return
		}
		var bitcoinNode sl.BitcoinNode
		for _, node := range slnetwork.BitcoinNodes {
			if node.GetName() == nodeName {
				bitcoinNode = node
				continue
			}
		}
		allNames := []string{}
		for _, node := range slnetwork.BitcoinNodes {
			allNames = append(allNames, node.GetName())
		}
		if bitcoinNode.Name == "" {
			fmt.Printf(
				"Can't find node with name %v, here are the nodes that are running: %v\n",
				nodeName,
				allNames,
			)
		}

		generateRes, err := bitcoinNode.Generate(numOfBlocks)
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

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().
		StringP("node", "n", "", "The name of the node to generate blocks on")
	generateCmd.MarkFlagRequired("node")

	generateCmd.Flags().
		Uint32P("blocks", "b", 50, "How many blocks to generate")

}
