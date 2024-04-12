package network

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"github.com/scaling-lightning/scaling-lightning/pkg/bitcoinnode"
	"github.com/scaling-lightning/scaling-lightning/pkg/kube"
	"github.com/scaling-lightning/scaling-lightning/pkg/lightningnode"
	stdbitcoinclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/bitcoin"
	stdcommonclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/common"
	stdlightningclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/lightning"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools/grpc_helpers"
	"github.com/scaling-lightning/scaling-lightning/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

const DefaultNamespace = "sl"
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

type LightningNodeInterface interface {
	GetName() string
	SendToAddress(
		client stdcommonclient.CommonClient,
		address string,
		amount types.Amount,
	) (string, error)
	GetNewAddress(client stdcommonclient.CommonClient) (string, error)
	GetPubKey(client stdlightningclient.LightningClient) (types.PubKey, error)
	GetWalletBalance(client stdcommonclient.CommonClient) (types.Amount, error)
	ConnectPeer(
		client stdlightningclient.LightningClient,
		pubkey types.PubKey,
		nodeName string,
	) error
	OpenChannel(
		client stdlightningclient.LightningClient,
		pubkey types.PubKey,
		localAmt types.Amount,
	) (types.ChannelPoint, error)
	GetConnectionFiles(network string, kubeConfig string) (*lightningnode.ConnectionFiles, error)
	WriteAuthFilesToDirectory(network string, kubeConfig string, dir string) error
	GetConnectionPort(kubeConfig string) (uint16, error)
	CreateInvoice(client stdlightningclient.LightningClient, amountSats uint64) (string, error)
	PayInvoice(client stdlightningclient.LightningClient, invoice string) (string, error)
	ChannelBalance(client stdlightningclient.LightningClient) (types.Amount, error)
}

type BitcoinNodeInterface interface {
	GetName() string
	Generate(
		client stdbitcoinclient.BitcoinClient,
		commonClient stdcommonclient.CommonClient,
		numBlocks uint32,
	) (hashes []string, err error)
	GetWalletBalance(client stdcommonclient.CommonClient) (types.Amount, error)
	SendToAddress(
		client stdcommonclient.CommonClient,
		address string,
		amount types.Amount,
	) (TxId string, err error)
	GetNewAddress(client stdcommonclient.CommonClient) (string, error)
	GetConnectionPorts(kubeConfig string) ([]bitcoinnode.ConnectionPorts, error)
}

type NodeInterface interface {
	GetName() string
	GetWalletBalance(client stdcommonclient.CommonClient) (types.Amount, error)
	SendToAddress(
		client stdcommonclient.CommonClient,
		address string,
		amount types.Amount,
	) (TxId string, err error)
	GetNewAddress(client stdcommonclient.CommonClient) (string, error)
}

type SLNetwork struct {
	BitcoinNodes   []BitcoinNodeInterface
	LightningNodes []LightningNodeInterface
	kubeConfig     string
	helmfile       string
	ApiHost        string
	ApiPort        uint16
	Network        NetworkType
	RetryCommands  bool
	RetryDelay     time.Duration
	RetryTimeout   time.Duration
	Namespace      string
}

