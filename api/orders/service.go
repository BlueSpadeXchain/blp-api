package orderHandler

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

	user "github.com/BlueSpadeXchain/blp-api/api/user"
	"github.com/BlueSpadeXchain/blp-api/pkg/db"
	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"github.com/supabase-community/supabase-go"
)

// each of these three request type add to requests pool to be processed

// func OpenOrderRequest(r *http.Request, parameters ...*interface{}) (interface{}, error) {
// 	return nil, nil
// }

// func CloseOrderRequest(r *http.Request, parameters ...*interface{}) (interface{}, error) {
// 	return nil, nil
// }

func GetOrdersRequest(r *http.Request, parameters ...*interface{}) (interface{}, error) {
	return nil, nil
}

// we need to sign with something that makes this specific tx unique
// for now we will just have the user sign of tx
func CreateOrderSignedRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*CreateOrderSignedRequestParams) (interface{}, error) {
	var params *CreateOrderSignedRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &CreateOrderSignedRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	txHash, _ := hex.DecodeString(utils.RemoveHex0xPrefix(""))
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

// type OrderRequestParams struct {
// 	Signer       string `query:"addr"`  // Signer address
// 	PerpId       string `query:"perp"`  // Target perpetual
// 	Collateral   string `query:"value"` // Collateral amount in USD
// 	EntryPrice   string `query:"entry"` // Entry price in USD
// 	Slippage     string `query:"slip"`  // Max slippage (basis points, out of 10,000)
// 	Leverage     string `query:"lev"`   // Leverage multiplier
// 	PositionType string `query:"type"`  // "long" or "short"
// 	Signature    string `query:"sig"`   // Signature over id . value . entry . slip . lev . type
// }

func CreateOrderUnsignedRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*CreateOrderUnsignedRequestParams) (interface{}, error) {
	var params *CreateOrderUnsignedRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &CreateOrderUnsignedRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	// var users []User
	// err := json.Unmarshal([]byte(response), &users)
	// if err != nil {
	// 	return nil, fmt.Errorf("error unmarshalling db.rpc response: %v", err)
	// }

	// if len(users) == 0 {
	// 	return nil, fmt.Errorf("db error: %v of type %v could not create a user", walletAddress, walletType)
	// }

	userData, err := user.GetUserByAddressRequest(r, supabaseClient, &user.GetUserByAddressRequestParams{
		Address:     params.WalletAddress,
		AddressType: params.WalletType})
	if err != nil {
		logrus.Error("GetUserByAddressRequest error:", err.Error())
		return nil, utils.ErrInternal(fmt.Sprintf("GetUserByAddressRequest error: %v", err.Error()))
	}

	balance := userData.(db.UserResponse).Balance
	fmt.Print(balance >= params.Collateral)

	// validate balance
	// validate pair exists in const mapping
	// fetch current value o

	// validate signature to verify backend query

	var userId string
	var orderType string
	var leverage string
	var pair string
	var amount string
	var entryPrice string
	var markPrice string

	if err := db.CreateOrder(supabaseClient, userId, orderType, leverage, pair, amount, entryPrice, markPrice); err != nil {
		return nil, utils.ErrInternal(fmt.Sprintf("Failed to add order: %v", err.Error()))
	}

	return nil, nil
}

func CloseOrder(r *http.Request, parameters ...*OrderCloseParams) (interface{}, error) {
	var params *OrderCloseParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &OrderCloseParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	// if err := validateOrderClose(params); err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
