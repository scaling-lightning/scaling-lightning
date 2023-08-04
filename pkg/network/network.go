package network

import (
	"os"
	"os/exec"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"github.com/scaling-lightning/scaling-lightning/pkg/types"

	"gopkg.in/yaml.v3"
)

type SLNetwork struct {
	BitcoinNodes   []BitcoinNode
	LightningNodes []LightningNode
	kubeConfig     string
	helmfile       string
	endpoint       string
}

type Node interface {
	GetNewAddress() (types.Address, error)
	Send(Node, types.Amount) error
	GetName() string
	GetWalletBalance() (types.Amount, error)
}

func (n *SLNetwork) CheckDependencies() error {
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

	// podNameCmdStr := `get pods -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx --field-selector=status.phase==Running -o jsonpath='{.items[0].metadata.name}'`
	// podNameCmd := exec.Command("kubectl", strings.Split(podNameCmdStr, " ")...)
	// podNameOut, err := podNameCmd.Output()
	// if err != nil {
	// 	log.Debug().Err(err).Msgf("pod name output was: %v", string(podNameOut))
	// 	return errors.Wrap(err, "Getting pod name for ingress-nginx")
	// }

	// podNameStr := strings.ReplaceAll(string(podNameOut), "'", "")
	// nginxIngressCmdStr := `exec -n ingress-nginx -it ` + string(
	// 	podNameStr,
	// ) + ` -- /nginx-ingress-controller --version`
	// nginxIngressCmd := exec.Command("kubectl", strings.Split(nginxIngressCmdStr, " ")...)
	// nginxIngressOut, err := nginxIngressCmd.Output()
	// if err != nil {
	// 	log.Debug().
	// 		Err(err).
	// 		Msgf("nginx-ingress-controller version output was: %v", string(nginxIngressOut))
	// 	return errors.Wrap(err, "Getting nginx-ingress-controller version")
	// }
	// if !strings.Contains(strings.ToLower(string(nginxIngressOut)), "nginx ingress controller") {
	// 	log.Debug().
	// 		Err(err).
	// 		Msgf("nginx-ingress-controller version output was: %v", string(nginxIngressOut))
	// 	return errors.New("Ingress nginx not installed")
	// }

	return nil
}

func NewSLNetwork(helmfile string, kubeConfig string) SLNetwork {
	return SLNetwork{
		helmfile:   helmfile,
		kubeConfig: kubeConfig,
	}
}

type helmListOutput struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Revision   string `json:"revision"`
	Updated    string `json:"updated"`
	Status     string `json:"status"`
	Chart      string `json:"chart"`
	AppVersion string `json:"app_version"`
}

func DiscoverStartedNetwork(kubeconfig string) (*SLNetwork, error) {
	helmCmd := exec.Command("helm", "list", "-o", "json")
	helmOut, err := helmCmd.Output()
	if err != nil {
		log.Debug().Err(err).Msgf("helm output was: %v", string(helmOut))
		return nil, errors.Wrap(err, "Running helmfile list command")
	}
	helmListOutput := []helmListOutput{}
	err = yaml.Unmarshal(helmOut, &helmListOutput)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshalling helm list output")
	}
	slnetwork := NewSLNetwork("", kubeconfig)
	for _, release := range helmListOutput {
		if strings.Contains(release.Chart, "bitcoin") {
			bitcoinNode := BitcoinNode{Name: release.Name, SLNetwork: &slnetwork}
			slnetwork.BitcoinNodes = append(slnetwork.BitcoinNodes, bitcoinNode)
		} else if strings.Contains(release.Chart, "lnd") || strings.Contains(release.Chart, "cln") {
			lightningNode := LightningNode{Name: release.Name, SLNetwork: &slnetwork}
			slnetwork.LightningNodes = append(slnetwork.LightningNodes, lightningNode)
		}
	}
	return &slnetwork, nil
}

func (n *SLNetwork) Start() error {
	log.Debug().Msg("Starting network")
	if err := n.CheckDependencies(); err != nil {
		return errors.Wrap(err, "Checking dependencies")
	}
	helmfileCmd := exec.Command("helmfile", "apply", "-f", n.helmfile)
	helmfileOut, err := helmfileCmd.Output()
	if err != nil {
		log.Debug().Err(err).Msgf("helmfile output was: %v", string(helmfileOut))
		return errors.Wrap(err, "Running helmfile apply command")
	}

	lightningNodes, bitcoinNodes, err := parseHelmfileForNodes(n)
	if err != nil {
		return errors.Wrap(err, "Parsing helmfile for nodes")
	}
	n.LightningNodes = lightningNodes
	n.BitcoinNodes = bitcoinNodes

	return nil
}

func (n *SLNetwork) Stop() error {
	log.Debug().Msg("Stopping network")
	helmfileCmd := exec.Command("helmfile", "destroy", "-f", n.helmfile)
	helmfileOut, err := helmfileCmd.Output()
	if err != nil {
		log.Debug().Err(err).Msgf("helmfile output was: %v", string(helmfileOut))
		return errors.Wrap(err, "Running helmfile destroy command")
	}
	return nil
}

func (n *SLNetwork) GetBitcoinNode(name string) (*BitcoinNode, error) {
	for _, node := range n.BitcoinNodes {
		if node.Name == name {
			return &node, nil
		}
	}
	return nil, errors.New("Bitcoin node not found")
}

func (n *SLNetwork) GetLightningNode(name string) (*LightningNode, error) {
	for _, node := range n.LightningNodes {
		if node.Name == name {
			return &node, nil
		}
	}
	return nil, errors.New("Lightning node not found")
}

func (n *SLNetwork) GetAllNodes() []Node {
	nodes := []Node{}
	for i := range n.BitcoinNodes {
		nodes = append(nodes, &n.BitcoinNodes[i])
	}
	for i := range n.LightningNodes {
		nodes = append(nodes, &n.LightningNodes[i])
	}
	return nodes
}

type helmFile struct {
	Releases []struct {
		Name  string `yaml:"name"`
		Chart string `yaml:"chart"`
	}
}

func parseHelmfileForNodes(
	slnetwork *SLNetwork,
) (lightningNodes []LightningNode, bitcoinNodes []BitcoinNode, err error) {
	bytes, err := os.ReadFile(slnetwork.helmfile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "OS Reading helmfile")
	}

	helmFile := &helmFile{}
	err = yaml.Unmarshal(bytes, helmFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Unmarshalling helmfile")
	}

	for _, release := range helmFile.Releases {
		if strings.Contains(release.Chart, "lnd") || strings.Contains(release.Chart, "cln") {
			lightningNodes = append(
				lightningNodes,
				LightningNode{Name: release.Name, SLNetwork: slnetwork},
			)
		}
		if strings.Contains(release.Chart, "bitcoind") {
			bitcoinNodes = append(
				bitcoinNodes,
				BitcoinNode{Name: release.Name, SLNetwork: slnetwork},
			)
		}
	}

	return lightningNodes, bitcoinNodes, err

}
