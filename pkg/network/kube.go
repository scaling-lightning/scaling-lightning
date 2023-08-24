package network

import (
	"encoding/json"
	"os/exec"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
)

func kubeCP(kubeconfig string, source string, destination string) error {
	kubectlCmd := exec.Command(
		"kubectl",
		"--kubeconfig",
		kubeconfig,
		"-n",
		mainNamespace,
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

type ingressRouteTCPData struct {
	Spec struct {
		EntryPoints []string `json:"entryPoints"`
	} `json:"spec"`
}

type endpointsData struct {
	Subsets []struct {
		Ports []struct {
			Name     string `json:"name"`
			Port     uint16 `json:"port"`
			Protocol string `json:"protocol"`
		} `json:"ports"`
	} `json:"subsets"`
}

func getEndpointForNode(kubeconfig string, nodeName string) (uint16, error) {
	kubectlCmd := exec.Command(
		"kubectl",
		"--kubeconfig",
		kubeconfig,
		"-n",
		mainNamespace,
		"get",
		"ingressroutetcps.traefik.containo.us",
		nodeName+"-direct-grpc",
		"-o",
		"json",
	)
	kubectlOut, err := kubectlCmd.Output()
	log.Debug().Msgf("kubectl output was: %v", string(kubectlOut))
	if err != nil {
		log.Error().Err(err).Msgf("kubectl output was: %v", string(kubectlOut))
		return 0, errors.Wrap(err, "Running kubectl get endpoints command")
	}
	ingressRouteTCPData := ingressRouteTCPData{}
	err = json.Unmarshal(kubectlOut, &ingressRouteTCPData)
	if err != nil {
		return 0, errors.Wrap(err, "Unmarshalling ingressRouteTCPData")
	}
	if len(ingressRouteTCPData.Spec.EntryPoints) != 1 {
		return 0, errors.New("Expected 1 entrypoint")
	}
	entrypoint := ingressRouteTCPData.Spec.EntryPoints[0]

	kubectlCmd = exec.Command(
		"kubectl",
		"--kubeconfig",
		kubeconfig,
		"-n",
		traefikNamespace,
		"get",
		"endpoints",
		"traefik",
		"-o",
		"json",
	)

	kubectlOut, err = kubectlCmd.Output()
	log.Debug().Msgf("kubectl output was: %v", string(kubectlOut))
	if err != nil {
		log.Error().Err(err).Msgf("kubectl output was: %v", string(kubectlOut))
		return 0, errors.Wrap(err, "Running kubectl get endpoints command")
	}

	endpointsData := endpointsData{}
	err = json.Unmarshal(kubectlOut, &endpointsData)
	if err != nil {
		return 0, errors.Wrap(err, "Unmarshalling endpointsData")
	}
	for _, subset := range endpointsData.Subsets {
		for _, port := range subset.Ports {
			if port.Name == entrypoint {
				return port.Port, nil
			}
		}
	}
	return 0, errors.New("Couldn't find port")
}
