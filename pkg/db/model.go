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

type OrderAndUserResponse struct {
	Order OrderResponse `json:"order"`
	User  UserResponse  `json:"user"`
}

type OrderAndUserResponse2 struct {
	Order OrderResponse2 `json:"order"`
	User  UserResponse   `json:"user"`
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

type UnsignedCreateOrderResponse struct {
	Order         OrderResponse2 `json:"order"`
	SignatureId   string         `json:"signature_id"`
	SignatureHash string         `json:"signature_hash"`
	ExpiryTime    string         `json:"expiry_time"`
}

type UnsignedCloseOrderResponse struct {
	OrderId       string `json:"order_id"`
	SignatureId   string `json:"signature_id"`
	SignatureHash string `json:"signature_hash"`
	ExpiryTime    string `json:"expiry_time"`
}

type UnsignedCancelOrderResponse struct {
	OrderId       string `json:"order_id"`
	SignatureId   string `json:"signature_id"`
	SignatureHash string `json:"signature_hash"`
	ExpiryTime    string `json:"expiry_time"`
}

type SignedCancelOrderResponse struct {
	Order        OrderResponse2 `json:"order"`
	IsValid      bool           `json:"is_valid"`
	ErrorMessage string         `json:"error_message"`
}

type SignedCloseOrderResponse struct {
	Order        OrderResponse2 `json:"order"`
	IsValid      bool           `json:"is_valid"`
	ErrorMessage string         `json:"error_message"`
}

type GetSignatureValidationHashResponse struct {
	Hash string `json:"signature_hash"`
}

type GlobalStateResponse struct {
	Key       string  `json:"key"`
	Value     float64 `json:"value"`
	UpdatedAt string  `json:"updated_at"`
}

type GetSignatureHashResponse struct {
	Hash string `json:"signature_hash"`
}
