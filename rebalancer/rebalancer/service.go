package rebalancer

import (
	"bufio"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/BlueSpadeXchain/blp-api/rebalancer/pkg/db"
	"github.com/BlueSpadeXchain/blp-api/rebalancer/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/supabase-community/supabase-go"
)

func getFeeScalingFactor() float64 {
	return 0.3
}

func getBaseFee() float64 {
	return 0.001
}

func getPerHourFee() float64 {
	return 0.0001
}

func dynamicLeverageFee(leverage float64) float64 {
	//fee percent = 1/ (1+ scaling factor * log(leverage)) * base fee / 100
	return 1 / (1 + getFeeScalingFactor()*math.Log(leverage)) * getBaseFee()
}

func dynamicUtilizationFee(startTimestamp time.Time, globalBorrowed, globalLiquidity float64) float64 {
	elapsedTime := time.Since(startTimestamp.UTC()).Seconds()

	return getPerHourFee() * (elapsedTime / 3600) * globalBorrowed / globalLiquidity
}

func getCurrentBorrowAndLiquidity(supabaseClient *supabase.Client) (float64, float64, error) {
	result, err := db.GetGlobalStateMetrics(supabaseClient, []string{"current_borrowed", "current_liquidity"})
	if err != nil {
		logrus.Error(err.Error())
		return 0, 0, err
	}
	if result == nil || len(*result) != 2 {
		return 0, 0, fmt.Errorf("unexpected response from GetGlobalStateMetrics: %v", err)
	}

	var currentBorrowed, currentLiquidity float64
	for _, metric := range *result {
		switch metric.Key {
		case "current_borrowed":
			currentBorrowed = metric.Value
		case "current_liquidity":
			currentLiquidity = metric.Value
		}
	}
	return currentBorrowed, currentLiquidity, nil
}

func printProcessedOrder(order db.OrderResponse, orderUpdate db.OrderUpdate) {
	utils.LogInfo("Processed order after iteration", utils.FormatKeyValueLogs([][2]string{
		{"old status", fmt.Sprint(order.OrderStatus)},
		{"Order ID", fmt.Sprint(order.ID)},
		{"User Id", fmt.Sprint(order.UserID)},
		{"Status", fmt.Sprint(orderUpdate.Status)},
		{"EntryPrice", fmt.Sprint(orderUpdate.EntryPrice)},
		{"ClosePrice", fmt.Sprint(orderUpdate.ClosePrice)},
		{"TpValue", fmt.Sprint(orderUpdate.TpValue)},
		{"Pnl", fmt.Sprint(orderUpdate.Pnl)},
		{"Collateral", fmt.Sprint(orderUpdate.Collateral)},
		{"BalanceChange", fmt.Sprint(orderUpdate.BalanceChange)},
		{"EscrowBalanceChange", fmt.Sprint(orderUpdate.EscrowBalanceChange)},
	}))
}

// order type assumed (not checked) short/long
func processOrderTakeProfit(globalBorrowed, globalLiquidity, payout, closeFee *float64, order *db.OrderResponse, orderUpdate *db.OrderUpdate) {
	logrus.Info(fmt.Sprintf("processing %s take profit order", order.OrderType))

	typeMultiplier := map[bool]float64{true: 1, false: -1}[order.OrderType == "long"]

	value := order.TakeProfitCollateral * order.Leverage * (1 + (order.TakeProfitPrice-order.EntryPrice)*typeMultiplier/order.EntryPrice)
	*closeFee = order.TakeProfitCollateral * (dynamicLeverageFee(order.Leverage) + dynamicUtilizationFee(order.StartedAt, *globalBorrowed, *globalLiquidity))
	*payout += value - *closeFee - order.TakeProfitCollateral*(order.Leverage-1)
	order.TakeProfitValue = 0 // reset tpValue, indication of tp fill
	orderUpdate.Status = "pending"
	orderUpdate.EntryPrice = order.EntryPrice
	orderUpdate.ClosePrice = 0
	orderUpdate.TpValue = 0
	orderUpdate.Pnl += *payout
	orderUpdate.Collateral = order.Collateral
	orderUpdate.TakeProfitAt = time.Now()

	*globalBorrowed -= order.TakeProfitCollateral * (order.Leverage - 1)
	*globalLiquidity -= value

	orderUpdate.OrderGlobalUpdate.CurrentBorrowed -= order.TakeProfitCollateral * (order.Leverage - 1)
	orderUpdate.OrderGlobalUpdate.CurrentLiquidity -= value
	orderUpdate.OrderGlobalUpdate.TotalPnlProfits += *payout
	orderUpdate.OrderGlobalUpdate.TotalRevenue += *closeFee
	orderUpdate.OrderGlobalUpdate.TreasuryBalance += *closeFee * 0.1
	orderUpdate.OrderGlobalUpdate.TotalTreasuryProfits += *closeFee * 0.1
	orderUpdate.OrderGlobalUpdate.VaultBalance += *closeFee * 0.1
	orderUpdate.OrderGlobalUpdate.TotalVaultProfits += *closeFee * 0.1
	orderUpdate.OrderGlobalUpdate.TotalLiquidityRewards += *closeFee * 0.5
	orderUpdate.OrderGlobalUpdate.TotalStakeRewards += *closeFee * 0.3

	printProcessedOrder(*order, *orderUpdate)
}

