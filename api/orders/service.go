package orderHandler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	user "github.com/BlueSpadeXchain/blp-api/api/user"
	db "github.com/BlueSpadeXchain/blp-api/pkg/db"
	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
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

func GetOrdersByUserAddressRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*GetOrdersByUserAddressRequestParams) (interface{}, error) {
	var params *GetOrdersByUserAddressRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetOrdersByUserAddressRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	orders, err := db.GetOrdersByUserAddress(supabaseClient, params.WalletAddress, params.WalletType)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	return orders, nil
}

func GetOrdersByUserIdRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*GetOrdersByUserIdRequestParams) (interface{}, error) {
	var params *GetOrdersByUserIdRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetOrdersByUserIdRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	orders, err := db.GetOrdersByUserId(supabaseClient, params.UserId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	return orders, nil
}

func GetOrderByIdRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*GetOrdersByIdRequestParams) (interface{}, error) {
	var params *GetOrdersByIdRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetOrdersByIdRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	order, err := db.GetOrderById(supabaseClient, params.OrderId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	return order, nil
}

// we need to sign with something that makes this specific tx unique
// for now we will just have the user sign of tx
func SignedOrderRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*SignedOrderRequestParams) (interface{}, error) {
	var params *SignedOrderRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &SignedOrderRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	// the user
	order, err := db.GetOrderById(supabaseClient, params.OrderId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}

	// fmt.Printf("\n order id: %v", response.ID)
	// orderIdBytes := []byte(response.ID)
	// orderIdHash := crypto.Keccak256(orderIdBytes)

	// return UnsignedOrderRequestResponse{
	// 	Order: *response,
	// 	Hash:  hex.EncodeToString(orderIdHash),
	// }, nil

	// orderIdBytesTest := []byte("ae96c97c-4e0d-4b13-bae6-54efffc72859")
	// orderIdHashTest := crypto.Keccak256(orderIdBytesTest)
	// fmt.Printf("\n generated from input: %v", hex.EncodeToString(orderIdHashTest))
	// fmt.Printf("\n generated from api:   %v", "8cebec0419712a2ed98cb3ddd8aec8e92cc71f1677af693cfb99732e287d5902")
	// tempr, _ := hex.DecodeString("2fbf4a2ac97a29b60e3b56bc6392805e2a398ade1d9a31e1a80961347764a49f")
	// temps, _ := hex.DecodeString("6a9c55b2242e00d371169cb650c89130eb58a81b352bd1146735ae7557488e48")
	// tempv, _ := hex.DecodeString("00")
	// signatureBytesTest := append(tempr, temps...)
	// signatureBytesTest = append(signatureBytesTest, tempv...)
	// if ok, err := utils.ValidateEvmEcdsaSignature(orderIdHashTest, signatureBytesTest, common.HexToAddress("0xaf73d6bc4017518f45106c4eeb896204b99fd0e9")); !ok || err != nil {
	// 	if err != nil {
	// 		utils.LogError("error validating signature", err.Error())
	// 		return nil, utils.ErrInternal(fmt.Sprintf("error validating signature: %v", err.Error()))
	// 	} else {
	// 		utils.LogError("signature validation failed", "invaid signature")
	// 		return nil, utils.ErrInternal("Signature validation failed: invalid signature")
	// 	}
	// }
	// fmt.Print("\n concluded test")

	orderIdBytes := []byte(params.OrderId)
	orderIdHash := crypto.Keccak256(orderIdBytes)

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

	utils.LogInfo("Signature details", utils.FormatKeyValueLogs([][2]string{
		{"address", order.User.WalletAddress},
		{"hash", hex.EncodeToString(orderIdHash)},
		{"signature", hex.EncodeToString(signatureBytes)},
		{"module", "signature-validation"},
	}))

	// if ok, err := utils.ValidateEvmEcdsaSignature(orderIdHash, signatureBytes, common.HexToAddress("0x"+order.User.WalletAddress)); !ok || err != nil {
	// 	if err != nil {
	// 		utils.LogError("error validating signature", err.Error())
	// 		return nil, utils.ErrInternal(fmt.Sprintf("error validating signature: %v", err.Error()))
	// 	} else {
	// 		utils.LogError("signature validation failed", "invaid signature")
	// 		return nil, utils.ErrInternal("Signature validation failed: invalid signature")
	// 	}
	// }

	orderResponse, err := db.SignOrder(supabaseClient, params.OrderId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}

	return orderResponse, nil
}

// GetCurrentPriceData queries the API for the current price data for the given pair ID.
func GetCurrentPriceData(pair string) (PriceUpdate, error) {
	baseURL := "https://hermes.pyth.network/v2/updates/price/latest"

	// Create the request with query parameters
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return PriceUpdate{}, fmt.Errorf("error parsing URL: %v", err)
	}

	q := reqURL.Query()
	q.Add("ids[]", pair)
	reqURL.RawQuery = q.Encode()

	resp, err := http.Get(reqURL.String())
	if err != nil {
		return PriceUpdate{}, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PriceUpdate{}, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PriceUpdate{}, fmt.Errorf("error reading response body: %v", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return PriceUpdate{}, fmt.Errorf("error unmarshaling response JSON: %v", err)
	}

	LogResponse(reqURL.String(), response.Parsed[0])

	return response.Parsed[0], nil
}

func UnsignedOrderRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*UnsignedOrderRequestParams) (interface{}, error) {
	var params *UnsignedOrderRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &UnsignedOrderRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			return nil, utils.ErrInternal(fmt.Sprintf("failed to parse params: %s", err.Error()))
		}
	}

	userData, err := user.GetUserByUserIdRequest(r, supabaseClient, &user.GetUserByUserIdRequestParams{
		UserId: params.UserId,
	})
	if err != nil {
		logrus.Error("GetUserByIdRequest error:", err.Error())
		return nil, utils.ErrInternal(fmt.Sprintf("GetUserByIdRequest error: %v", err.Error()))
	}

	collateral, err := strconv.ParseFloat(params.Collateral, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid collateral value: %w", err)
	}

	pair, err := getPair(params.Pair)
	if err != nil {
		return nil, err
	}

	priceData, err := GetCurrentPriceData(pair)
	if err != nil {
		return nil, err
	}
	markPrice, _ := strconv.ParseFloat(priceData.Price.Price, 64)
	exponent := priceData.Price.Expo

	var slippage float64
	if params.Slippage != "" {
		var err error
		slippage, err = strconv.ParseFloat(params.Slippage, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid slippage value: %w", err)
		}
	}

	// If exponent is negative, divide by 10^abs(Expo)
	if exponent < 0 {
		// Convert exponent to int64 for proper comparison
		for i := int64(0); i < -int64(exponent); i++ {
			markPrice /= 10
		}
	} else { // If exponent is positive, multiply by 10^Expo
		for i := int64(0); i < int64(exponent); i++ {
			markPrice *= 10
		}
	}

	balance := userData.(*db.UserResponse).Balance
	if balance < collateral {
		return nil, utils.ErrInternal(fmt.Sprintf("user %v insufficent balance: expected >=%v, found %v", params.UserId, params.Collateral, balance))
	}

	leverage, err := strconv.ParseFloat(params.Leverage, 64)
	if err != nil {
		return nil, utils.ErrInternal(fmt.Sprintf("invalid leverage value: %v", err.Error()))
	}

	// /api/order?query=create-order-unsigned&user-id=04b89ffbb4f53a4e&pair=ethusd&value=1&entry=3505&slip=500&lev=1&position-type=long

	// Calculate liquidation price
	var liqPrice float64
	switch params.PositionType {
	case "long":
		liqPrice = markPrice * (1 - (1 / leverage))
	case "short":
		liqPrice = markPrice * (1 + (1 / leverage))
	default:
		return nil, fmt.Errorf("invalid position type: %s", params.PositionType)
	}

	if liqPrice <= 0 {
		return nil, fmt.Errorf("invalid liquidation price calculated")
	}

	var maxSlippage = 0.05 // 5% slippage threshold
	if slippage != 0 {
		maxSlippage = slippage
	}
	slippageThreshold := markPrice * maxSlippage

	var entryPrice float64
	// Validate that the entryPrice is within acceptable slippage from the markPrice
	if params.EntryPrice != "0" && params.EntryPrice != "" {
		var err error
		entryPrice, err = strconv.ParseFloat(params.EntryPrice, 64)
		if err != nil {
			return nil, utils.ErrInternal(fmt.Sprintf("invalid entry price: %v", err.Error()))
		}

		if params.PositionType == "long" && (entryPrice-markPrice) > slippageThreshold {
			return nil, fmt.Errorf("long position: entry price exceeds 5%% slippage from the mark price")
		} else if params.PositionType == "short" && (markPrice-entryPrice) > slippageThreshold {
			return nil, fmt.Errorf("short position: entry price exceeds 5%% slippage from the mark price")
		}
	}

	if params.PositionType == "long" && (markPrice <= liqPrice) {
		return nil, fmt.Errorf("long position: mark price in under liquidation price")
	} else if params.PositionType == "short" && (markPrice >= liqPrice) {
		return nil, fmt.Errorf("short position: mark price in over liquidation price")
	}

	response, err := db.CreateOrder(supabaseClient, params.UserId, params.PositionType, leverage, pair, collateral, markPrice, liqPrice)
	if err != nil {
		return nil, utils.ErrInternal(fmt.Sprintf("db post response: %v", err.Error()))
	}

	LogCreateOrderResponse("Supabase create_order response", *response)
	LogBeforeCreateOrderResponse(params.UserId, params.Pair, pair, collateral, entryPrice, markPrice, liqPrice, leverage, params.PositionType, "unsigned")

	// still need to make the hash that the user will sign
	// thinking of just taking the keccak256 of the string order_.id

	fmt.Printf("\n order id: %v", response.ID)
	orderIdBytes := []byte(response.ID)
	orderIdHash := crypto.Keccak256(orderIdBytes)

	return UnsignedOrderRequestResponse{
		Order: *response,
		Hash:  hex.EncodeToString(orderIdHash),
	}, nil
}

// f22cb3fe-3514-4b5d-a763-4c16e6b3330b
// http://localhost:8080/api/order?query=create-unsigned-order&user-id=f22cb3fe-3514-4b5d-a763-4c16e6b3330b&pair=ethusd&value=1000&entry=33867498&slip=500&lev=10&position-type=long

func CloseOrder(r *http.Request, parameters ...*OrderCloseParams) (interface{}, error) {
	var params *OrderCloseParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &OrderCloseParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			return nil, utils.ErrInternal(fmt.Sprintf("failed to parse params: %s", err.Error()))
		}
	}

	// if err := validateOrderClose(params); err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
