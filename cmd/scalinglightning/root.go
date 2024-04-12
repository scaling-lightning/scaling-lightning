package scalinglightning

import (
	"github.com/scaling-lightning/scaling-lightning/pkg/network"
	"os"
	"os/user"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var kubeConfigPath string //nolint:gochecknoglobals
var apiHost string        //nolint:gochecknoglobals
var apiPort uint16        //nolint:gochecknoglobals
var namespace string      //nolint:gochecknoglobals

var rootCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "sl",
	Short: "A CLI for interacting with the scaling-lightning network",
	Long:  ``,
}

func processDebugFlag(cmd *cobra.Command) {
	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid debug flag")
	}
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug logging enabled")
	} else {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not get current user")
	}

	kubeConfigPath = path.Join(currentUser.HomeDir, ".kube", "config")

	rootCmd.PersistentFlags().
		StringVarP(&kubeConfigPath, "kubeconfig", "k", kubeConfigPath, "Location of Kubernetes config file")

	rootCmd.PersistentFlags().
		StringVarP(&apiHost, "host", "H", "", "Host of the scaling-lightning API")
	rootCmd.PersistentFlags().
		Uint16VarP(&apiPort, "port", "p", 0, "Port of the scaling-lightning API")
	rootCmd.PersistentFlags().
		StringVarP(&namespace, "namespace", "N", network.DefaultNamespace, "Kubernetes namespace for the network")

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug logging")
}
