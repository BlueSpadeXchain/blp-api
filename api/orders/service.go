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

// func CreateOrder2(
// 	client *supabase.Client,
// 	userId, orderType, pair string,
// 	leverage, collateral, entryPrice, liquidationPrice, maxPrice, limitPrice, stopLossPrice, takeProfitPrice, takeProfitAmount float64) (*OrderResponse2, error) {
// 	// Convert chainID, block, and depositNonce to string for TEXT type in the database
// 	params := map[string]interface{}{
// 		"user_id":     userId,
// 		"order_type":  orderType,
// 		"leverage":    leverage,
// 		"pair":        pair,
// 		"collateral":  collateral,
// 		"entry_price": entryPrice,
// 		"liq_price":   liquidationPrice,
// 		"max_price":   maxPrice,
// 		"lim_price":   limitPrice,
// 		"stop_price":  stopLossPrice,
// 		"tp_price":    takeProfitPrice,
// 		"tp_amount":   takeProfitAmount,
// 	}

type UnsignedOrder2RequestParams struct {
	UserId            string `query:"user-id" optional:"true"`    // implied user has an existing account if to have collateral
	Pair              string `query:"pair"`                       // Target perpetual, expects "BTC/USD", "ETH/USD", etc
	Collateral        string `query:"value" optional:"true"`      // Collateral amount in USD
	EntryPrice        string `query:"entry" optional:"true"`      // Entry price in USD
	Slippage          string `query:"slip" optional:"true"`       // Max slippage (basis points, out of 10,000)
	Leverage          string `query:"lev" optional:"true"`        // Leverage multiplier
	PositionType      string `query:"order-type" optional:"true"` // "long" or "short"
	LimitPrice        string `query:"lim-price" optional:"true"`
	StopLossPrice     string `query:"stop-price" optional:"true"`
	TakeProfitPrice   string `query:"tp-price" optional:"true"`
	TakeProfitPercent string `query:"tp-percent" optional:"true"` // percent to close the position for take profit, when achieved the tp_price and tp_value are set to null
}

// http://localhost:8080/api/order?query=create-unsigned-order&user-id=1d2664a39eee6098&pair=ethusd&value=1000&entry=33867498&lev=5&order-type=long
// http://localhost:8080/api/order?query=create-order-unsigned2&user-id=1d2664a39eee6098&pair=ethusd&value=1000&lev=5&order-type=short&lim-price=4000