func processOrderFill(globalBorrowed, globalLiquidity, borrowed, payout, closeFee *float64, order *db.OrderResponse, orderUpdate *db.OrderUpdate) {

	typeMultiplier := map[bool]float64{true: 1, false: -1}[order.OrderType == "long"]

	var value float64
	var liquidityChange float64
	if order.TakeProfitValue == 0 && order.TakeProfitCollateral != 0 {
		logrus.Info(fmt.Sprintf("processing %s fill order, after %s take profit", order.OrderType, order.OrderType))
		value = (order.Collateral - order.TakeProfitCollateral) * order.Leverage * (1 + (order.MaxPrice-order.EntryPrice)*typeMultiplier/order.EntryPrice)
		*closeFee = (order.Collateral - order.TakeProfitCollateral) * (dynamicLeverageFee(order.Leverage) + dynamicUtilizationFee(order.StartedAt, *globalBorrowed, *globalLiquidity))
		liquidityChange = order.Collateral - order.TakeProfitCollateral
		*globalBorrowed -= liquidityChange * (order.Leverage - 1)
		orderUpdate.OrderGlobalUpdate.CurrentBorrowed -= liquidityChange * (order.Leverage - 1)
		orderUpdate.TpValue = 0
		*borrowed = (order.Collateral - order.TakeProfitCollateral) * (order.Leverage - 1)

	} else { // if there is no tp collateral (implying not set)
		logrus.Info(fmt.Sprintf("processing %s fill order", order.OrderType))
		value = order.Collateral * order.Leverage * (1 + (order.MaxPrice-order.EntryPrice)*typeMultiplier/order.EntryPrice)
		*closeFee = order.Collateral * dynamicLeverageFee(order.Leverage)
		*globalBorrowed -= order.Collateral * (order.Leverage - 1)
		orderUpdate.OrderGlobalUpdate.CurrentBorrowed -= order.Collateral * (order.Leverage - 1)
		orderUpdate.TpValue = order.TakeProfitValue
		*borrowed = order.Collateral * (order.Leverage - 1)
	}

	*payout += value - *closeFee - *borrowed

	orderUpdate.Status = "filled"
	orderUpdate.EntryPrice = order.EntryPrice
	orderUpdate.ClosePrice = order.MaxPrice
	orderUpdate.Pnl += *payout
	orderUpdate.Collateral = order.Collateral
	*globalLiquidity -= value
	orderUpdate.OrderGlobalUpdate.CurrentLiquidity -= value
	orderUpdate.OrderGlobalUpdate.CurrentOrdersActive = -1
	orderUpdate.OrderGlobalUpdate.CurrentOrdersPending = -1
	orderUpdate.OrderGlobalUpdate.TotalOrdersFilled = 1
	orderUpdate.OrderGlobalUpdate.TotalPnlProfits += *payout
	orderUpdate.OrderGlobalUpdate.TotalRevenue += *closeFee
	orderUpdate.OrderGlobalUpdate.TreasuryBalance += *closeFee * 0.1
	orderUpdate.OrderGlobalUpdate.TotalTreasuryProfits += *closeFee * 0.1
	orderUpdate.OrderGlobalUpdate.VaultBalance += *closeFee * 0.1
	orderUpdate.OrderGlobalUpdate.TotalVaultProfits += *closeFee * 0.1
	orderUpdate.OrderGlobalUpdate.TotalLiquidityRewards += *closeFee * 0.5
	orderUpdate.OrderGlobalUpdate.TotalStakeRewards += *closeFee * 0.3

	printProcessedOrder(*order, *orderUpdate)
}

