package userHandler

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/BlueSpadeXchain/blp-api/pkg/db"
	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"github.com/supabase-community/supabase-go"
)

func DespositRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*DespositRequestParams) (interface{}, error) {
	var params *DespositRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &DespositRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	// validate signature to verify backend query
	txHash, _ := hex.DecodeString(utils.RemoveHex0xPrefix(params.TxHash))
	fmt.Printf("\n txhash: %v", txHash)
	signature, _ := hex.DecodeString(params.Signature)
	pubkey := os.Getenv("EVM_ADDRESS")
	if pubkey == "" {
		logrus.Fatal("EVM_ADDRESS is not set")
	}
	if ok, err := utils.ValidateEvmEcdsaSignature(crypto.Keccak256(txHash), signature, common.HexToAddress(pubkey)); !ok || err != nil {
		if err != nil {
			utils.LogError("error validating isgnature", err.Error())
			return nil, utils.ErrInternal(fmt.Sprintf("error validating signature: %v", err.Error()))
		} else {
			utils.LogError("signature validation failed", "invaid signature")
			return nil, utils.ErrInternal("Signature validation failed: invalid signature")
		}
	}

	// parse value or deposit (1 eth = 3000 usd, 1 token = 1 usd)
	var value string
	amount, ok := new(big.Int).SetString(params.Amount, 10) // Convert amount to big.Int
	if !ok {
		return nil, utils.ErrInternal("Invalid amount format")
	}

	if utils.RemoveHex0xPrefix(params.Asset) == "0000000000000000000000000000000000000000" {
		// If address(0), assume 18 decimals
		// 1 * 10^18 tokens = 3000 USD
		tokensPerUSD := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil) // 10^18
		usdValue := new(big.Float).Quo(new(big.Float).SetInt(amount), new(big.Float).SetInt(tokensPerUSD))
		usdValue.Mul(usdValue, big.NewFloat(3000)) // Multiply by 3000 USD
		value = fmt.Sprintf("%.9f", usdValue)
	} else {
		// For non-address(0), assume 9 decimals
		// 1 * 10^9 tokens = 1 USD
		tokensPerUSD := new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil) // 10^9
		usdValue := new(big.Float).Quo(new(big.Float).SetInt(amount), new(big.Float).SetInt(tokensPerUSD))
		value = fmt.Sprintf("%.9f", usdValue)
	}

	if err := db.AddUserDeposit(
		supabaseClient,
		utils.RemoveHex0xPrefix(params.Receiver),
		"ecdsa",
		params.ChainId,
		params.Block,
		utils.RemoveHex0xPrefix(params.BlockHash),
		utils.RemoveHex0xPrefix(params.TxHash),
		utils.RemoveHex0xPrefix(params.Sender),
		params.DepositNonce,
		utils.RemoveHex0xPrefix(params.Asset),
		params.Amount,
		value); err != nil {
		return nil, utils.ErrInternal(fmt.Sprintf("Failed to add deposit: %v", err.Error()))
	}

	return nil, nil
}

func GetDepositsByUserAddressRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*GetDepositsByUserAddressRequestParams) (interface{}, error) {
	var params *GetDepositsByUserAddressRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetDepositsByUserAddressRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	deposits, err := db.GetDepositsByUserAddress(supabaseClient, params.WalletAddress, params.WalletType)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	return deposits, nil
}

func GetDepositsByUserIdRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*GetDepositsByUserIdRequestParams) (interface{}, error) {
	var params *GetDepositsByUserIdRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetDepositsByUserIdRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	deposits, err := db.GetDepositsByUserId(supabaseClient, params.UserId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	return deposits, nil
}

func UserDataRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*UserDataRequestParams) (interface{}, error) {
	var params *UserDataRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &UserDataRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	return nil, nil
}

func GetUserByUserIdRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*GetUserByUserIdRequestParams) (interface{}, error) {
	var params *GetUserByUserIdRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetUserByUserIdRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	user, err := db.GetUserByUserId(supabaseClient, params.UserId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}

	return user, nil
}

func GetUserByUserAddressRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*GetUserByUserAddressRequestParams) (interface{}, error) {
	var params *GetUserByUserAddressRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetUserByUserAddressRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	user, err := db.GetOrCreateUser(supabaseClient, params.Address, params.AddressType)
	if err != nil {
		utils.LogError("db GetOrCreateUser failed", err.Error())
		return nil, utils.ErrInternal(err.Error())
	}

	return user, nil
}

func AddAuthorizedWalletRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*AddAuthorizedWalletRequestParams) (interface{}, error) {
	var params *AddAuthorizedWalletRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &AddAuthorizedWalletRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	return nil, nil
}

func RemoveAuthorizedWalletRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*RemoveAuthorizedWalletRequestParams) (interface{}, error) {
	var params *RemoveAuthorizedWalletRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &RemoveAuthorizedWalletRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	return nil, nil
}

func UnsignedStakeFromBalanceRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*DespositRequestParams) (interface{}, error) {
	var params *DespositRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &DespositRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	// validate signature to verify backend query
	txHash, _ := hex.DecodeString(utils.RemoveHex0xPrefix(params.TxHash))
	fmt.Printf("\n txhash: %v", txHash)
	signature, _ := hex.DecodeString(params.Signature)
	pubkey := os.Getenv("EVM_ADDRESS")
	if pubkey == "" {
		logrus.Fatal("EVM_ADDRESS is not set")
	}
	if ok, err := utils.ValidateEvmEcdsaSignature(crypto.Keccak256(txHash), signature, common.HexToAddress(pubkey)); !ok || err != nil {
		if err != nil {
			utils.LogError("error validating isgnature", err.Error())
			return nil, utils.ErrInternal(fmt.Sprintf("error validating signature: %v", err.Error()))
		} else {
			utils.LogError("signature validation failed", "invaid signature")
			return nil, utils.ErrInternal("Signature validation failed: invalid signature")
		}
	}

	// parse value or deposit (1 eth = 3000 usd, 1 token = 1 usd)
	var value string
	amount, ok := new(big.Int).SetString(params.Amount, 10) // Convert amount to big.Int
	if !ok {
		return nil, utils.ErrInternal("Invalid amount format")
	}

	if utils.RemoveHex0xPrefix(params.Asset) == "0000000000000000000000000000000000000000" {
		// If address(0), assume 18 decimals
		// 1 * 10^18 tokens = 3000 USD
		tokensPerUSD := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil) // 10^18
		usdValue := new(big.Float).Quo(new(big.Float).SetInt(amount), new(big.Float).SetInt(tokensPerUSD))
		usdValue.Mul(usdValue, big.NewFloat(3000)) // Multiply by 3000 USD
		value = fmt.Sprintf("%.9f", usdValue)
	} else {
		// For non-address(0), assume 9 decimals
		// 1 * 10^9 tokens = 1 USD
		tokensPerUSD := new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil) // 10^9
		usdValue := new(big.Float).Quo(new(big.Float).SetInt(amount), new(big.Float).SetInt(tokensPerUSD))
		value = fmt.Sprintf("%.9f", usdValue)
	}

	if err := db.AddUserDeposit(
		supabaseClient,
		utils.RemoveHex0xPrefix(params.Receiver),
		"ecdsa",
		params.ChainId,
		params.Block,
		utils.RemoveHex0xPrefix(params.BlockHash),
		utils.RemoveHex0xPrefix(params.TxHash),
		utils.RemoveHex0xPrefix(params.Sender),
		params.DepositNonce,
		utils.RemoveHex0xPrefix(params.Asset),
		params.Amount,
		value); err != nil {
		return nil, utils.ErrInternal(fmt.Sprintf("Failed to add deposit: %v", err.Error()))
	}

	return nil, nil
}

func StakeFromBalanceRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*DespositRequestParams) (interface{}, error) {
	var params *DespositRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &DespositRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	// validate signature to verify backend query
	txHash, _ := hex.DecodeString(utils.RemoveHex0xPrefix(params.TxHash))
	fmt.Printf("\n txhash: %v", txHash)
	signature, _ := hex.DecodeString(params.Signature)
	pubkey := os.Getenv("EVM_ADDRESS")
	if pubkey == "" {
		logrus.Fatal("EVM_ADDRESS is not set")
	}
	if ok, err := utils.ValidateEvmEcdsaSignature(crypto.Keccak256(txHash), signature, common.HexToAddress(pubkey)); !ok || err != nil {
		if err != nil {
			utils.LogError("error validating isgnature", err.Error())
			return nil, utils.ErrInternal(fmt.Sprintf("error validating signature: %v", err.Error()))
		} else {
			utils.LogError("signature validation failed", "invaid signature")
			return nil, utils.ErrInternal("Signature validation failed: invalid signature")
		}
	}

	// parse value or deposit (1 eth = 3000 usd, 1 token = 1 usd)
	var value string
	amount, ok := new(big.Int).SetString(params.Amount, 10) // Convert amount to big.Int
	if !ok {
		return nil, utils.ErrInternal("Invalid amount format")
	}

	if utils.RemoveHex0xPrefix(params.Asset) == "0000000000000000000000000000000000000000" {
		// If address(0), assume 18 decimals
		// 1 * 10^18 tokens = 3000 USD
		tokensPerUSD := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil) // 10^18
		usdValue := new(big.Float).Quo(new(big.Float).SetInt(amount), new(big.Float).SetInt(tokensPerUSD))
		usdValue.Mul(usdValue, big.NewFloat(3000)) // Multiply by 3000 USD
		value = fmt.Sprintf("%.9f", usdValue)
	} else {
		// For non-address(0), assume 9 decimals
		// 1 * 10^9 tokens = 1 USD
		tokensPerUSD := new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil) // 10^9
		usdValue := new(big.Float).Quo(new(big.Float).SetInt(amount), new(big.Float).SetInt(tokensPerUSD))
		value = fmt.Sprintf("%.9f", usdValue)
	}

	if err := db.AddUserDeposit(
		supabaseClient,
		utils.RemoveHex0xPrefix(params.Receiver),
		"ecdsa",
		params.ChainId,
		params.Block,
		utils.RemoveHex0xPrefix(params.BlockHash),
		utils.RemoveHex0xPrefix(params.TxHash),
		utils.RemoveHex0xPrefix(params.Sender),
		params.DepositNonce,
		utils.RemoveHex0xPrefix(params.Asset),
		params.Amount,
		value); err != nil {
		return nil, utils.ErrInternal(fmt.Sprintf("Failed to add deposit: %v", err.Error()))
	}

	return nil, nil
}

func EoaStakeRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*DespositRequestParams) (interface{}, error) {
	var params *DespositRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &DespositRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	// validate signature to verify backend query
	txHash, _ := hex.DecodeString(utils.RemoveHex0xPrefix(params.TxHash))
	fmt.Printf("\n txhash: %v", txHash)
	signature, _ := hex.DecodeString(params.Signature)
	pubkey := os.Getenv("EVM_ADDRESS")
	if pubkey == "" {
		logrus.Fatal("EVM_ADDRESS is not set")
	}
	if ok, err := utils.ValidateEvmEcdsaSignature(crypto.Keccak256(txHash), signature, common.HexToAddress(pubkey)); !ok || err != nil {
		if err != nil {
			utils.LogError("error validating isgnature", err.Error())
			return nil, utils.ErrInternal(fmt.Sprintf("error validating signature: %v", err.Error()))
		} else {
			utils.LogError("signature validation failed", "invaid signature")
			return nil, utils.ErrInternal("Signature validation failed: invalid signature")
		}
	}

	amount, ok := new(big.Int).SetString(params.Amount, 10) // Convert amount to big.Int
	if !ok {
		return nil, utils.ErrInternal("Invalid amount format")
	}

	bluAddress, usdcAddress, err := getNetworkAddresses()
	if err != nil {
		utils.LogError("getNetworkAddresses error", err.Error())
		return nil, utils.ErrInternal(err.Error())
	}

	assetAddress := utils.RemoveHex0xPrefix(params.Asset)

	logrus.Warning("assetAddress", assetAddress)
	logrus.Warning("bluAddress", bluAddress)
	logrus.Warning("usdcAddress", usdcAddress)

	// need to fetch the current price of eth
	switch assetAddress {
	case "0000000000000000000000000000000000000000": // stake blp: from eth
		priceData, err := utils.GetCurrentPriceData("ff61491a931112ddf1bd8147cd1b641375f79f5825126d665480874634fd0ace")
		if err != nil {
			return "", utils.ErrInternal(err.Error())
		}

		markPrice, err := strconv.ParseFloat(priceData.Price.Price, 64)
		if err != nil {
			return "", utils.ErrInternal(fmt.Sprintf("failed to parse mark price: %v", err))
		}

		value, err := calculatePriceValue(amount, markPrice, int32(priceData.Price.Expo), 18)
		if err != nil {
			return "", utils.ErrInternal(fmt.Sprintf("failed to calculate price value: %v", err))
		}

		if err := db.ProcessDepositAndStake(
			supabaseClient,
			utils.RemoveHex0xPrefix(params.Receiver),
			"ecdsa",
			params.ChainId,
			params.Block,
			utils.RemoveHex0xPrefix(params.BlockHash),
			utils.RemoveHex0xPrefix(params.TxHash),
			utils.RemoveHex0xPrefix(params.Sender),
			params.DepositNonce,
			utils.RemoveHex0xPrefix(params.Asset),
			params.Amount,
			value,
			"BLP"); err != nil {
			return nil, utils.ErrInternal(fmt.Sprintf("Failed to add deposit: %v", err.Error()))
		}
		break
	case bluAddress: // staked blu
		value, err := calculatePriceValue(amount, 1, 0, 18)
		if err != nil {
			return "", utils.ErrInternal(fmt.Sprintf("failed to calculate price value: %v", err))
		}

		if err := db.ProcessDepositAndStake(
			supabaseClient,
			utils.RemoveHex0xPrefix(params.Receiver),
			"ecdsa",
			params.ChainId,
			params.Block,
			utils.RemoveHex0xPrefix(params.BlockHash),
			utils.RemoveHex0xPrefix(params.TxHash),
			utils.RemoveHex0xPrefix(params.Sender),
			params.DepositNonce,
			utils.RemoveHex0xPrefix(params.Asset),
			params.Amount,
			value,
			"BLU"); err != nil {
			return nil, utils.ErrInternal(fmt.Sprintf("Failed to add deposit: %v", err.Error()))
		}
		break
	case usdcAddress: // stake blp: from usdc
		value, err := calculatePriceValue(amount, 1, 0, 6)
		if err != nil {
			return "", utils.ErrInternal(fmt.Sprintf("failed to calculate price value: %v", err))
		}

		if err := db.ProcessDepositAndStake(
			supabaseClient,
			utils.RemoveHex0xPrefix(params.Receiver),
			"ecdsa",
			params.ChainId,
			params.Block,
			utils.RemoveHex0xPrefix(params.BlockHash),
			utils.RemoveHex0xPrefix(params.TxHash),
			utils.RemoveHex0xPrefix(params.Sender),
			params.DepositNonce,
			utils.RemoveHex0xPrefix(params.Asset),
			params.Amount,
			value,
			"BLP"); err != nil {
			return nil, utils.ErrInternal(fmt.Sprintf("Failed to add deposit: %v", err.Error()))
		}
		break
	default:
		return utils.ErrInternal("Stake address not recognized"), nil
	}

	// if assetAddress == "0000000000000000000000000000000000000000" {
	// 	// If address(0), assume 18 decimals
	// 	// 1 * 10^18 tokens = 3000 USD

	// } else {
	// 	// For non-address(0), assume 9 decimals
	// 	// 1 * 10^9 tokens = 1 USD
	// 	tokensPerUSD := new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil) // 10^9
	// 	usdValue := new(big.Float).Quo(new(big.Float).SetInt(amount), new(big.Float).SetInt(tokensPerUSD))
	// 	value = fmt.Sprintf("%.9f", usdValue)
	// }

	// if utils.RemoveHex0xPrefix(params.Asset) == "0000000000000000000000000000000000000000" {
	// 	// If address(0), assume 18 decimals
	// 	// 1 * 10^18 tokens = 3000 USD
	// 	tokensPerUSD := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil) // 10^18
	// 	usdValue := new(big.Float).Quo(new(big.Float).SetInt(amount), new(big.Float).SetInt(tokensPerUSD))
	// 	usdValue.Mul(usdValue, big.NewFloat(3000)) // Multiply by 3000 USD
	// 	value = fmt.Sprintf("%.9f", usdValue)
	// } else {
	// 	// For non-address(0), assume 9 decimals
	// 	// 1 * 10^9 tokens = 1 USD
	// 	tokensPerUSD := new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil) // 10^9
	// 	usdValue := new(big.Float).Quo(new(big.Float).SetInt(amount), new(big.Float).SetInt(tokensPerUSD))
	// 	value = fmt.Sprintf("%.9f", usdValue)
	// }

	return nil, nil
}

