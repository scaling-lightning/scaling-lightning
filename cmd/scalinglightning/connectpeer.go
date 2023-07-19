package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

var connectpeerFromName string
var connectpeerToName string

var connectpeerCmd = &cobra.Command{
	Use:   "connectpeer",
	Short: "Connect peers",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := sl.ConnectPeer(connectpeerFromName, connectpeerToName)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Connect peer command received")
	},
}

func init() {
	rootCmd.AddCommand(connectpeerCmd)

	connectpeerCmd.Flags().
		StringVarP(&connectpeerFromName, "from", "f", "", "Name of the node to connect from")
	connectpeerCmd.MarkFlagRequired("from")

	connectpeerCmd.Flags().
		StringVarP(&connectpeerToName, "to", "t", "", "Name of the node to connect from")
	connectpeerCmd.MarkFlagRequired("from")
}
