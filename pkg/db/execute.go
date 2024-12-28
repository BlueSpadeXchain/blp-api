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

func CreateOrder(client *supabase.Client, orderData map[string]interface{}) error {
	_, _, err := client.From("orders").Insert(orderData, false, "", "minimal", "").Execute()
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		return err
	}
	return nil
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

// %!(EXTRA string={"code":"PGRST202","details":"Searched for the function public.get_user_by_userid with parameter userid or with a single unnamed json/jsonb parameter, but no matches were found in the schema cache.","hint":"Perhaps you meant to call the function public.get_user_by_userid(user_id)","message":"Could not find the function public.get_user_by_userid(userid) in the schema cache"})
func GetUserByUserId(client *supabase.Client, userId string) (*User, error) {
	params := map[string]interface{}{
		"userid": userId,
	}
	response := client.Rpc("get_user_by_userid", "exact", params)

	// First try to unmarshal as error response
	var errResp SupabaseError
	if err := json.Unmarshal([]byte(response), &errResp); err == nil {
		if errResp.Code != "" { // If we successfully unmarshaled an error
			return nil, fmt.Errorf("supabase error: %s - %s", errResp.Code, errResp.Message)
		}
	}

	// If no error, try to unmarshal as user response
	var users []User
	if err := json.Unmarshal([]byte(response), &users); err != nil {
		return nil, fmt.Errorf("error unmarshalling user response: %v", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no user found for userId: %s", userId)
	}

	return &users[0], nil
}

func GetOrCreateUser(client *supabase.Client, walletAddress, walletType string) (*User, error) {
	params := map[string]interface{}{
		"wallet_addr": walletAddress,
		"wallet_t":    walletType,
	}

	response := client.Rpc("get_or_create_user", "exact", params)

	var users []User
	err := json.Unmarshal([]byte(response), &users)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("db error: %v of type %v could not create a user", walletAddress, walletType)
	}

	return &users[0], nil
}