type ConnectionDetails struct {
	NodeName string
	Type     string
	Host     string
	Port     uint16
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

// NewSLNetworkWithoutNamespace should be only used when creating a new network from helmfile
func NewSLNetworkWithoutNamespace(helmfile string, kubeConfig string, Network NetworkType) (SLNetwork, error) {
	// allow empty namespace in creation, since it will be read from the helmfile
	return newSLNetworkInternal(helmfile, kubeConfig, Network, "")
}

// NewSLNetwork creates a new SLNetwork object that is used to interact with an existing network
func NewSLNetwork(helmfile string, kubeConfig string, Network NetworkType, namespace string) (SLNetwork, error) {
	if namespace == "" {
		return SLNetwork{}, errors.New("Must provide a namespace for network")
	}
	if namespace == traefikNamespace {
		return SLNetwork{}, errors.New("Namespace cannot be same as the traefik namespace")
	}

	return newSLNetworkInternal(helmfile, kubeConfig, Network, namespace)
}

func newSLNetworkInternal(helmfile string, kubeConfig string, Network NetworkType, namespace string) (SLNetwork, error) {
	return SLNetwork{
		helmfile:      helmfile,
		kubeConfig:    kubeConfig,
		Network:       Network,
		RetryCommands: true,
		RetryDelay:    time.Second * 10,
		RetryTimeout:  time.Minute * 3,
		Namespace:     namespace,
	}, nil
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

func DiscoverRunningNetwork(
	kubeconfig string,
	overrideAPIHost string,
	overrideAPIPort uint16,
	namespace string,
) (*SLNetwork, error) {
	// TODO: santise inputs here
	if namespace == "" {
		return nil, errors.New("Namespace cannot be empty")
	}

	helmCmd := exec.Command( //nolint:gosec
		"helm",
		"--kubeconfig",
		kubeconfig,
		"list",
		"-n",
		namespace,
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
	slnetwork, err := NewSLNetwork("", kubeconfig, Regtest, namespace)
	if err != nil {
		return nil, errors.Wrap(err, "Making SLNetwork")
	}

	for _, release := range helmListOutput {
		switch {
		case strings.Contains(release.Chart, "bitcoin"):
			bitcoinNode := bitcoinnode.BitcoinNode{Name: release.Name, Namespace: namespace}
			slnetwork.BitcoinNodes = append(slnetwork.BitcoinNodes, &bitcoinNode)
		case strings.Contains(release.Chart, "lnd"):
			lightningNode := lightningnode.LightningNode{
				Name:      release.Name,
				Impl:      lightningnode.LND,
				Namespace: namespace,
			}
			slnetwork.LightningNodes = append(slnetwork.LightningNodes, &lightningNode)
		case strings.Contains(release.Chart, "cln"):
			lightningNode := lightningnode.LightningNode{
				Name:      release.Name,
				Impl:      lightningnode.CLN,
				Namespace: namespace,
			}
			slnetwork.LightningNodes = append(slnetwork.LightningNodes, &lightningNode)
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

func (n *SLNetwork) CreateAndStart() error {
	log.Debug().Msg("Creating network")
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

	lightningNodes, bitcoinNodes, ns, err := parseHelmfileForNodes(n)
	if err != nil {
		return errors.Wrap(err, "Parsing helmfile for nodes")
	}
	n.LightningNodes = lightningNodes
	n.BitcoinNodes = bitcoinNodes

	if n.Namespace != "" && n.Namespace != ns {
		log.Warn().Msgf("Namespace was set to %v, but helmfile specifies namespace %v. "+
			"Replacing the namespace with the one in helmfile.", n.Namespace, ns)
	}

	n.Namespace = ns

	log.Info().Msg("Discovering ingress connection details, may take a few minutes")

	err = tools.Retry(func(cancel context.CancelFunc) error {
		log.Info().Msg("waiting...")
		return n.discoverConnectionDetails()
	}, n.RetryDelay, n.RetryTimeout)
	if err != nil {
		return errors.Wrap(err, "Discovering connection details timeout")
	}

	err = n.discoverConnectionDetails()
	if err != nil {
		return errors.Wrap(err, "Discovering connection details")
	}

	fmt.Printf("Network started in namespace '%v' with %v lightning node(s) "+
		"and %v bitcoin node(s).\n", n.Namespace, len(n.LightningNodes), len(n.BitcoinNodes))

	return nil
}

func (n *SLNetwork) Start() error {
	log.Debug().Msgf("Starting network in namespace '%v'", n.Namespace)
	for _, node := range n.GetAllNodes() {
		err := n.StartNode(node.GetName())
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *SLNetwork) StartNode(nodeName string) error {
	log.Debug().Msg("Starting node")
	err := kube.Scale(n.kubeConfig, nodeName, "statefulset", 1, n.Namespace)
	if err != nil {
		return errors.Wrapf(err, "Scaling node %v to 1 replicas", nodeName)
	}
	return nil
}

func (n *SLNetwork) Stop() error {
	log.Debug().Msg("Stopping network")
	for _, node := range n.GetAllNodes() {
		err := n.StopNode(node.GetName())
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *SLNetwork) StopNode(nodeName string) error {
	log.Debug().Msg("Stopping node")
	err := kube.Scale(n.kubeConfig, nodeName, "statefulset", 0, n.Namespace)
	if err != nil {
		return errors.Wrapf(err, "Scaling node %v to 0 replicas", nodeName)
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

func (n *SLNetwork) Destroy() error {
	log.Debug().Msg("Stopping network")

	err := kube.DeleteMainNamespace(n.kubeConfig, n.Namespace)
	if err != nil {
		return errors.Wrapf(err, "Deleting main namespace: %s", n.Namespace)
	}
	return nil
}

func (n *SLNetwork) IsNodeRunning(nodeName string) (bool, error) {
	log.Debug().Msg("Checking if node is running")
	scale, err := kube.GetScale(n.kubeConfig, nodeName, n.Namespace)
	if err != nil {
		return false, errors.Wrapf(err, "Getting scale for %v", nodeName)
	}
	return scale > 0, nil
}

func (n *SLNetwork) GetBitcoinNode(name string) (BitcoinNodeInterface, error) {
	for _, node := range n.BitcoinNodes {
		if node.GetName() == name {
			return node, nil
		}
	}
	return nil, errors.New("Bitcoin node not found")
}

func (n *SLNetwork) GetLightningNode(name string) (LightningNodeInterface, error) {
	for _, node := range n.LightningNodes {
		if node.GetName() == name {
			return node, nil
		}
	}
	return nil, errors.New("Lightning node not found")
}

func (n *SLNetwork) GetNode(name string) (NodeInterface, error) {
	for _, node := range n.BitcoinNodes {
		if node.GetName() == name {
			return node, nil
		}
	}
	for _, node := range n.LightningNodes {
		if node.GetName() == name {
			return node, nil
		}
	}
	return nil, errors.New("Node not found")
}

func (n *SLNetwork) GetAllNodes() []NodeInterface {
	nodes := []NodeInterface{}
	for _, node := range n.BitcoinNodes {
		nodes = append(nodes, node)
	}
	for _, node := range n.LightningNodes {
		nodes = append(nodes, node)
	}
	return nodes
}

type helmFile struct {
	Releases []struct {
		Name      string `yaml:"name"`
		Chart     string `yaml:"chart"`
		Namespace string `yaml:"namespace"`
	}
}

func parseHelmfileForNodes(
	slnetwork *SLNetwork,
) (lightningNodes []LightningNodeInterface, bitcoinNodes []BitcoinNodeInterface, namespace string, err error) {
	bytes, err := os.ReadFile(slnetwork.helmfile)
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "OS Reading helmfile")
	}

	helmFile := &helmFile{}
	err = yaml.Unmarshal(bytes, helmFile)
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "Unmarshalling helmfile")
	}

	// check that all the namespaces are the same, and read the namespaces from the helmfile
	namespace = ""

	for _, release := range helmFile.Releases {
		if namespace == "" {
			namespace = release.Namespace
		}
		if release.Namespace != namespace {
			return nil, nil, "", errors.New("All nodes must be in the same namespace, check the helmfile")
		}
	}

	for _, release := range helmFile.Releases {
		if strings.Contains(release.Chart, "lnd") {
			lightningNodes = append(
				lightningNodes,
				&lightningnode.LightningNode{
					Name:      release.Name,
					Impl:      lightningnode.LND,
					Namespace: release.Namespace,
				},
			)
		}
		if strings.Contains(release.Chart, "cln") {
			lightningNodes = append(
				lightningNodes,
				&lightningnode.LightningNode{
					Name:      release.Name,
					Impl:      lightningnode.CLN,
					Namespace: release.Namespace,
				},
			)
		}
		if strings.Contains(release.Chart, "bitcoind") {
			bitcoinNodes = append(
				bitcoinNodes,
				&bitcoinnode.BitcoinNode{
					Name:      release.Name,
					Namespace: release.Namespace,
				},
			)
		}
	}

	return lightningNodes, bitcoinNodes, namespace, nil
}

func (n *SLNetwork) GetWalletBalance(nodeName string) (types.Amount, error) {
	if !n.RetryCommands {
		return n.getWalletBalance(nodeName)
	}
	balance, err := tools.RetryWithReturn(func(cancel context.CancelFunc) (types.Amount, error) {
		balance, err := n.getWalletBalance(nodeName)
		if err != nil {
			return types.NewAmountSats(0), err
		}
		return balance, nil
	}, n.RetryDelay, n.RetryTimeout)
	if err != nil {
		return types.NewAmountSats(0), errors.Wrap(err, "Getting wallet balance")
	}
	return balance, nil
}

func (n *SLNetwork) getWalletBalance(nodeName string) (types.Amount, error) {
	conn, err := connectToGRPCServer(n.ApiHost, n.ApiPort, nodeName)
	if err != nil {
		return types.NewAmountSats(
				0,
			), errors.Wrapf(
				err,
				"Connecting to gRPC for %v's client",
				nodeName,
			)
	}
	defer conn.Close()
	client := stdcommonclient.NewCommonClient(conn)

	node, err := n.GetNode(nodeName)
	if err != nil {
		return types.NewAmountSats(0), err
	}

	walletBalance, err := node.GetWalletBalance(client)
	if err != nil {
		return types.NewAmountSats(
				0,
			), errors.Wrapf(
				err,
				"Getting wallet balance for %v",
				node.GetName(),
			)
	}
	return walletBalance, nil
}

func (n *SLNetwork) Generate(nodeName string, numBlocks uint32) (hashes []string, err error) {
	if !n.RetryCommands {
		return n.generate(nodeName, numBlocks)
	}
	hashes, err = tools.RetryWithReturn(func(cancel context.CancelFunc) ([]string, error) {
		hashes, err := n.generate(nodeName, numBlocks)
		if err != nil {
			return []string{}, err
		}
		return hashes, nil
	}, n.RetryDelay, n.RetryTimeout)
	if err != nil {
		return []string{}, errors.Wrap(err, "Generating blocks")
	}
	return hashes, nil
}

func (n *SLNetwork) generate(nodeName string, numBlocks uint32) (hashes []string, err error) {
	conn, err := connectToGRPCServer(n.ApiHost, n.ApiPort, nodeName)
	if err != nil {
		return []string{}, errors.Wrapf(err, "Connecting to gRPC for %v's client", nodeName)
	}
	defer conn.Close()

	client := stdbitcoinclient.NewBitcoinClient(conn)
	commonClient := stdcommonclient.NewCommonClient(conn)

	node, err := n.GetBitcoinNode(nodeName)
	if err != nil {
		return []string{}, err
	}

	hashes, err = node.Generate(client, commonClient, numBlocks)
	if err != nil {
		log.Error().Err(err).Msgf("Generating blocks for %v", node.GetName())
	}
	return hashes, nil
}

func (n *SLNetwork) GetNewAddress(nodeName string) (string, error) {
	if !n.RetryCommands {
		return n.getNewAddress(nodeName)
	}
	address, err := tools.RetryWithReturn(func(cancel context.CancelFunc) (string, error) {
		address, err := n.getNewAddress(nodeName)
		if err != nil {
			return "", err
		}
		return address, nil
	}, n.RetryDelay, n.RetryTimeout)
	if err != nil {
		return "", errors.Wrap(err, "Getting new address")
	}
	return address, nil
}
func (n *SLNetwork) getNewAddress(nodeName string) (string, error) {
	conn, err := connectToGRPCServer(n.ApiHost, n.ApiPort, nodeName)
	if err != nil {
		return "", errors.Wrapf(err, "Connecting to gRPC for %v's client", nodeName)
	}
	defer conn.Close()
	client := stdcommonclient.NewCommonClient(conn)

	node, err := n.GetNode(nodeName)
	if err != nil {
		return "", err
	}
	newAddress, err := node.GetNewAddress(client)
	if err != nil {
		return "", errors.Wrapf(err, "Getting new address for %v", node.GetName())
	}
	return newAddress, nil
}

func (n *SLNetwork) GetConnectionDetailsForAllNodes() ([]ConnectionDetails, error) {
	connectionDetails := []ConnectionDetails{}
	getDetailForNode := func(nodeName string) error {
		connDetails, err := n.GetConnectionDetails(nodeName)
		if err != nil {
			return err
		}
		connectionDetails = append(connectionDetails, connDetails...)
		return nil
	}
	for _, node := range n.BitcoinNodes {
		if err := getDetailForNode(node.GetName()); err != nil {
			return nil, errors.Wrapf(err, "Getting connection details for %v", node.GetName())
		}
	}
	for _, node := range n.LightningNodes {
		if err := getDetailForNode(node.GetName()); err != nil {
			return nil, errors.Wrapf(err, "Getting connection details for %v", node.GetName())
		}
	}
	return connectionDetails, nil
}

func (n *SLNetwork) GetConnectionDetails(nodeName string) ([]ConnectionDetails, error) {
	bitcoinNode, _ := n.GetBitcoinNode(nodeName)
	if bitcoinNode != nil {
		connectionPorts, err := bitcoinNode.GetConnectionPorts(n.kubeConfig)
		if err != nil {
			return nil, errors.Wrapf(err, "Getting connection ports for %v", nodeName)
		}
		connectionDetails := []ConnectionDetails{}
		for _, connectionPort := range connectionPorts {
			connectionDetails = append(
				connectionDetails,
				ConnectionDetails{
					NodeName: bitcoinNode.GetName(),
					Type:     connectionPort.Name,
					Host:     n.ApiHost,
					Port:     connectionPort.Port,
				},
			)
		}
		return connectionDetails, nil
	}

	node, err := n.GetLightningNode(nodeName)
	if err != nil {
		return nil, errors.Newf("Node not found. Available nodes: %v", n.listNodes(false, true))
	}

	port, err := node.GetConnectionPort(n.kubeConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "Getting grpc endpoint for %v", nodeName)
	}

	return []ConnectionDetails{
		{Type: "grpc", Host: n.ApiHost, Port: port, NodeName: node.GetName()},
	}, nil
}

func (n *SLNetwork) Send(
	fromNodeName string,
	toNodeName string,
	amountSats uint64,
) (string, error) {
	if !n.RetryCommands {
		return n.send(fromNodeName, toNodeName, amountSats)
	}
	txid, err := tools.RetryWithReturn(func(cancel context.CancelFunc) (string, error) {
		txid, err := n.send(fromNodeName, toNodeName, amountSats)
		if err != nil {
			return "", err
		}
		return txid, nil
	}, n.RetryDelay, n.RetryTimeout)
	if err != nil {
		return "", errors.Wrap(err, "Sending onchain")
	}
	return txid, nil
}

func (n *SLNetwork) send(
	fromNodeName string,
	toNodeName string,
	amountSats uint64,
) (string, error) {
	log.Debug().Msgf("Sending %v from %v to %v", amountSats, fromNodeName, toNodeName)

	address, err := n.GetNewAddress(toNodeName)
	if err != nil {
		return "", errors.Wrapf(err, "Getting new address for %v", toNodeName)
	}

	amount := types.NewAmountSats(amountSats)
	conn, err := connectToGRPCServer(n.ApiHost, n.ApiPort, fromNodeName)
	if err != nil {
		return "", errors.Wrapf(err, "Connecting to gRPC for %v's client", fromNodeName)
	}
	defer conn.Close()

	client := stdcommonclient.NewCommonClient(conn)

	node, err := n.GetNode(fromNodeName)
	if err != nil {
		return "", errors.Newf("From node not found. Available nodes: %v", n.listNodes(true, true))
	}

	txid, err := node.SendToAddress(client, address, amount)
	if err != nil {
		return "", errors.Wrapf(err, "Sending to addres")
	}

	_, err = n.Generate("bitcoind", 10)
	if err != nil {
		return "", errors.Wrapf(
			err,
			"TxId was: %v but problem generating blocks on %v",
			txid,
			"bitcoind",
		)
	}
	return txid, nil
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

func (n *SLNetwork) GetPubKey(nodeName string) (types.PubKey, error) {
	if !n.RetryCommands {
		return n.getPubKey(nodeName)
	}
	pubkey, err := tools.RetryWithReturn(func(cancel context.CancelFunc) (types.PubKey, error) {
		pubkey, err := n.getPubKey(nodeName)
		if err != nil {
			return types.PubKey{}, err
		}
		return pubkey, nil
	}, n.RetryDelay, n.RetryTimeout)
	if err != nil {
		return types.PubKey{}, errors.Wrap(err, "Getting pubkey")
	}
	return pubkey, nil
}

func (n *SLNetwork) getPubKey(nodeName string) (types.PubKey, error) {
	conn, err := connectToGRPCServer(n.ApiHost, n.ApiPort, nodeName)
	if err != nil {
		return types.PubKey{}, errors.Wrapf(err, "Connecting to gRPC for %v's client", nodeName)
	}
	defer conn.Close()

	client := stdlightningclient.NewLightningClient(conn)

	node, err := n.GetLightningNode(nodeName)
	if err != nil {
		return types.PubKey{}, errors.Newf(
			"Node not found. Available nodes: %v",
			n.listNodes(false, true),
		)
	}

	pubkey, err := node.GetPubKey(client)
	if err != nil {
		return types.PubKey{}, errors.Wrapf(err, "Getting pubkey for %v", node.GetName())
	}
	return pubkey, nil
}

func (n *SLNetwork) ConnectPeer(fromNodeName string, toNodeName string) error {
	if !n.RetryCommands {
		return n.connectPeer(fromNodeName, toNodeName)
	}
	err := tools.Retry(func(cancel context.CancelFunc) error {
		err := n.connectPeer(fromNodeName, toNodeName)
		if err != nil {
			return err
		}
		return nil
	}, n.RetryDelay, n.RetryTimeout)
	if err != nil {
		return errors.Wrap(err, "Connecting peer")
	}
	return nil
}

func (n *SLNetwork) connectPeer(fromNodeName string, toNodeName string) error {
	log.Debug().Msgf("Connecting %v to %v", fromNodeName, toNodeName)

	conn, err := connectToGRPCServer(n.ApiHost, n.ApiPort, fromNodeName)
	if err != nil {
		return errors.Wrapf(err, "Connecting to gRPC for %v's client", fromNodeName)
	}
	defer conn.Close()

	client := stdlightningclient.NewLightningClient(conn)

	toPubKey, err := n.GetPubKey(toNodeName)
	if err != nil {
		return errors.Wrapf(err, "Getting pubkey for peer: %v", toNodeName)
	}

	node, err := n.GetLightningNode(fromNodeName)
	if err != nil {
		return errors.Newf("Node not found. Available nodes: %v", n.listNodes(false, true))
	}

	err = node.ConnectPeer(client, toPubKey, toNodeName)
	if err != nil {
		return errors.Wrapf(err, "Connecting to %v", toNodeName)
	}
	return nil
}

func (n *SLNetwork) OpenChannel(
	fromNodeName string,
	toNodeName string,
	localAmountSats uint64,
) (types.ChannelPoint, error) {
	if !n.RetryCommands {
		return n.openChannel(fromNodeName, toNodeName, localAmountSats)
	}
	chanPoint, err := tools.RetryWithReturn(
		func(cancel context.CancelFunc) (types.ChannelPoint, error) {
			chanPoint, err := n.openChannel(fromNodeName, toNodeName, localAmountSats)
			if err != nil {
				return types.ChannelPoint{}, err
			}
			return chanPoint, nil
		},
		n.RetryDelay,
		n.RetryTimeout,
	)
	if err != nil {
		return types.ChannelPoint{}, errors.Wrap(err, "Opening channel")
	}
	return chanPoint, nil
}

func (n *SLNetwork) openChannel(
	fromNodeName string,
	toNodeName string,
	localAmountSats uint64,
) (types.ChannelPoint, error) {
	log.Debug().
		Msgf("Opening channel from %v to %v for %d sats", fromNodeName, toNodeName, localAmountSats)

	conn, err := connectToGRPCServer(n.ApiHost, n.ApiPort, fromNodeName)
	if err != nil {
		return types.ChannelPoint{}, errors.Wrapf(
			err,
			"Connecting to gRPC for %v's client",
			fromNodeName,
		)
	}
	defer conn.Close()

	client := stdlightningclient.NewLightningClient(conn)

	amount := types.NewAmountSats(localAmountSats)
	toPubKey, err := n.GetPubKey(toNodeName)
	if err != nil {
		return types.ChannelPoint{}, errors.Wrapf(err, "Getting pubkey for peer: %v", toNodeName)
	}

	node, err := n.GetLightningNode(fromNodeName)
	if err != nil {
		return types.ChannelPoint{}, errors.Newf(
			"Node not found. Available nodes: %v",
			n.listNodes(false, true),
		)
	}

	channelPoint, err := node.OpenChannel(client, toPubKey, amount)
	if err != nil {
		return types.ChannelPoint{}, errors.Wrapf(err, "Opening channel to %v", toNodeName)
	}
	_, err = n.Generate("bitcoind", 10)
	if err != nil {
		return types.ChannelPoint{}, errors.Wrapf(err,
			"Generating blocks after opening channel with channel point: %v", channelPoint)
	}
	return channelPoint, nil
}

func (n *SLNetwork) ChannelBalance(nodeName string) (types.Amount, error) {
	if !n.RetryCommands {
		return n.channelBalance(nodeName)
	}
	amount, err := tools.RetryWithReturn(func(cancel context.CancelFunc) (types.Amount, error) {
		amount, err := n.channelBalance(nodeName)
		if err != nil {
			return types.Amount{}, err
		}
		return amount, nil
	}, n.RetryDelay, n.RetryTimeout)
	if err != nil {
		return types.Amount{}, errors.Wrap(err, "Getting channel balance")
	}
	return amount, nil
}

func (n *SLNetwork) channelBalance(nodeName string) (types.Amount, error) {
	conn, err := connectToGRPCServer(n.ApiHost, n.ApiPort, nodeName)
	if err != nil {
		return types.NewAmountSats(
				0,
			), errors.Wrapf(
				err,
				"Connecting to gRPC for %v's client",
				nodeName,
			)
	}
	defer conn.Close()

	client := stdlightningclient.NewLightningClient(conn)

	node, err := n.GetLightningNode(nodeName)
	if err != nil {
		return types.NewAmountSats(
				0,
			), errors.Newf(
				"Node not found. Available nodes: %v",
				n.listNodes(false, true),
			)
	}

	channelBalance, err := node.ChannelBalance(client)
	if err != nil {
		return types.NewAmountSats(
				0,
			), errors.Wrapf(
				err,
				"Getting channel balance for %v",
				node.GetName(),
			)
	}
	return channelBalance, nil
}

func (n *SLNetwork) CreateInvoice(nodeName string, amountSats uint64) (string, error) {
	if !n.RetryCommands {
		return n.createInvoice(nodeName, amountSats)
	}
	invoice, err := tools.RetryWithReturn(func(cancel context.CancelFunc) (string, error) {
		invoice, err := n.createInvoice(nodeName, amountSats)
		if err != nil {
			return "", err
		}
		return invoice, nil
	}, n.RetryDelay, n.RetryTimeout)
	if err != nil {
		return "", errors.Wrap(err, "Creating invoice")
	}
	return invoice, nil
}

func (n *SLNetwork) createInvoice(nodeName string, amountSats uint64) (string, error) {
	conn, err := connectToGRPCServer(n.ApiHost, n.ApiPort, nodeName)
	if err != nil {
		return "", errors.Wrapf(err, "Connecting to gRPC for %v's client", nodeName)
	}
	defer conn.Close()

	client := stdlightningclient.NewLightningClient(conn)

	node, err := n.GetLightningNode(nodeName)
	if err != nil {
		return "", errors.Newf("Node not found. Available nodes: %v", n.listNodes(false, true))
	}

	invoice, err := node.CreateInvoice(client, amountSats)
	if err != nil {
		return "", errors.Wrapf(err, "Creating invoice for %v", node.GetName())
	}
	return invoice, nil
}

func (n *SLNetwork) PayInvoice(nodeName string, invoice string) (preimage string, err error) {
	if !n.RetryCommands {
		return n.payInvoice(nodeName, invoice)
	}
	preimage, err = tools.RetryWithReturn(func(cancel context.CancelFunc) (string, error) {
		preimage, err := n.payInvoice(nodeName, invoice)
		if err != nil {
			return "", err
		}
		return preimage, nil
	}, n.RetryDelay, n.RetryTimeout)
	if err != nil {
		return "", errors.Wrap(err, "Paying invoice")
	}
	return preimage, nil
}

func (n *SLNetwork) payInvoice(nodeName string, invoice string) (preimage string, err error) {
	conn, err := connectToGRPCServer(n.ApiHost, n.ApiPort, nodeName)
	if err != nil {
		return "", errors.Wrapf(err, "Connecting to gRPC for %v's client", nodeName)
	}
	defer conn.Close()

	client := stdlightningclient.NewLightningClient(conn)

	node, err := n.GetLightningNode(nodeName)
	if err != nil {
		return "", errors.Newf("Node not found. Available nodes: %v", n.listNodes(false, true))
	}

	preimage, err = node.PayInvoice(client, invoice)
	if err != nil {
		return "", errors.Wrapf(err, "Paying invoice for %v", node.GetName())
	}
	return preimage, nil
}

func (n *SLNetwork) listNodes(bitcoin bool, lightning bool) (nodes []string) {
	if bitcoin {
		for _, node := range n.BitcoinNodes {
			nodes = append(nodes, node.GetName())
		}
	}
	if lightning {
		for _, node := range n.LightningNodes {
			nodes = append(nodes, node.GetName())
		}
	}
	return nodes
}
