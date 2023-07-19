package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

var sendFromName string
var sendToName string
var sendAmount uint64

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send on chain funds betwen nodes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := sl.Send(sendFromName, sendToName, sendAmount)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Sent funds")
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringVarP(&sendFromName, "from", "f", "", "Name of node to send from")
	sendCmd.MarkFlagRequired("from")

	sendCmd.Flags().StringVarP(&sendToName, "to", "t", "", "Name of node to send to")
	sendCmd.MarkFlagRequired("to")

	sendCmd.Flags().Uint64VarP(&sendAmount, "amount", "a", 0, "Amount of satoshis to send")
	sendCmd.MarkFlagRequired("amount")

}
