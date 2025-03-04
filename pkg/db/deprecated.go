package db

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

func GetOrdersByUserId_old(client *supabase.Client, userId string) (*[]OrderResponse_old, error) {
	params := map[string]interface{}{
		"user_id": userId,
	}
	response := client.Rpc("get_orders_by_userid_deprecated", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	var orders []OrderResponse_old
	if err := json.Unmarshal([]byte(response), &orders); err != nil {
		return nil, fmt.Errorf("error unmarshalling user response: %v", err)
	}

	return &orders, nil
}

func GetOrdersByUserAddress_old(client *supabase.Client, walletAddress, walletType string) (*[]OrderResponse_old, error) {
	fmt.Print("\n inside of GetOrdersByUserAddress")
	params := map[string]interface{}{
		"wallet_addr": walletAddress,
		"wallet_t":    walletType,
	}

	response := client.Rpc("get_orders_by_address_deprecated", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	fmt.Printf("\n response: %v", response)
	var orders []OrderResponse_old
	err := json.Unmarshal([]byte(response), &orders)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}
	fmt.Printf("\n response: %v", &orders)

	return &orders, nil
}

func GetOrderById_old(client *supabase.Client, id string) (*OrderAndUserResponse_old, error) {
	fmt.Printf("\n this is where i really am")
	params := map[string]interface{}{
		"id_": id,
	}
	response := client.Rpc("get_order_by_id_deprecated", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	var order OrderAndUserResponse_old
	if err := json.Unmarshal([]byte(response), &order); err != nil {
		return nil, fmt.Errorf("error unmarshalling user response: %v", err)
	}

	return &order, nil
}

type OrderAndUserResponse_old struct {
	Order OrderResponse_old `json:"order"`
	User  UserResponse      `json:"user"`
}

type OrderResponse_old struct {
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

func CreateOrder_old(client *supabase.Client, userId, orderType string, leverage float64, pair string, collateral, entryPrice, liquidationPrice float64) (*OrderResponse_old, error) {
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
	response := client.Rpc("create_order_deprecated", "exact", params)

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

	var order OrderResponse_old
	err := json.Unmarshal([]byte(response), &order)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	return &order, nil
}

func SignOrder_old(client *supabase.Client, orderId string) (*OrderResponse_old, error) {
	_, err := uuid.Parse(orderId)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID format: %v", err)
	}

	params := map[string]interface{}{
		"order_id": orderId,
	}

	// Execute the RPC call
	response := client.Rpc("sign_order_deprecated", "estimate", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	if response == "" {
		return nil, fmt.Errorf("db error: failed to execute create_order")
	}

	var order OrderResponse_old
	if err := json.Unmarshal([]byte(response), &order); err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	return &order, nil
}
