package network

import (
	"context"
	"fmt"
	"path"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	stdcommonclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/common"
	stdlightningclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/lightning"
	basictypes "github.com/scaling-lightning/scaling-lightning/pkg/types"
)

type NodeImpl int

const (
	LND NodeImpl = iota
	CLN
	LDK
	Eclair
)

type LightningNode struct {
	Name        string
	BitcoinNode *BitcoinNode
	SLNetwork   *SLNetwork
	Impl        NodeImpl
}

type ConnectionDetails struct {
	Host string
	Port uint16
	LND  *LNDConnectionDetails
	CLN  *CLNConnectionDetails
}

type LNDConnectionDetails struct {
	TLSCert  []byte
	Macaroon []byte
}

type CLNConnectionDetails struct {
	LightningNode
	ClientCert []byte
	ClientKey  []byte
	CACert     []byte
}

func (n *LightningNode) GetName() string {
	return n.Name
}

func (n *LightningNode) Send(to Node, amount basictypes.Amount) (string, error) {
	log.Debug().Msgf("Sending %v from %v to %v", amount, n.Name, to)

	var toNode Node
	toNode, err := n.SLNetwork.GetNode(to.GetName())
	if err != nil {
		return "", errors.Wrapf(err, "Looking up lightning node %v", to.GetName())
	}
	address, err := toNode.GetNewAddress()
	if err != nil {
		return "", errors.Wrapf(err, "Getting new address for %v", to.GetName())
	}

	txid, err := n.SendToAddress(address, amount)
	if err != nil {
		return "", errors.Wrapf(err, "Sending %v from %v to %v", amount, n.Name, to.GetName())
	}

	_, err = n.BitcoinNode.Generate(50)
	if err != nil {
		return "", errors.Wrapf(err, "Generating blocks for %v", "bitcoind")
	}

	return txid, nil
}

func (n *LightningNode) SendToAddress(
	address string,
	amount basictypes.Amount,
) (string, error) {

	conn, err := connectToGRPCServer(n.SLNetwork.ApiHost, n.SLNetwork.ApiPort, n.Name)
	if err != nil {
		return "", errors.Wrapf(err, "Connecting to gRPC for %v's client", n.Name)
	}
	defer conn.Close()
	client := stdcommonclient.NewCommonClient(conn)

	sendRes, err := client.Send(
		context.Background(),
		&stdcommonclient.SendRequest{
			Address: address,
			Amount:  amount.AsSats(),
		},
	)
	if err != nil {
		return "", errors.Wrapf(err, "Sending %v to %v", amount, address)
	}
	return sendRes.TxId, nil
}

func (n *LightningNode) GetNewAddress() (string, error) {
	conn, err := connectToGRPCServer(n.SLNetwork.ApiHost, n.SLNetwork.ApiPort, n.Name)
	if err != nil {
		return "", errors.Wrapf(err, "Connecting to gRPC for %v's client", n.Name)
	}
	defer conn.Close()
	client := stdcommonclient.NewCommonClient(conn)

	newAddress, err := client.NewAddress(
		context.Background(),
		&stdcommonclient.NewAddressRequest{},
	)
	if err != nil {
		return "", errors.Wrapf(err, "Getting new address for %v", n.Name)
	}

	return newAddress.Address, nil
}

func (n *LightningNode) GetPubKey() (basictypes.PubKey, error) {
	conn, err := connectToGRPCServer(n.SLNetwork.ApiHost, n.SLNetwork.ApiPort, n.Name)
	if err != nil {
		return basictypes.PubKey{}, errors.Wrapf(err, "Connecting to gRPC for %v's client", n.Name)
	}
	defer conn.Close()
	client := stdlightningclient.NewLightningClient(conn)

	pubKeyRes, err := client.PubKey(context.Background(), &stdlightningclient.PubKeyRequest{})
	if err != nil {
		return basictypes.PubKey{}, errors.Wrapf(err, "Getting pubkey for %v", n.Name)
	}

	return basictypes.NewPubKeyFromByte(pubKeyRes.PubKey), nil
}

func (n *LightningNode) GetWalletBalance() (basictypes.Amount, error) {
	conn, err := connectToGRPCServer(n.SLNetwork.ApiHost, n.SLNetwork.ApiPort, n.Name)
	if err != nil {
		return basictypes.Amount{}, errors.Wrapf(err, "Connecting to gRPC for %v's client", n.Name)
	}
	defer conn.Close()
	client := stdcommonclient.NewCommonClient(conn)
	walletBalance, err := client.WalletBalance(
		context.Background(),
		&stdcommonclient.WalletBalanceRequest{},
	)
	if err != nil {
		return basictypes.Amount{}, errors.Wrapf(err, "Getting wallet balance for %v", n.Name)
	}
	return basictypes.NewAmountSats(walletBalance.Balance), nil
}

