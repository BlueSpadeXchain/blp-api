package userHandler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

func getNetworkAddresses() (string, string, error) {
	var bluAddress, usdcAddress string

	// Check if mainnet is enabled
	mainnetEnabled := os.Getenv("MAINNET_ENABLED") == "true"

	if mainnetEnabled {
		bluAddress = os.Getenv("MAINNET_BLU")
		usdcAddress = os.Getenv("MAINNET_USDC")
	} else {
		bluAddress = os.Getenv("TESTNET_BLU")
		usdcAddress = os.Getenv("TESTNET_USDC")
	}

	fmt.Printf("\n mainnet_endable: %v", mainnetEnabled)

	// Validate addresses
	if !common.IsHexAddress(bluAddress) {
		return "", "", fmt.Errorf("invalid BLU address: %s", bluAddress)
	}
	if !common.IsHexAddress(usdcAddress) {
		return "", "", fmt.Errorf("invalid USDC address: %s", usdcAddress)
	}

	return bluAddress, usdcAddress, nil
}

func ConvertStructToQuery(params interface{}) (string, error) {
	fmt.Printf("\n params: %v", params)
	v := reflect.ValueOf(params)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("expected a pointer to a struct")
	}

	query := url.Values{}

	t := v.Elem().Type()
	for i := 0; i < v.Elem().NumField(); i++ {
		field := v.Elem().Field(i)
		fieldName := t.Field(i).Tag.Get("query")

		if field.IsZero() {
			continue
		}

		query.Set(fieldName, fmt.Sprintf("%v", field.Interface()))
	}
	fmt.Printf("\n params end: %v", query.Encode())
	return query.Encode(), nil
}

// func sendRequest(api, query, body string) {
// 	url := fmt.Sprintf("%v?query=%v&%v", api, query, body)
// 	logrus.Info("Request forwarded: ", url)
// 	req, err := http.Get(url)
// 	if err != nil {
// 		logrus.Error("Request creation error: ", err)
// 		return
// 	}

// 	defer req.Body.Close()
// }

func sendRequest(api, query, body string) {
	url := fmt.Sprintf("%v?query=%v&%v", api, query, body)
	logrus.Info("Request forwarded: ", url)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error("Request creation error: ", err)
		return
	}

	// Set a whitelisted origin
	// This needs to be one of the origins that the receiving API accepts
	whitelistedOrigin := "https://blp-api-vercel.vercel.app"
	req.Header.Set("Origin", whitelistedOrigin)
	logrus.Info("Setting Origin header to: ", whitelistedOrigin)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("Request error: ", err)
		return
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("Failed to read response body: ", err)
		return
	}

	logrus.Infof("Response Status: %d", resp.StatusCode)
	logrus.Infof("Response Body: %s", string(respBody))
}
