package types

type NewAddressRes struct {
	Address string `json:"address"`
}

type PubKeyRes struct {
	PubKey string `json:"pubkey"`
}

type ConnectPeerReq struct {
	PubKey string `json:"pubKey"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}

type OpenChannelReq struct {
	PubKey       string `json:"pubKey"`
	LocalAmtSats uint64 `json:"localAmountSats"`
}

type OpenChannelRes struct {
	FundingTx   string `json:"fundingTx"`
	OutputIndex uint32 `json:"outputIndex"`
}

type SendToAddressReq struct {
	Address string `json:"address"`
	AmtSats uint64 `json:"amountSats"`
}

type GenerateToAddressReq struct {
	Address     string `json:"address"`
	NumOfBlocks uint64 `json:"numberOfBlocks"`
}
