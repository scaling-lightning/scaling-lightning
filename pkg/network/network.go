package network

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"github.com/scaling-lightning/scaling-lightning/pkg/bitcoinnode"
	"github.com/scaling-lightning/scaling-lightning/pkg/lightningnode"
	stdbitcoinclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/bitcoin"
	stdcommonclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/common"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools/grpc_helpers"
	"github.com/scaling-lightning/scaling-lightning/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

const mainNamespace = "sl"
const traefikNamespace = "sl-traefik"
const clientgRPCPort = 28100

type NetworkType int

const (
	Mainnet NetworkType = iota
	Testnet
	Signet
	Regtest
	Simnet
)

func (n NetworkType) String() string {
	switch n {
	case Mainnet:
		return "mainnet"
	case Testnet:
		return "testnet"
	case Signet:
		return "signet"
	case Regtest:
		return "regtest"
	case Simnet:
		return "simnet"
	default:
		return "unknown"
	}
}

//go:generate mockery --name SLNetworkInterface --exported
type SLNetworkInterface interface {
	Start() error
	Stop() error
	GetBitcoinNode(name string) (*bitcoinnode.BitcoinNode, error)
	GetLightningNode(name string) (*lightningnode.LightningNode, error)
	GetNode(name string) (Node, error)
	GetAllNodes() []Node
}

type SLNetwork struct {
	BitcoinNodes   []bitcoinnode.BitcoinNode
	LightningNodes []lightningnode.LightningNode
	kubeConfig     string
	helmfile       string
	ApiHost        string
	ApiPort        uint16
	Network        NetworkType
}

type ConnectionDetails struct {
	Name string
	Host string
	Port uint16
}

type ConnectionFiles struct {
	LND *LNDConnectionFiles
	CLN *CLNConnectionFiles
}

type LNDConnectionFiles struct {
	TLSCert  []byte
	Macaroon []byte
}

type CLNConnectionFiles struct {
	ClientCert []byte
	ClientKey  []byte
	CACert     []byte
}