func calculatePriceValue(amount *big.Int, price float64, priceExponent int32, tokenDecimals int64) (string, error) {
	if amount == nil {
		return "", fmt.Errorf("amount cannot be nil")
	}

	// Adjust mark price based on its exponent
	adjustedPrice := adjustPriceByExponent(price, priceExponent)

	// Convert amount from token decimals to regular value
	tokensScale := new(big.Int).Exp(big.NewInt(10), big.NewInt(tokenDecimals), nil)
	usdValue := new(big.Float).Quo(new(big.Float).SetInt(amount), new(big.Float).SetInt(tokensScale))

	// Multiply by adjusted mark price
	usdValue.Mul(usdValue, big.NewFloat(adjustedPrice))

	return fmt.Sprintf("%.9f", usdValue), nil
}

func adjustPriceByExponent(price float64, exponent int32) float64 {
	return price * math.Pow10(int(exponent))
}

func StakeRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*DespositRequestParams) (interface{}, error) {
	var params *DespositRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &DespositRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	// validate signature to verify backend query
	txHash, _ := hex.DecodeString(utils.RemoveHex0xPrefix(params.TxHash))
	fmt.Printf("\n txhash: %v", txHash)
	signature, _ := hex.DecodeString(params.Signature)
	pubkey := os.Getenv("EVM_ADDRESS")
	if pubkey == "" {
		logrus.Fatal("EVM_ADDRESS is not set")
	}
	if ok, err := utils.ValidateEvmEcdsaSignature(crypto.Keccak256(txHash), signature, common.HexToAddress(pubkey)); !ok || err != nil {
		if err != nil {
			utils.LogError("error validating isgnature", err.Error())
			return nil, utils.ErrInternal(fmt.Sprintf("error validating signature: %v", err.Error()))
		} else {
			utils.LogError("signature validation failed", "invaid signature")
			return nil, utils.ErrInternal("Signature validation failed: invalid signature")
		}
	}

	// parse value or deposit (1 eth = 3000 usd, 1 token = 1 usd)
	var value string
	amount, ok := new(big.Int).SetString(params.Amount, 10) // Convert amount to big.Int
	if !ok {
		return nil, utils.ErrInternal("Invalid amount format")
	}

	if utils.RemoveHex0xPrefix(params.Asset) == "0000000000000000000000000000000000000000" {
		// If address(0), assume 18 decimals
		// 1 * 10^18 tokens = 3000 USD
		tokensPerUSD := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil) // 10^18
		usdValue := new(big.Float).Quo(new(big.Float).SetInt(amount), new(big.Float).SetInt(tokensPerUSD))
		usdValue.Mul(usdValue, big.NewFloat(3000)) // Multiply by 3000 USD
		value = fmt.Sprintf("%.9f", usdValue)
	} else {
		// For non-address(0), assume 9 decimals
		// 1 * 10^9 tokens = 1 USD
		tokensPerUSD := new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil) // 10^9
		usdValue := new(big.Float).Quo(new(big.Float).SetInt(amount), new(big.Float).SetInt(tokensPerUSD))
		value = fmt.Sprintf("%.9f", usdValue)
	}

	if err := db.AddUserDeposit(
		supabaseClient,
		utils.RemoveHex0xPrefix(params.Receiver),
		"ecdsa",
		params.ChainId,
		params.Block,
		utils.RemoveHex0xPrefix(params.BlockHash),
		utils.RemoveHex0xPrefix(params.TxHash),
		utils.RemoveHex0xPrefix(params.Sender),
		params.DepositNonce,
		utils.RemoveHex0xPrefix(params.Asset),
		params.Amount,
		value); err != nil {
		return nil, utils.ErrInternal(fmt.Sprintf("Failed to add deposit: %v", err.Error()))
	}

	return nil, nil
}

func GetStakesByUserIdRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*GetStakesByUserIdRequestParams) (interface{}, error) {
	var params *GetStakesByUserIdRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetStakesByUserIdRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	deposits, err := db.GetStakesByUserId(supabaseClient, params.UserId, params.StakeType, 0)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	return deposits, nil
}

func GetStakesByUserAddressRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*GetStakesByUserAddressRequestParams) (interface{}, error) {
	var params *GetStakesByUserAddressRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetStakesByUserAddressRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	deposits, err := db.GetStakesByUserAddress(supabaseClient, params.WalletAddress, params.WalletType, params.StakeType, 0)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	return deposits, nil
}

type WithdrawBluRequestParams struct {
	PendingWithdrawlId string `query:"pending-withdrawl-id"`
	Amount             string `query:"amount"`
	WalletAddress      string `query:"wallet-address"`
}

func UnstakeRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*UnstakeRequestParams) (interface{}, error) {
	var params *UnstakeRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &UnstakeRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	amount, err := strconv.ParseFloat(params.Amount, 64)
	if err != nil {
		return nil, utils.ErrInternal(fmt.Sprintf("invalid amount input: %v", err.Error()))
	}

	stakeType := strings.ToUpper(params.StakeType)

	response, err := db.GetUserByUserId(supabaseClient, params.UserId)
	if err != nil {
		return nil, utils.ErrInternal(fmt.Sprintf("fetch user failed: %v", err.Error()))
	}

	switch stakeType {
	case "BLU":
		if (response.BluStakeBalance + response.BluStakePending) < amount {
			return nil, utils.ErrInternal(fmt.Sprintf("insufficent BLU balance: %v", response.BluStakeBalance+response.BluStakePending))
		}
	case "BLP":
		if (response.BlpStakeBalance + response.BlpStakePending) < amount {
			return nil, utils.ErrInternal(fmt.Sprintf("insufficent BLP balance: %v", response.BlpStakeBalance+response.BlpStakePending))
		}
	default:
		return nil, utils.ErrInternal(fmt.Sprintf("invalid stake-type found: %v", stakeType))
	}

	withdrawlApi := os.Getenv("WITHDRAW_API")
	if withdrawlApi == "" {
		logrus.Fatal("WITHDRAW_API is not set")
	}

	unstakeResponse, err := db.Unstake(supabaseClient, params.UserId, stakeType, amount)
	if err != nil {
		return nil, utils.ErrInternal(fmt.Sprintf("db unstake error: %v", err.Error()))
	}

	if unstakeResponse.PendingWithdrawl.TokenType == "BLU" {
		// call the withdraw api
		request := &WithdrawBluRequestParams{
			PendingWithdrawlId: unstakeResponse.PendingWithdrawl.ID,
			Amount:             fmt.Sprint(unstakeResponse.PendingWithdrawl.Amount),
			WalletAddress:      unstakeResponse.PendingWithdrawl.WalletAddress,
		}
		body, _ := ConvertStructToQuery(request)
		logrus.Info("body: ", body)
		logrus.Warning("withdraw-blu was triggered")
		sendRequest(withdrawlApi, "withdraw-blu", body)
	}

	// the api upon a good response, will call on chain
	return unstakeResponse, nil
}

