package db

import (
	"encoding/json"
	"fmt"

	"github.com/supabase-community/supabase-go"
)

func GetUserByUserId(client *supabase.Client, userId string) (*UserResponse, error) {
	fmt.Printf("\n inside of getuser")
	fmt.Printf("\n inside of getuser: %v", userId)
	params := map[string]interface{}{
		"user_id": userId,
	}
	response := client.Rpc("get_user_by_userid", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}
	fmt.Printf("\n inside of getuser: %v", response)

	var users []UserResponse
	if err := json.Unmarshal([]byte(response), &users); err != nil {
		return nil, fmt.Errorf("error unmarshalling user response: %v", err)
	}
	fmt.Printf("\n inside of getuser: %v", users)

	if len(users) == 0 {
		return nil, fmt.Errorf("no user found for userId: %s", userId)
	}

	return &users[0], nil
}

func GetDepositsByUserId(client *supabase.Client, userId string) (*[]DepositResponse, error) {
	params := map[string]interface{}{
		"user_id": userId,
	}
	response := client.Rpc("get_deposits_by_userid", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	var deposits []DepositResponse
	if err := json.Unmarshal([]byte(response), &deposits); err != nil {
		return nil, fmt.Errorf("error unmarshalling user response: %v", err)
	}

	return &deposits, nil
}

func GetDepositsByUserAddress(client *supabase.Client, walletAddress, walletType string) (*[]DepositResponse, error) {
	params := map[string]interface{}{
		"wallet_addr": walletAddress,
		"wallet_t":    walletType,
	}

	response := client.Rpc("get_deposits_by_address", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	var deposits []DepositResponse
	err := json.Unmarshal([]byte(response), &deposits)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	return &deposits, nil
}

func GetOrderById(client *supabase.Client, id string) (*OrderAndUserResponse, error) {
	fmt.Printf("\n this is where i really am")
	params := map[string]interface{}{
		"id_": id,
	}
	response := client.Rpc("get_order_by_id", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	var order OrderAndUserResponse
	if err := json.Unmarshal([]byte(response), &order); err != nil {
		return nil, fmt.Errorf("error unmarshalling user response: %v", err)
	}

	return &order, nil
}

func GetOrderById2(client *supabase.Client, id string) (*OrderAndUserResponse2, error) {
	fmt.Printf("\n this is where i really am")
	params := map[string]interface{}{
		"id_": id,
	}
	response := client.Rpc("get_order_by_id2", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	var order OrderAndUserResponse2
	if err := json.Unmarshal([]byte(response), &order); err != nil {
		return nil, fmt.Errorf("error unmarshalling user response: %v", err)
	}

	return &order, nil
}

func GetOrdersByUserId(client *supabase.Client, userId string) (*[]OrderResponse, error) {
	params := map[string]interface{}{
		"user_id": userId,
	}
	response := client.Rpc("get_orders_by_userid", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	var orders []OrderResponse
	if err := json.Unmarshal([]byte(response), &orders); err != nil {
		return nil, fmt.Errorf("error unmarshalling user response: %v", err)
	}

	return &orders, nil
}

func GetOrdersByUserId2(client *supabase.Client, userId string) (*[]OrderResponse2, error) {
	params := map[string]interface{}{
		"user_id": userId,
	}
	response := client.Rpc("get_orders_by_userid2", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	var orders []OrderResponse2
	if err := json.Unmarshal([]byte(response), &orders); err != nil {
		return nil, fmt.Errorf("error unmarshalling user response: %v", err)
	}

	return &orders, nil
}

func GetOrdersByAddress(client *supabase.Client, walletAddress, walletType string) (*[]OrderResponse, error) {
	params := map[string]interface{}{
		"wallet_addr": walletAddress,
		"wallet_t":    walletType,
	}
	response := client.Rpc("get_orders_by_address", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	var orders []OrderResponse
	if err := json.Unmarshal([]byte(response), &orders); err != nil {
		return nil, fmt.Errorf("error unmarshalling user response: %v", err)
	}

	return &orders, nil
}

func GetOrdersByAddress2(client *supabase.Client, walletAddress, walletType string) (*[]OrderResponse2, error) {
	params := map[string]interface{}{
		"wallet_addr": walletAddress,
		"wallet_t":    walletType,
	}
	response := client.Rpc("get_orders_by_address2", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	var orders []OrderResponse2
	if err := json.Unmarshal([]byte(response), &orders); err != nil {
		return nil, fmt.Errorf("error unmarshalling user response: %v", err)
	}

	return &orders, nil
}

func GetOrdersByUserAddress(client *supabase.Client, walletAddress, walletType string) (*[]OrderResponse, error) {
	fmt.Print("\n inside of GetOrdersByUserAddress")
	params := map[string]interface{}{
		"wallet_addr": walletAddress,
		"wallet_t":    walletType,
	}

	response := client.Rpc("get_orders_by_address", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	fmt.Printf("\n response: %v", response)
	var orders []OrderResponse
	err := json.Unmarshal([]byte(response), &orders)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}
	fmt.Printf("\n response: %v", &orders)

	return &orders, nil
}

func GetSignatureValidationHash(client *supabase.Client, SignatureId string) (*GetSignatureValidationHashResponse, error) {
	params := map[string]interface{}{
		"p_signature_id": SignatureId,
	}
	response := client.Rpc("get_signature_hash", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	var hash_ GetSignatureValidationHashResponse
	if err := json.Unmarshal([]byte(response), &hash_); err != nil {
		return nil, fmt.Errorf("error unmarshalling user response: %v", err)
	}

	return &hash_, nil
}

func GetGlobalStateMetrics(client *supabase.Client, metrics []string) (*[]GlobalStateResponse, error) {
	params := map[string]interface{}{
		"metrics": metrics,
	}

	// Call the RPC function in Supabase
	response := client.Rpc("get_global_state_metrics", "exact", params)

	// Check for Supabase errors
	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	// Parse the response into the GlobalStateResponse slice
	var metricsResponse []GlobalStateResponse
	err := json.Unmarshal([]byte(response), &metricsResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	return &metricsResponse, nil
}

func GetSignatureHash(client *supabase.Client, signatureId string) (*GetSignatureHashResponse, error) {
	fmt.Printf("\n this is where i really am")
	params := map[string]interface{}{
		"p_signature_id": signatureId,
	}
	response := client.Rpc("get_signature_hash", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	var order GetSignatureHashResponse
	if err := json.Unmarshal([]byte(response), &order); err != nil {
		return nil, fmt.Errorf("error unmarshalling user response: %v", err)
	}

	return &order, nil
}
