package db

type SupabaseError struct {
	Code    string `json:"code"`
	Details string `json:"details"`
	Hint    string `json:"hint"`
	Message string `json:"message"`
}

type OrderResponse2 struct {
	// Pair                  string  `json:"pair"`
	ID                   string  `json:"id"`
	UserID               string  `json:"userid"`
	OrderType            string  `json:"order_type"`
	Leverage             float64 `json:"leverage"`
	PairId               string  `json:"pair_id"`
	OrderStatus          string  `json:"status"`
	Collateral           float64 `json:"collateral"`
	EntryPrice           float64 `json:"entry_price"`
	LiquidationPrice     float64 `json:"liq_price"`
	MaxPrice             float64 `json:"max_price"`
	MaxValue             float64 `json:"max_value"`
	LimitPrice           float64 `json:"limit_price"`
	StopLossPrice        float64 `json:"stop_price"`
	TakeProfitPrice      float64 `json:"tp_price"`
	TakeProfitValue      float64 `json:"tp_value"`
	TakeProfitCollateral float64 `json:"tp_collateral"`
	CreatedAt            string  `json:"created_at"`
	SignedAt             string  `json:"signed_at"`
	StartedAt            string  `json:"started_at"`
	EndedAt              string  `json:"ended_at"`
}
