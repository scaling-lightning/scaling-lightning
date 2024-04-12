package scalinglightning

import (
	"fmt"
	"log"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

func init() {

	var payInvoiceCmd = &cobra.Command{
		Use:   "payinvoice",
		Short: "pay lightning invoice",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			processDebugFlag(cmd)
			nodeName := cmd.Flag("node").Value.String()
			invoice := cmd.Flag("invoice").Value.String()
			slnetwork, err := sl.DiscoverRunningNetwork(kubeConfigPath, apiHost, apiPort, namespace)
			if err != nil {
				fmt.Printf(
					"Problem with network discovery, is there a network running? Error: %v\n",
					err.Error(),
				)
				return
			}
			preimage, err := slnetwork.PayInvoice(nodeName, invoice)
			if err != nil {
				fmt.Printf("Problem paying the invoice: %v\n", err.Error())
				return
			}
			fmt.Printf("preimage: %v\n", preimage)
		},
	}

	rootCmd.AddCommand(payInvoiceCmd)

	payInvoiceCmd.Flags().
		StringP("node", "n", "", "The name of the node that will pay the invoice")
	err := payInvoiceCmd.MarkFlagRequired("node")
	if err != nil {
		log.Fatalf("Problem marking node flag as required: %v", err.Error())
	}

	payInvoiceCmd.Flags().StringP("invoice", "i", "", "The invoice to pay")
	err = payInvoiceCmd.MarkFlagRequired("invoice")
	if err != nil {
		log.Fatalf("Problem marking invoice flag as required: %v", err.Error())
	}
}
