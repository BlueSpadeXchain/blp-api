package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	yourpackage "path/to/your/generated/go/bindings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/supabase-community/postgrest-go"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mainnetURL := os.Getenv("MAINNET_URL")
	sepoliaURL := os.Getenv("SEPOLIA_URL")
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	escrowAddress := os.Getenv("TESTNET_ESCROW")

	// Ensure URLs are set
	if mainnetURL == "" || sepoliaURL == "" || supabaseURL == "" || supabaseKey == "" {
		log.Fatal("Required environment variables are not set")
	}

	// Connect to the Ethereum network
	client, err := ethclient.Dial(sepoliaURL) // Change to correct network
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Ethereum network")

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", supabaseKey),
	}
	supabaseClient := postgrest.NewClient(supabaseURL, "public", headers)

	// Address of the deployed contract
	contractAddress := common.HexToAddress(escrowAddress)
	contract, err := yourpackage.NewYourContract(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// Subscribe to Deposit events
	depositEvent := make(chan *yourpackage.YourContractDeposit)
	sub, err := contract.WatchDeposit(&bind.WatchOpts{Context: context.Background()}, depositEvent, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening for Deposit events...")

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case event := <-depositEvent:
			fmt.Printf("Deposit event: User=%s Token=%s Amount=%s\n", event.User.Hex(), event.Token.Hex(), event.Amount.String())
			// Update the user balance in the Supabase database
			updateUserBalance(supabaseClient, event.User.Hex(), event.Token.Hex(), event.Amount)
			// Store the event log in the Supabase database
			storeEventLog(supabaseClient, event)
		}
	}
}

func updateUserBalance(client *postgrest.Client, userAddress string, tokenAddress string, amount *big.Int) {
	// Check if the user already exists
	resp, statusCode, err := client.From("users").Select("*", "", false).Eq("signer", userAddress).Single().Execute()
	if err != nil {
		if statusCode == http.StatusNotFound {
			// User does not exist, create a new user row
			newUser := map[string]interface{}{
				"signer":     userAddress,
				"created_on": time.Now(),
				// Initialize balance fields
				fmt.Sprintf("%s_balance", tokenAddress): amount.String(),
			}
			_, _, err := client.From("users").Insert(newUser, false, "", "", "").Execute()
			if err != nil {
				log.Printf("Failed to create new user: %v\n", err)
				return
			}
			log.Printf("Created new user with address: %s\n", userAddress)
			return
		}
		log.Printf("Error fetching user: %v\n", err)
		return
	}

	// Parse the user data
	var user map[string]interface{}
	if err := json.Unmarshal(resp, &user); err != nil {
		log.Printf("Failed to parse user data: %v\n", err)
		return
	}

	// Update the user's balance
	balanceField := fmt.Sprintf("%s_balance", tokenAddress)
	newBalance := big.NewInt(0)
	if existingBalance, ok := user[balanceField]; ok && existingBalance != nil {
		existingBalanceBigInt, _ := new(big.Int).SetString(existingBalance.(string), 10)
		newBalance.Add(existingBalanceBigInt, amount)
	} else {
		newBalance.Set(amount)
	}

	updateData := map[string]interface{}{
		balanceField: newBalance.String(),
	}
	_, _, err = client.From("users").Update(updateData, "", "").Eq("signer", userAddress).Execute()
	if err != nil {
		log.Printf("Failed to update user balance: %v\n", err)
		return
	}

	log.Printf("Updated user balance for %s: %s %s\n", userAddress, newBalance.String(), tokenAddress)
}

func storeEventLog(client *postgrest.Client, event *yourpackage.YourContractDeposit) {
	eventData := map[string]interface{}{
		"token":  event.Token.Hex(),
		"user":   event.User.Hex(),
		"amount": event.Amount.String(),
	}

	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		log.Printf("Failed to marshal event data: %v\n", err)
		return
	}

	logEntry := map[string]interface{}{
		"tx_hash":         event.Raw.TxHash.Hex(),
		"block_timestamp": time.Unix(int64(event.Raw.BlockNumber), 0), // Adjust the conversion as needed
		"signer":          event.User.Hex(),
		"event_data":      string(eventJSON),
	}

	_, _, err = client.From("emit_logs").Insert(logEntry, false, "", "", "").Execute()
	if err != nil {
		log.Printf("Failed to store event log: %v\n", err)
		return
	}

	log.Printf("Stored event log for transaction: %s\n", event.Raw.TxHash.Hex())
}
