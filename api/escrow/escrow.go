package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/BlueSpadeXchain/blp-api/bindings"
	"github.com/ethereum/go-ethereum/common"

	// "github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

type User struct {
	Signer      string `json:"signer"`
	DateCreated string `json:"date_created"`
	Balances    string `json:"balances"` // JSONB data as a string
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mainnetWS := os.Getenv("MAINNET_WS")
	sepoliaWS := os.Getenv("SEPOLIA_WS")
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	escrowAddress := os.Getenv("TESTNET_ESCROW")

	// Ensure URLs are set
	if mainnetWS == "" || sepoliaWS == "" || supabaseURL == "" || supabaseKey == "" {
		log.Fatal("Required environment variables are not set")
	}

	// Connect to the Ethereum network
	// client, err := ethclient.Dial(sepoliaWS) // Change to correct network
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("Connected to Ethereum network")

	// Initialize Supabase client
	options := &supabase.ClientOptions{}
	supabaseClient, err := supabase.NewClient(supabaseURL, supabaseKey, options)
	if err != nil {
		log.Fatalf("cannot initialize client: %v", err)
	}

	// Fetch data from the 'users' table
	data, count, err := supabaseClient.From("users").Select("*", "exact", false).Execute()
	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
	}

	// Print fetched data
	fmt.Printf("Data: %v\n", data)
	fmt.Printf("Count: %d\n", count)
	fmt.Println("Connected to Supabase")

	// Address of the deployed contract
	contractAddress := common.HexToAddress(escrowAddress)
	// contract, err := bindings.NewBindings(contractAddress, client)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Printf("contract: %s\n", contractAddress)

	// Example deposit event data
	signer := "0xB9c28690C461C68190cc44739b82863F73dA9E22"
	dateCreated := "2024-08-05T16:40:06"
	assetAddress := "0x0000000000000000000000000000000000000000"
	assetAmount := "10000000000000"

	// Create the user with the deposit event data
	createUser(supabaseClient, signer, dateCreated, assetAddress, assetAmount)

	signer = "0xB9c28690C461C68190cc44739b82863F73dA9E22"
	dateCreated = "2024-08-05T16:40:06"
	assetAddress = "0x0000000000000000000000000000000000000002"
	assetAmount = "20000000000000"

	log.Println("User inserted successfully")

	// Subscribe to Deposit events
	// depositEvent := make(chan *bindings.BindingsDeposit)
	// sub, err := contract.WatchDeposit(&bind.WatchOpts{Context: context.Background()}, depositEvent)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Listening for Deposit events...")

	// for {
	// 	select {
	// 	case err := <-sub.Err():
	// 		log.Fatal(err)
	// 	case event := <-depositEvent:
	// 		fmt.Printf("Deposit event: Asset=%s From=%s Amount=%s\n", event.Asset.Hex(), event.From.Hex(), event.Amount.String())
	// 		// Update the user balance in the Supabase database
	// 		updateUserBalance(supabaseClient, event.From.Hex(), event.Asset.Hex(), event.Amount)
	// 		// Store the event log in the Supabase database
	// 		storeEventLog(supabaseClient, event)
	// 	}
	// }
}

func createUser(client *supabase.Client, signer string, dateCreated string, assetAddress string, assetAmount string) {
	newUser := map[string]interface{}{
		"signer":       signer,
		"date_created": dateCreated,
		"balances": map[string]interface{}{
			fmt.Sprintf("%s_balance", assetAddress): assetAmount,
		},
	}

	_, _, err := client.From("users").
		Insert(newUser, true, "", "*", "").
		Execute()
	if err != nil {
		log.Printf("Failed to create new user: %v\n", err)
		return
	}

	log.Printf("Created new user with address: %s\n", signer)
}

func updateBalance(user map[string]interface{}, assetAddress string, assetAmount string) {
	// Step 1: Parse the JSONB to a map
	var balances map[string]string
	if balanceData, ok := user["balances"]; ok {
		if err := json.Unmarshal(balanceData.([]byte), &balances); err != nil {
			log.Printf("Failed to unmarshal balances: %v\n", err)
			return
		}
	} else {
		balances = make(map[string]string)
	}

	// Step 2: Convert the string integer values to int64
	balanceField := fmt.Sprintf("%s_balance", assetAddress)
	newBalance := int64(0)
	if existingBalanceStr, ok := balances[balanceField]; ok {
		existingBalance, err := strconv.ParseInt(existingBalanceStr, 10, 64)
		if err != nil {
			log.Printf("Failed to parse existing balance: %v\n", err)
			return
		}
		newBalance = existingBalance + assetAmount
	} else {
		newBalance = assetAmount
	}

	// Step 4: Convert the result back to a string
	balances[balanceField] = strconv.FormatInt(newBalance, 10)

	// Step 5: Marshal the map back to JSONB
	updatedBalanceData, err := json.Marshal(balances)
	if err != nil {
		log.Printf("Failed to marshal updated balances: %v\n", err)
		return
	}

	// Update the user's balance
	user["balances"] = updatedBalanceData
}

func storeEventLog(client *supabase.Client, event *bindings.BindingsDeposit) {
	eventData := map[string]interface{}{
		"asset":  event.Asset.Hex(),
		"from":   event.From.Hex(),
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
		"signer":          event.From.Hex(),
		"event_data":      string(eventJSON),
	}

	_, _, err = client.From("emit_logs").Insert(logEntry, false, "", "", "1").Execute()
	if err != nil {
		log.Printf("Failed to store event log: %v\n", err)
		return
	}

	log.Printf("Stored event log for transaction: %s\n", event.Raw.TxHash.Hex())
}
