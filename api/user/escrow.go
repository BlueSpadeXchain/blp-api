package userHandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/BlueSpadeXchain/blp-api/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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

	// Initialize Supabase client
	options := &supabase.ClientOptions{}
	supabaseClient, err := supabase.NewClient(supabaseURL, supabaseKey, options)
	if err != nil {
		log.Fatalf("cannot initialize client: %v", err)
	}

	for {
		err := startListening(supabaseClient, sepoliaWS, escrowAddress)
		if err != nil {
			fmt.Printf("Connection lost. Attempting to reconnect in 5 seconds...\n")
			time.Sleep(5 * time.Second)
		}
	}
}

func startListening(supabaseClient *supabase.Client, sepoliaWS, escrowAddress string) error {
	// Connect to the Ethereum network
	client, err := ethclient.Dial(sepoliaWS)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum network: %v", err)
	}
	defer client.Close()

	fmt.Println("Connected to Ethereum network")

	// Address of the deployed contract
	contractAddress := common.HexToAddress(escrowAddress)
	contract, err := bindings.NewBindings(contractAddress, client)
	if err != nil {
		return fmt.Errorf("failed to bind to deployed contract: %v", err)
	}
	fmt.Printf("contract: %s\n", contractAddress)

	// Subscribe to Deposit events
	depositEvent := make(chan *bindings.BindingsDeposit)
	sub, err := contract.WatchDeposit(&bind.WatchOpts{Context: context.Background()}, depositEvent)
	if err != nil {
		return fmt.Errorf("failed to subscribe to Deposit events: %v", err)
	}
	defer sub.Unsubscribe()

	fmt.Println("Listening for Deposit events...")

	for {
		select {
		case err := <-sub.Err():
			log.Printf("WebSocket error: %v", err)
			return fmt.Errorf("WebSocket error: %v", err)
		case event := <-depositEvent:
			fmt.Printf("Deposit event: Asset=%s From=%s Amount=%s\n", event.Asset.Hex(), event.From.Hex(), event.Amount.String())
			// Update the user balance in the Supabase database
			signer := event.From.Hex()
			fmt.Printf("signer found: %s\n", signer)

			data, _, err := supabaseClient.From("users").Select("*", "exact", false).Filter("signer", "eq", signer).Limit(1, "").Execute()
			if err != nil {
				fmt.Printf("Failed to execute query: %v", err)
				continue
			}

			if string(data) != "[]" {
				var results2 []map[string]interface{}
				if err := json.Unmarshal(data, &results2); err != nil {
					log.Fatalf("Failed to unmarshal data: %v", err)
				}
				for _, userData := range results2 {
					if balanceData, ok := userData["balances"].(string); ok {
						var balances map[string]string

						if err := json.Unmarshal([]byte(balanceData), &balances); err != nil {
							log.Printf("Failed to unmarshal balances: %v\n", err)
						}
						err := updateBalance(supabaseClient, userData, event.Asset.Hex(), event.Amount.String())
						if err != nil {
							fmt.Printf("Error update balance: %v\n", err)
						}
					}
				}
			} else {
				err = createUser(supabaseClient, event.From.Hex(), time.Unix(int64(event.Raw.BlockNumber), 0).Format(time.RFC3339), event.Asset.Hex(), event.Amount.String())
				if err != nil {
					fmt.Printf("Error create user: %v\n", err)
				}
			}

			// Store the event log in the Supabase database
			storeEventLog(supabaseClient, event)
		}
	}
}

func createUser(client *supabase.Client, signer string, dateCreated string, assetAddress string, assetAmount string) error {
	balances := map[string]string{
		fmt.Sprintf("%s_balance", assetAddress): assetAmount,
	}
	balancesJSON, err := json.Marshal(balances)
	if err != nil {
		return fmt.Errorf("failed to marshal balances: %v", err)
	}

	newUser := map[string]interface{}{
		"signer":       signer,
		"date_created": dateCreated,
		"balances":     string(balancesJSON),
	}

	_, _, err = client.From("users").
		Insert(newUser, true, "", "*", "").
		Execute()
	if err != nil {
		return fmt.Errorf("failed to create new user: %v", err)
	}

	log.Printf("Created new user with address: %s\n", signer)

	return nil
}

// current goal is to add a balance to the jsonb string
func updateBalance(client *supabase.Client, user map[string]interface{}, assetAddress string, assetAmount string) error {
	if assetAmount == "0" {
		return fmt.Errorf("asset amount cannot be zero")
	}

	signer, ok := user["signer"].(string)
	if !ok {
		return fmt.Errorf("signer field is missing or not a string")
	}

	// ctx := context.Background()

	var balances map[string]string
	if balanceData, ok := user["balances"].(string); ok {

		if err := json.Unmarshal([]byte(balanceData), &balances); err != nil {
			return fmt.Errorf("failed to unmarshal balances: %v", err)
		}
		if len(balances) == 0 { // should never happen
			return fmt.Errorf("balances field is empty")
		}
	} else { // should never happen
		return fmt.Errorf("balances field does not exist")
	}

	balanceField := fmt.Sprintf("%s_balance", assetAddress)
	assetAmountInt, err := strconv.ParseInt(assetAmount, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse asset amount: %v", err)
	}
	if existingBalanceStr, ok := balances[balanceField]; ok {
		existingBalance, err := strconv.ParseInt(existingBalanceStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse existing balance: %v", err)
		}
		newBalance := existingBalance + assetAmountInt
		balances[balanceField] = strconv.FormatInt(newBalance, 10)
	} else {
		balances[balanceField] = strconv.FormatInt(assetAmountInt, 10)
	}

	updatedBalanceData, err := json.Marshal(balances)
	if err != nil {
		return fmt.Errorf("failed to marshal updated balances: %v", err)
	}

	_, _, err = client.From("users").
		Update(map[string]interface{}{"balances": string(updatedBalanceData)}, "", "").
		Filter("signer", "eq", signer).
		Execute()
	if err != nil {
		return fmt.Errorf("error updating user record: %v", err)
	}

	fmt.Printf("User %s balances updated!\n", signer)
	fmt.Println("Balances: %s\n", string(updatedBalanceData))

	return nil
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