func processStopLoss(globalBorrowed, globalLiquidity, borrowed, payout, closeFee *float64, order *db.OrderResponse, orderUpdate *db.OrderUpdate) {

	typeMultiplier := map[bool]float64{true: 1, false: -1}[order.OrderType == "long"]

	var value float64
	var liquidityChange float64
	if order.TakeProfitValue == 0 && order.TakeProfitCollateral != 0 {
		logrus.Info(fmt.Sprintf("processing %s stop loss order, after %s take profit", order.OrderType, order.OrderType))
		value = (order.Collateral - order.TakeProfitCollateral) * (1 + order.Leverage*(order.StopLossPrice-order.EntryPrice)*typeMultiplier/order.EntryPrice)
		*closeFee = (order.Collateral - order.TakeProfitCollateral) * (dynamicLeverageFee(order.Leverage) + dynamicUtilizationFee(order.StartedAt, *globalBorrowed, *globalLiquidity))
		*borrowed = (order.Collateral - order.TakeProfitCollateral) * (order.Leverage - 1)
		*payout = value - *closeFee - *borrowed
		liquidityChange = order.Collateral - order.TakeProfitCollateral
		*closeFee += liquidityChange - value
		*globalBorrowed -= liquidityChange * (order.Leverage - 1)
		orderUpdate.OrderGlobalUpdate.CurrentBorrowed -= liquidityChange * (order.Leverage - 1)
		orderUpdate.TpValue = 0
	} else { // if there is no tp collateral (implying not set)
		logrus.Info(fmt.Sprintf("processing %s stop loss order", order.OrderType))
		value = order.Collateral * (1 + order.Leverage*(order.StopLossPrice-order.EntryPrice)*typeMultiplier/order.EntryPrice)
		*closeFee = order.Collateral * (dynamicLeverageFee(order.Leverage) + dynamicUtilizationFee(order.StartedAt, *globalBorrowed, *globalLiquidity))
		*borrowed = order.Collateral * (order.Leverage - 1)
		*payout = value - *closeFee - *borrowed
		*closeFee += order.Collateral - value
		*globalBorrowed -= order.Collateral * (order.Leverage - 1)
		orderUpdate.OrderGlobalUpdate.CurrentBorrowed -= order.Collateral * (order.Leverage - 1)
		orderUpdate.TpValue = order.TakeProfitValue
	}

	orderUpdate.Status = "stopped"
	orderUpdate.EntryPrice = order.EntryPrice
	orderUpdate.ClosePrice = order.StopLossPrice
	orderUpdate.Pnl -= (liquidityChange - *payout)
	orderUpdate.Collateral = order.Collateral
	*globalLiquidity -= value
	orderUpdate.OrderGlobalUpdate.CurrentLiquidity -= value
	orderUpdate.OrderGlobalUpdate.CurrentOrdersActive = -1
	orderUpdate.OrderGlobalUpdate.CurrentOrdersPending = -1
	orderUpdate.OrderGlobalUpdate.TotalOrdersStopped = 1
	orderUpdate.OrderGlobalUpdate.TotalPnlLosses -= (liquidityChange - *payout)
	orderUpdate.OrderGlobalUpdate.TotalRevenue += *closeFee
	orderUpdate.OrderGlobalUpdate.TreasuryBalance += *closeFee * 0.1
	orderUpdate.OrderGlobalUpdate.TotalTreasuryProfits += *closeFee * 0.1
	orderUpdate.OrderGlobalUpdate.VaultBalance += *closeFee * 0.1
	orderUpdate.OrderGlobalUpdate.TotalVaultProfits += *closeFee * 0.1
	orderUpdate.OrderGlobalUpdate.TotalLiquidityRewards += *closeFee * 0.5
	orderUpdate.OrderGlobalUpdate.TotalStakeRewards += *closeFee * 0.3

	printProcessedOrder(*order, *orderUpdate)
}

