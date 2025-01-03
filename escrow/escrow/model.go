package escrow

type DespositRequestParams struct {
	ChainId      string `query:"chain-id"`
	Block        string `query:"block"`
	BlockHash    string `query:"block-hash"`
	TxHash       string `query:"tx-hash"`
	Sender       string `query:"sender"`
	Receiver     string `query:"receiver"`
	DepositNonce string `query:"nonce"`
	Asset        string `query:"asset"`
	Amount       string `query:"amount"`
	Signature    string `query:"signature"`
}

type StakeRequestParams struct {
	ChainId      string `query:"chain-id"`
	Block        string `query:"block"`
	BlockHash    string `query:"block-hash"`
	TxHash       string `query:"tx-hash"`
	Sender       string `query:"sender"`
	Receiver     string `query:"receiver"`
	DepositNonce string `query:"nonce"`
	Asset        string `query:"asset"`
	Amount       string `query:"amount"`
	Signature    string `query:"signature"`
}
