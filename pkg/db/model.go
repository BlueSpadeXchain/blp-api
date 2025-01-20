package db

type UserResponse struct {
	ID            string  `json:"id"`
	UserID        string  `json:"userid"`
	WalletAddress string  `json:"wallet_address"`
	WalletType    string  `json:"wallet_type"`
	Balance       float64 `json:"balance"`
	PerpBalance   float64 `json:"perp_balance"`
	EscrowBalance float64 `json:"escrow_balance"`
	StakeBalance  float64 `json:"stake_balance"`
	FrozenBalance float64 `json:"frozen_balance"`
	TotalBalance  float64 `json:"total_balance"`
	CreatedAt     string  `json:"created_at"`
}

type OrderResponse struct {
	ID         string  `json:"id"`
	UserID     string  `json:"userid"`
	OrderType  string  `json:"order_type"`
	Leverage   float64 `json:"leverage"`
	PairId     string  `json:"pair"`
	Status     string  `json:"status"`
	EntryPrice float64 `json:"entry_price"`
	LiqPrice   float64 `json:"liq_price"`
	CreatedAt  string  `json:"created_at"`
	EndedAt    string  `json:"ended_at"`
	Collateral float64 `json:"collateral"`
}

type OrderResponse2 struct {
	ID              string  `json:"id"`
	UserID          string  `json:"userid"`
	OrderType       string  `json:"order_type"`
	Leverage        float64 `json:"leverage"`
	PairId          string  `json:"pair"`
	Status          string  `json:"status"`
	EntryPrice      float64 `json:"entry_price"`
	LiqPrice        float64 `json:"liq_price"`
	LimitPrice      float64 `json:"limit_price"`
	StopLossPrice   float64 `json:"stop_price"`
	TakeProfitPrice float64 `json:"tp_price"`
	MaxPrice        float64 `json:"max_price"`
	CreatedAt       string  `json:"created_at"`
	EndedAt         string  `json:"ended_at"`
	Collateral      float64 `json:"collateral"`
}

type OrderAndUserResponse struct {
	Order OrderResponse `json:"order"`
	User  UserResponse  `json:"user"`
}

type DepositResponse struct {
	ID            string  `json:"id"`
	UserID        string  `json:"userid"`
	WalletAddress string  `json:"wallet_address"`
	WalletType    string  `json:"wallet_type"`
	ChainID       string  `json:"chain_id"`
	Block         string  `json:"block"`
	BlockHash     string  `json:"block_hash"`
	TxHash        string  `json:"tx_hash"`
	Sender        string  `json:"sender"`
	DepositNonce  string  `json:"deposit_nonce"`
	Asset         string  `json:"asset"`
	Amount        string  `json:"amount"`
	Value         float64 `json:"value"`
	CreatedAt     string  `json:"created_at"`
}

type SupabaseError struct {
	Code    string `json:"code"`
	Details string `json:"details"`
	Hint    string `json:"hint"`
	Message string `json:"message"`
}