func processLiquidation(globalBorrowed, globalLiquidity, borrowed, payout, closeFee *float64, order *db.OrderResponse, orderUpdate *db.OrderUpdate) {

	typeMultiplier := map[bool]float64{true: 1, false: -1}[order.OrderType == "long"]

	var value float64
	var liquidityChange float64
	if order.TakeProfitValue == 0 && order.TakeProfitCollateral != 0 {
		logrus.Info(fmt.Sprintf("processing %s liquidate order, after %s take profit", order.OrderType, order.OrderType))
		value = (order.Collateral - order.TakeProfitCollateral) * (1 + order.Leverage*(order.LiquidationPrice-order.EntryPrice)*typeMultiplier/order.EntryPrice)
		*closeFee = (order.Collateral - order.TakeProfitCollateral) * (dynamicLeverageFee(order.Leverage) + dynamicUtilizationFee(order.StartedAt, *globalBorrowed, *globalLiquidity))
		*borrowed = (order.Collateral - order.TakeProfitCollateral) * (order.Leverage - 1)
		*payout = value - *closeFee - *borrowed
		liquidityChange = order.Collateral - order.TakeProfitCollateral
		*closeFee += liquidityChange - value
		if *closeFee > liquidityChange {
			*payout = 0
			*closeFee = liquidityChange
		}
		*globalBorrowed -= liquidityChange * (order.Leverage - 1)
		orderUpdate.OrderGlobalUpdate.CurrentBorrowed -= liquidityChange * (order.Leverage - 1)
		orderUpdate.TpValue = 0
	} else { // if there is no tp collateral (implying not set)
		logrus.Info(fmt.Sprintf("processing %s liquidate order", order.OrderType))
		value = order.Collateral * (1 + order.Leverage*(order.LiquidationPrice-order.EntryPrice)*typeMultiplier/order.EntryPrice)
		*closeFee = order.Collateral * (dynamicLeverageFee(order.Leverage) + dynamicUtilizationFee(order.StartedAt, *globalBorrowed, *globalLiquidity))
		*borrowed = order.Collateral * (order.Leverage - 1)
		*payout = value - *closeFee - *borrowed
		*closeFee += order.Collateral - value
		if *closeFee > order.Collateral {
			*payout = 0
			*closeFee = order.Collateral
		}
		*globalBorrowed -= order.Collateral * (order.Leverage - 1)
		orderUpdate.OrderGlobalUpdate.CurrentBorrowed -= order.Collateral * (order.Leverage - 1)

		orderUpdate.TpValue = order.TakeProfitValue
	}

	orderUpdate.Status = "liquidated"
	orderUpdate.EntryPrice = order.EntryPrice
	orderUpdate.ClosePrice = order.LiquidationPrice
	orderUpdate.Pnl -= (liquidityChange - *payout)
	orderUpdate.Collateral = order.Collateral
	*globalLiquidity -= order.TakeProfitCollateral
	orderUpdate.OrderGlobalUpdate.CurrentLiquidity -= order.TakeProfitCollateral
	orderUpdate.OrderGlobalUpdate.CurrentOrdersActive = -1
	orderUpdate.OrderGlobalUpdate.CurrentOrdersPending = -1
	orderUpdate.OrderGlobalUpdate.TotalOrdersLiquidated = 1
	orderUpdate.OrderGlobalUpdate.TotalPnlLosses -= (liquidityChange - *payout)
	orderUpdate.OrderGlobalUpdate.TotalOrdersFilled = 1
	orderUpdate.OrderGlobalUpdate.TotalRevenue += *closeFee
	orderUpdate.OrderGlobalUpdate.TreasuryBalance += *closeFee * 0.1
	orderUpdate.OrderGlobalUpdate.TotalTreasuryProfits += *closeFee * 0.1
	orderUpdate.OrderGlobalUpdate.VaultBalance += *closeFee * 0.1
	orderUpdate.OrderGlobalUpdate.TotalVaultProfits += *closeFee * 0.1
	orderUpdate.OrderGlobalUpdate.TotalLiquidityRewards += *closeFee * 0.5
	orderUpdate.OrderGlobalUpdate.TotalStakeRewards += *closeFee * 0.3

	printProcessedOrder(*order, *orderUpdate)
}

