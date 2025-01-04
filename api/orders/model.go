package orderHandler

type OrderRaw struct {
	Signer    string    `json:"signer"`
	CreatedOn string    `json:"createdOn"`
	ChainId   string    `json:"chainId"`
	Order     OrderData `json:"order"`
	MessageId string    `json:"messageId"`
	Signature string    `json:"signature"`
	Nonce     int64     `json:"nonce"`
}

type CreateOrderUnsignedRequestParams struct {
	WalletAddress string `query:"wallet-adress"`
	WalletType    string `query:"wallet-type"`   // ecdsa / eddsa / edd25519 / secp256r1 / etc
	Pair          string `query:"pair"`          // Target perpetual, expects "BTC/USD", "ETH/USD", etc
	Collateral    string `query:"value"`         // Collateral amount in USD
	EntryPrice    string `query:"entry"`         // Entry price in USD
	Slippage      string `query:"slip"`          // Max slippage (basis points, out of 10,000)
	Leverage      string `query:"lev"`           // Leverage multiplier
	PositionType  string `query:"position-type"` // "long" or "short"
	Signature     string `query:"sig"`           // Signature over id . value . entry . slip . lev . type
}

type CreateOrderSignedRequestParams struct {
	Wallet       string `query:"wallet-adress"`
	WalletType   string `query:"wallet-type"`   // ecdsa / eddsa / edd25519 / secp256r1 / etc
	PerpId       string `query:"perp"`          // Target perpetual
	Collateral   string `query:"value"`         // Collateral amount in USD
	EntryPrice   string `query:"entry"`         // Entry price in USD
	Slippage     string `query:"slip"`          // Max slippage (basis points, out of 10,000)
	Leverage     string `query:"lev"`           // Leverage multiplier
	PositionType string `query:"position-type"` // "long" or "short"
	Signature    string `query:"sig"`           // Signature over id . value . entry . slip . lev . type
}

type OrderCloseParams struct {
	Signer     string `query:"addr"`
	OrderId    string `query:"order"` // from db
	PerpId     string `query:"perp"`  // target of perp
	Percentage string `query:"percent"`
	Signature  string `query:"sig"` // sign over order . id
}

type OrderRequestResponse struct {
	OrderId string `json:"order-id"`
}

// actually we don't need to store the hash or signature on the order, only for message validity

type OrderData struct {
	OrderId          string `json:"orderId"`
	NetValue         string `json:"netValue"`
	Amount           string `json:"amount"`
	Collateral       string `json:"collateral"`
	MarkPrice        string `json:"markPrice"`
	EntryPrice       string `json:"EntryPrice"`
	LiquidationPrice string `json:"liquidationPrice"`
}

type OrderResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Data    *OrderData `json:"data,omitempty"` // this will assign a new orderId if not already applied
}
