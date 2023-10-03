package lightningnode

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"path"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"

	"github.com/scaling-lightning/scaling-lightning/pkg/kube"
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

type LightningNode struct {
	Name string
	Impl NodeImpl
}

func (n *LightningNode) GetName() string {
	return n.Name
}

func (n *LightningNode) SendToAddress(
	client stdcommonclient.CommonClient,
	address string,
	amount basictypes.Amount,
) (string, error) {

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

func (n *LightningNode) GetNewAddress(client stdcommonclient.CommonClient) (string, error) {
	newAddress, err := client.NewAddress(
		context.Background(),
		&stdcommonclient.NewAddressRequest{},
	)
	if err != nil {
		return "", errors.Wrapf(err, "Getting new address for %v", n.Name)
	}

	return newAddress.Address, nil
}

func (n *LightningNode) GetPubKey(client stdlightningclient.LightningClient) (basictypes.PubKey, error) {
	pubKeyRes, err := client.PubKey(context.Background(), &stdlightningclient.PubKeyRequest{})
	if err != nil {
		return basictypes.PubKey{}, errors.Wrapf(err, "Getting pubkey for %v", n.Name)
	}

	return basictypes.NewPubKeyFromByte(pubKeyRes.PubKey), nil
}

func (n *LightningNode) GetWalletBalance(client stdcommonclient.CommonClient) (basictypes.Amount, error) {
	walletBalance, err := client.WalletBalance(
		context.Background(),
		&stdcommonclient.WalletBalanceRequest{},
	)
	if err != nil {
		return basictypes.Amount{}, errors.Wrapf(err, "Getting wallet balance for %v", n.Name)
	}
	return basictypes.NewAmountSats(walletBalance.BalanceSats), nil
}

func (n *LightningNode) ConnectPeer(client stdlightningclient.LightningClient, pubkey basictypes.PubKey, nodeName string) error {
	_, err := client.ConnectPeer(
		context.Background(),
		&stdlightningclient.ConnectPeerRequest{
			PubKey: pubkey.AsBytes(),
			Host:   nodeName,
			Port:   9735,
		},
	)
	if err != nil {
		return errors.Wrapf(err, "Connecting %v to %v", n.Name, pubkey.AsHexString())
	}

	return nil
}

func (n *LightningNode) OpenChannel(
	client stdlightningclient.LightningClient,
	pubkey basictypes.PubKey,
	localAmt basictypes.Amount,
) (basictypes.ChannelPoint, error) {
	log.Debug().
		Msgf("Opening channel from %v to %v for %d sats", n.Name, pubkey.AsHexString(), localAmt.AsSats())

	openChannelRes, err := client.OpenChannel(
		context.Background(),
		&stdlightningclient.OpenChannelRequest{
			PubKey:       pubkey.AsBytes(),
			LocalAmtSats: localAmt.AsSats(),
		},
	)
	if err != nil {
		return basictypes.ChannelPoint{}, errors.Wrapf(
			err,
			"Opening channel from %v to %v for %d sats",
			n.Name,
			pubkey.AsHexString(),
			localAmt.AsSats(),
		)
	}
	return basictypes.ChannelPoint{
		FundingTx:   basictypes.NewTransactionFromByte(openChannelRes.FundingTxId),
		OutputIndex: uint(openChannelRes.FundingTxOutputIndex),
	}, nil
}

func (n *LightningNode) GetConnectionFiles(network string, kubeConfig string) (*ConnectionFiles, error) {
	dir, err := os.MkdirTemp("", "slconnectionfiles")
	if err != nil {
		return nil, errors.Wrap(err, "Creating temp dir")
	}
	defer os.RemoveAll(dir)

	err = n.WriteAuthFilesToDirectory(network, kubeConfig, dir)
	if err != nil {
		return nil, errors.Wrap(err, "Writing auth files")
	}

	switch n.Impl {
	case LND:
		tlsCert, err := os.ReadFile(path.Join(dir, "tls.cert"))
		if err != nil {
			return nil, errors.Wrap(err, "Reading tls.cert")
		}
		macaroon, err := os.ReadFile(path.Join(dir, "admin.macaroon"))
		if err != nil {
			return nil, errors.Wrap(err, "Reading admin.macaroon")
		}
		return &ConnectionFiles{LND: &LNDConnectionFiles{
			TLSCert:  tlsCert,
			Macaroon: macaroon,
		}}, nil
	case CLN:
		clientCert, err := os.ReadFile(path.Join(dir, "client.pem"))
		if err != nil {
			return nil, errors.Wrap(err, "Reading client.pem")
		}
		clientKey, err := os.ReadFile(path.Join(dir, "client-key.pem"))
		if err != nil {
			return nil, errors.Wrap(err, "Reading client-key.pem")
		}
		caCert, err := os.ReadFile(path.Join(dir, "ca.pem"))
		if err != nil {
			return nil, errors.Wrap(err, "Reading ca.pem")
		}
		return &ConnectionFiles{CLN: &CLNConnectionFiles{
			ClientCert: clientCert,
			ClientKey:  clientKey,
			CACert:     caCert,
		}}, nil
	default:
		return nil, errors.Newf("Unknown node implementation: %v", n.Impl)
	}
}

func (n *LightningNode) WriteAuthFilesToDirectory(network string, kubeConfig string, dir string) error {
	switch n.Impl {
	case LND:
		err := kube.KubeCP(
			kubeConfig,
			fmt.Sprintf("%v-0:root/.lnd/tls.cert", n.Name),
			path.Join(dir, "tls.cert"),
		)
		if err != nil {
			return errors.Wrap(err, "KubeCP LND's tls.cert")
		}
		err = kube.KubeCP(
			kubeConfig,
			fmt.Sprintf("%v-0:root/.lnd/data/chain/bitcoin/%v/admin.macaroon", n.Name, network),
			path.Join(dir, "admin.macaroon"),
		)
		if err != nil {
			return errors.Wrap(err, "KubeCP LND's admin.macaroon")
		}
	case CLN:
		err := kube.KubeCP(
			kubeConfig,
			fmt.Sprintf("%v-0:root/.lightning/%v/client.pem", n.Name, network),
			path.Join(dir, "client.pem"),
		)
		if err != nil {
			return errors.Wrap(err, "KubeCP CLN's client.pem")
		}
		err = kube.KubeCP(
			kubeConfig,
			fmt.Sprintf("%v-0:root/.lightning/%v/client-key.pem", n.Name, network),
			path.Join(dir, "client-key.pem"),
		)
		if err != nil {
			return errors.Wrap(err, "KubeCP CLN's client-key.pem")
		}
		err = kube.KubeCP(
			kubeConfig,
			fmt.Sprintf("%v-0:root/.lightning/%v/ca.pem", n.Name, network),
			path.Join(dir, "ca.pem"),
		)
		if err != nil {
			return errors.Wrap(err, "KubeCP CLN's ca.pem")
		}
	}
	return nil
}

func (n *LightningNode) GetConnectionPort(kubeConfig string) (uint16, error) {
	port, err := kube.GetEndpointForNode(kubeConfig, n.Name+"-direct-grpc", kube.ModeTCP)
	if err != nil {
		return 0, errors.Wrapf(err, "Getting endpoint for %v", n.Name)
	}
	return port, nil
}

func (n *LightningNode) CreateInvoice(client stdlightningclient.LightningClient, amountSats uint64) (string, error) {
	invoiceRes, err := client.CreateInvoice(
		context.Background(),
		&stdlightningclient.CreateInvoiceRequest{
			AmtSats: amountSats,
		},
	)
	if err != nil {
		return "", errors.Wrapf(err, "Creating invoice for %v", n.Name)
	}
	return invoiceRes.Invoice, nil
}

func (n *LightningNode) PayInvoice(client stdlightningclient.LightningClient, invoice string) (string, error) {
	payRes, err := client.PayInvoice(
		context.Background(),
		&stdlightningclient.PayInvoiceRequest{
			Invoice: invoice,
		},
	)
	if err != nil {
		return "", errors.Wrapf(err, "Paying invoice for %v", n.Name)
	}

	preImageStr := hex.EncodeToString(payRes.PaymentPreimage)
	return preImageStr, nil
}

func (n *LightningNode) ChannelBalance(client stdlightningclient.LightningClient) (basictypes.Amount, error) {
	balanceRes, err := client.ChannelBalance(
		context.Background(),
		&stdlightningclient.ChannelBalanceRequest{},
	)
	if err != nil {
		return basictypes.Amount{}, errors.Wrapf(err, "Getting channel balance for %v", n.Name)
	}
	return basictypes.NewAmountSats(balanceRes.BalanceSats), nil
}