// this is specifically for withdrawling from the users balance, thus the user must exist and have a balance
func UnsignedWithdrawRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*UnsignedWithdrawRequestParams) (interface{}, error) {
	var params *UnsignedWithdrawRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &UnsignedWithdrawRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	amount, err := strconv.ParseFloat(params.Amount, 64)
	if err != nil {
		return nil, utils.ErrInternal(fmt.Sprintf("invalid amount input: %v", err.Error()))
	}

	if response, err := db.GetUserByUserId(supabaseClient, params.UserId); err != nil {
		return nil, utils.ErrInternal(err.Error())
	} else {
		if response.Balance < amount {
			return nil, utils.ErrInternal(fmt.Sprintf("insufficent balance: %v", response.Balance))
		}
	}

	response, err := db.Withdraw(supabaseClient, params.UserId, amount)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	return response, nil
}

func SignedWithdrawRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*SignedWithdrawRequestParams) (interface{}, error) {
	var params *SignedWithdrawRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &SignedWithdrawRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	signatureV, err := strconv.ParseUint(params.V, 16, 64) // the value from raw metamask is messed up
	if err != nil {
		return nil, utils.ErrInternal(fmt.Sprintf("invalid v value: %v", err.Error()))
	}

	signatureR, err := hex.DecodeString(params.R)
	if err != nil {
		utils.LogError("invalid sig-s value", err.Error())
		return nil, utils.ErrInternal(fmt.Sprintf("invalid sig-r value: %v", err.Error()))
	}

	signatureS, err := hex.DecodeString(params.S)
	if err != nil {
		utils.LogError("invalid sig-s value", err.Error())
		return nil, utils.ErrInternal(fmt.Sprintf("invalid sig-s value: %v", err.Error()))
	}

	if signatureV >= 27 {
		signatureV -= 27
	}

	signatureBytes := append(signatureR, signatureS...)
	signatureBytes = append(signatureBytes, byte(signatureV))

	if response, err := db.GetSignatureValidationHash(supabaseClient, params.SignatureId); err != nil {
		return nil, utils.ErrInternal(err.Error())
	} else {
		hash_, _ := hex.DecodeString(response.Hash)
		logrus.Info(fmt.Sprintf("hash to evaluate: %v", hash_))
		// if ok, err := utils.ValidateEvmEcdsaSignature(orderIdHash, signatureBytes, common.HexToAddress("0x"+order.User.WalletAddress)); !ok || err != nil {
		// 	if err != nil {
		// 		utils.LogError("error validating signature", err.Error())
		// 		return nil, utils.ErrInternal(fmt.Sprintf("error validating signature: %v", err.Error()))
		// 	} else {
		// 		utils.LogError("signature validation failed", "invaid signature")
		// 		return nil, utils.ErrInternal("Signature validation failed: invalid signature")
		// 	}
		// }
	}

	cancelResponse, err := db.SignWithdraw(supabaseClient, params.WithdrawId, params.SignatureId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	if !cancelResponse.IsValid {
		utils.LogError("sign cancel order error", err.Error())
		return nil, utils.ErrInternal(fmt.Sprintf("invalid sig-s value: %v", cancelResponse.ErrorMessage))
	}
	return cancelResponse, nil
}
