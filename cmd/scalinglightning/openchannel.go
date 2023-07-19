package scalinglightning

import (
	"fmt"

	"github.com/rs/zerolog/log"
	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

var openchannelFromName string
var openchannelToName string
var openchannelLocalAmt uint64

var openchannelCmd = &cobra.Command{
	Use:   "openchannel",
	Short: "Open a channel between two nodes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := sl.OpenChannel(openchannelFromName, openchannelToName, openchannelLocalAmt)
		if err != nil {
			log.Error().Err(err).Msg("Error opening channel")
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Open channel command received")
	},
}

func init() {
	rootCmd.AddCommand(openchannelCmd)

	openchannelCmd.Flags().
		StringVarP(&openchannelFromName, "from", "f", "", "Name of node to open channel from")
	openchannelCmd.MarkFlagRequired("from")

	openchannelCmd.Flags().
		StringVarP(&openchannelToName, "to", "t", "", "Name of node to open channel to")
	openchannelCmd.MarkFlagRequired("to")

	openchannelCmd.Flags().
		Uint64VarP(&openchannelLocalAmt, "amount", "a", 0, "Amount of satoshis to put into channel from the opening side")
	openchannelCmd.MarkFlagRequired("amount")

}
