package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"github.com/scaling-lightning/scaling-lightning/pkg/standardclient/types"
)

func CheckDependencies() error {
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

	podNameCmdStr := `get pods -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx --field-selector=status.phase==Running -o jsonpath='{.items[0].metadata.name}'`
	podNameCmd := exec.Command("kubectl", strings.Split(podNameCmdStr, " ")...)
	podNameOut, err := podNameCmd.Output()
	if err != nil {
		log.Debug().Err(err).Msgf("pod name output was: %v", string(podNameOut))
		return errors.Wrap(err, "Getting pod name for ingress-nginx")
	}

	podNameStr := strings.ReplaceAll(string(podNameOut), "'", "")
	nginxIngressCmdStr := `exec -n ingress-nginx -it ` + string(
		podNameStr,
	) + ` -- /nginx-ingress-controller --version`
	nginxIngressCmd := exec.Command("kubectl", strings.Split(nginxIngressCmdStr, " ")...)
	nginxIngressOut, err := nginxIngressCmd.Output()
	if err != nil {
		log.Debug().
			Err(err).
			Msgf("nginx-ingress-controller version output was: %v", string(nginxIngressOut))
		return errors.Wrap(err, "Getting nginx-ingress-controller version")
	}
	if !strings.Contains(strings.ToLower(string(nginxIngressOut)), "nginx ingress controller") {
		log.Debug().
			Err(err).
			Msgf("nginx-ingress-controller version output was: %v", string(nginxIngressOut))
		return errors.New("Ingress nginx not installed")
	}

	return nil
}

func StartViaHelmfile(helmfilePath string) error {
	log.Debug().Msg("Starting network")
	if err := CheckDependencies(); err != nil {
		return errors.Wrap(err, "Checking dependencies")
	}
	helmfileCmd := exec.Command("helmfile", "apply", "-f", helmfilePath)
	helmfileOut, err := helmfileCmd.Output()
	if err != nil {
		log.Debug().Err(err).Msgf("helmfile output was: %v", string(helmfileOut))
		return errors.Wrap(err, "Running helmfile apply command")
	}
	return nil
}

func StopViaHelmfile(helmfilePath string) error {
	log.Debug().Msg("Stopping network")
	helmfileCmd := exec.Command("helmfile", "destroy", "-f", helmfilePath)
	helmfileOut, err := helmfileCmd.Output()
	if err != nil {
		log.Debug().Err(err).Msgf("helmfile output was: %v", string(helmfileOut))
		return errors.Wrap(err, "Running helmfile apply command")
	}
	return nil
}

func Send(from string, to string, amount uint64) error {
	log.Debug().Msgf("Sending %v from %v to %v", amount, from, to)

	address, err := GetNewAddress(to)
	if err != nil {
		return errors.Wrapf(err, "Getting new address for %v", to)
	}

	err = SendToAddress(from, address, amount)
	if err != nil {
		return errors.Wrapf(err, "Sending %v from %v to %v", amount, from, to)
	}

	err = Generate(from, 50)
	if err != nil {
		return errors.Wrapf(err, "Generating blocks for %v", from)
	}

	return nil
}

func Generate(name string, numBlocks uint64) error {
	address, err := GetNewAddress(name)
	if err != nil {
		return errors.Wrapf(err, "Getting new address for %v", name)
	}
	req := types.GenerateToAddressReq{Address: address, NumOfBlocks: numBlocks}
	postBody, _ := json.Marshal(req)
	postBuf := bytes.NewBuffer(postBody)
	resp, err := http.Post(
		fmt.Sprintf("http://localhost/%v/generatetoaddress", name),
		"application/json",
		postBuf,
	)
	if err != nil {
		return errors.Wrapf(err, "Sending POST request to %v/generatetoaddress", name)
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
			name,
			resp.StatusCode,
		)
	}
	return nil
}

