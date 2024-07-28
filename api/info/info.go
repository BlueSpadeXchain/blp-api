package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Version struct {
	Version string `json:"version"`
}

// Error data structure
type Error struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
}

// Asset info
type AssetInfo struct {
	ChainId   string `json:"chain-id"`
	ChainName string `json:"chain-name"`
	Address   string `json:"address"`
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	Decimals  string `json:"decimals"`
}

// Balance Response
type BalanceInfoResponse struct {
	Address string    `json:"address"`
	Balance string    `json:"balance"`
	Info    AssetInfo `json:"info"`
}

// Balances Response
type BalancesResponse struct {
	Address0 string    `json:"address0"`
	Balance0 string    `json:"balance0"`
	Info0    AssetInfo `json:"info0"`
	Address1 string    `json:"address1"`
	Balance1 string    `json:"balance1"`
	Info1    AssetInfo `json:"info1"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	handlerWithCORS := EnableCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// capture query variables
		query := r.URL.Query()
		switch query.Get("query") {
		case "version":
			response, err := http.Get("https://lubab-api-vercel.vercel.app/api/info?query=version")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer response.Body.Close()
			var crosschainVersion Version
			if err := json.NewDecoder(response.Body).Decode(&crosschainVersion); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			version := Version{Version: "BLP Order API v0.0.1"}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(version); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			version := "Hello, World!"
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(version); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}))

	handlerWithCORS.ServeHTTP(w, r)
}

func errUnsupportedChain(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(&Error{
		Code:    0,
		Message: "Chain not currently supported",
	})
}

func errMalformedRequest(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(&Error{
		Code:    400,
		Message: "Malformed request",
	})
}

func errInternal(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(&Error{
		Code:    500,
		Message: "Internal server error",
	})
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		fmt.Printf("Method: %s, URL: %s", r.Method, r.URL)

		next.ServeHTTP(w, r)
	})
}
