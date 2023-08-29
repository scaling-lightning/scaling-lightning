package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

var createInvoiceCmd = &cobra.Command{
	Use:   "createinvoice",
	Short: "Create lightning invoice",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		processDebugFlag(cmd)
		nodeName := cmd.Flag("node").Value.String()
		amount, err := cmd.Flags().GetUint64("amount")
		if err != nil {
			fmt.Println("Amount must be a valid number")
			return
		}
		slnetwork, err := sl.DiscoverStartedNetwork(kubeConfigPath)
		if err != nil {
			fmt.Printf(
				"Problem with network discovery, is there a network running? Error: %v\n",
				err.Error(),
			)
			return
		}
		allNodes := slnetwork.LightningNodes
		for _, node := range allNodes {
			if node.GetName() == nodeName {
				invoice, err := node.CreateInvoice(amount)
				if err != nil {
					fmt.Printf("Problem generating invoice: %v\n", err.Error())
					return
				}
				fmt.Printf("bolt11: %v\n", invoice)
				return
			}
		}
		allNames := []string{}
		for _, node := range allNodes {
			allNames = append(allNames, node.GetName())
		}
		fmt.Printf(
			"Can't find node with name %v, here are the nodes that are running: %v\n",
			nodeName,
			allNames,
		)
	},
}

func init() {
	rootCmd.AddCommand(createInvoiceCmd)

	createInvoiceCmd.Flags().
		StringP("node", "n", "", "The name of the node to generate the invoice")
	createInvoiceCmd.MarkFlagRequired("node")

	createInvoiceCmd.Flags().Uint64P("amount", "a", 0, "Amount of satoshis to invoice for")
	createInvoiceCmd.MarkFlagRequired("amount")
}
