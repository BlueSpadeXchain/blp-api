package db

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/supabase-community/supabase-go"
)

func CreateWithdrawal(client *supabase.Client, userID string, amount int64, status string, txHash string) error {
	withdrawData := map[string]interface{}{
		"userid":  userID,
		"amount":  amount,
		"status":  status,
		"tx_hash": txHash,
	}

	_, _, err := client.From("withdrawals").Insert(withdrawData, false, "", "minimal", "").Execute()
	if err != nil {
		log.Printf("Failed to create withdrawal: %v", err)
		return err
	}
	return nil
}

func ModifyWithdrawalStatus(client *supabase.Client, withdrawalID string, status string) error {
	updateData := map[string]interface{}{
		"status": status,
	}

	_, _, err := client.From("withdrawals").Update(updateData, "", "").Eq("id", withdrawalID).Execute()
	if err != nil {
		log.Printf("Failed to modify withdrawal status: %v", err)
		return err
	}
	return nil
}

func ModifyUserBalance(client *supabase.Client, userID string, newBalance int64) error {
	updateData := map[string]interface{}{
		"balance": newBalance,
	}

	_, _, err := client.From("users").Update(updateData, "", "").Eq("userid", userID).Execute()
	if err != nil {
		log.Printf("Failed to modify user balance: %v", err)
		return err
	}
	return nil
}

func SignOrder(client *supabase.Client, orderId string) (*OrderResponse, error) {
	params := map[string]interface{}{
		"order_id": orderId,
	}

	// Execute the RPC call
	response := client.Rpc("sign_order", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	if response == "" {
		return nil, fmt.Errorf("db error: failed to execute create_order")
	}

	var order OrderResponse
	err := json.Unmarshal([]byte(response), &order)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	return &order, nil
}

func CreateOrder(client *supabase.Client, userId, orderType string, leverage float64, pair string, collateral, entryPrice, liquidationPrice float64) (*OrderResponse, error) {
	// Convert chainID, block, and depositNonce to string for TEXT type in the database
	params := map[string]interface{}{
		"user_id":     userId,
		"order_type":  orderType,
		"leverage":    leverage,
		"pair":        pair,
		"collateral":  collateral,
		"entry_price": entryPrice,
		"liq_price":   liquidationPrice,
	}

	// Execute the RPC call
	response := client.Rpc("create_order", "exact", params)

	// Check for any Supabase errors
	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	// If no response or an error, return
	if response == "" {
		return nil, fmt.Errorf("db error: failed to execute create_order for user ID %v", userId)
	}

	var order OrderResponse
	err := json.Unmarshal([]byte(response), &order)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	return &order, nil
}

func ModifyOrder(client *supabase.Client, orderID string, updatedData map[string]interface{}) error {
	_, _, err := client.From("orders").Update(updatedData, "", "").Eq("id", orderID).Execute()
	if err != nil {
		log.Printf("Failed to modify order: %v", err)
		return err
	}
	return nil
}

func CloseOrder(client *supabase.Client, orderID string) error {
	updateData := map[string]interface{}{
		"status": "canceled",
	}

	_, _, err := client.From("orders").Update(updateData, "", "").Eq("id", orderID).Execute()
	if err != nil {
		log.Printf("Failed to close order: %v", err)
		return err
	}
	return nil
}

func GetOrCreateUser(client *supabase.Client, walletAddress, walletType string) (*UserResponse, error) {
	params := map[string]interface{}{
		"wallet_addr": walletAddress,
		"wallet_t":    walletType,
	}

	response := client.Rpc("get_or_create_user", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	var users []UserResponse
	err := json.Unmarshal([]byte(response), &users)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("db error: %v of type %v could not create a user", walletAddress, walletType)
	}

	return &users[0], nil
}

func AddUserDeposit(client *supabase.Client, walletAddress, walletType, chainID, block, blockHash, txHash, sender, depositNonce, asset, amount, value string) error {
	// Convert chainID, block, and depositNonce to string for TEXT type in the database
	params := map[string]interface{}{
		"wallet_addr":   walletAddress,
		"wallet_t":      walletType,
		"chain":         chainID,
		"blk":           block,
		"blk_hash":      blockHash,
		"tx_hash":       txHash,
		"sndr":          sender,
		"deposit_nonce": depositNonce,
		"asset_addr":    asset,
		"amt":           amount,
		"val":           value,
	}

	// Execute the RPC call
	response := client.Rpc("add_user_deposit", "exact", params)

	// Check for any Supabase errors
	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	// If no response or an error, return
	if response == "" {
		return fmt.Errorf("db error: failed to execute add_user_deposit for wallet %v", walletAddress)
	}

	return nil
}