func processLimit(globalBorrowed, globalLiquidity *float64, order *db.OrderResponse, orderUpdate *db.OrderUpdate) {

	logrus.Info(fmt.Sprintf("processing %s limit order", order.OrderType))
	openFee := order.Collateral * (dynamicLeverageFee(order.Leverage) + dynamicUtilizationFee(order.StartedAt, *globalBorrowed, *globalLiquidity))

	order.OrderStatus = "pending"
	orderUpdate.Status = "pending"
	orderUpdate.EntryPrice = order.LimitPrice
	orderUpdate.ClosePrice = 0
	orderUpdate.Pnl = 0
	order.Collateral = order.Collateral - openFee
	orderUpdate.Collateral = order.Collateral - openFee
	*globalBorrowed += order.Collateral * (order.Leverage - 1)
	*globalLiquidity += order.Collateral
	orderUpdate.OrderGlobalUpdate.CurrentOrdersActive += 1
	orderUpdate.OrderGlobalUpdate.CurrentOrdersPending += 1
	orderUpdate.OrderGlobalUpdate.CurrentOrdersLimit -= 1
	orderUpdate.OrderGlobalUpdate.TotalBorrowed += order.Collateral * (order.Leverage - 1)
	orderUpdate.OrderGlobalUpdate.TotalOrdersActive += 1
	orderUpdate.OrderGlobalUpdate.TotalRevenue += openFee
	orderUpdate.OrderGlobalUpdate.TreasuryBalance += openFee * 0.1
	orderUpdate.OrderGlobalUpdate.TotalTreasuryProfits += openFee * 0.1
	orderUpdate.OrderGlobalUpdate.VaultBalance += openFee * 0.1
	orderUpdate.OrderGlobalUpdate.TotalVaultProfits += openFee * 0.1
	orderUpdate.OrderGlobalUpdate.TotalLiquidityRewards += openFee * 0.5
	orderUpdate.OrderGlobalUpdate.TotalStakeRewards += openFee * 0.3

	printProcessedOrder(*order, *orderUpdate)
}

