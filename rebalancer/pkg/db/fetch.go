package db

import (
	"encoding/json"
	"fmt"

	"github.com/supabase-community/supabase-go"
)

func GetOrdersParsingRange(client *supabase.Client, pairId string, minPrice, maxPrice float64) (*[]OrderResponse2, error) {
	params := map[string]interface{}{
		"pair_id_":   pairId,
		"min_price_": minPrice,
		"max_price_": maxPrice,
	}
	response := client.Rpc("get_orders_parsing_range", "estimate", params) // parse, so we already know the count

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	// edge case: can return empty array
	var orders *[]OrderResponse2
	if err := json.Unmarshal([]byte(response), &orders); err != nil {
		return nil, fmt.Errorf("error unmarshalling user response: %v", err)
	}

	return orders, nil
}
