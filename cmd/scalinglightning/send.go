package scalinglightning

import (
	"fmt"
	"log"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

func init() {
	var sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Send on chain funds betwen nodes",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			processDebugFlag(cmd)
			sendFromName := cmd.Flag("from").Value.String()
			sendToName := cmd.Flag("to").Value.String()
			sendAmount, err := cmd.Flags().GetUint64("amount")
			if err != nil {
				fmt.Println("Amount must be a valid number")
				return
			}

			slnetwork, err := sl.DiscoverRunningNetwork(kubeConfigPath, apiHost, apiPort)
			if err != nil {
				fmt.Printf(
					"Problem with network discovery, is there a network running? Error: %v\n",
					err.Error(),
				)
				return
			}

			sendRes, err := slnetwork.Send(sendFromName, sendToName, sendAmount)
			if err != nil {
				fmt.Printf("Problem sending funds: %v\n", err.Error())
				return
			}

			fmt.Printf("Sent funds, txid: %v\n", sendRes)
		},
	}

	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringP("from", "f", "", "Name of node to send from")
	err := sendCmd.MarkFlagRequired("from")
	if err != nil {
		log.Fatalf("Problem marking from flag required: %v\n", err.Error())
	}

	sendCmd.Flags().StringP("to", "t", "", "Name of node to send to")
	err = sendCmd.MarkFlagRequired("to")
	if err != nil {
		log.Fatalf("Problem marking to flag required: %v\n", err.Error())
	}

	sendCmd.Flags().Uint64P("amount", "a", 0, "Amount of satoshis to send")
	err = sendCmd.MarkFlagRequired("amount")
	if err != nil {
		log.Fatalf("Problem marking amount flag required: %v\n", err.Error())
	}

}
