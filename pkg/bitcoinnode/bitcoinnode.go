package bitcoinnode

import (
	"context"

	"github.com/cockroachdb/errors"
	stdbitcoinclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/bitcoin"
	stdcommonclient "github.com/scaling-lightning/scaling-lightning/pkg/standardclient/common"

	basictypes "github.com/scaling-lightning/scaling-lightning/pkg/types"
)

//go:generate mockery --name LightningNodeInterface --exported
type BitcoinNodeInterface interface {
	GetName() string
	Generate(numBlocks uint32) (hashes []string, err error)
	GetWalletBalance() (basictypes.Amount, error)
	SendToAddress(address string, amount basictypes.Amount) (string, error)
	GetNewAddress() (string, error)
}

type BitcoinNode struct {
	Name string
	stdbitcoinclient.BitcoinClient
}

func (n *BitcoinNode) GetName() string {
	return n.Name
}

func (n *BitcoinNode) Generate(client stdbitcoinclient.BitcoinClient, commonClient stdcommonclient.CommonClient, numBlocks uint32) (hashes []string, err error) {

	address, err := n.GetNewAddress(commonClient)
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

func (n *BitcoinNode) GetWalletBalance(client stdcommonclient.CommonClient) (basictypes.Amount, error) {
	walletBalance, err := client.WalletBalance(
		context.Background(),
		&stdcommonclient.WalletBalanceRequest{},
	)
	if err != nil {
		return basictypes.Amount{}, errors.Wrapf(err, "Getting wallet balance for %v", n.Name)
	}
	return basictypes.NewAmountSats(walletBalance.BalanceSats), nil
}

func (n *BitcoinNode) SendToAddress(client stdcommonclient.CommonClient, address string, amount basictypes.Amount) (TxId string, err error) {
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

func (n *BitcoinNode) GetNewAddress(client stdcommonclient.CommonClient) (string, error) {
	newAddress, err := client.NewAddress(
		context.Background(),
		&stdcommonclient.NewAddressRequest{},
	)
	if err != nil {
		return "", errors.Wrapf(err, "Getting new address for %v", n.Name)
	}

	return newAddress.Address, nil
}
