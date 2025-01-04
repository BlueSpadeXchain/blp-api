package db

import (
	"encoding/json"
	"fmt"

	"github.com/supabase-community/supabase-go"
)

func GetUserByUserId(client *supabase.Client, userId string) (*UserResponse, error) {
	params := map[string]interface{}{
		"user_id": userId,
	}
	response := client.Rpc("get_user_by_userid", "exact", params)

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	var users []UserResponse
	if err := json.Unmarshal([]byte(response), &users); err != nil {
		return nil, fmt.Errorf("error unmarshalling user response: %v", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no user found for userId: %s", userId)
	}

	return &users[0], nil
}

func GetDepositsByUserId(client *supabase.Client, userId string) (*[]DepositResponse, error) {
	params := map[string]interface{}{
		"user_id": userId,
	}
	response := client.Rpc("get_user_by_userid", "exact", params)

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

	response := client.Rpc("get_or_create_user", "exact", params)

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

func GetOrdersByUserId(client *supabase.Client, userId string) (*[]OrderResponse, error) {
	params := map[string]interface{}{
		"user_id": userId,
	}
	response := client.Rpc("get_user_by_userid", "exact", params)

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

func GetOrdersByUserAddress(client *supabase.Client, walletAddress, walletType string) (*[]OrderResponse, error) {
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

	var orders []OrderResponse
	err := json.Unmarshal([]byte(response), &orders)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	return &orders, nil
}
