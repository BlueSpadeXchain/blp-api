package db

import (
	"fmt"
	"strings"
	"time"
)

type UserResponse struct {
	ID              string  `json:"id"`
	UserID          string  `json:"userid"`
	WalletAddress   string  `json:"wallet_address"`
	WalletType      string  `json:"wallet_type"`
	Balance         float64 `json:"balance"`
	PerpBalance     float64 `json:"perp_balance"`
	EscrowBalance   float64 `json:"escrow_balance"`
	FrozenBalance   float64 `json:"frozen_balance"`
	BluStakeBalance float64 `json:"blu_stake_balance"`
	BlpStakeBalance float64 `json:"blp_stake_balance"`
	BluStakePending float64 `json:"blu_stake_pending"`
	BlpStakePending float64 `json:"blp_stake_pending"`
	TotalBalance    float64 `json:"total_balance"`
	CreatedAt       string  `json:"created_at"`
}

type OrderResponse struct {
	ID                   string    `json:"id"`
	UserID               string    `json:"userid"`
	OrderType            string    `json:"order_type"`
	Leverage             float64   `json:"leverage"`
	Pair                 string    `json:"pair"`
	PairId               string    `json:"pair_id"`
	OrderStatus          string    `json:"status"`
	Collateral           float64   `json:"collateral"`
	EntryPrice           float64   `json:"entry_price"`
	ClosePrice           float64   `json:"close_price"`
	LiquidationPrice     float64   `json:"liq_price"`
	MaxPrice             float64   `json:"max_price"`
	MaxValue             float64   `json:"max_value"`
	LimitPrice           float64   `json:"limit_price"`
	StopLossPrice        float64   `json:"stop_price"`
	TakeProfitPrice      float64   `json:"tp_price"`
	TakeProfitValue      float64   `json:"tp_value"`
	TakeProfitCollateral float64   `json:"tp_collateral"`
	CreatedAt            time.Time `json:"created_at"`
	SignedAt             time.Time `json:"signed_at"`
	StartedAt            time.Time `json:"started_at"`
	ModifiedAt           time.Time `json:"modified_at"`
	EndedAt              time.Time `json:"ended_at"`
	TakeProfitAt         time.Time `json:"tp_at"`
	ProfitAndLoss        float64   `json:"pnl"`
	OpenFee              float64   `json:"open_fee"`
	CloseFee             float64   `json:"close_fee"`
}

type StakeResponse struct {
	ID        string  `json:"id"`
	UserID    string  `json:"userid"`
	StakeType string  `json:"stake_type"`
	Amount    float64 `json:"amount"`
	CreatedAt string  `json:"created_at"`
}

type StakesAndUserResponse struct {
	Stakes []StakeResponse `json:"stakes"`
	User   UserResponse    `json:"user"`
}

type OrderAndUserResponse struct {
	Order OrderResponse `json:"order"`
	User  UserResponse  `json:"user"`
}

type WithdrawalAndUserResponse struct {
	Withdrawal WithdrawalResponse `json:"pending_withdrawal"`
	User       UserResponse       `json:"user"`
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
	Order         OrderResponse `json:"order"`
	SignatureId   string        `json:"signature_id"`
	SignatureHash string        `json:"signature_hash"`
	ExpiryTime    string        `json:"expiry_time"`
}

type SignOrderResponse struct {
	Order OrderResponse `json:"order"`
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
	Order        OrderResponse `json:"order"`
	IsValid      bool          `json:"is_valid"`
	ErrorMessage string        `json:"error_message"`
}

type SignedCloseOrderResponse struct {
	Order        OrderResponse `json:"order"`
	IsValid      bool          `json:"is_valid"`
	ErrorMessage string        `json:"error_message"`
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

type UnsignedWithdrawalResponse struct {
	WithdrawalId  string `json:"pending_withdrawal_id"`
	SignatureId   string `json:"signature_id"`
	SignatureHash string `json:"signature_hash"`
	ExpiryTime    string `json:"expiry_time"`
}

type SignedWithdrawalResponse struct {
	Withdrawal   WithdrawalResponse `json:"pending_withdrawal"`
	IsValid      bool               `json:"is_valid"`
	ErrorMessage string             `json:"error_message"`
}

type WithdrawalResponse struct {
	ID            string    `json:"id"`
	UserID        string    `json:"userid"`
	Amount        float64   `json:"amount"`
	TokenType     string    `json:"token_type"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	TxHash        string    `json:"tx_hash"`
	WalletAddress string    `json:"wallet_address"`
}

type StakeDepositResponse struct {
	ID        string     `json:"id"`
	Userid    string     `json:"userid"`
	StakeType string     `json:"stake_type"`
	Amount    float64    `json:"amount"`
	CreatedAt CustomTime `json:"created_at"`
}

type PendingWithdrawalResponse struct {
	ID            string     `json:"id"`
	Userid        string     `json:"userid"`
	Amount        float64    `json:"amount"`
	TokenType     string     `json:"token_type"`
	Status        string     `json:"status"`
	CreatedAt     CustomTime `json:"created_at"`
	UpdatedAt     CustomTime `json:"updated_at"`
	TxHash        string     `json:"tx_hash"`
	WalletAddress string     `json:"wallet_address"`
}

type ProcessUnstakeResponse struct {
	StakeDeposit      StakeDepositResponse      `json:"stake_deposit"`
	PendingWithdrawal PendingWithdrawalResponse `json:"pending_withdrawal"`
}

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)
	if str == "null" || str == "" {
		return nil
	}

	// Define possible time formats
	formats := []string{
		time.RFC3339Nano,             // "2025-02-18T11:14:37.858677Z"
		"2006-01-02T15:04:05.999999", // "2025-02-18T11:14:37.858677" (no timezone)
	}

	var err error
	for _, layout := range formats {
		var t time.Time
		t, err = time.Parse(layout, str)
		if err == nil {
			ct.Time = t
			return nil
		}
	}

	return fmt.Errorf("error parsing time: %v", err)
}
