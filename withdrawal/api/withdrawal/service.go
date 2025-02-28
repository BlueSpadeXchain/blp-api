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

	"github.com/BlueSpadeXchain/blp-api/withdrawal/pkg/db"
	"github.com/BlueSpadeXchain/blp-api/withdrawal/pkg/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/supabase-community/supabase-go"
)

type WithdrawBluRequestResponse struct {
}

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

	withdrawalApiKey := os.Getenv("WITHDRAWAL_API_KEY")

	if params.ApiKey != withdrawalApiKey {
		err_ := utils.ErrInternal("invalid api key")
		utils.LogError("", err_.Error())
		db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "FAILED", "Invalid withdrawal api key")
		return nil, err_
	}

	// Start blockchain processing in a goroutine
	go func() {
		// Setup context with 2-minute timeout
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		// Log start of processing
		utils.LogInfo("Starting blockchain withdrawal", fmt.Sprintf("UnstakeID: %s, Amount: %s", params.PendingWithdrawalId, params.Amount))

		// Get environment variables based on mainnet flag
		isMainnetEnabled := os.Getenv("MAINNET_ENABLED") == "true"
		var rpcURL, escrowAddress, bluAddress string

		if isMainnetEnabled {
			rpcURL = os.Getenv("MAINNET_JSON_RPC")
			escrowAddress = os.Getenv("MAINNET_ESCROW")
			bluAddress = os.Getenv("MAINNET_BLU")
		} else {
			rpcURL = os.Getenv("TESTNET_JSON_RPC")
			escrowAddress = os.Getenv("TESTNET_ESCROW")
			bluAddress = os.Getenv("TESTNET_BLU")
		}

		// Connect to blockchain
		client, err := ethclient.Dial(rpcURL)
		if err != nil {
			utils.LogError("failed to connect to blockchain", err.Error())
			db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "failure", "Blockchain connection failed")
			return
		}

		// Get user wallet address
		userAddress := common.HexToAddress(params.WalletAddress)

		// Parse ABI for escrow contract
		parsedEscrowABI, err := abi.JSON(strings.NewReader(escrowContractABI))
		if err != nil {
			utils.LogError("failed to parse escrow ABI", err.Error())
			db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "failure", "Contract configuration error")
			return
		}

		// Convert amount to big.Int with 18 decimals
		amountFloat, err := strconv.ParseFloat(params.Amount, 64)
		if err != nil {
			utils.LogError("invalid amount format", err.Error())
			db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "FAILED", "Invalid amount format")
			return
		}

		// Convert to wei (assuming 18 decimals)
		amountInWei := toWei(amountFloat, 18)

		// Get contract addresses
		escrowAddr := common.HexToAddress(escrowAddress)
		bluAddr := common.HexToAddress(bluAddress)

		utils.LogInfo("execute_transfer_function params", utils.StringifyStructFields(map[string]interface{}{
			"function":       "transfer",
			"escrow_address": escrowAddr,
			"user_address":   userAddress,
			"blu_address":    bluAddr,
			"amount_in_wei":  amountInWei.String(),
		}, ""))

		// Execute the transfer function
		txResponse, err := ExecuteFunction(
			*client,
			escrowAddr,
			parsedEscrowABI,
			"transfer",
			common.Big0, // No ETH value being sent
			userAddress, // Destination address
			bluAddr,     // Asset address
			amountInWei, // Amount in wei
		)
		//10 00000000 0000000000

		if err != nil {
			utils.LogError("blockchain transaction failed", err.Error())
			db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "failure", "")
			return
		}

		// Wait for up to 90 seconds for receipt
		receiptCtx, receiptCancel := context.WithTimeout(ctx, 90*time.Second)
		defer receiptCancel()

		receipt, err := waitForReceipt(receiptCtx, client, txResponse.TxHash)
		if err != nil {
			if err == context.DeadlineExceeded {
				utils.LogInfo("Transaction pending", "Receipt not available within timeout")
				db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "failure", txResponse.TxHash.String())
			} else {
				utils.LogError("failed to get receipt", err.Error())
				db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "failure", txResponse.TxHash.String())
			}
			return
		}

		// Check receipt status
		if receipt.Status == 1 {
			db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "success", txResponse.TxHash.String())
			utils.LogInfo("Blockchain withdrawal completed", fmt.Sprintf("UnstakeID: %s, TxHash: %s", params.PendingWithdrawalId, txResponse.TxHash.String()))
		} else {
			db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "failure", txResponse.TxHash.String())
			utils.LogError("transaction reverted", txResponse.TxHash.String())
		}
	}()

	// Return immediate response
	return map[string]string{
		"status":     "ACCEPTED",
		"message":    "Withdrawal BLU process initiated",
		"unstake_id": params.PendingWithdrawalId,
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

func WithdrawBalanceRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*WithdrawBalanceRequestParams) (interface{}, error) {
	var params *WithdrawBalanceRequestParams
	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &WithdrawBalanceRequestParams{}
	}
	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	withdrawalApiKey := os.Getenv("WITHDRAWAL_API_KEY")

	if params.ApiKey != withdrawalApiKey {
		err_ := utils.ErrInternal("invalid api key")
		utils.LogError("", err_.Error())
		db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "FAILED", "Invalid withdrawal api key")
		return nil, err_
	}

	// Start blockchain processing in a goroutine
	go func() {
		// Setup context with 2-minute timeout
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		// Log start of processing
		utils.LogInfo("Starting blockchain withdrawal", fmt.Sprintf("UnstakeID: %s, Amount: %s", params.PendingWithdrawalId, params.Amount))

		// Get environment variables based on mainnet flag
		isMainnetEnabled := os.Getenv("MAINNET_ENABLED") == "true"
		var rpcURL, escrowAddress, usdcAddress string

		if isMainnetEnabled {
			rpcURL = os.Getenv("MAINNET_JSON_RPC")
			escrowAddress = os.Getenv("MAINNET_ESCROW")
			usdcAddress = os.Getenv("MAINNET_USDC")
		} else {
			rpcURL = os.Getenv("TESTNET_JSON_RPC")
			escrowAddress = os.Getenv("TESTNET_ESCROW")
			usdcAddress = os.Getenv("MAINNET_USDC")
		}

		// Connect to blockchain
		client, err := ethclient.Dial(rpcURL)
		if err != nil {
			utils.LogError("failed to connect to blockchain", err.Error())
			db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "failure", "Blockchain connection failed")
			return
		}

		// Get user wallet address
		userAddress := common.HexToAddress(params.WalletAddress)

		// Parse ABI for escrow contract
		parsedEscrowABI, err := abi.JSON(strings.NewReader(escrowContractABI))
		if err != nil {
			utils.LogError("failed to parse escrow ABI", err.Error())
			db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "failure", "Contract configuration error")
			return
		}

		// Convert amount to big.Int with 18 decimals
		amountFloat, err := strconv.ParseFloat(params.Amount, 64)
		if err != nil {
			utils.LogError("invalid amount format", err.Error())
			db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "FAILED", "Invalid amount format")
			return
		}

		// Convert to wei (assuming 18 decimals)
		amountInWei := toWei(amountFloat*0.997, 18)

		// Get contract addresses
		escrowAddr := common.HexToAddress(escrowAddress)
		usdcAddr := common.HexToAddress(usdcAddress)

		utils.LogInfo("execute_transfer_function params", utils.StringifyStructFields(map[string]interface{}{
			"function":       "transfer",
			"escrow_address": escrowAddr,
			"user_address":   userAddress,
			"blu_address":    usdcAddr,
			"amount_in_wei":  amountInWei.String(),
		}, ""))

		// Execute the transfer function
		txResponse, err := ExecuteFunction(
			*client,
			escrowAddr,
			parsedEscrowABI,
			"transfer",
			common.Big0, // No ETH value being sent
			userAddress, // Destination address
			usdcAddr,    // Asset address
			amountInWei, // Amount in wei
		)
		//10 00000000 0000000000

		if err != nil {
			utils.LogError("blockchain transaction failed", err.Error())
			db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "failure", "")
			return
		}

		// Wait for up to 90 seconds for receipt
		receiptCtx, receiptCancel := context.WithTimeout(ctx, 90*time.Second)
		defer receiptCancel()

		receipt, err := waitForReceipt(receiptCtx, client, txResponse.TxHash)
		if err != nil {
			if err == context.DeadlineExceeded {
				utils.LogInfo("Transaction pending", "Receipt not available within timeout")
				db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "failure", txResponse.TxHash.String())
			} else {
				utils.LogError("failed to get receipt", err.Error())
				db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "failure", txResponse.TxHash.String())
			}
			return
		}

		// Check receipt status
		if receipt.Status == 1 {
			db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "success", txResponse.TxHash.String())
			utils.LogInfo("Blockchain withdrawal completed", fmt.Sprintf("UnstakeID: %s, TxHash: %s", params.PendingWithdrawalId, txResponse.TxHash.String()))
		} else {
			db.UpdateWithdrawalStatus(supabaseClient, params.PendingWithdrawalId, "failure", txResponse.TxHash.String())
			utils.LogError("transaction reverted", txResponse.TxHash.String())
		}
	}()

	// Return immediate response
	return map[string]string{
		"status":     "ACCEPTED",
		"message":    "Withdrawal balance process initiated",
		"unstake_id": params.PendingWithdrawalId,
	}, nil
}