func (n *LightningNode) ConnectPeer(to Node) error {
	log.Debug().Msgf("Connecting %v to %v", n.Name, to)
	conn, err := connectToGRPCServer(n.SLNetwork.ApiHost, n.SLNetwork.ApiPort, n.Name)
	if err != nil {
		return errors.Wrapf(err, "Connecting to gRPC for %v's client", n.Name)
	}
	defer conn.Close()
	client := stdlightningclient.NewLightningClient(conn)

	toNode, err := n.SLNetwork.GetLightningNode(to.GetName())
	if err != nil {
		return errors.Wrapf(err, "Looking up lightning node for %v", to.GetName())
	}
	toPubKey, err := toNode.GetPubKey()
	if err != nil {
		return errors.Wrapf(err, "Getting pubkey for %v", to)
	}

	_, err = client.ConnectPeer(
		context.Background(),
		&stdlightningclient.ConnectPeerRequest{
			PubKey: toPubKey.AsBytes(),
			Host:   to.GetName(),
			Port:   9735,
		},
	)
	if err != nil {
		return errors.Wrapf(err, "Connecting %v to %v", n.Name, to)
	}

	return nil
}

func (n *LightningNode) OpenChannel(
	to *LightningNode,
	localAmt basictypes.Amount,
) (basictypes.ChannelPoint, error) {
	log.Debug().
		Msgf("Opening channel from %v to %v for %d sats", n.Name, to.GetName(), localAmt.AsSats())

	conn, err := connectToGRPCServer(n.SLNetwork.ApiHost, n.SLNetwork.ApiPort, n.Name)
	if err != nil {
		return basictypes.ChannelPoint{}, errors.Wrapf(
			err,
			"Connecting to gRPC for %v's client",
			n.Name,
		)
	}
	defer conn.Close()
	client := stdlightningclient.NewLightningClient(conn)

	toNode, err := n.SLNetwork.GetLightningNode(to.GetName())
	if err != nil {
		return basictypes.ChannelPoint{}, errors.Wrapf(
			err,
			"Looking up lightning node for %v",
			to.GetName(),
		)
	}
	toPubKey, err := toNode.GetPubKey()
	if err != nil {
		return basictypes.ChannelPoint{}, errors.Wrapf(err, "Getting pubkey for %v", to.GetName())
	}

	openChannelRes, err := client.OpenChannel(
		context.Background(),
		&stdlightningclient.OpenChannelRequest{
			PubKey:       toPubKey.AsBytes(),
			LocalAmtSats: localAmt.AsSats(),
		},
	)
	if err != nil {
		return basictypes.ChannelPoint{}, errors.Wrapf(
			err,
			"Opening channel from %v to %v for %d sats",
			n.Name,
			to.GetName(),
			localAmt.AsSats(),
		)
	}
	_, err = n.BitcoinNode.Generate(50)
	if err != nil {
		return basictypes.ChannelPoint{}, errors.Wrapf(err, "Generating blocks for %v", "bitcoind")
	}
	return basictypes.ChannelPoint{
		FundingTx:   basictypes.NewTransactionFromByte(openChannelRes.FundingTxId),
		OutputIndex: uint(openChannelRes.FundingTxOutputIndex),
	}, nil
}

func (n *LightningNode) WriteAuthFilesToDirectory(dir string) error {
	network := n.SLNetwork.Network.String()
	switch n.Impl {
	case LND:
		err := kubeCP(
			n.SLNetwork.kubeConfig,
			fmt.Sprintf("%v-0:root/.lnd/tls.cert", n.Name),
			path.Join(dir, "tls.cert"),
		)
		if err != nil {
			return errors.Wrap(err, "KubeCP LND's tls.cert")
		}
		err = kubeCP(
			n.SLNetwork.kubeConfig,
			fmt.Sprintf("%v-0:root/.lnd/data/chain/bitcoin/%v/admin.macaroon", n.Name, network),
			path.Join(dir, "admin.macaroon"),
		)
		if err != nil {
			return errors.Wrap(err, "KubeCP LND's admin.macaroon")
		}
	case CLN:
		err := kubeCP(
			n.SLNetwork.kubeConfig,
			fmt.Sprintf("%v-0:root/.lightning/%v/client.pem", n.Name, network),
			path.Join(dir, "client.pem"),
		)
		if err != nil {
			return errors.Wrap(err, "KubeCP CLN's client.pem")
		}
		err = kubeCP(
			n.SLNetwork.kubeConfig,
			fmt.Sprintf("%v-0:root/.lightning/%v/client-key.pem", n.Name, network),
			path.Join(dir, "client-key.pem"),
		)
		if err != nil {
			return errors.Wrap(err, "KubeCP CLN's client-key.pem")
		}
		err = kubeCP(
			n.SLNetwork.kubeConfig,
			fmt.Sprintf("%v-0:root/.lightning/%v/ca.pem", n.Name, network),
			path.Join(dir, "ca.pem"),
		)
		if err != nil {
			return errors.Wrap(err, "KubeCP CLN's ca.pem")
		}
	}
	return nil
}

func (n *LightningNode) GetConnectionDetails() ConnectionDetails {
	return ConnectionDetails{Host: n.SLNetwork.ApiHost, Port: 12345}
}
