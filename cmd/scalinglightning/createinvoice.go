package scalinglightning

import (
	"fmt"
	"log"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

func init() {

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
			slnetwork, err := sl.DiscoverRunningNetwork(kubeConfigPath, apiHost, apiPort, namespace)
			if err != nil {
				fmt.Printf(
					"Problem with network discovery, is there a network running? Error: %v\n",
					err.Error(),
				)
				return
			}
			invoice, err := slnetwork.CreateInvoice(nodeName, amount)
			if err != nil {
				fmt.Printf("Problem generating invoice: %v\n", err.Error())
				return
			}
			fmt.Printf("bolt11: %v\n", invoice)
		},
	}

	rootCmd.AddCommand(createInvoiceCmd)

	createInvoiceCmd.Flags().
		StringP("node", "n", "", "The name of the node to generate the invoice")
	err := createInvoiceCmd.MarkFlagRequired("node")
	if err != nil {
		log.Fatalf("Problem marking node flag as required: %v", err.Error())
	}

	createInvoiceCmd.Flags().Uint64P("amount", "a", 0, "Amount of satoshis to invoice for")
	err = createInvoiceCmd.MarkFlagRequired("amount")
	if err != nil {
		log.Fatalf("Problem marking amount flag as required: %v", err.Error())
	}
}