// todo
// supabase client need to be called for the metrics and we need to select where to call to reduce the number of calls
// but we also need to consider the change in liquidity and borrow state from the incoming order changes
func processOrders(supabaseClient *supabase.Client, pairId string, priceMap []float64, minPrice, maxPrice float64) {
	orders, err := db.GetOrdersParsingRange(supabaseClient, pairId, minPrice, maxPrice)
	if err != nil {
		logrus.Error(fmt.Sprintf("could not fetch orders using pair id %v, minPrice %v, maxPrice %v: %v", pairId, minPrice, maxPrice, err))
	}

	globalBorrowed, globalLiquidity, err := getCurrentBorrowAndLiquidity(supabaseClient)

	orderUpdates_ := []db.OrderUpdate{}
	OrderGlobalUpdate_ := db.OrderGlobalUpdate{}
	if orders == nil || len(*orders) == 0 {
		logrus.Error(fmt.Sprintf("No orders returned for pair id %v, minPrice %v, maxPrice %v", pairId, minPrice, maxPrice))
		return
	}
	for _, order := range *orders {
		if order.OrderStatus == "unsigned" {
			continue
		}
		LogCreateOrderResponse(order)
		orderUpdate_ := db.OrderUpdate{}
		orderUpdate_.OrderID = order.ID
		orderUpdate_.UserID = order.UserID
		var payout float64
		var borrowed float64
		// add utilitization fee to order liquidation
		for _, markPrice := range priceMap {
			var closeFee float64
			// assume the order collateral is the exact, fees are already taken
			// collateral_ := order.Collateral * 0.99975
			if order.OrderType == "long" && order.EndedAt.IsZero() {
				if order.OrderStatus == "pending" {
					// profits
					if order.TakeProfitPrice <= markPrice && order.TakeProfitValue > 0 {
						processOrderTakeProfit(&globalBorrowed, &globalLiquidity, &payout, &closeFee, &order, &orderUpdate_)
					}
					if order.MaxPrice <= markPrice {
						processOrderFill(&globalBorrowed, &globalLiquidity, &borrowed, &payout, &closeFee, &order, &orderUpdate_)
						break
					}
					// losses
					if order.StopLossPrice >= markPrice && order.StopLossPrice > 0 {
						processStopLoss(&globalBorrowed, &globalLiquidity, &borrowed, &payout, &closeFee, &order, &orderUpdate_)
						break
					}
					// assume liquidations occur where value is non zero
					if order.LiquidationPrice >= markPrice || markPrice <= 0 {
						processLiquidation(&globalBorrowed, &globalLiquidity, &borrowed, &payout, &closeFee, &order, &orderUpdate_)
						break
					}
				} else if order.OrderStatus == "limit" {
					// assuming order.LimitPrice != 0
					if order.LimitPrice > order.EntryPrice && markPrice >= order.LimitPrice {
						processLimit(&globalBorrowed, &globalLiquidity, &order, &orderUpdate_)
					}
				} else {
					continue
				}
			} else if order.OrderType == "short" && order.EndedAt.IsZero() {
				if order.OrderStatus == "pending" {
					// profits
					if order.TakeProfitPrice >= markPrice && order.TakeProfitValue > 0 {
						processOrderTakeProfit(&globalBorrowed, &globalLiquidity, &payout, &closeFee, &order, &orderUpdate_)
					}
					if order.MaxPrice >= markPrice {
						processOrderFill(&globalBorrowed, &globalLiquidity, &borrowed, &payout, &closeFee, &order, &orderUpdate_)
						break
					}
					// losses
					if order.StopLossPrice <= markPrice && order.StopLossPrice > 0 {
						processStopLoss(&globalBorrowed, &globalLiquidity, &borrowed, &payout, &closeFee, &order, &orderUpdate_)
						break
					}
					// assume liquidations occur where value is non zero
					if order.LiquidationPrice <= markPrice {
						processLiquidation(&globalBorrowed, &globalLiquidity, &borrowed, &payout, &closeFee, &order, &orderUpdate_)
						break
					}
				} else if order.OrderStatus == "limit" {
					// assuming order.LimitPrice != 0
					if order.LimitPrice < order.EntryPrice && markPrice <= order.LimitPrice {
						processLimit(&globalBorrowed, &globalLiquidity, &order, &orderUpdate_)
					}
				} else {
					continue
				}
			}

			OrderGlobalUpdate_.CurrentBorrowed += orderUpdate_.OrderGlobalUpdate.CurrentBorrowed
			OrderGlobalUpdate_.CurrentLiquidity += orderUpdate_.OrderGlobalUpdate.CurrentLiquidity
			OrderGlobalUpdate_.CurrentOrdersActive += orderUpdate_.OrderGlobalUpdate.CurrentOrdersActive
			OrderGlobalUpdate_.CurrentOrdersLimit += orderUpdate_.OrderGlobalUpdate.CurrentOrdersLimit
			OrderGlobalUpdate_.CurrentOrdersPending += orderUpdate_.OrderGlobalUpdate.CurrentOrdersPending
			OrderGlobalUpdate_.TotalBorrowed += orderUpdate_.OrderGlobalUpdate.TotalBorrowed
			OrderGlobalUpdate_.TotalLiquidations += orderUpdate_.OrderGlobalUpdate.TotalLiquidations
			OrderGlobalUpdate_.TotalOrdersActive += orderUpdate_.OrderGlobalUpdate.TotalOrdersActive
			OrderGlobalUpdate_.TotalOrdersFilled += orderUpdate_.OrderGlobalUpdate.TotalOrdersFilled
			OrderGlobalUpdate_.TotalOrdersLimit += orderUpdate_.OrderGlobalUpdate.TotalOrdersLimit
			OrderGlobalUpdate_.TotalOrdersLiquidated += orderUpdate_.OrderGlobalUpdate.TotalOrdersLiquidated
			OrderGlobalUpdate_.TotalOrdersStopped += orderUpdate_.OrderGlobalUpdate.TotalOrdersStopped
			OrderGlobalUpdate_.TotalPnlLosses += orderUpdate_.OrderGlobalUpdate.TotalPnlLosses
			OrderGlobalUpdate_.TotalPnlProfits += orderUpdate_.OrderGlobalUpdate.TotalPnlProfits
			OrderGlobalUpdate_.TotalRevenue += orderUpdate_.OrderGlobalUpdate.TotalRevenue
			OrderGlobalUpdate_.TreasuryBalance += orderUpdate_.OrderGlobalUpdate.TreasuryBalance
			OrderGlobalUpdate_.TotalTreasuryProfits += orderUpdate_.OrderGlobalUpdate.TotalTreasuryProfits
			OrderGlobalUpdate_.VaultBalance += orderUpdate_.OrderGlobalUpdate.VaultBalance
			OrderGlobalUpdate_.TotalVaultProfits += orderUpdate_.OrderGlobalUpdate.TotalVaultProfits
			OrderGlobalUpdate_.TotalLiquidityRewards += orderUpdate_.OrderGlobalUpdate.TotalLiquidityRewards
			OrderGlobalUpdate_.TotalStakeRewards += orderUpdate_.OrderGlobalUpdate.TotalStakeRewards
		}

		orderUpdates_ = append(orderUpdates_, orderUpdate_)
	}

	if len(orderUpdates_) > 0 {
		if err := db.ProcessBatchOrders(supabaseClient, time.Now(), orderUpdates_, OrderGlobalUpdate_); err != nil {
			logrus.Error(fmt.Sprintf("Error processing batch orders: %v", err.Error()))
		}
	} else {
		logrus.Info("No order processed")
	}
}

