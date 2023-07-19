package scalinglightning

import (
	"fmt"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

var pubkeyNodeName string

var pubkeyCmd = &cobra.Command{
	Use:   "pubkey",
	Short: "Get the pubkey of a node",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		pubkey, err := sl.GetPubKey(pubkeyNodeName)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(pubkey)
	},
}

func init() {
	rootCmd.AddCommand(pubkeyCmd)

	pubkeyCmd.Flags().
		StringVarP(&pubkeyNodeName, "node", "n", "", "The name of the node to get the wallet balance of")
	pubkeyCmd.MarkFlagRequired("node")
}
