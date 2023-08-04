package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/types"

	basictypes "github.com/scaling-lightning/scaling-lightning/pkg/types"
)

type BitcoinNode struct {
	Name      string
	SLNetwork *SLNetwork
}

func (n *BitcoinNode) GetName() string {
	return n.Name
}

func (n *BitcoinNode) Send(to Node, amount basictypes.Amount) error {
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

	err = n.Generate(50)
	if err != nil {
		return errors.Wrapf(err, "Generating blocks for %v", "bitcoind")
	}

	return nil
}

func (n *BitcoinNode) Generate(numBlocks uint64) error {
	address, err := n.GetNewAddress()
	if err != nil {
		return errors.Wrapf(err, "Getting new address for %v", n.Name)
	}
	req := types.GenerateToAddressReq{Address: address.AsBase58String(), NumOfBlocks: numBlocks}
	postBody, _ := json.Marshal(req)
	postBuf := bytes.NewBuffer(postBody)
	resp, err := http.Post(
		fmt.Sprintf("http://localhost/%v/generatetoaddress", n.Name),
		"application/json",
		postBuf,
	)
	if err != nil {
		return errors.Wrapf(err, "Sending POST request to %v/generatetoaddress", n.Name)
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			log.Debug().
				Msgf("Response body to failed generatetoaddress request was: %v", string(body))
		}
		return errors.Newf(
			"Got non-200 status code from %v/generatetoaddress: %v",
			n.Name,
			resp.StatusCode,
		)
	}
	return nil
}

func (n *BitcoinNode) SendToAddress(address basictypes.Address, amount basictypes.Amount) error {
	req := types.SendToAddressReq{Address: address.AsBase58String(), AmtSats: amount.AsSats()}
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

func (n *BitcoinNode) GetNewAddress() (basictypes.Address, error) {
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
