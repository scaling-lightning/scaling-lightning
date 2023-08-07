package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/lightning"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/types"
	"github.com/scaling-lightning/scaling-lightning/pkg/tools/grpc_helpers"
	basictypes "github.com/scaling-lightning/scaling-lightning/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LightningNode struct {
	Name        string
	Host        string
	Port        int
	BitcoinNode *BitcoinNode
	SLNetwork   *SLNetwork
}

func (n *LightningNode) GetName() string {
	return n.Name
}

func (n *LightningNode) Send(to Node, amount basictypes.Amount) error {
	log.Debug().Msgf("Sending %v from %v to %v", amount, n.Name, to)

	var toNode Node
	toNode, err := n.SLNetwork.GetLightningNode(to.GetName())
	if err != nil {
		toNode, err = n.SLNetwork.GetBitcoinNode(to.GetName())
		if err != nil {
			return errors.Wrapf(err, "Looking up lightning node %v", to.GetName())
		}
	}
	address, err := toNode.GetNewAddress()
	if err != nil {
		return errors.Wrapf(err, "Getting new address for %v", to.GetName())
	}

	err = n.SendToAddress(address, amount)
	if err != nil {
		return errors.Wrapf(err, "Sending %v from %v to %v", amount, n.Name, to.GetName())
	}

	err = n.BitcoinNode.Generate(50)
	if err != nil {
		return errors.Wrapf(err, "Generating blocks for %v", "bitcoind")
	}

	return nil
}

func (n *LightningNode) SendToAddress(
	toAddress basictypes.Address,
	amount basictypes.Amount,
) error {
	req := types.SendToAddressReq{Address: toAddress.AsBase58String(), AmtSats: amount.AsSats()}
	postBody, _ := json.Marshal(req)
	postBuf := bytes.NewBuffer(postBody)
	resp, err := http.Post(
		fmt.Sprintf("http://localhost/%v/sendtoaddress", n.Name),
		"application/json",
		postBuf,
	)
	if err != nil {
		return errors.Wrapf(err, "Sending POST request to %v/sendtoaddress", n.Name)
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			log.Debug().Msgf("Response body to failed sendtoaddress request was: %v", string(body))
		}
		return errors.Newf(
			"Got non-200 status code from %v/sendtoaddress: %v",
			n.Name,
			resp.StatusCode,
		)
	}
	return nil
}

func (n *LightningNode) GetNewAddress() (basictypes.Address, error) {
	resp, err := http.Post(
		fmt.Sprintf("http://localhost/%v/newaddress", n.Name),
		"application/json",
		nil,
	)
	if err != nil {
		return basictypes.Address{}, errors.Wrapf(
			err,
			"Sending POST request to %v/newaddress",
			n.Name,
		)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return basictypes.Address{}, errors.Wrapf(
			err,
			"Reading response body from %v/newaddress",
			n.Name,
		)
	}
	var newAddress types.NewAddressRes
	err = json.Unmarshal(body, &newAddress)
	if err != nil {
		fmt.Println("error:", err)
	}
	return basictypes.NewAddressFromBase58String(newAddress.Address), nil
}

func (n *LightningNode) GetPubKey() (basictypes.PubKey, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost/%v/pubkey", n.Name))
	if err != nil {
		return basictypes.PubKey{}, errors.Wrapf(err, "Sending GET request to %v/pubkey", n.Name)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return basictypes.PubKey{}, errors.Wrapf(
			err,
			"Reading response body from %v/pubkey",
			n.Name,
		)
	}
	var pubKeyRes types.PubKeyRes
	err = json.Unmarshal(body, &pubKeyRes)
	if err != nil {
		return basictypes.PubKey{}, errors.Wrapf(
			err,
			"Unmarshalling pubkey response body from %v",
			n.Name,
		)
	}

	pubKey, err := basictypes.NewPubKeyFromHexString(pubKeyRes.PubKey)
	if err != nil {
		return basictypes.PubKey{}, errors.Wrapf(err, "Converting pubkey from hex string")
	}

	return pubKey, nil
}

func (n *LightningNode) GetWalletBalance() (basictypes.Amount, error) {

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_helpers.ClientInterceptor(n.Name)),
	}
	conn, err := grpc.Dial("localhost:80", opts...)
	if err != nil {
		return basictypes.Amount{}, errors.Wrapf(err, "Connecting to gRPC for %v's client", n.Name)
	}
	defer conn.Close()
	client := lightning.NewCommonClient(conn)
	walletBalance, err := client.WalletBalance(
		context.Background(),
		&lightning.WalletBalanceRequest{},
	)
	if err != nil {
		return basictypes.Amount{}, errors.Wrapf(err, "Getting wallet balance for %v", n.Name)
	}
	return basictypes.NewAmountSats(walletBalance.Balance), nil

}

func (n *LightningNode) ConnectPeer(to Node) error {
	log.Debug().Msgf("Connecting %v to %v", n.Name, to)
	toNode, err := n.SLNetwork.GetLightningNode(to.GetName())
	if err != nil {
		return errors.Wrapf(err, "Looking up lightning node for %v", to.GetName())
	}
	toPubKey, err := toNode.GetPubKey()
	if err != nil {
		return errors.Wrapf(err, "Getting pubkey for %v", to)
	}
	req := types.ConnectPeerReq{PubKey: toPubKey.AsHexString(), Host: to.GetName(), Port: 9735}
	postBody, _ := json.Marshal(req)
	postBuf := bytes.NewBuffer(postBody)
	resp, err := http.Post(
		fmt.Sprintf("http://localhost/%v/connectpeer", n.Name),
		"application/json",
		postBuf,
	)
	if err != nil {
		return errors.Wrapf(err, "Sending POST request to %v/connectpeer", n.Name)
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "Status code was not 200, and error reading response body")
		}
		if strings.Contains(string(body), "already connected") {
			return nil
		}
		return errors.Newf(
			"Problem calling %v/connectpeer: %v",
			n.Name,
			string(body),
		)
	}
	return nil
}

func (n *LightningNode) OpenChannel(to *LightningNode, localAmt basictypes.Amount) error {
	log.Debug().
		Msgf("Opening channel from %v to %v for %d sats", n.Name, to.GetName(), localAmt.AsSats())

	toNode, err := n.SLNetwork.GetLightningNode(to.GetName())
	if err != nil {
		return errors.Wrapf(err, "Looking up lightning node for %v", to.GetName())
	}
	toPubKey, err := toNode.GetPubKey()
	if err != nil {
		return errors.Wrapf(err, "Getting pubkey for %v", to.GetName())
	}
	req := types.OpenChannelReq{PubKey: toPubKey.AsHexString(), LocalAmtSats: localAmt.AsSats()}
	postBody, _ := json.Marshal(req)
	postBuf := bytes.NewBuffer(postBody)
	resp, err := http.Post(
		fmt.Sprintf("http://localhost/%v/openchannel", n.Name),
		"application/json",
		postBuf,
	)
	if err != nil {
		return errors.Wrapf(err, "Sending POST request to %v/openchannel", n.Name)
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "Status code was not 200, and error reading response body")
		}
		return errors.Newf(
			"Problem calling %v/openchannel: %v",
			n.Name,
			string(body),
		)
	}
	err = n.BitcoinNode.Generate(50)
	if err != nil {
		return errors.Wrapf(err, "Generating blocks for %v", "bitcoind")
	}
	return nil
}
