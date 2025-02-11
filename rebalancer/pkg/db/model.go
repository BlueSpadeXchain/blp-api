package db

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SupabaseError struct {
	Code    string `json:"code"`
	Details string `json:"details"`
	Hint    string `json:"hint"`
	Message string `json:"message"`
}

type OrderResponse struct {
	ID                   uuid.UUID `json:"id"`
	UserID               string    `json:"userid"`
	OrderType            string    `json:"order_type"`
	Leverage             float64   `json:"leverage"`
	Pair                 string    `json:"pair"`
	PairID               string    `json:"pair_id"`
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
	TakeProfitAt         time.Time `json:"tp_at"`
	EndedAt              time.Time `json:"ended_at"`
	ProfitAndLoss        float64   `json:"pnl"`
	OpenFee              float64   `json:"open_fee"`
	CloseFee             float64   `json:"close_fee"`
}

type OrderGlobalUpdate struct {
	CurrentBorrowed       float64 `json:"current_borrowed"`
	CurrentLiquidity      float64 `json:"current_liquidity"`
	CurrentOrdersActive   float64 `json:"current_orders_active"`
	CurrentOrdersLimit    float64 `json:"current_orders_limit"`
	CurrentOrdersPending  float64 `json:"current_orders_pending"`
	TotalBorrowed         float64 `json:"total_borrowed"`
	TotalLiquidations     float64 `json:"total_liquidations"`
	TotalOrdersActive     float64 `json:"total_orders_active"`
	TotalOrdersFilled     float64 `json:"total_orders_filled"`
	TotalOrdersLimit      float64 `json:"total_orders_limit"`
	TotalOrdersLiquidated float64 `json:"total_orders_liquidated"`
	TotalOrdersStopped    float64 `json:"total_orders_stopped"`
	TotalPnlLosses        float64 `json:"total_pnl_losses"`
	TotalPnlProfits       float64 `json:"total_pnl_profits"`
	TotalRevenue          float64 `json:"total_revenue"`
	TreasuryBalance       float64 `json:"treasury_balance"`
	TotalTreasuryProfits  float64 `json:"total_treasury_profits"`
	VaultBalance          float64 `json:"vault_balance"`
	TotalVaultProfits     float64 `json:"total_vault_profits"`
	TotalLiquidityRewards float64 `json:"total_liquidity_rewards"`
	TotalStakeRewards     float64 `json:"total_stake_rewards"`
}

// OrderUpdate represents the PostgreSQL order_update type
type OrderUpdate struct {
	OrderID             uuid.UUID         `json:"order_id"`
	UserID              string            `json:"userid"`
	Status              string            `json:"status"`
	EntryPrice          float64           `json:"entry_price"`
	ClosePrice          float64           `json:"close_price"`
	TpValue             float64           `json:"tp_value"`
	Pnl                 float64           `json:"pnl"`
	Collateral          float64           `json:"collateral"`
	TakeProfitAt        time.Time         `json:"tp_at"`
	BalanceChange       float64           `json:"balance_change"`
	EscrowBalanceChange float64           `json:"escrow_balance_change"`
	OrderGlobalUpdate   OrderGlobalUpdate `json:"order_global_update"`
}

func (ou OrderUpdate) Value() (driver.Value, error) {
	return json.Marshal(ou)
}

// Scan implements the sql.Scanner interface for OrderUpdate
func (ou *OrderUpdate) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, &ou)
}

type GlobalStateResponse struct {
	Key       string  `json:"key"`
	Value     float64 `json:"value"`
	UpdatedAt string  `json:"updated_at"`
}
