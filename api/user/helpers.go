package userHandler

import (
	"fmt"
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

func sendRequest(api, query, body string, r *http.Request) {
	url := fmt.Sprintf("%v?query=%v&%v", api, query, body)
	logrus.Info("Request forwarded: ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error("Request creation error: ", err)
		return
	}

	// Get the host from the incoming request
	if r != nil {
		// Construct the full origin with scheme and host
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		origin := fmt.Sprintf("%s://%s", scheme, r.Host)

		// Set the Origin header dynamically
		req.Header.Set("Origin", origin)
		logrus.Info("Setting Origin header to: ", origin)
	} else {
		// Fallback if no request context is available
		req.Header.Set("Origin", "http://localhost:8080")
		logrus.Info("Using default Origin header: http://localhost:8080")
	}

	go func() {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logrus.Error("Request error: ", err)
			return
		}
		defer resp.Body.Close()
	}()
}