func processPrices(supabaseClient *supabase.Client, priceMap map[string][]float64) {

	// Process the collected prices every 3 seconds
	for id, prices := range priceMap {
		var minPrice, maxPrice float64
		minPrice = 2000000 // set to a large enough value to cover all subvalues
		if len(prices) == 0 {
			continue
		}

		for _, price := range prices {
			if price < minPrice {
				minPrice = price
			}
			if price > maxPrice {
				maxPrice = price
			}
		}

		logrus.Warning(fmt.Sprintf("max price: %v, min price: %v", maxPrice, minPrice))
		logrus.Warning(fmt.Sprintf("price map: %v", priceMap))
		processOrders(supabaseClient, id, prices, minPrice, maxPrice)
	}
}

func SubscribeToPriceStream(supabaseClient *supabase.Client, url string, ids []string) {
	var markPriceMap = make(map[string][]float64)
	var mu sync.Mutex

	go func() {
		for {
			time.Sleep(3 * time.Second) // Periodically process prices

			mu.Lock()
			processPrices(supabaseClient, markPriceMap)
			for k := range markPriceMap {
				delete(markPriceMap, k)
			}
			mu.Unlock()
		}
	}()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Fatal("Error creating request:", err)
		return
	}

	q := req.URL.Query()
	for _, id := range ids {
		q.Add("ids[]", id)
	}
	req.URL.RawQuery = q.Encode()

	// Open a connection to the SSE stream
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Fatal("Error connecting to stream:", err)
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		// logrus.Info("Received line:", line)

		response, err := processSSELine(line)
		if err != nil {
			logrus.Error(err.Error())
		}

		parsedResponse, ok := response.(Response)
		if ok && err == nil {
			for _, priceUpdate := range parsedResponse.Parsed {
				// assume exponent always negative
				markPrice, _ := strconv.ParseFloat(priceUpdate.Price.Price, 64)
				exponent := priceUpdate.Price.Expo
				for i := int64(0); i < -int64(exponent); i++ {
					markPrice /= 10
				}

				mu.Lock()
				markPriceMap[priceUpdate.ID] = append(markPriceMap[priceUpdate.ID], markPrice)
				mu.Unlock()

				//logrus.Infof("Received Price Update - ID: %s, MarkPrice: %.6f", priceUpdate.ID, markPrice)

			}
		}
	}

	if err := scanner.Err(); err != nil {
		logrus.Error("Error reading from SSE stream:", err)
		// this should restart and log panic but fix later
	}

	select {}

}

