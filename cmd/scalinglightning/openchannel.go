package scalinglightning

import (
	"fmt"
	"log"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

func init() {

	var openchannelCmd = &cobra.Command{
		Use:   "openchannel",
		Short: "Open a channel between two nodes",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			processDebugFlag(cmd)
			openchannelFromName := cmd.Flag("from").Value.String()
			openchannelToName := cmd.Flag("to").Value.String()
			openchannelLocalAmt, err := cmd.Flags().GetUint64("amount")
			if err != nil {
				fmt.Println("Amount must be a valid number")
				return
			}

			slnetwork, err := sl.DiscoverStartedNetwork(kubeConfigPath, apiHost, apiPort)
			if err != nil {
				fmt.Printf(
					"Problem with network discovery, is there a network running? Error: %v\n",
					err.Error(),
				)
				return
			}
			chanPoint, err := slnetwork.OpenChannel(
				openchannelFromName,
				openchannelToName,
				openchannelLocalAmt,
			)
			if err != nil {
				fmt.Printf("Problem opening channel: %v\n", err.Error())
				return
			}

			fmt.Printf(
				"Open channel command received.\nTxid: %v\nOutputIndex: %d\n",
				chanPoint.FundingTx.IdAsHexString(),
				chanPoint.OutputIndex,
			)
		},
	}

	rootCmd.AddCommand(openchannelCmd)

	openchannelCmd.Flags().
		StringP("from", "f", "", "Name of node to open channel from")
	err := openchannelCmd.MarkFlagRequired("from")
	if err != nil {
		log.Fatalf("Problem marking from flag as required: %v\n", err.Error())
	}

	openchannelCmd.Flags().
		StringP("to", "t", "", "Name of node to open channel to")
	err = openchannelCmd.MarkFlagRequired("to")
	if err != nil {
		log.Fatalf("Problem marking to flag as required: %v\n", err.Error())
	}

	openchannelCmd.Flags().
		Uint64P("amount", "a", 0, "Amount of satoshis to put into channel from the opening side")
	err = openchannelCmd.MarkFlagRequired("amount")
	if err != nil {
		log.Fatalf("Problem marking amount flag as required: %v\n", err.Error())
	}

}
