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
	ID                   string     `json:"id"`
	UserID               string     `json:"userid"`
	OrderType            string     `json:"order_type"`
	Leverage             float64    `json:"leverage"`
	Pair                 string     `json:"pair"`
	PairId               string     `json:"pair_id"`
	OrderStatus          string     `json:"status"`
	Collateral           float64    `json:"collateral"`
	EntryPrice           float64    `json:"entry_price"`
	ClosePrice           float64    `json:"close_price"`
	LiquidationPrice     float64    `json:"liq_price"`
	MaxPrice             float64    `json:"max_price"`
	MaxValue             float64    `json:"max_value"`
	LimitPrice           float64    `json:"limit_price"`
	StopLossPrice        float64    `json:"stop_price"`
	TakeProfitPrice      float64    `json:"tp_price"`
	TakeProfitValue      float64    `json:"tp_value"`
	TakeProfitCollateral float64    `json:"tp_collateral"`
	CreatedAt            CustomTime `json:"created_at"`
	SignedAt             CustomTime `json:"signed_at"`
	StartedAt            CustomTime `json:"started_at"`
	ModifiedAt           CustomTime `json:"modified_at"`
	EndedAt              CustomTime `json:"ended_at"`
	TakeProfitAt         CustomTime `json:"tp_at"`
	ProfitAndLoss        float64    `json:"pnl"`
	OpenFee              float64    `json:"open_fee"`
	CloseFee             float64    `json:"close_fee"`
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
	ID            string     `json:"id"`
	UserID        string     `json:"userid"`
	Amount        float64    `json:"amount"`
	TokenType     string     `json:"token_type"`
	Status        string     `json:"status"`
	CreatedAt     CustomTime `json:"created_at"`
	UpdatedAt     CustomTime `json:"updated_at"`
	TxHash        string     `json:"tx_hash"`
	WalletAddress string     `json:"wallet_address"`
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

type MetricsResponse struct {
	TotalRevenue          float64 `json:"total_revenue"`
	VaultBalance          float64 `json:"vault_balance"`
	TotalBorrowed         float64 `json:"total_borrowed"`
	TotalDeposits         float64 `json:"total_deposits"`
	CurrentBalance        float64 `json:"current_balance"`
	TotalLiquidity        float64 `json:"total_liquidity"`
	TotalWithdrawn        float64 `json:"total_witdhrawn"`
	CurrentBorrowed       float64 `json:"current_borrowed"`
	OpenFeePercent        float64 `json:"open_fee_percent"`
	TotalPnlLosses        float64 `json:"total_pnl_losses"`
	TreasuryBalance       float64 `json:"treasury_balance"`
	CloseFeePercent       float64 `json:"close_fee_percent"`
	CurrentLiquidity      float64 `json:"current_liquidity"`
	FeePercentStake       float64 `json:"fee_percent_stake"`
	FeePercentVault       float64 `json:"fee_percent_vault"`
	TotalBlpRewards       float64 `json:"total_blp_rewards"`
	TotalBluRewards       float64 `json:"total_blu_rewards"`
	TotalLiquidations     float64 `json:"total_liquidations"`
	TotalPnlProfits       float64 `json:"total_pnl_profits"`
	CurrentBlpStaked      float64 `json:"current_blp_staked"`
	CurrentBluStaked      float64 `json:"current_blu_staked"`
	TotalBlpUnstaked      float64 `json:"total_blp_unstaked"`
	TotalBluUnstaked      float64 `json:"total_blu_unstaked"`
	TotalOrdersLimit      float64 `json:"total_orders_limit"`
	CurrentBlpPending     float64 `json:"current_blp_pending"`
	CurrentBluPending     float64 `json:"current_blu_pending"`
	CurrentBlpRewards     float64 `json:"current_blp_rewards"`
	CurrentBluRewards     float64 `json:"current_blu_rewards"`
	TotalOrdersActive     float64 `json:"total_orders_active"`
	TotalOrdersClosed     float64 `json:"total_orders_closed"`
	TotalOrdersFilled     float64 `json:"total_orders_filled"`
	TotalOrdersSigned     float64 `json:"total_orders_signed"`
	TotalVaultProfits     float64 `json:"total_vault_profits"`
	CurrentOrdersLimit    float64 `json:"current_orders_limit"`
	FeePercentTreasury    float64 `json:"fee_percent_treasury"`
	TotalOrdersCreated    float64 `json:"total_orders_created"`
	TotalOrdersStopped    float64 `json:"total_orders_stopped"`
	CurrentOrdersActive   float64 `json:"current_order_active"`
	FeePercentLiquidity   float64 `json:"fee_percent_liquidity"`
	TotalOrdersCanceled   float64 `json:"total_orders_canceled"`
	CurrentOrdersPending  float64 `json:"current_orders_pending"`
	TotalTreasuryProfits  float64 `json:"total_treasury_profits"`
	TotalOrdersLiquidated float64 `json:"total_orders_liquidated"`
}

type GetLatestMetricSnapshotResponse struct {
	SnapshotTime CustomTime      `json:"snapshot_time"`
	Metrics      MetricsResponse `json:"metrics"`
}
