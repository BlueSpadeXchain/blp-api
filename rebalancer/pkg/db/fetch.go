package db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

func GetOrdersParsingRange(client *supabase.Client, pairId string, minPrice, maxPrice float64) (*[]OrderResponse, error) {
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
	var orders *[]OrderResponse
	if err := json.Unmarshal([]byte(response), &orders); err != nil {
		return nil, fmt.Errorf("error unmarshalling user response: %v", err)
	}

	return orders, nil
}

func (o *OrderResponse) UnmarshalJSON(data []byte) error {
	type Alias OrderResponse // Create alias to avoid recursion

	// Temporary struct with string fields for parsing
	temp := struct {
		ID         string `json:"id"`
		CreatedAt  string `json:"created_at"`
		SignedAt   string `json:"signed_at"`
		StartedAt  string `json:"started_at"`
		ModifiedAt string `json:"modified_at"`
		EndedAt    string `json:"ended_at"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Parse UUIDs
	if id, err := uuid.Parse(temp.ID); err == nil {
		o.ID = id
	}

	// Parse timestamps
	layout := time.RFC3339
	if t, err := time.Parse(layout, temp.CreatedAt); err == nil {
		o.CreatedAt = t
	}
	if t, err := time.Parse(layout, temp.SignedAt); err == nil {
		o.SignedAt = t
	}
	if t, err := time.Parse(layout, temp.StartedAt); err == nil {
		o.StartedAt = t
	}
	if t, err := time.Parse(layout, temp.ModifiedAt); err == nil {
		o.ModifiedAt = t
	}
	if t, err := time.Parse(layout, temp.EndedAt); err == nil {
		o.EndedAt = t
	}

	return nil
}
