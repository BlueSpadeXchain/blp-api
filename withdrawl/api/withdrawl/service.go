package withdrawHandler

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BlueSpadeXchain/blp-api/pkg/db"
	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/supabase-community/supabase-go"
)

func WithdrawBluRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*WithdrawBluRequestParams) (interface{}, error) {
	var params *WithdrawBluRequestParams
	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &WithdrawBluRequestParams{}
	}
	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	// Start blockchain processing in a goroutine
	go func() {
		// Setup context with 2-minute timeout
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		// Log start of processing
		utils.LogInfo("Starting blockchain withdrawal", fmt.Sprintf("UnstakeID: %s, Amount: %s", params.PendingWithdrawlId, params.Amount))

		// Get environment variables based on mainnet flag
		isMainnetEnabled := os.Getenv("MAINNET_ENABLED") == "true"
		var rpcURL, escrowAddress, bluAddress string

		if isMainnetEnabled {
			rpcURL = os.Getenv("MAINNET_JSON_RPC")
			escrowAddress = os.Getenv("MAINNET_ESCROW_ADDRESS")
			bluAddress = os.Getenv("MAINNET_BLU")
		} else {
			rpcURL = os.Getenv("TESTNET_JSON_RPC")
			escrowAddress = os.Getenv("TESTNET_ESCROW_ADDRESS")
			bluAddress = os.Getenv("TESTNET_BLU")
		}

		// Connect to blockchain
		client, err := ethclient.Dial(rpcURL)
		if err != nil {
			utils.LogError("failed to connect to blockchain", err.Error())
			db.UpdateWithdrawlStatus(supabaseClient, params.PendingWithdrawlId, "failure", "Blockchain connection failed")
			return
		}

		// Get user wallet address
		userAddress := common.HexToAddress(params.WalletAddress)

		// Parse ABI for escrow contract
		parsedEscrowABI, err := abi.JSON(strings.NewReader(escrowContractABI))
		if err != nil {
			utils.LogError("failed to parse escrow ABI", err.Error())
			db.UpdateWithdrawlStatus(supabaseClient, params.PendingWithdrawlId, "failure", "Contract configuration error")
			return
		}

		// Convert amount to big.Int with 18 decimals
		amountFloat, err := strconv.ParseFloat(params.Amount, 64)
		if err != nil {
			utils.LogError("invalid amount format", err.Error())
			db.UpdateWithdrawlStatus(supabaseClient, params.PendingWithdrawlId, "FAILED", "Invalid amount format")
			return
		}

		// Convert to wei (assuming 18 decimals)
		amountInWei := toWei(amountFloat, 18)

		// Get contract addresses
		escrowAddr := common.HexToAddress(escrowAddress)
		bluAddr := common.HexToAddress(bluAddress)

		// Execute the transfer function
		txHash, err := ExecuteFunction(
			*client,
			escrowAddr,
			parsedEscrowABI,
			"transfer",
			common.Big0, // No ETH value being sent
			userAddress, // Destination address
			bluAddr,     // Asset address
			amountInWei, // Amount in wei
		)

		if err != nil {
			utils.LogError("blockchain transaction failed", err.Error())
			updateWithdrawStatusInDB(supabaseClient, params.PendingWithdrawlId, "FAILED", fmt.Sprintf("Transaction failed: %v", err))
			return
		}

		// Update status with transaction hash
		updateWithdrawStatusInDB(supabaseClient, params.PendingWithdrawlId, "PENDING", txHash.String())

		// Wait for up to 90 seconds for receipt
		receiptCtx, receiptCancel := context.WithTimeout(ctx, 90*time.Second)
		defer receiptCancel()

		receipt, err := waitForReceipt(receiptCtx, client, txHash)
		if err != nil {
			if err == context.DeadlineExceeded {
				utils.LogInfo("Transaction pending", "Receipt not available within timeout")
				// Keep status as PENDING
			} else {
				utils.LogError("failed to get receipt", err.Error())
				updateWithdrawStatusInDB(supabaseClient, params.UnstakeID, "FAILED", "Failed to verify transaction")
			}
			return
		}

		// Check receipt status
		if receipt.Status == 1 {
			updateWithdrawStatusInDB(supabaseClient, params.UnstakeID, "COMPLETED", txHash.String())
			utils.LogInfo("Blockchain withdrawal completed", fmt.Sprintf("UnstakeID: %s, TxHash: %s", params.UnstakeID, txHash.String()))
		} else {
			updateWithdrawStatusInDB(supabaseClient, params.UnstakeID, "FAILED", "Transaction reverted on chain")
			utils.LogError("transaction reverted", txHash.String())
		}
	}()

	// Return immediate response
	return map[string]string{
		"status":     "ACCEPTED",
		"message":    "Withdrawal process initiated",
		"unstake_id": params.UnstakeID,
	}, nil
}

// Helper function to convert float to wei with specified decimals
func toWei(amount float64, decimals int) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(amount)

	// Create the multiplier (10^decimals)
	multiplier := new(big.Float)
	multiplier.SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))

	// Multiply
	bigval.Mul(bigval, multiplier)

	// Convert to big.Int
	result := new(big.Int)
	bigval.Int(result)
	return result
}

// Helper to wait for transaction receipt with timeout
func waitForReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			receipt, err := client.TransactionReceipt(ctx, txHash)
			if err == nil {
				return receipt, nil
			}
			if err != ethereum.NotFound {
				return nil, err
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func WithdrawBlpRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*WithdrawBlpRequestParams) (interface{}, error) {
	var params *WithdrawBlpRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &WithdrawBlpRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	return nil, nil
}
