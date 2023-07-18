package scalinglightning

import (
	"fmt"

	"github.com/spf13/cobra"
)

// walletbalanceCmd represents the walletbalance command
var walletbalanceCmd = &cobra.Command{
	Use:   "walletbalance",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("walletbalance called")
	},
}

func init() {
	rootCmd.AddCommand(walletbalanceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	walletbalanceCmd.PersistentFlags().String("node", "", "The name of the node to get the wallet balance of")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// walletbalanceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
