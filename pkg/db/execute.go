package db

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
	"github.com/google/uuid"
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

func SignOrder(client *supabase.Client, orderId string) (*SignOrderResponse, error) {
	_, err := uuid.Parse(orderId)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID format: %v", err)
	}

	params := map[string]interface{}{
		"order_id": orderId,
	}

	utils.LogInfo("sign_order params", utils.StringifyStructFields(params, ""))

	response := client.Rpc("sign_order", "estimate", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	if response == "" {
		return nil, fmt.Errorf("db error: failed to execute create_order")
	}

	var order SignOrderResponse
	if err := json.Unmarshal([]byte(response), &order); err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	return &order, nil
}

func CreateOrder(
	client *supabase.Client,
	userId, orderType, pair, pairId string,
	leverage, collateral, entryPrice, liquidationPrice, maxPrice, limitPrice, stopLossPrice, takeProfitPrice, takeProfitValue, takeProfitCollateral, openFee float64) (*UnsignedCreateOrderResponse, error) {
	// Convert chainID, block, and depositNonce to string for TEXT type in the database
	params := map[string]interface{}{
		"user_id":     userId,
		"order_type":  orderType,
		"leverage":    leverage,
		"pair":        pair,
		"pair_id":     pairId,
		"collateral":  collateral,
		"entry_price": entryPrice,
		"liq_price":   liquidationPrice,
		"max_price":   maxPrice,
		"open_fee":    openFee,
	}

	if limitPrice != 0 {
		params["lim_price"] = limitPrice
	}

	if stopLossPrice != 0 {
		params["stop_price"] = stopLossPrice
	}

	if takeProfitPrice != 0 && takeProfitValue != 0 && takeProfitCollateral != 0 {
		params["tp_price"] = takeProfitPrice
		params["tp_value"] = takeProfitValue
		params["tp_collateral"] = takeProfitCollateral
	}

	utils.LogInfo("create_order params", utils.StringifyStructFields(params, ""))

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

	var order UnsignedCreateOrderResponse
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

func CloseOrder(client *supabase.Client, orderID string) (*UnsignedCloseOrderResponse, error) {
	params := map[string]interface{}{
		"order_id": orderID,
	}

	utils.LogInfo("unsigned_close_order params", utils.StringifyStructFields(params, ""))

	// Execute the RPC call
	response := client.Rpc("unsigned_close_order", "exact", params)

	// Check for any Supabase errors
	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	// If no response or an error, return
	if response == "" {
		return nil, fmt.Errorf("db error: failed to execute close_order for order ID %v", orderID)
	}

	var order UnsignedCloseOrderResponse
	err := json.Unmarshal([]byte(response), &order)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	return &order, nil
}

func SignCloseOrder(client *supabase.Client, orderId, signatureId string, remainingCollateral, payoutValue, closeFee, closePrice float64) (*SignedCloseOrderResponse, error) {
	params := map[string]interface{}{
		"order_id":             orderId,
		"signature_id":         signatureId,
		"remaining_collateral": remainingCollateral,
		"payout_value":         payoutValue,
		"close_fee_":           closeFee,
		"close_price_":         closePrice,
	}

	utils.LogInfo("signed_close_order params", utils.StringifyStructFields(params, ""))

	// Execute the RPC call
	response := client.Rpc("signed_close_order", "exact", params)

	// Check for any Supabase errors
	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	// If no response or an error, return
	if response == "" {
		return nil, fmt.Errorf("db error: failed to execute close_order for order ID %v", orderId)
	}

	var order SignedCloseOrderResponse
	err := json.Unmarshal([]byte(response), &order)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	return &order, nil
}

func CancelOrder(client *supabase.Client, orderID string) (*UnsignedCancelOrderResponse, error) {
	params := map[string]interface{}{
		"order_id": orderID,
	}

	utils.LogInfo("unsigned_cancel_order params", utils.StringifyStructFields(params, ""))

	response := client.Rpc("unsigned_cancel_order", "exact", params)

	// Check for any Supabase errors
	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	// If no response or an error, return
	if response == "" {
		return nil, fmt.Errorf("db error: failed to execute cancel_order for order ID %v", orderID)
	}

	var order UnsignedCancelOrderResponse
	err := json.Unmarshal([]byte(response), &order)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	return &order, nil
}

func SignCancelOrder(client *supabase.Client, orderId, signatureId string) (*SignedCancelOrderResponse, error) {
	params := map[string]interface{}{
		"order_id":     orderId,
		"signature_id": signatureId,
	}

	utils.LogInfo("signed_cancel_order params", utils.StringifyStructFields(params, ""))

	response := client.Rpc("signed_cancel_order", "exact", params)

	// Check for any Supabase errors
	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	// If no response or an error, return
	if response == "" {
		return nil, fmt.Errorf("db error: failed to execute cancel_order for order ID %v", orderId)
	}

	var order SignedCancelOrderResponse
	err := json.Unmarshal([]byte(response), &order)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	return &order, nil
}

func GetOrCreateUser(client *supabase.Client, walletAddress, walletType string) (*UserResponse, error) {
	params := map[string]interface{}{
		"wallet_addr": walletAddress,
		"wallet_t":    walletType,
	}

	utils.LogInfo("get_or_create_user params", utils.StringifyStructFields(params, ""))

	response := client.Rpc("get_or_create_user", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	var users UserResponse
	err := json.Unmarshal([]byte(response), &users)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	// if len(users) == 0 {
	// 	return nil, fmt.Errorf("db error: %v of type %v could not create a user", walletAddress, walletType)
	// }

	return &users, nil
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

	utils.LogInfo("add_user_deposit params", utils.StringifyStructFields(params, ""))

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

func ProcessDepositAndStake(client *supabase.Client, walletAddress, walletType, chainID, block, blockHash, txHash, sender, depositNonce, asset, amount, value, stakeType string) error {
	// Convert chainID, block, and depositNonce to string for TEXT type in the database
	params := map[string]interface{}{
		"wallet_addr":      walletAddress,
		"wallet_t":         walletType,
		"chain":            chainID,
		"blk":              block,
		"blk_hash":         blockHash,
		"tx_hash":          txHash,
		"sndr":             sender,
		"deposit_nonce":    depositNonce,
		"asset_addr":       asset,
		"amt":              amount,
		"val":              value,
		"stake_type_param": stakeType,
	}

	utils.LogInfo("process_deposit_and_stake params", utils.StringifyStructFields(params, ""))

	response := client.Rpc("process_deposit_and_stake", "exact", params)

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

// CREATE OR REPLACE FUNCTION process_deposit_and_stake(
// 	-- Deposit parameters
// 	wallet_addr VARCHAR,
// 	wallet_t VARCHAR,
// 	chain TEXT,
// 	blk TEXT,
// 	blk_hash VARCHAR,
// 	tx_hash VARCHAR,
// 	sndr VARCHAR,
// 	deposit_nonce TEXT,
// 	asset_addr VARCHAR,
// 	amt TEXT,
// 	val NUMERIC(78, 9),
// 	-- Stake parameters
// 	stake_type_param stake_type
// )

// withdraw function and signed withdraw function the generation of teh signing needs to be over
// the amount + withdraw id
// and validate the user/signer balance
// also affect the user frozen balance
func Withdraw(client *supabase.Client, withdrawId string, amount float64) (*UnsignedWithdrawResponse, error) {
	params := map[string]interface{}{
		"withdraw_id": withdrawId,
		"amount":      amount,
	}

	utils.LogInfo("unsigned_withdraw params", utils.StringifyStructFields(params, ""))

	// Execute the RPC call
	response := client.Rpc("unsigned_withdraw", "exact", params)

	// Check for any Supabase errors
	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	// If no response or an error, return
	if response == "" {
		return nil, fmt.Errorf("db error: failed to execute unsigned_wthdraw for withdraw ID %v", withdrawId)
	}

	var withdrawl UnsignedWithdrawResponse
	err := json.Unmarshal([]byte(response), &withdrawl)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	return &withdrawl, nil
}

func SignWithdraw(client *supabase.Client, withdrawId, signatureId string) (*SignedWithdrawResponse, error) {
	params := map[string]interface{}{
		"withdraw_id":  withdrawId,
		"signature_id": signatureId,
	}

	utils.LogInfo("signed_withdraw params", utils.StringifyStructFields(params, ""))

	// Execute the RPC call
	response := client.Rpc("signed_withdraw", "exact", params)

	// Check for any Supabase errors
	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	// If no response or an error, return
	if response == "" {
		return nil, fmt.Errorf("db error: failed to execute signed_withdraw for withdraw ID %v", withdrawId)
	}

	var withdrawl SignedWithdrawResponse
	err := json.Unmarshal([]byte(response), &withdrawl)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	return &withdrawl, nil
}

func Unstake(client *supabase.Client, userId, stakeType string, amount float64) (*ProcessUnstakeResponse, error) {
	params := map[string]interface{}{
		"user_id":    userId,
		"stake_type": stakeType,
		"amount":     amount,
	}

	utils.LogInfo("process_unstake_deposit params", utils.StringifyStructFields(params, ""))

	// Execute the RPC call
	response := client.Rpc("process_unstake_deposit", "exact", params)

	// Check for any Supabase errors
	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	// If no response or an error, return
	if response == "" {
		return nil, fmt.Errorf("db error: failed to execute process_unstake_deposit for user ID %v", userId)
	}

	var unstake ProcessUnstakeResponse
	err := json.Unmarshal([]byte(response), &unstake)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	return &unstake, nil
}

func UpdateWithdrawlStatus(client *supabase.Client, withdrawlId, status, txHash string) (*ProcessUnstakeResponse, error) {
	params := map[string]interface{}{
		"p_withdrawl_id": withdrawlId,
		"p_status":       status,
		"p_tx_hash":      txHash,
	}

	utils.LogInfo("update_withdrawal_status params", utils.StringifyStructFields(params, ""))

	// Execute the RPC call
	response := client.Rpc("update_withdrawal_status", "exact", params)

	// Check for any Supabase errors
	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	// If no response or an error, return
	if response == "" {
		return nil, fmt.Errorf("db error: failed to execute update_withdrawal_status for withdraw ID %v", withdrawlId)
	}

	var unstake ProcessUnstakeResponse
	err := json.Unmarshal([]byte(response), &unstake)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	return &unstake, nil
}
