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
	PubKey   string `json:"pubKey"`
	LocalAmt int64  `json:"localAmount"`
}

type OpenChannelRes struct {
	FundingTx   string `json:"fundingTx"`
	OutputIndex uint32 `json:"outputIndex"`
}
