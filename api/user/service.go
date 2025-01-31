package userHandler

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"os"

	"github.com/BlueSpadeXchain/blp-api/pkg/db"
	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"github.com/supabase-community/supabase-go"
)

func WithdrawRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*WithdrawRequestParams) (interface{}, error) {
	var params *WithdrawRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &WithdrawRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	return nil, nil
}

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
