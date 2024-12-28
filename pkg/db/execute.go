package db

import (
	"encoding/json"
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

func CreateUser(client *supabase.Client, email, username string, walletAddresses []map[string]string) error {
	userData := map[string]interface{}{
		"email":            email,
		"username":         username,
		"wallet_addresses": walletAddresses,
	}

	_, _, err := client.From("users").Insert(userData, false, "", "minimal", "").Execute()
	if err != nil {
		log.Printf("Failed to create user: %v", err)
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

func AddWalletToUser(client *supabase.Client, userID string, wallet map[string]string) error {
	// Retrieve current wallet_addresses
	response, _, err := client.From("users").Select("wallet_addresses", "", "").Eq("userid", userID).Execute()
	if err != nil {
		log.Printf("Failed to fetch user wallet_addresses: %v", err)
		return err
	}

	var users []map[string]interface{}
	err = json.Unmarshal(response, &users)
	if err != nil || len(users) == 0 {
		log.Printf("Failed to parse user wallet_addresses: %v", err)
		return err
	}

	currentWallets := users[0]["wallet_addresses"].([]map[string]string)
	currentWallets = append(currentWallets, wallet)

	updateData := map[string]interface{}{
		"wallet_addresses": currentWallets,
	}

	_, _, err = client.From("users").Update(updateData, "", "").Eq("userid", userID).Execute()
	if err != nil {
		log.Printf("Failed to add wallet to user: %v", err)
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
