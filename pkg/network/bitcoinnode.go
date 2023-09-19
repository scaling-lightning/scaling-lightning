package network

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	stdbitcoinclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/bitcoin"
	stdcommonclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/common"

	basictypes "github.com/scaling-lightning/scaling-lightning/pkg/types"
)

type BitcoinNode struct {
	Name      string
	SLNetwork *SLNetwork
}

func (n *BitcoinNode) GetName() string {
	return n.Name
}

func (n *BitcoinNode) Send(to Node, amount basictypes.Amount) (string, error) {
	log.Debug().Msgf("Sending %v from %v to %v", amount, n.Name, to)

	toNode, err := n.SLNetwork.GetNode(to.GetName())
	if err != nil {
		return "", errors.Wrapf(err, "Looking up lightning node %v", to.GetName())
	}
	address, err := toNode.GetNewAddress()
	if err != nil {
		return "", errors.Wrapf(err, "Getting new address for %v", to.GetName())
	}

	addressRes, err := n.SendToAddress(address, amount)
	if err != nil {
		return "", errors.Wrapf(err, "Sending %v from %v to %v", amount, n.Name, to.GetName())
	}

	_, err = n.Generate(50)
	if err != nil {
		return "", errors.Wrapf(err, "Generating blocks for %v", "bitcoind")
	}

	return addressRes, nil
}

func (n *BitcoinNode) Generate(numBlocks uint32) (hashes []string, err error) {
	conn, err := connectToGRPCServer(n.SLNetwork.ApiHost, n.SLNetwork.ApiPort, n.Name)
	if err != nil {
		return []string{}, errors.Wrapf(err, "Connecting to gRPC for %v's client", n.Name)
	}
	defer conn.Close()
	client := stdbitcoinclient.NewBitcoinClient(conn)

	address, err := n.GetNewAddress()
	if err != nil {
		return []string{}, errors.Wrapf(err, "Getting new address for %v", n.Name)
	}
	generateRes, err := client.GenerateToAddress(
		context.Background(),
		&stdbitcoinclient.GenerateToAddressRequest{
			Address:     address,
			NumOfBlocks: numBlocks,
		},
	)
	if err != nil {
		return []string{}, errors.Wrapf(err, "Generating %v blocks for %v", numBlocks, n.Name)
	}

	return generateRes.Hashes, nil
}

func (n *BitcoinNode) GetWalletBalance() (basictypes.Amount, error) {
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
	return basictypes.NewAmountSats(walletBalance.BalanceSats), nil

}

func (n *BitcoinNode) SendToAddress(address string, amount basictypes.Amount) (string, error) {

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

func (n *BitcoinNode) GetNewAddress() (string, error) {
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

func (n *BitcoinNode) GetConnectionDetails() ([]ConnectionDetails, error) {
	rpcPort, err := getEndpointForNode(n.SLNetwork.kubeConfig, n.Name+"-direct-rpc", modeHTTP)
	if err != nil {
		return nil, errors.Wrapf(err, "Getting endpoint for %v", n.Name)
	}
	zmqBlockPort, err := getEndpointForNode(n.SLNetwork.kubeConfig, n.Name+"-direct-zmq-pub-block", modeTCP)
	if err != nil {
		return nil, errors.Wrapf(err, "Getting endpoint for %v", n.Name)
	}
	zmqTxPort, err := getEndpointForNode(n.SLNetwork.kubeConfig, n.Name+"-direct-zmq-pub-tx", modeTCP)
	if err != nil {
		return nil, errors.Wrapf(err, "Getting endpoint for %v", n.Name)
	}
	return []ConnectionDetails{
		{Name: "rpc", Host: n.SLNetwork.ApiHost, Port: rpcPort},
		{Name: "zmq blocks", Host: n.SLNetwork.ApiHost, Port: zmqBlockPort},
		{Name: "zmp txs", Host: n.SLNetwork.ApiHost, Port: zmqTxPort},
	}, err
}
