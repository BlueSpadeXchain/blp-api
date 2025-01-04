package db

type UserResponse struct {
	ID            string `json:"id"`
	UserID        string `json:"userid"`
	WalletAddress string `json:"wallet_address"`
	WalletType    string `json:"wallet_type"`
	Balance       int64  `json:"balance"`
	PerpBalance   int64  `json:"perp_balance"`
	EscrowBalance int64  `json:"escrow_balance"`
	StakeBalance  int64  `json:"stake_balance"`
	FrozenBalance int64  `json:"frozen_balance"`
	TotalBalance  int64  `json:"total_balance"`
	CreatedAt     string `json:"created_at"`
}

type OrderResponse struct {
	ID         string  `json:"id"`
	UserID     string  `json:"userid"`
	OrderType  string  `json:"order_type"`
	Leverage   float64 `json:"leverage"`
	Pair       string  `json:"pair"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	EntryPrice float64 `json:"entry_price"`
	MarkPrice  float64 `json:"mark_price"`
	LiqPrice   float64 `json:"liq_price"`
	CreatedAt  string  `json:"created_at"`
	EndedAt    string  `json:"ended_at"`
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