func SendToAddress(fromName string, toAddress string, amount uint64) error {
	req := types.SendToAddressReq{Address: toAddress, AmtSats: amount}
	postBody, _ := json.Marshal(req)
	postBuf := bytes.NewBuffer(postBody)
	resp, err := http.Post(
		fmt.Sprintf("http://localhost/%v/sendtoaddress", fromName),
		"application/json",
		postBuf,
	)
	if err != nil {
		return errors.Wrapf(err, "Sending POST request to %v/sendtoaddress", fromName)
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			log.Debug().Msgf("Response body to failed sendtoaddress request was: %v", string(body))
		}
		return errors.Newf(
			"Got non-200 status code from %v/sendtoaddress: %v",
			fromName,
			resp.StatusCode,
		)
	}
	return nil
}

func GetNewAddress(name string) (string, error) {
	resp, err := http.Post(
		fmt.Sprintf("http://localhost/%v/newaddress", name),
		"application/json",
		nil,
	)
	if err != nil {
		return "", errors.Wrapf(err, "Sending POST request to %v/newaddress", name)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "Reading response body from %v/newaddress", name)
	}
	var newAddress types.NewAddressRes
	err = json.Unmarshal(body, &newAddress)
	if err != nil {
		fmt.Println("error:", err)
	}
	return newAddress.Address, nil
}

func GetPubKey(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost/%v/pubkey", name))
	if err != nil {
		return "", errors.Wrapf(err, "Sending GET request to %v/pubkey", name)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "Reading response body from %v/pubkey", name)
	}
	var pubKey types.PubKeyRes
	err = json.Unmarshal(body, &pubKey)
	if err != nil {
		fmt.Println("error:", err)
	}
	return pubKey.PubKey, nil
}

func GetWalletBalanceSats(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost/%v/walletbalace", name))
	if err != nil {
		return "", errors.Wrapf(err, "Sending GET request to %v/walletbalance", name)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "Reading response body from %v/walletbalance", name)
	}
	return string(body), nil
}

func ConnectPeer(fromName string, toName string) error {
	log.Debug().Msgf("Connecting %v to %v", fromName, toName)
	toPubKey, err := GetPubKey(toName)
	if err != nil {
		return errors.Wrapf(err, "Getting pubkey for %v", toName)
	}
	req := types.ConnectPeerReq{PubKey: toPubKey, Host: toName, Port: 9735}
	postBody, _ := json.Marshal(req)
	postBuf := bytes.NewBuffer(postBody)
	resp, err := http.Post(
		fmt.Sprintf("http://localhost/%v/connectpeer", fromName),
		"application/json",
		postBuf,
	)
	if err != nil {
		return errors.Wrapf(err, "Sending POST request to %v/connectpeer", fromName)
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			if strings.Contains(string(body), "already connected") {
				return nil
			}
			log.Debug().Msgf("Response body to failed connectpeer request was: %v", string(body))
		}
		return errors.Newf(
			"Got non-200 status code from %v/connectpeer: %v",
			fromName,
			resp.StatusCode,
		)
	}
	return nil
}

func OpenChannel(fromName string, toName string, localAmtSats uint64) error {
	log.Debug().Msgf("Opening channel from %v to %v for %d sats", fromName, toName, localAmtSats)
	toPubKey, err := GetPubKey(toName)
	if err != nil {
		return errors.Wrapf(err, "Getting pubkey for %v", toName)
	}
	req := types.OpenChannelReq{PubKey: toPubKey, LocalAmtSats: localAmtSats}
	postBody, _ := json.Marshal(req)
	postBuf := bytes.NewBuffer(postBody)
	resp, err := http.Post(
		fmt.Sprintf("http://localhost/%v/openchannel", fromName),
		"application/json",
		postBuf,
	)
	if err != nil {
		return errors.Wrapf(err, "Sending POST request to %v/openchannel", fromName)
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			log.Debug().Msgf("Response body to failed openchanel request was: %v", string(body))
		}
		return errors.Newf(
			"Got non-200 status code from %v/openchannel: %v",
			fromName,
			resp.StatusCode,
		)
	}
	return nil
}
