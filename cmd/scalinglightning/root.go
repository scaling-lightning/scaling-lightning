package scalinglightning

import (
	"os"
	"os/user"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var kubeConfigPath string

var rootCmd = &cobra.Command{
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

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug logging")
}
