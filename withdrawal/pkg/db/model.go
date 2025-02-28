package db

import (
	"fmt"
	"strings"
	"time"
)

type SupabaseError struct {
	Code    string `json:"code"`
	Details string `json:"details"`
	Hint    string `json:"hint"`
	Message string `json:"message"`
}

type Error struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
	Origin  string `json:"origin"`
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
