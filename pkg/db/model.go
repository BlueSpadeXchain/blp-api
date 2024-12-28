package db

type User struct {
	ID            string `json:"id"`
	UserID        string `json:"userid"`
	WalletAddress string `json:"wallet_address"`
	WalletType    string `json:"wallet_type"`
	Balance       int64  `json:"balance"`
	PerpBalance   int64  `json:"perp_balance"`
	EscrowBalance int64  `json:"escrow_balance"`
	StakeBalance  int64  `json:"stake_balance"`
	FrozenBalance int64  `json:"frozen_balance"`
	CreatedAt     string `json:"created_at"`
}

type SupabaseError struct {
	Code    string `json:"code"`
	Details string `json:"details"`
	Hint    string `json:"hint"`
	Message string `json:"message"`
}
