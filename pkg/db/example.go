package db

import (
	"log"
	"os"

	"github.com/supabase-community/supabase-go"
)

func example() {
	// Initialize Supabase client
	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	supabaseClient, err := supabase.NewClient(supabaseUrl, supabaseKey, nil)

	// Example: Log an error
	err = LogError(supabaseClient, nil, "Example error message", map[string]interface{}{"example": "context"})
	if err != nil {
		log.Printf("Error logging example error: %v", err)
	}

	// Example: Modify user balance
	err = ModifyUserBalance(supabaseClient, "example_user_id", 100000000) // Update balance to 100 nano-USD
	if err != nil {
		log.Printf("Error modifying user balance: %v", err)
	}

	// Example: Create a withdrawal
	err = CreateWithdrawal(supabaseClient, "example_user_id", 50000000, "pending", "0xabc123...txhash")
	if err != nil {
		log.Printf("Error creating withdrawal: %v", err)
	}

	// Example: Modify withdrawal status
	err = ModifyWithdrawalStatus(supabaseClient, "example_withdrawal_id", "completed")
	if err != nil {
		log.Printf("Error modifying withdrawal status: %v", err)
	}

	// Example: Create an order
	order := map[string]interface{}{
		"userid":     "example_user_id",
		"order_type": "long",
		"leverage":   5.5,
		"pair":       "ETH-USD",
		"amount":     200000000, // 200 nano-USD
		"status":     "pending",
	}
	err = CreateOrder(supabaseClient, order)
	if err != nil {
		log.Printf("Error creating order: %v", err)
	}

	// Example: Modify an order
	orderUpdate := map[string]interface{}{
		"status": "filled",
		"amount": 150000000, // Adjusted to 150 nano-USD
	}
	err = ModifyOrder(supabaseClient, "example_order_id", orderUpdate)
	if err != nil {
		log.Printf("Error modifying order: %v", err)
	}

	// Example: Close an order
	err = CloseOrder(supabaseClient, "example_order_id")
	if err != nil {
		log.Printf("Error closing order: %v", err)
	}
}
