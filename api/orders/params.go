package orderHandler

type UnsignedOrderRequestParams struct {
	UserId       string `query:"user-id" optional:"true"`       // implied user has an existing account if to have collateral
	Pair         string `query:"pair"`                          // Target perpetual, expects "BTC/USD", "ETH/USD", etc
	Collateral   string `query:"value" optional:"true"`         // Collateral amount in USD
	EntryPrice   string `query:"entry" optional:"true"`         // Entry price in USD
	Slippage     string `query:"slip" optional:"true"`          // Max slippage (basis points, out of 10,000)
	Leverage     string `query:"lev" optional:"true"`           // Leverage multiplier
	PositionType string `query:"position-type" optional:"true"` // "long" or "short"
}

type SignedOrderRequestParams struct {
	OrderId string `query:"order-id"`
	R       string `query:"r"`
	S       string `query:"s"`
	V       string `query:"v"`
}

type ModifyOrderRequestParams struct {
	OrderId   string `query:"order-id"`
	NewStatus string `query:"new-status"`
}

type GetOrdersByUserAddressRequestParams struct {
	WalletAddress string `query:"wallet-address"`
	WalletType    string `query:"wallet-type"`
	OrderType     string `query:"order-type" optional:"true"`   // 'long', 'short'
	OrderStatus   string `query:"order-status" optional:"true"` // 'unsigned', 'pending', 'filled', 'canceled', 'closed', 'liquidated'
}

type GetOrdersByUserIdRequestParams struct {
	UserId      string `query:"user-id"`
	OrderType   string `query:"order-type" optional:"true"` // filter results
	OrderStatus string `query:"order-status" optional:"true"`
}

type GetOrdersByIdRequestParams struct {
	OrderId string `query:"order-id"`
}

type UnsignedCancelOrderRequestParams struct {
	OrderId string `query:"order-id"`
}

type UnsignedCloseOrderRequestParams struct {
	OrderId string `query:"order-id"`
}

// set the ended_at timestamp
// the signature ov
type SignedCloseOrderRequestParams struct {
	OrderId     string `query:"order-id"`
	SignatureId string `query:"signature-id"`
	R           string `query:"r" optional:"true"`
	S           string `query:"s" optional:"true"`
	V           string `query:"v" optional:"true"`
}

type SignedCancelOrderRequestParams struct {
	OrderId     string `query:"order-id"`
	SignatureId string `query:"signature-id"`
	R           string `query:"r" optional:"true"`
	S           string `query:"s" optional:"true"`
	V           string `query:"v" optional:"true"`
}

type CreateOrderRequestParams struct {
	UserId            string `query:"user-id" optional:"true"`    // implied user has an existing account if to have collateral
	Pair              string `query:"pair"`                       // Target perpetual, expects "BTC/USD", "ETH/USD", etc
	Collateral        string `query:"value" optional:"true"`      // Collateral amount in USD
	EntryPrice        string `query:"entry" optional:"true"`      // Entry price in USD
	Slippage          string `query:"slip" optional:"true"`       // Max slippage (basis points, out of 10,000)
	Leverage          string `query:"lev" optional:"true"`        // Leverage multiplier
	PositionType      string `query:"order-type" optional:"true"` // "long" or "short"
	LimitPrice        string `query:"lim-price" optional:"true"`
	StopLossPrice     string `query:"stop-price" optional:"true"`
	TakeProfitPrice   string `query:"tp-price" optional:"true"`
	TakeProfitPercent string `query:"tp-percent" optional:"true"` // percent to close the position for take profit, when achieved the tp_price and tp_value are set to null
}
