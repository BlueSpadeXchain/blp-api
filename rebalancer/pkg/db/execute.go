package db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/supabase-community/supabase-go"
)

func ProcessBatchOrders(client *supabase.Client, batchTimestamp time.Time, orderUpdates []OrderUpdate, globalUpdates OrderGlobalUpdate) error {
	orderGlobalUpdateTuple := fmt.Sprintf(
		"(%f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f)",
		globalUpdates.CurrentBorrowed,
		globalUpdates.CurrentLiquidity,
		globalUpdates.CurrentOrdersActive,
		globalUpdates.CurrentOrdersLimit,
		globalUpdates.CurrentOrdersPending,
		globalUpdates.TotalBorrowed,
		globalUpdates.TotalLiquidations,
		globalUpdates.TotalOrdersActive,
		globalUpdates.TotalOrdersFilled,
		globalUpdates.TotalOrdersLimit,
		globalUpdates.TotalOrdersLiquidated,
		globalUpdates.TotalOrdersStopped,
		globalUpdates.TotalPnlLosses,
		globalUpdates.TotalPnlProfits,
		globalUpdates.TotalRevenue,
		globalUpdates.TreasuryBalance,
		globalUpdates.TotalTreasuryProfits,
		globalUpdates.VaultBalance,
		globalUpdates.TotalVaultProfits,
		globalUpdates.TotalLiquidityRewards,
		globalUpdates.TotalStakeRewards,
	)

	params := map[string]interface{}{
		"batch_timestamp":      batchTimestamp.Format(time.RFC3339),
		"order_updates":        orderUpdates,
		"order_global_update_": orderGlobalUpdateTuple, // Pass as a string
	}
	response := client.Rpc("process_batch_orders", "estimate", params) // parse, so we already know the count

	var supabaseError SupabaseError
	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
		LogSupabaseError(supabaseError)
		return fmt.Errorf("supabase error: %v", supabaseError.Message)
	}

	return nil
}

// func ProcessBatchOrders(client *supabase.Client, batchTimestamp time.Time, orderUpdates []OrderUpdate, globalUpdates OrderGlobalUpdate) error {
// 	// Build array of order updates as ROW expressions
// 	orderUpdatesArray := make([]string, len(orderUpdates))
// 	for i, o := range orderUpdates {
// 		// Format the nested OrderGlobalUpdate as a ROW
// 		globalUpdateStr := fmt.Sprintf("ROW(%f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f)",
// 			o.OrderGlobalUpdate.CurrentBorrowed,
// 			o.OrderGlobalUpdate.CurrentLiquidity,
// 			o.OrderGlobalUpdate.CurrentOrdersActive,
// 			o.OrderGlobalUpdate.CurrentOrdersLimit,
// 			o.OrderGlobalUpdate.CurrentOrdersPending,
// 			o.OrderGlobalUpdate.TotalBorrowed,
// 			o.OrderGlobalUpdate.TotalLiquidations,
// 			o.OrderGlobalUpdate.TotalOrdersActive,
// 			o.OrderGlobalUpdate.TotalOrdersFilled,
// 			o.OrderGlobalUpdate.TotalOrdersLimit,
// 			o.OrderGlobalUpdate.TotalOrdersLiquidated,
// 			o.OrderGlobalUpdate.TotalOrdersStopped,
// 			o.OrderGlobalUpdate.TotalPnlLosses,
// 			o.OrderGlobalUpdate.TotalPnlProfits,
// 			o.OrderGlobalUpdate.TotalRevenue,
// 			o.OrderGlobalUpdate.TreasuryBalance,
// 			o.OrderGlobalUpdate.TotalTreasuryProfits,
// 			o.OrderGlobalUpdate.VaultBalance,
// 			o.OrderGlobalUpdate.TotalVaultProfits,
// 			o.OrderGlobalUpdate.TotalLiquidityRewards,
// 			o.OrderGlobalUpdate.TotalStakeRewards,
// 		)

// 		// Format the complete OrderUpdate as a ROW
// 		orderUpdatesArray[i] = fmt.Sprintf("ROW('%s', '%s', '%s', %f, %f, %f, %f, %f, %f, %s)::order_update",
// 			o.OrderID,
// 			o.UserID,
// 			o.Status,
// 			o.EntryPrice,
// 			o.ClosePrice,
// 			o.TpValue,
// 			o.Pnl,
// 			o.BalanceChange,
// 			o.EscrowBalanceChange,
// 			globalUpdateStr+"::order_global_update",
// 		)
// 	}

// 	// Join all ROWs into an array
// 	orderUpdatesSQL := "ARRAY[" + strings.Join(orderUpdatesArray, ", ") + "]::order_update[]"

// 	// Format the global update parameter
// 	globalUpdateSQL := fmt.Sprintf("ROW(%f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f)::order_global_update",
// 		globalUpdates.CurrentBorrowed,
// 		globalUpdates.CurrentLiquidity,
// 		globalUpdates.CurrentOrdersActive,
// 		globalUpdates.CurrentOrdersLimit,
// 		globalUpdates.CurrentOrdersPending,
// 		globalUpdates.TotalBorrowed,
// 		globalUpdates.TotalLiquidations,
// 		globalUpdates.TotalOrdersActive,
// 		globalUpdates.TotalOrdersFilled,
// 		globalUpdates.TotalOrdersLimit,
// 		globalUpdates.TotalOrdersLiquidated,
// 		globalUpdates.TotalOrdersStopped,
// 		globalUpdates.TotalPnlLosses,
// 		globalUpdates.TotalPnlProfits,
// 		globalUpdates.TotalRevenue,
// 		globalUpdates.TreasuryBalance,
// 		globalUpdates.TotalTreasuryProfits,
// 		globalUpdates.VaultBalance,
// 		globalUpdates.TotalVaultProfits,
// 		globalUpdates.TotalLiquidityRewards,
// 		globalUpdates.TotalStakeRewards,
// 	)

// 	params := map[string]interface{}{
// 		"batch_timestamp":      batchTimestamp.Format(time.RFC3339),
// 		"order_updates":        orderUpdatesSQL,
// 		"order_global_update_": globalUpdateSQL,
// 	}

// 	response := client.Rpc("process_batch_orders", "estimate", params) // parse, so we already know the count

// 	var supabaseError SupabaseError
// 	if err := json.Unmarshal([]byte(response), &supabaseError); err == nil && supabaseError.Message != "" {
// 		LogSupabaseError(supabaseError)
// 		return fmt.Errorf("supabase error: %v", supabaseError.Message)
// 	}

// 	return nil
// }
