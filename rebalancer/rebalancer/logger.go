package rebalancer

import (
	"fmt"

	"github.com/BlueSpadeXchain/blp-api/rebalancer/pkg/db"
	"github.com/sirupsen/logrus"
)

func LogCreateOrderResponse2(response db.OrderResponse2) {
	if logrus.GetLevel() < logrus.InfoLevel {
		return
	}

	message := fmt.Sprintf(
		"Order Returned: \033[1m\033[0m\n"+
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