func UnsignedOrder2Request(r *http.Request, supabaseClient *supabase.Client, parameters ...*UnsignedOrder2RequestParams) (interface{}, error) {
	var params *UnsignedOrder2RequestParams
	var markPrice, entryPrice, limitPrice, stopLossPrice, tpPrice, tpValue, tpCollateral, maxProfitPrice, liqPrice float64 // init as zero

	fmt.Print("\n I got this far")

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &UnsignedOrder2RequestParams{}
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

	if params.LimitPrice == "" || params.EntryPrice == "" {
		priceData, err := GetCurrentPriceData(pair)
		if err != nil {
			return nil, err
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
	}

	// skip mark price evaluation, if limit order
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

	// limit price
	if params.LimitPrice != "" && params.LimitPrice != "0" {
		var err error
		limitPrice, err = strconv.ParseFloat(params.LimitPrice, 64)
		if err != nil {
			return nil, utils.ErrInternal(fmt.Errorf("invalid limit price value: %w", err).Error())
		}
		markPrice = limitPrice
	}

	balance := userData.(*db.UserResponse).Balance
	if balance < collateral {
		return nil, utils.ErrInternal(fmt.Sprintf("user %v insufficent balance: expected >=%v, found %v", params.UserId, params.Collateral, balance))
	}

	leverage, err := strconv.ParseFloat(params.Leverage, 64)
	if err != nil {
		return nil, utils.ErrInternal(fmt.Sprintf("invalid leverage value: %v", err.Error()))
	}

	// Calculate liquidation price
	switch params.PositionType {
	case "long":
		liqPrice = markPrice * (1 - (1 / leverage))
		maxProfitPrice = markPrice * (1 + 10/leverage)
	case "short":
		liqPrice = markPrice * (1 + (1 / leverage))
		maxProfitPrice = markPrice * (1 - 10/leverage)
	default:
		return nil, utils.ErrInternal(fmt.Sprintf("invalid position type: %s", params.PositionType))
	}

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
				return nil, utils.ErrInternal(fmt.Sprintf("stop loss price too low, liquidation price: %v", liqPrice))
			}
			if markPrice <= stopLossPrice {
				return nil, utils.ErrInternal(fmt.Sprintf("stop loss price too high, entry price: %v", entryPrice))
			}
		case "short":
			if liqPrice <= stopLossPrice {
				return nil, utils.ErrInternal(fmt.Sprintf("stop loss price too high, liquidation price: %v", liqPrice))
			}
			if markPrice >= stopLossPrice {
				return nil, utils.ErrInternal(fmt.Sprintf("stop loss price too hilowgh, entry price: %v", entryPrice))
			}
		default:
			return nil, utils.ErrInternal(fmt.Sprintf("invalid position type: %s", params.PositionType))
		}

	}

	if params.TakeProfitPrice != "" && params.TakeProfitPrice != "0" {

		tpPrice_, err := strconv.ParseFloat(params.TakeProfitPrice, 64)
		if err != nil {
			return nil, utils.ErrInternal(fmt.Sprintf("invalid take profit price: %v", err.Error()))
		}
		tpPrice = tpPrice_
		if tpPrice <= 0 {
			return nil, utils.ErrInternal(fmt.Sprintf("invalid take profit price: %v", err.Error()))
		}
		tpPercent, err := strconv.ParseFloat(params.TakeProfitPercent, 64)
		if err != nil {
			return nil, utils.ErrInternal(fmt.Sprintf("invalid take profit price: %v", err.Error()))
		}
		tpCollateral = collateral * tpPercent / 100
		if params.PositionType == "long" {
			// For long positions: Profit when tpPrice > entryPrice
			if tpPrice <= entryPrice {
				return nil, utils.ErrInternal("Take profit price must be greater than the entry price for long positions")
			}
			tpValue = tpCollateral * leverage * (1 + (tpPrice-markPrice)/markPrice)
		} else if params.PositionType == "short" {
			// For short positions: Profit when tpPrice < entryPrice
			if tpPrice >= markPrice {
				return nil, utils.ErrInternal("Take profit price must be less than the entry price for short positions")
			}
			tpValue = tpCollateral * leverage * (1 + (markPrice-tpPrice)/markPrice)
		} else {
			return nil, utils.ErrInternal("Invalid order type")
		}

		if tpPrice >= maxProfitPrice {
			return nil, utils.ErrInternal(fmt.Sprintf("take profit price cannot exceed max price: %v", maxProfitPrice))
		}
		if tpPrice <= markPrice {
			return nil, utils.ErrInternal(fmt.Sprintf("take profit price must exceed entry price: %v", markPrice))
		}
	}

	response, err := db.CreateOrder2(
		supabaseClient,
		params.UserId,
		params.PositionType,
		pair,
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

	LogCreateOrderResponse2("Supabase create_order response", *response)
	//LogBeforeCreateOrderResponse(params.UserId, params.Pair, pair, collateral, entryPrice, entryPrice, liqPrice, leverage, params.PositionType, "unsigned")

	// still need to make the hash that the user will sign
	// thinking of just taking the keccak256 of the string order_.id

	// fmt.Printf("\n order id: %v", response.ID)
	orderIdBytes := []byte(response.ID)
	orderIdHash := crypto.Keccak256(orderIdBytes)

	return UnsignedOrderRequestResponse2{
		Order: *response,
		Hash:  hex.EncodeToString(orderIdHash),
	}, nil
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

// type SignedOrderRequestParams struct {
// 	OrderId string `query:"order-id"`
// 	R       string `query:"r"`
// 	S       string `query:"s"`
// 	V       string `query:"v"`
// }

func canModifyOrder(status string) error {
	invalid := []string{"filled", "canceled", "closed", "liquidated"}
	for _, i := range invalid {
		if status == i {
			return fmt.Errorf("orders of status %v cannot be mutated", status)
		}
	}
	return nil
}

func canCancelOrder(status string) error {
	invalid := []string{"pending", "filled", "canceled", "closed", "liquidated"}
	for _, i := range invalid {
		if status == i {
			return fmt.Errorf("orders of status %v cannot be mutated", status)
		}
	}
	return nil
}

func canCloseOrder(status string) error {
	invalid := []string{"unsigned", "filled", "canceled", "closed", "liquidated"}
	for _, i := range invalid {
		if status == i {
			return fmt.Errorf("orders of status %v cannot be mutated", status)
		}
	}
	return nil
}

func CloseOrderRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*ModifyOrderRequestParams) (interface{}, error) {
	var params *ModifyOrderRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &ModifyOrderRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	// need to fetch the order
	order, err := db.GetOrderById(supabaseClient, params.OrderId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}

	// order status
	// need to make sure that the order status is only pending
	// the user should not be able to mutate an unsigned, closed, filed, or cancelled order

	// also need add stop loss and limit order for the leverage positions
	// this means that the stop loss and limit are within the range of the closed positions
	// additionally a field for 10x needs to be added as a closing price
	if err := canCloseOrder(order.Order.Status); err != nil {
		return nil, utils.ErrInternal(err.Error())
	}

	// params.NewStatus,
	// need to capture mark price
	// order.Order.
	// 	orders, err := db.GetOrdersByUserAddress(supabaseClient, params.WalletAddress, params.WalletType)
	// if err != nil {
	// 	return nil, utils.ErrInternal(err.Error())
	// }
	return nil, nil
}

func CancelRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*GetOrdersByUserAddressRequestParams) (interface{}, error) {
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

func StopLossOrder() {}
func LimitOrder()    {}
