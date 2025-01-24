package orderHandler

import "github.com/BlueSpadeXchain/blp-api/pkg/db"

type OrderRaw struct {
	Signer    string    `json:"signer"`
	CreatedOn string    `json:"createdOn"`
	ChainId   string    `json:"chainId"`
	Order     OrderData `json:"order"`
	MessageId string    `json:"messageId"`
	Signature string    `json:"signature"`
	Nonce     int64     `json:"nonce"`
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

type UnsignedOrderRequestResponse struct {
	Order db.OrderResponse `json:"order"` // created unsigned position, so it has no affect on balances
	Hash  string           `json:"hash"`  // Hash in hex to be signed by the user
}

type UnsignedOrderRequestResponse2 struct {
	Order db.OrderResponse2 `json:"order"` // created unsigned position, so it has no affect on balances
	Hash  string            `json:"hash"`  // Hash in hex to be signed by the user
}

type BinaryData struct {
	Encoding string   `json:"encoding"`
	Data     []string `json:"data"`
}

type PriceUpdate struct {
	EmaPrice Price  `json:"ema_price"`
	ID       string `json:"id"`
	Metadata Meta   `json:"metadata"`
	Price    Price  `json:"price"`
}

type Price struct {
	Conf        string `json:"conf"`
	Expo        int    `json:"expo"`
	Price       string `json:"price"`
	PublishTime int64  `json:"publish_time"`
}

type Meta struct {
	PrevPublishTime    int64 `json:"prev_publish_time"`
	ProofAvailableTime int64 `json:"proof_available_time"`
	Slot               int64 `json:"slot"`
}

type Response struct {
	Binary BinaryData    `json:"binary"`
	Parsed []PriceUpdate `json:"parsed"`
}
