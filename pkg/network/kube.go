package network

import (
	"os/exec"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
)

func kubeCP(kubeconfig string, source string, destination string) error {
	kubectlCmd := exec.Command(
		"kubectl",
		"--kubeconfig",
		kubeconfig,
		"cp",
		source,
		destination,
	)
	kubectlOut, err := kubectlCmd.CombinedOutput()
	log.Debug().Msgf("kubectl output was: %v", string(kubectlOut))
	if err != nil {
		log.Error().Err(err).Msgf("kubectl output was: %v", string(kubectlOut))
		return errors.Wrap(err, "Running kubectl cp command")
	}
	return nil
}
