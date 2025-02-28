package db

import (
	"encoding/json"
	"fmt"

	"github.com/BlueSpadeXchain/blp-api/withdrawal/pkg/utils"
	"github.com/supabase-community/supabase-go"
)

func UpdateWithdrawalStatus(client *supabase.Client, withdrawalId, status, txHash string) (*ProcessUnstakeResponse, error) {
	params := map[string]interface{}{
		"p_withdrawal_id": withdrawalId,
		"p_tx_hash":       txHash,
	}

	if status != "success" || txHash == "" {
		params["p_status"] = "failure"
	} else {
		params["p_status"] = "success"
	}

	utils.LogInfo("update_withdrawal_status params", utils.StringifyStructFields(params, ""))

	// Execute the RPC call
	response := client.Rpc("update_withdrawal_status", "exact", params)

	// Check for any Supabase errors
	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return nil, fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	// If no response or an error, return
	if response == "" {
		return nil, fmt.Errorf("db error: failed to execute update_withdrawal_status for withdraw ID %v", withdrawalId)
	}

	var unstake ProcessUnstakeResponse
	err := json.Unmarshal([]byte(response), &unstake)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	}

	return &unstake, nil
}
