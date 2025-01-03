package userHandler

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"

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

	txHash, _ := hex.DecodeString(RemoveHex0xPrefix(params.TxHash))
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

	return nil, nil
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

func GetUserByIdRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*GetUserByIdRequestParams) (interface{}, error) {
	var params *GetUserByIdRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetUserByIdRequestParams{}
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

func GetUserByAddressRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*GetUserByAddressRequestParams) (interface{}, error) {
	var params *GetUserByAddressRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetUserByAddressRequestParams{}
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

func RemoveHex0xPrefix(hex string) string {
	if strings.HasPrefix(hex, "0x") || strings.HasPrefix(hex, "0X") {
		return hex[2:]
	}
	return hex
}
