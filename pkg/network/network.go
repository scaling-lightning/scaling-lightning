package network

import (
	"os/exec"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
)

func CheckDependencies() error {
	_, err := exec.LookPath("helmfile")
	if err != nil {
		return errors.Wrap(err, "Looking for helmfile executable on system")
	}
	_, err = exec.LookPath("helm")
	if err != nil {
		return errors.Wrap(err, "Looking for helm executable on system")
	}
	_, err = exec.LookPath("kubectl")
	if err != nil {
		return errors.Wrap(err, "Looking for kubectl executable on system")
	}
	hplCmd := exec.Command("helm", "plugin", "list")
	hplOut, err := hplCmd.Output()
	if err != nil {
		return errors.Wrap(err, "Running helm plugin list command")
	}
	if !strings.Contains(string(hplOut), "diff\t") {
		log.Debug().Err(err).Msgf("helm plugin list output was: %v", string(hplOut))
		return errors.New("Helm plugin diff not installed")
	}

	podNameCmdStr := `get pods -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx --field-selector=status.phase==Running -o jsonpath='{.items[0].metadata.name}'`
	podNameCmd := exec.Command("kubectl", strings.Split(podNameCmdStr, " ")...)
	podNameOut, err := podNameCmd.Output()
	if err != nil {
		log.Debug().Err(err).Msgf("pod name output was: %v", string(podNameOut))
		return errors.Wrap(err, "Getting pod name for ingress-nginx")
	}

	podNameStr := strings.ReplaceAll(string(podNameOut), "'", "")
	nginxIngressCmdStr := `exec -n ingress-nginx -it ` + string(
		podNameStr,
	) + ` -- /nginx-ingress-controller --version`
	nginxIngressCmd := exec.Command("kubectl", strings.Split(nginxIngressCmdStr, " ")...)
	nginxIngressOut, err := nginxIngressCmd.Output()
	if err != nil {
		log.Debug().
			Err(err).
			Msgf("nginx-ingress-controller version output was: %v", string(nginxIngressOut))
		return errors.Wrap(err, "Getting nginx-ingress-controller version")
	}
	if !strings.Contains(strings.ToLower(string(nginxIngressOut)), "nginx ingress controller") {
		log.Debug().
			Err(err).
			Msgf("nginx-ingress-controller version output was: %v", string(nginxIngressOut))
		return errors.New("Ingress nginx not installed")
	}

	return nil
}

func StartViaHelmfile(helmfilePath string) error {
	if err := CheckDependencies(); err != nil {
		return errors.Wrap(err, "Checking dependencies")
	}
	helmfileCmd := exec.Command("helmfile", "apply", "-f", helmfilePath)
	helmfileOut, err := helmfileCmd.Output()
	if err != nil {
		log.Debug().Msgf("helmfile output was: %v", string(helmfileOut))
		return errors.Wrap(err, "Running helmfile apply command")
	}
	log.Debug().Err(err).Msgf("helmfile output was: %v", string(helmfileOut))
	return nil
}

func Stop() {
	log.Info().Msg("Stopping network")
}