/*
data:{
	"binary":{
		"encoding":"hex",
		"data":[
			"504e41550100000003b801000000040d00e351a260da76aec683d5b6179b6f8e817d8d7d0727e84ceaf6dac0deca05a4fc0ce55d0a5addaf6a871ed4aaacb785752db5d24b4a8bfa90158962a2517cc1f400024cb0a66e9f4adcbab8568d8f6032360dbf93def497ccbd02402a8abb7169e3ed1de3fa827f86ba7b3992679dcd7c6136719121254ce96d78309e1bcfe263d00400035876e6d45911aa1c73028be8d41a5412701732c3f7005a86440500310b3eaa0f635a3dfcd9d3757a31f78ad5b3e7c4c7f1561b04c7e1af3d360dee650030ade100049613ef5f6ddbe9271ffb7bd41ee781cbe8c55180bca45b794318a9d199e229f32bf8f44b8da144430d09826406861d13cbc8b9a53f50bd30a616c74d6849d7c101061cd7547a31bc897fc7e2443e9cbee426ca992048d064712f5596a954fa5689374741e6ec2842532fdd480604d34292336c6744150a144658bd9092bf9e87434b000852205ad844251a1349b6fbe1e7d9f27b8681d0628f47f7769142cc421b0b890433dc5ed49285b2717169edfa5b07544f8d53eec122079b5ceb22ccefbd4b63ef000a231ae39d67ed1f15b199ea34018660bfebb4e8f275037667f594868c74a64b362b548d8649136b230159342a6b2761bbca4c4996a6ce21d3504f1ff935bd670c010bb9e2a7ec75a714a5f77760968e50f32d9c8a70295f265e71aff8dc2c52a2d1d73753f36979cf2b2ec788961743b729976675305e7725a1a4862741bc7f9aea33000c5554b32e7c51e4f5383319af3345ab57a434bb6df4251d126ba6b8667fec55366e83b8acc085bbd5c75bcc1d99338f4bf1c75aaf08fe9d007593042cc183aedd000d717e0a95bc5f4db1a8f5a0795258f030b5190d059d6cff387b1b0783cb1012242c27f23aee7ccd09729a5feb72fb22d4b8bccd9d25bc5d6f4e49892026975031000f18fa460fd4967520c9976682c721d9dd0b532bf6d6e286af38998d41975f1eb4562be5a05911d14def552318667eb0d9a8e0e8b9307be23bf112aa700b446030011069b26caf4a5b2afad642c2c2e6efb9d5ee57dc37f39988f22fe7dc2f8001f6fb3fa03dfca0b7bf03dbbf9b02b2885a4785186d552de621b826ed4ef16a87741d00118e4b2346adcfa23e76e089b74753d49fc27e5c041be6c627959d5a51a5d9b8f706a943fa815ebac68606c9807bfa3fa4daec36a4d5a7aebb4fa02ce3229b8a73016773daea00000000001ae101faedac5851e32b9b23b5f9411a8c2bac4aae3ed4dd7b811dd1a72ea4aa7100000000062a89a7014155575600000000000b376d6500002710d86affe1de75b312f0a6f7383f0b0732330ea7c601005500e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b430000088e17d0360600000002305d11befffffff8000000006773daea000000006773daea0000088ac179e42000000001bc5439500b31be0d91a82eae8b2c080cf02a88e5f96a22a52fa70923bc214a8d0f3b0b4d5f4251d4fd655a2cc50ab6b8017b7ea1b2dea695791776ca14d5331491993a0e8576e8e4e5991983c3be76742f55e851f5d47879b0c8728abb48ff3042a66f2ab688922740cb33732f7b765f1c96c91a99a5a7ddb339c155f10ccd5fe1471fd63c6e8c59271b06cd7027380aee4d1a5605dc190821952944fc9168c58f7a5b04d2271cf84cf207de2550df8db4984bb7e473c269e452cf31930a62a0040e4edb21c6a6740d2619be0f0a92c4d15a3355ecc0050863529ff14422e9fbda"
		]
	},
	"parsed":[{
		"id":"e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43",
		"price":{
			"price":"9406377899526",
			"conf":"9401340350",
			"expo":-8,
			"publish_time":1735645930
		},
		"ema_price":{
			"price":"9392044500000",
			"conf":"7454603600",
			"expo":-8,
			"publish_time":1735645930
		},
		"metadata":{
			"slot":188181861,
			"proof_available_time":1735645932,
			"prev_publish_time":1735645930
		}
	}]
}
*/