type Node interface {
	GetNewAddress() (string, error)
	Send(Node, types.Amount) (string, error)
	GetName() string
	GetWalletBalance() (types.Amount, error)
	GetConnectionDetails() ([]ConnectionDetails, error)
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
	hplOut, err := hplCmd.CombinedOutput()
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

func NewSLNetwork(helmfile string, kubeConfig string, Network NetworkType) SLNetwork {
	return SLNetwork{
		helmfile:   helmfile,
		kubeConfig: kubeConfig,
		Network:    Network,
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

func DiscoverStartedNetwork(
	kubeconfig string,
	overrideAPIHost string,
	overrideAPIPort uint16,
) (*SLNetwork, error) {
	// TODO: santise inputs here
	helmCmd := exec.Command( //nolint:gosec
		"helm",
		"--kubeconfig",
		kubeconfig,
		"list",
		"-n",
		mainNamespace,
		"-o",
		"json",
	)
	helmOut, err := helmCmd.Output()
	log.Debug().Err(err).Msgf("helm output was: %v", string(helmOut))
	if err != nil {
		log.Error().Err(err).Msgf("helm output was: %v", string(helmOut))
		return nil, errors.Wrap(err, "Running helmfile list command")
	}
	helmListOutput := []helmListOutput{}
	err = yaml.Unmarshal(helmOut, &helmListOutput)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshalling helm list output")
	}
	slnetwork := NewSLNetwork("", kubeconfig, Regtest)
	for _, release := range helmListOutput {
		switch {
		case strings.Contains(release.Chart, "bitcoin"):
			bitcoinNode := bitcoinnode.BitcoinNode{Name: release.Name}
			slnetwork.BitcoinNodes = append(slnetwork.BitcoinNodes, bitcoinNode)
		case strings.Contains(release.Chart, "lnd"):
			lightningNode := lightningnode.LightningNode{
				Name: release.Name,
				Impl: lightningnode.LND,
			}
			slnetwork.LightningNodes = append(slnetwork.LightningNodes, lightningNode)
		case strings.Contains(release.Chart, "cln"):
			lightningNode := lightningnode.LightningNode{
				Name: release.Name,
				Impl: lightningnode.CLN,
			}
			slnetwork.LightningNodes = append(slnetwork.LightningNodes, lightningNode)
		}
	}

	if overrideAPIHost == "" {
		err = slnetwork.discoverConnectionDetails()
		if err != nil {
			return nil, errors.Wrap(err, "Discovering connection details")
		}
	} else {
		slnetwork.ApiHost = overrideAPIHost
		slnetwork.ApiPort = clientgRPCPort
	}

	if overrideAPIPort != 0 {
		slnetwork.ApiPort = overrideAPIPort
	}

	return &slnetwork, nil
}

type k8sService struct {
	Status struct {
		LoadBalancer struct {
			Ingress []struct {
				Hostname *string `json:"hostname"`
				IP       *string `json:"ip"`
			} `json:"ingress"`
		} `json:"loadBalancer"`
	} `json:"status"`
}

func GetLoadbalancerHostname(
	serviceName string,
	namespace string,
	kubeconfig string,
) (string, error) {
	// TODO: santise inputs here
	kubectlCmd := exec.Command( //nolint:gosec
		"kubectl",
		"--kubeconfig",
		kubeconfig,
		"get",
		"service",
		serviceName,
		"-n",
		namespace,
		"-o",
		"json",
	)
	kubectlOut, err := kubectlCmd.Output()
	log.Debug().Msgf("kubectl output was: %v", string(kubectlOut))
	if err != nil {
		log.Error().Err(err).Msgf("kubectl output was: %v", string(kubectlOut))
		return "", errors.Wrap(err, "Running kubectl get service command")
	}
	k8sService := k8sService{}
	err = json.Unmarshal(kubectlOut, &k8sService)
	if err != nil {
		return "", errors.Wrap(err, "Unmarshalling kubectl get service output")
	}
	if len(k8sService.Status.LoadBalancer.Ingress) == 0 {
		return "", errors.New("No loadbalancer ingress found")
	}
	var host string
	if k8sService.Status.LoadBalancer.Ingress[0].IP != nil {
		host = *k8sService.Status.LoadBalancer.Ingress[0].IP
	} else if k8sService.Status.LoadBalancer.Ingress[0].Hostname != nil {
		host = *k8sService.Status.LoadBalancer.Ingress[0].Hostname
	} else {
		return "", errors.New("No loadbalancer ingress found")
	}
	return host, nil
}

func (n *SLNetwork) Start() error {
	log.Debug().Msg("Starting network")
	if err := n.CheckDependencies(); err != nil {
		return errors.Wrap(err, "Checking dependencies")
	}
	// TODO: santise helmfile flag
	helmfileCmd := exec.Command( //nolint:gosec
		"helmfile",
		"apply",
		"-f",
		n.helmfile,
	)
	helmfileCmd.Env = append(os.Environ(), "KUBECONFIG="+n.kubeConfig)
	helmfileOut, err := helmfileCmd.CombinedOutput()
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

	log.Info().Msg("Discovering ingress connection details, may take a few minutes")
	err = tools.Retry(func() error {
		log.Info().Msg("waiting...")
		return n.discoverConnectionDetails()
	}, time.Second*30, time.Minute*5)
	if err != nil {
		return errors.Wrap(err, "Discovering connection details timeout")
	}

	err = n.discoverConnectionDetails()
	if err != nil {
		return errors.Wrap(err, "Discovering connection details")
	}

	return nil
}

func (n *SLNetwork) discoverConnectionDetails() error {
	loadbalancer, err := GetLoadbalancerHostname("traefik", traefikNamespace, n.kubeConfig)
	if err != nil {
		return errors.Wrap(err, "Getting loadbalancer hostname")
	}
	n.ApiHost = loadbalancer
	n.ApiPort = clientgRPCPort
	return nil
}

func (n *SLNetwork) Stop() error {
	log.Debug().Msg("Stopping network")

	// TODO: santise helmfile flag
	helmfileCmd := exec.Command( //nolint:gosec
		"helmfile",
		"destroy",
		"-f",
		n.helmfile,
	)
	helmfileCmd.Env = append(os.Environ(), "KUBECONFIG="+n.kubeConfig)
	helmfileOut, err := helmfileCmd.CombinedOutput()
	if err != nil {
		log.Debug().Err(err).Msgf("helmfile output was: %v", string(helmfileOut))
		return errors.Wrap(err, "Running helmfile destroy command")
	}
	return nil
}

func (n *SLNetwork) GetBitcoinNode(name string) (*bitcoinnode.BitcoinNode, error) {
	for _, node := range n.BitcoinNodes {
		if node.Name == name {
			return &node, nil
		}
	}
	return nil, errors.New("Bitcoin node not found")
}

func (n *SLNetwork) GetLightningNode(name string) (*lightningnode.LightningNode, error) {
	for _, node := range n.LightningNodes {
		if node.Name == name {
			return &node, nil
		}
	}
	return nil, errors.New("Lightning node not found")
}

// func (n *SLNetwork) GetNode(name string) (Node, error) {
// 	var node Node
// 	node, err := n.GetLightningNode(name)
// 	if err != nil {
// 		node, err = n.GetBitcoinNode(name)
// 		if err != nil {
// 			return nil, errors.Wrapf(err, "Looking up node %v", name)
// 		}
// 	}
// 	if node.GetName() != name {
// 		return nil, errors.New("Node not found")
// 	}

// 	return node, nil
// }

// func (n *SLNetwork) GetAllNodes() []Node {
// 	nodes := []Node{}
// 	for i := range n.BitcoinNodes {
// 		nodes = append(nodes, &n.BitcoinNodes[i])
// 	}
// 	for i := range n.LightningNodes {
// 		nodes = append(nodes, &n.LightningNodes[i])
// 	}
// 	return nodes
// }

type helmFile struct {
	Releases []struct {
		Name  string `yaml:"name"`
		Chart string `yaml:"chart"`
	}
}

func parseHelmfileForNodes(
	slnetwork *SLNetwork,
) (lightningNodes []lightningnode.LightningNode, bitcoinNodes []bitcoinnode.BitcoinNode, err error) {
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
		if strings.Contains(release.Chart, "lnd") {
			lightningNodes = append(
				lightningNodes,
				lightningnode.LightningNode{
					Name: release.Name,
					Impl: lightningnode.LND,
				},
			)
		}
		if strings.Contains(release.Chart, "cln") {
			lightningNodes = append(
				lightningNodes,
				lightningnode.LightningNode{
					Name: release.Name,
					Impl: lightningnode.CLN,
				},
			)
		}
		if strings.Contains(release.Chart, "bitcoind") {
			bitcoinNodes = append(
				bitcoinNodes,
				bitcoinnode.BitcoinNode{Name: release.Name},
			)
		}
	}

	return lightningNodes, bitcoinNodes, nil
}

func (n *SLNetwork) GetWalletBalance(nodeName string) (types.Amount, error) {
	conn, err := connectToGRPCServer(n.ApiHost, n.ApiPort, nodeName)
	if err != nil {
		return types.NewAmountSats(0), errors.Wrapf(err, "Connecting to gRPC for %v's client", nodeName)
	}
	defer conn.Close()
	client := stdcommonclient.NewCommonClient(conn)

	for _, node := range n.BitcoinNodes {
		if node.Name != nodeName {
			continue
		}
		walletBalance, err := node.GetWalletBalance(client)
		if err != nil {
			return types.NewAmountSats(0), errors.Wrapf(err, "Getting wallet balance for %v", node.Name)
		}
		return walletBalance, nil
	}
	for _, node := range n.LightningNodes {
		if node.Name != nodeName {
			continue
		}
		walletBalance, err := node.GetWalletBalance(client)
		if err != nil {
			return types.NewAmountSats(0), errors.Wrapf(err, "Getting wallet balance for %v", node.Name)
		}
		return walletBalance, nil
	}
	return types.NewAmountSats(0), errors.New("Node not found")
}

func (n *SLNetwork) Generate(nodeName string) (hashes []string, err error) {
	conn, err := connectToGRPCServer(n.ApiHost, n.ApiPort, nodeName)
	if err != nil {
		return []string{}, errors.Wrapf(err, "Connecting to gRPC for %v's client", nodeName)
	}
	defer conn.Close()

	client := stdbitcoinclient.NewBitcoinClient(conn)
	commonClient := stdcommonclient.NewCommonClient(conn)

	for _, node := range n.BitcoinNodes {
		if node.Name != nodeName {
			continue
		}
		hashes, err := node.Generate(client, commonClient, 50)
		if err != nil {
			log.Error().Err(err).Msgf("Generating blocks for %v", node.Name)
		}
		return hashes, nil
	}
	return []string{}, errors.New("Node not found")
}

func (n *SLNetwork) GetNewAddress(nodeName string) (string, error) {
	conn, err := connectToGRPCServer(n.ApiHost, n.ApiPort, nodeName)
	if err != nil {
		return "", errors.Wrapf(err, "Connecting to gRPC for %v's client", nodeName)
	}
	defer conn.Close()
	client := stdcommonclient.NewCommonClient(conn)

	for _, node := range n.BitcoinNodes {
		if node.Name != nodeName {
			continue
		}
		newAddress, err := node.GetNewAddress(client)
		if err != nil {
			return "", errors.Wrapf(err, "Getting new address for %v", node.Name)
		}
		return newAddress, nil
	}
	for _, node := range n.LightningNodes {
		if node.Name != nodeName {
			continue
		}
		newAddress, err := node.GetNewAddress(client)
		if err != nil {
			return "", errors.Wrapf(err, "Getting new address for %v", node.Name)
		}
		return newAddress, nil
	}
	return "", errors.New("Node not found")
}

func (n *SLNetwork) GetConnectionDetails(nodeName string) ([]ConnectionDetails, error) {

	for _, node := range n.BitcoinNodes {
		if node.Name != nodeName {
			continue
		}
		connectionPorts, err := node.GetConnectionPorts(n.kubeConfig)
		if err != nil {
			return nil, errors.Wrapf(err, "Getting connection ports for %v", nodeName)
		}
		connectionDetails := []ConnectionDetails{}
		for _, connectionPort := range connectionPorts {
			connectionDetails = append(
				connectionDetails,
				ConnectionDetails{
					Name: connectionPort.Name,
					Host: n.ApiHost,
					Port: connectionPort.Port,
				},
			)
		}
		return connectionDetails, nil
	}

	for _, node := range n.LightningNodes {
		if node.Name != nodeName {
			continue
		}
		port, err := node.GetConnectionPort(n.kubeConfig)
		if err != nil {
			return nil, errors.Wrapf(err, "Getting grpc endpoint for %v", nodeName)
		}
		return []ConnectionDetails{
			{Name: "grpc", Host: n.ApiHost, Port: port},
		}, nil
	}

	return nil, errors.New("Node not found")
}

func (n *SLNetwork) Send(fromNodeName string, toNodeName string, amount types.Amount) (string, error) {
	log.Debug().Msgf("Sending %v from %v to %v", amount, fromNodeName, toNodeName)

	address, err := n.GetNewAddress(toNodeName)
	if err != nil {
		return "", errors.Wrapf(err, "Getting new address for %v", toNodeName)
	}

	conn, err := connectToGRPCServer(n.ApiHost, n.ApiPort, fromNodeName)
	if err != nil {
		return "", errors.Wrapf(err, "Connecting to gRPC for %v's client", fromNodeName)
	}
	defer conn.Close()

	client := stdcommonclient.NewCommonClient(conn)
	txid := ""
	for _, node := range n.BitcoinNodes {
		if node.Name != fromNodeName {
			continue
		}
		txid, err = node.SendToAddress(client, address, amount)
		if err != nil {
			return "", errors.Wrapf(err, "Sending to addres")
		}
	}

	if txid == "" {
		for _, node := range n.LightningNodes {
			if node.Name != fromNodeName {
				continue
			}
			txid, err = node.SendToAddress(client, address, amount)
			if err != nil {
				return "", errors.Wrapf(err, "Sending to addres")
			}
		}
	}

	if txid != "" {
		// TODO: Look this up
		_, err = n.Generate("bitcoind")
		if err != nil {
			return "", errors.Wrapf(err, "Generating blocks for %v", "bitcoind")
		}
	}

	return "", errors.New("Node not found")
}

func connectToGRPCServer(host string, port uint16, nodeName string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_helpers.ClientInterceptor(nodeName)),
	}
	if host == "" {
		host = "localhost"
	}
	if port == 0 {
		port = clientgRPCPort
	}
	conn, err := grpc.Dial(fmt.Sprintf("%v:%d", host, port), opts...)
	if err != nil {
		return nil, errors.Wrap(err, "Connecting to gRPC")
	}
	return conn, nil
}
