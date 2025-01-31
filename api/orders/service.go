package orderHandler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	user "github.com/BlueSpadeXchain/blp-api/api/user"
	db "github.com/BlueSpadeXchain/blp-api/pkg/db"
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

func GetOrdersByUserIdRequest_old(r *http.Request, supabaseClient *supabase.Client, parameters ...*GetOrdersByUserIdRequestParams) (interface{}, error) {
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

	orders, err := db.GetOrdersByUserId2(supabaseClient, params.UserId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	return orders, nil
}

func GetOrderByIdRequest_old(r *http.Request, supabaseClient *supabase.Client, parameters ...*GetOrdersByIdRequestParams) (interface{}, error) {
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

	order, err := db.GetOrderById2(supabaseClient, params.OrderId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	return order, nil
}

// deprecated
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

// http://localhost:8080/api/order?query=create-unsigned-order&user-id=1d2664a39eee6098&pair=ethusd&value=1000&entry=33867498&lev=5&order-type=long
// http://localhost:8080/api/order?query=create-order-unsigned2&user-id=1d2664a39eee6098&pair=ethusd&value=1000&lev=5&order-type=short&lim-price=4000

func CreateOrderRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*CreateOrderRequestParams) (interface{}, error) {
	var params *CreateOrderRequestParams
	var markPrice, entryPrice, limitPrice, stopLossPrice, tpPrice, tpValue, tpCollateral, maxProfitPrice, liqPrice float64 // init as zero

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &CreateOrderRequestParams{}
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
	pairId, err := getPairId(params.Pair)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	priceData, err := GetCurrentPriceData(pairId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	markPrice, _ = strconv.ParseFloat(priceData.Price.Price, 64)
	exponent := priceData.Price.Expo
	if exponent < 0 {
		for i := int64(0); i < -int64(exponent); i++ {
			markPrice /= 10
		}
	} else {
		for i := int64(0); i < int64(exponent); i++ {
			markPrice *= 10
		}
	}
	fmt.Printf("\n mark priceL: %v", markPrice)
	// skip mark price evaluation, if limit order
	fmt.Print("\n got here")
	if params.LimitPrice == "" && params.EntryPrice != "" {
		var err error
		entryPrice, err = strconv.ParseFloat(params.EntryPrice, 64)
		if err != nil {
			return nil, utils.ErrInternal(fmt.Sprintf("invalid entry price: %v", err.Error()))
		}

		var slippage float64
		if params.Slippage != "" {
			slippage, err = strconv.ParseFloat(params.Slippage, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid slippage value: %w", err)
			}
		}

		var maxSlippage = 0.05 // 5% slippage threshold
		if slippage != 0 {
			maxSlippage = slippage
		}
		slippageThreshold := markPrice * maxSlippage

		// Validate that the entryPrice is within acceptable slippage from the markPrice
		if params.PositionType == "long" && (entryPrice-markPrice) > slippageThreshold {
			return nil, fmt.Errorf("long position: entry price exceeds 5%% slippage from the mark price %v", markPrice)
		} else if params.PositionType == "short" && (markPrice-entryPrice) > slippageThreshold {
			return nil, fmt.Errorf("short position: entry price exceeds 5%% slippage from the mark price")
		}

		entryPrice = markPrice
	} else {
		entryPrice = markPrice
	}
	fmt.Print("\n got here")

	// limit price
	if params.LimitPrice != "" && params.LimitPrice != "0" {
		var err error
		limitPrice, err = strconv.ParseFloat(params.LimitPrice, 64)
		if err != nil {
			return nil, utils.ErrInternal(fmt.Errorf("invalid limit price value: %w", err).Error())
		}
		markPrice = limitPrice
	}
	fmt.Print("\n got here0")

	balance := userData.(*db.UserResponse).Balance
	if balance < collateral {
		return nil, utils.ErrInternal(fmt.Sprintf("user %v insufficent balance: expected >=%v, found %v", params.UserId, params.Collateral, balance))
	}
	fmt.Print("\n got here1")

	leverage, err := strconv.ParseFloat(params.Leverage, 64)
	if err != nil {
		return nil, utils.ErrInternal(fmt.Sprintf("invalid leverage value: %v", err.Error()))
	}
	fmt.Print("\n got here")

	// Calculate liquidation price
	switch params.PositionType {
	case "long":
		liqPrice = markPrice * (1 - (1 / leverage))
		maxProfitPrice = markPrice * (1 + 10/leverage)
	case "short":
		liqPrice = markPrice * (1 + (1 / leverage))
		maxProfitPrice = markPrice * (1 - 10/leverage)
	default:
		return nil, utils.ErrInternal(fmt.Sprintf("invalid position type: %v", params.PositionType))
	}
	fmt.Print("\n got here 45678")

	if liqPrice <= 0 {
		return nil, utils.ErrInternal(fmt.Sprintf("invalid liquidation price calculated %v", liqPrice))
	}

	// if maxProfitPrice <= 0 {
	// 	return nil, utils.ErrInternal(fmt.Sprintf("invalid max profit price calculated %v", maxProfitPrice))
	// }

	if params.PositionType == "long" && (markPrice <= liqPrice) {
		return nil, fmt.Errorf("long position: entry price in under liquidation price")
	} else if params.PositionType == "short" && (markPrice >= liqPrice) {
		return nil, fmt.Errorf("short position: entry price in over liquidation price")
	}

	// stop loss price
	if params.StopLossPrice != "" && params.StopLossPrice != "0" {
		var err error
		stopLossPrice, err = strconv.ParseFloat(params.StopLossPrice, 64)
		if err != nil {
			return nil, utils.ErrInternal(fmt.Errorf("invalid stop loss price value: %w", err).Error())
		}
		switch params.PositionType {
		case "long":
			if liqPrice >= stopLossPrice {
				return nil, utils.ErrInternal(fmt.Sprintf("stop loss %v price must exceed liquidation price: %v", params.PositionType, liqPrice))
			}
			if markPrice <= stopLossPrice {
				return nil, utils.ErrInternal(fmt.Sprintf("stop loss %v price cannot exceed entry price: %v", params.PositionType, entryPrice))
			}
		case "short":
			if liqPrice <= stopLossPrice {
				return nil, utils.ErrInternal(fmt.Sprintf("stop loss %v price cannot exceed liquidation price: %v", params.PositionType, liqPrice))
			}
			if markPrice >= stopLossPrice {
				return nil, utils.ErrInternal(fmt.Sprintf("stop loss %v price must exceed entry price: %v", params.PositionType, entryPrice))
			}
		default:
			return nil, utils.ErrInternal(fmt.Sprintf("invalid position type: %s", params.PositionType))
		}

	}
	fmt.Print("\n got here8")

	if params.TakeProfitPrice != "" && params.TakeProfitPrice != "0" {

		tpPrice_, err := strconv.ParseFloat(params.TakeProfitPrice, 64)
		if err != nil {
			return nil, utils.ErrInternal(fmt.Sprintf("invalid take profit price: %v", err.Error()))
		}
		tpPrice = tpPrice_
		if tpPrice <= 0 {
			return nil, utils.ErrInternal(utils.ErrInternal("invalid take profit price").Error())
		}
		tpPercent, err := strconv.ParseFloat(params.TakeProfitPercent, 64)
		if err != nil {
			return nil, utils.ErrInternal(fmt.Sprintf("invalid take profit price: %v", err.Error()))
		}
		tpCollateral = collateral * tpPercent / 100
		if params.PositionType == "long" {
			// For long positions: Profit when tpPrice > entryPrice
			if tpPrice <= entryPrice {
				return nil, utils.ErrInternal(fmt.Sprintf("take profit %v price must exceed entry price %v", params.PositionType, entryPrice))
			}
			if tpPrice >= maxProfitPrice {
				return nil, utils.ErrInternal(fmt.Sprintf("take profit %v price cannot exceed max price: %v", params.PositionType, maxProfitPrice))
			}
			if tpPrice <= markPrice {
				return nil, utils.ErrInternal(fmt.Sprintf("take profit %v price must exceed entry price: %v", params.PositionType, markPrice))
			}
			tpValue = tpCollateral * leverage * (1 + (tpPrice-markPrice)/markPrice)
		} else if params.PositionType == "short" {
			// For short positions: Profit when tpPrice < entryPrice
			if tpPrice >= markPrice {
				return nil, utils.ErrInternal(fmt.Sprintf("take profit %v price must exceed entry price %v", params.PositionType, entryPrice))
			}
			if tpPrice >= maxProfitPrice {
				return nil, utils.ErrInternal(fmt.Sprintf("take profit %v price cannot exceed max price: %v", params.PositionType, maxProfitPrice))
			}
			if tpPrice <= markPrice {
				return nil, utils.ErrInternal(fmt.Sprintf("take profit %v price must exceed entry price: %v", params.PositionType, markPrice))
			}
			tpValue = tpCollateral * leverage * (1 + (markPrice-tpPrice)/markPrice)
		} else {
			return nil, utils.ErrInternal("Invalid order type")
		}

	}
	fmt.Print("\n got here7")

	response, err := db.CreateOrder2(
		supabaseClient,
		params.UserId,
		params.PositionType,
		pairId,
		params.Pair,
		leverage,
		collateral,
		entryPrice,
		liqPrice,
		maxProfitPrice,
		limitPrice,
		stopLossPrice,
		tpPrice,
		tpValue,
		tpCollateral)
	if err != nil {
		return nil, utils.ErrInternal(fmt.Sprintf("db post response: %v", err.Error()))
	}

	LogCreateOrderResponse2("Supabase create_order response", response.Order)
	//LogBeforeCreateOrderResponse(params.UserId, params.Pair, pair, collateral, entryPrice, entryPrice, liqPrice, leverage, params.PositionType, "unsigned")

	// still need to make the hash that the user will sign
	// thinking of just taking the keccak256 of the string order_.id

	// fmt.Printf("\n order id: %v", response.ID)
	orderIdBytes := []byte(response.Order.ID)
	orderIdHash := crypto.Keccak256(orderIdBytes)

	return UnsignedOrderRequestResponse2{
		Order: response.Order,
		Hash:  hex.EncodeToString(orderIdHash),
	}, nil
}

func SignOrderRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*SignedOrderRequestParams) (interface{}, error) {
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
		err_ := utils.ErrInternal(fmt.Sprintf("invalid v value: %v", err.Error()))
		utils.ErrInternal(err.Error())
		utils.LogError(err_.Message, err_.Details)
		return nil, err_
	}

	signatureR, err := hex.DecodeString(params.R)
	if err != nil {
		utils.LogError("invalid sig-s value", err.Error())
		err_ := utils.ErrInternal(fmt.Sprintf("invalid sig-r value: %v", err.Error()))
		utils.ErrInternal(err.Error())
		utils.LogError(err_.Message, err_.Details)
		return nil, err_
	}

	signatureS, err := hex.DecodeString(params.S)
	if err != nil {
		utils.LogError("invalid sig-s value", err.Error())
		err_ := utils.ErrInternal(fmt.Sprintf("invalid sig-s value: %v", err.Error()))
		utils.ErrInternal(err.Error())
		utils.LogError(err_.Message, err_.Details)
		return nil, err_
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

	if ok, err := utils.ValidateEvmEcdsaSignature(orderIdHash, signatureBytes, common.HexToAddress("0x"+order.User.WalletAddress)); !ok || err != nil {
		if err != nil {
			utils.LogError("error validating signature", err.Error())
			// return nil, utils.ErrInternal(fmt.Sprintf("error validating signature: %v", err.Error()))
		} else {
			utils.LogError("signature validation failed", "invaid signature")
			// return nil, utils.ErrInternal("Signature validation failed: invalid signature")
		}
	}

	orderResponse, err := db.SignOrder2(supabaseClient, params.OrderId)
	if err != nil {
		err_ := utils.ErrInternal(err.Error())
		utils.LogError(err_.Message, err_.Details)
		return nil, err_
	}

	return orderResponse.Order, nil
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
// http://localhost:8080/api/order?query=create-unsigned-order&user-id=1d2664a39eee6098&pair=ethusd&collateral=1000&entry=33867498&slip=500&lev=10&position-type=long

func canModifyOrder(status string) error {
	invalid := []string{"filled", "canceled", "closed", "liquidated", "stopped"}
	for _, i := range invalid {
		if status == i {
			return fmt.Errorf("orders of status %v cannot be mutated", status)
		}
	}
	return nil
}

func canCancelOrder(status string) error {
	invalid := []string{"pending", "filled", "canceled", "closed", "liquidated", "stopped"}
	for _, i := range invalid {
		if status == i {
			return fmt.Errorf("orders of status %v cannot be mutated", status)
		}
	}
	return nil
}

func canCloseOrder(status string) error {
	invalid := []string{"unsigned", "filled", "canceled", "closed", "liquidated", "limit", "stopped"}
	for _, i := range invalid {
		if status == i {
			return fmt.Errorf("orders of status %v cannot be mutated", status)
		}
	}
	return nil
}

func UnsignedCloseOrderRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*UnsignedCloseOrderRequestParams) (interface{}, error) {
	var params *UnsignedCloseOrderRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &UnsignedCloseOrderRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	if response, err := db.GetOrderById(supabaseClient, params.OrderId); err != nil {
		return nil, utils.ErrInternal(err.Error())
	} else {
		if err := canCloseOrder(response.Order.Status); err != nil {
			return nil, utils.ErrInternal(err.Error())
		}
	}

	response, err := db.CloseOrder(supabaseClient, params.OrderId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	return response, nil
}

func SignedCloseOrderRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*SignedCloseOrderRequestParams) (interface{}, error) {
	var params *SignedCloseOrderRequestParams
	var markPrice, collateral, payoutValue, feeValue float64

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &SignedCloseOrderRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	// need to fetch the order
	response, err := db.GetOrderById2(supabaseClient, params.OrderId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}

	order_ := response.Order
	// user_ := response.User

	if err := canCloseOrder(order_.OrderStatus); err != nil {
		return nil, utils.ErrInternal(err.Error())
	}

	priceData, err := GetCurrentPriceData(order_.PairId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	markPrice, _ = strconv.ParseFloat(priceData.Price.Price, 64)
	exponent := priceData.Price.Expo

	if exponent < 0 {
		for i := int64(0); i < -int64(exponent); i++ {
			markPrice /= 10
		}
	} else {
		for i := int64(0); i < int64(exponent); i++ {
			markPrice *= 10
		}
	}

	// do we want deposits
	// do we want wthdraw fee (hmx does 0.3% for both)
	// we need a max utilization percent (hmx does 80%)
	// we need a delevergae buffer (20%)

	result, err := db.GetGlobalStateMetrics(supabaseClient, []string{"current_borrowed", "current_liquidity"})
	if result == nil || len(*result) != 2 {
		return nil, utils.ErrInternal("unexpected response from GetGlobalStateMetrics")
	}

	var totalBorrowed, totalLiquidity float64
	for _, metric := range *result {
		switch metric.Key {
		case "current_borrowed":
			totalBorrowed = metric.Value
		case "current_liquidity":
			totalLiquidity = metric.Value
		}
	}

	if order_.TakeProfitValue == 0 {
		collateral = order_.Collateral - order_.TakeProfitCollateral - order_.Collateral*0.00025
	} else {
		collateral = order_.Collateral - order_.Collateral*0.00025
	}

	createdAt, err := time.Parse(time.RFC3339Nano, order_.CreatedAt)
	if err != nil {
		return nil, utils.ErrInternal(fmt.Sprintf("invalid CreatedAt format: %v", err))
	}
	elapsedTime := time.Now().UTC().Sub(createdAt).Seconds()

	if totalLiquidity == 0 {
		return nil, utils.ErrInternal("total liquidity cannot be zero")
	}
	feePercent := (0.0001 * elapsedTime / 60 * totalBorrowed / totalLiquidity) + 0.00025

	// need to calculate the v
	switch order_.OrderType {
	case "long":
		payoutValue = collateral * order_.Leverage * (1 + (markPrice-order_.EntryPrice)/order_.EntryPrice)
	case "short":
		payoutValue = collateral * order_.Leverage * (1 + (order_.EntryPrice-markPrice)/order_.EntryPrice)
	default:
		return nil, utils.ErrInternal(fmt.Sprintf("unexpected order type: %v", order_.OrderType))
	}

	// Calculate feeValue and adjust value
	feeValue = payoutValue * feePercent
	payoutValue -= feeValue

	closeResponse, err := db.SignCloseOrder(supabaseClient, params.OrderId, params.SignatureId, payoutValue, feeValue)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	if !closeResponse.IsValid {
		return nil, utils.ErrInternal(closeResponse.ErrorMessage)
	}
	return closeResponse, nil
}

func UnsignedCancelOrderRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*UnsignedCancelOrderRequestParams) (interface{}, error) {
	var params *UnsignedCancelOrderRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &UnsignedCancelOrderRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	if response, err := db.GetOrderById(supabaseClient, params.OrderId); err != nil {
		return nil, utils.ErrInternal(err.Error())
	} else {
		if err := canCancelOrder(response.Order.Status); err != nil {
			return nil, utils.ErrInternal(err.Error())
		}
	}

	response, err := db.CancelOrder(supabaseClient, params.OrderId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	return response, nil
}

func SignedCancelOrderRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*SignedCancelOrderRequestParams) (interface{}, error) {
	var params *SignedCancelOrderRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &SignedCancelOrderRequestParams{}
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

	// given an order-id,
	// if response, err := db.GetOrderById(supabaseClient, params.OrderId); err != nil {
	// 	return nil, utils.ErrInternal(err.Error())
	// } else {
	// 	if err := canCancelOrder(response.Order.Status); err != nil {
	// 		return nil, utils.ErrInternal(err.Error())
	// 	}
	// }
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

	cancelResponse, err := db.SignCancelOrder(supabaseClient, params.OrderId, params.SignatureId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}
	if !cancelResponse.IsValid {
		return nil, utils.ErrInternal(cancelResponse.ErrorMessage)
	}
	return cancelResponse, nil
}
