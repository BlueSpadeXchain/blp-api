package orderHandler

import (
	"fmt"
	"reflect"

	"github.com/BlueSpadeXchain/blp-api/pkg/db"
	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
	"github.com/sirupsen/logrus"
)

func ParseStructToKeyValue(response interface{}, prefix string) [][2]string {
	var keyValuePairs [][2]string

	val := reflect.ValueOf(response)
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // Dereference pointer if necessary.
	}

	if val.Kind() != reflect.Struct {
		return keyValuePairs // Return empty for non-struct types.
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name
		fieldValue := val.Field(i)

		// Construct the full key name (including prefix for nested structs).
		fullName := fieldName
		if prefix != "" {
			fullName = fmt.Sprintf("%s.%s", prefix, fieldName)
		}

		// Check if the field is a nested struct.
		if fieldValue.Kind() == reflect.Struct {
			// Recursively process the nested struct.
			nestedKeyValuePairs := ParseStructToKeyValue(fieldValue.Interface(), fullName)
			keyValuePairs = append(keyValuePairs, nestedKeyValuePairs...)
		} else {
			// Convert the field value to a string representation.
			fieldValueStr := fmt.Sprintf("%v", fieldValue.Interface())
			keyValuePairs = append(keyValuePairs, [2]string{fullName, fieldValueStr})
		}
	}

	return keyValuePairs
}

func LogResponse(url string, response interface{}) {
	if logrus.GetLevel() < logrus.InfoLevel {
		return
	}

	keyValueData := ParseStructToKeyValue(response, "")
	message := utils.FormatKeyValueLogs(keyValueData)

	logrus.Info(fmt.Sprintf(
		"URL request: \033[1m%s\033[0m:\n%s",
		url,
		message,
	))
}

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
		"Request: \033[1m%s\033[0m\n"+
			"ID: %s\n"+
			"UserID: %s\n"+
			"OrderType: %s\n"+
			"Leverage: %.2f\n"+
			"Status: %s\n"+
			"EntryPrice: %.2f\n"+
			"LiqPrice: %.2f\n"+
			"Collateral: %.2f\n"+
			"CreatedAt: %s\n"+
			"EndedAt: %s\n",
		url,
		response.ID,
		response.UserID,
		response.OrderType,
		response.Leverage,
		response.Status,
		response.EntryPrice,
		response.LiqPrice,
		response.Collateral,
		response.CreatedAt,
		response.EndedAt,
	)

	logrus.Info(message)
}
