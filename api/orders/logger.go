package orderHandler

import (
	"fmt"

	"github.com/BlueSpadeXchain/blp-api/pkg/db"
	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
	"github.com/sirupsen/logrus"
)

func LogInfo(title string, message string) {
	utils.LogInfo(title, message)
}

func LogError(message string, errStr string) {
	utils.LogError(message, errStr)
}

func LogBeforeCreateOrderResponse(
	userId string,
	pair string,
	pairId string,
	collateral float64,
	entryPrice float64,
	markPrice float64,
	liqPrice float64,
	leverage float64,
	positionType string,
	status string,
) {
	if logrus.GetLevel() < logrus.InfoLevel {
		return
	}

	order := [][2]string{
		{"user_id", userId},
		{"pair", pair},
		{"pairId", pairId},
		{"collateral", fmt.Sprintf("%v", collateral)},
		{"entry_price", fmt.Sprintf("%v", entryPrice)},
		{"mark_price", fmt.Sprintf("%v", markPrice)},
		{"liq_price", fmt.Sprintf("%v", liqPrice)},
		{"leverage", fmt.Sprintf("%v", leverage)},
		{"position_type", positionType},
		{"status", status},
	}

	LogInfo("Unsigned order request created successfully: %+v", utils.FormatKeyValueLogs(order))
}

func LogCreateOrderResponse(url string, response db.OrderResponse) {
	if logrus.GetLevel() < logrus.InfoLevel {
		return
	}

	message := fmt.Sprintf(
		"Order Returned: \033[1m%s\033[0m\n"+
			"ID:                      %s\n"+
			"UserID:                  %s\n"+
			"Order Type:              %s\n"+
			"Leverage:                %.2f\n"+
			"Pair ID:                 %s\n"+
			"Order Status:            %s\n"+
			"Collateral:              %.2f\n"+
			"Entry Price:             %.2f\n"+
			"Liquidation Price:       %.2f\n"+
			"Limit Order Price:       %.2f\n"+
			"Max Price:               %.2f\n"+
			"Max Value:               %.2f\n"+
			"Stop Loss Price:         %.2f\n"+
			"Take Profit Price:       %.2f\n"+
			"Take Profit Value:       %.2f\n"+
			"Take Profit Collateral:  %.2f\n"+
			"Created At:              %s\n"+
			"Signed At:               %s\n"+
			"Started At:              %s\n"+
			"Ended At:                %s\n",
		url,
		response.ID,
		response.UserID,
		response.OrderType,
		response.Leverage,
		response.PairId,
		response.OrderStatus,
		response.Collateral,
		response.EntryPrice,
		response.LiquidationPrice,
		response.MaxPrice,
		response.MaxValue,
		response.LimitPrice,
		response.StopLossPrice,
		response.TakeProfitPrice,
		response.TakeProfitValue,
		response.TakeProfitCollateral,
		response.CreatedAt,
		response.SignedAt,
		response.StartedAt,
		response.EndedAt,
	)

	logrus.Info(message)
}
