package userHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BlueSpadeXchain/blp-api/pkg/db"
	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/supabase-community/supabase-go"
)

// handler handles signed and unsigned user requests
// login (maybe?), user data, escrow (withdrawals/deposits), staking
// 	deposits will have to be signed by our blockchain listener (logged regardless, to catch missed payout history)

func Handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("\nRecovered from panic: %v", rec)

			supabaseUrl := os.Getenv("SUPABASE_URL")
			supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
			supabaseClient, err := supabase.NewClient(supabaseUrl, supabaseKey, nil)
			if err == nil {
				logErr := db.LogPanic(supabaseClient, fmt.Sprintf("%v", rec), nil)
				if logErr != nil {
					log.Printf("\nFailed to log panic to Supabase: %v", logErr)
				}
			} else {
				log.Printf("\nFailed to create Supabase client for panic logging: %v", err)
			}

			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	handlerWithCORS := utils.EnableCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		var response interface{}
		var err error
		supabaseUrl := os.Getenv("SUPABASE_URL")
		supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
		supabaseClient, err := supabase.NewClient(supabaseUrl, supabaseKey, nil)
		if err != nil {
			http.Error(w, "Failed to create Supabase client", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		switch query.Get("query") {
		case "deposit":
			response, err = DespositRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "stake":
			response, err = StakeRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "eoa-stake":
			response, err = EoaStakeRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "unstake":
			response, err = UnstakeRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "user-data":
			response, err = UserDataRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "get-user-by-user-id":
			response, err = GetUserByUserIdRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
		case "get-user-by-user-address":
			response, err = GetUserByUserAddressRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "get-deposits-by-user-id":
			response, err = GetDepositsByUserIdRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "get-deposits-by-user-address":
			response, err = GetDepositsByUserAddressRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		// todo
		case "get-withdrawals-by-user-id":
			response, err = GetDepositsByUserIdRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		// todo
		case "get-withdrawals-by-user-address":
			response, err = GetDepositsByUserAddressRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "get-stakes-by-user-id":
			response, err = GetStakesByUserIdRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "get-stakes-by-user-address":
			response, err = GetStakesByUserAddressRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "withdraw":
			response, err = UnsignedWithdrawRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "sign-withdraw":
			response, err = SignedWithdrawRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		// case "get-orders-by-user-id":
		// 	response, err = GetOrdersByUserIdRequest(r, supabaseClient)
		// 	HandleResponse(w, r, supabaseClient, response, err)
		// 	return
		// case "get-orders-by-user-address":
		// 	response, err = GetOrdersByUserAddressRequest(r, supabaseClient)
		// 	HandleResponse(w, r, supabaseClient, response, err)
		// return
		case "add-wallet": // both need some type of connection token
			response, err = AddAuthorizedWalletRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "remove-wallet":
			response, err = RemoveAuthorizedWalletRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "test":
			callWithdrawalAPI()
			return
		default:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(utils.ErrMalformedRequest("Invalid query parameter"))
			return
		}
	}))

	handlerWithCORS.ServeHTTP(w, r)
}

func HandleResponse(w http.ResponseWriter, r *http.Request, supabaseClient *supabase.Client, response interface{}, err error) {
	if err != nil {
		if logErr := db.LogError(supabaseClient, err, r.URL.Query().Get("query"), response); logErr != nil {
			fmt.Printf("Failed to log error: %v\n", logErr.Error())
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}
}

func callWithdrawalAPI() {
	apiURL := "https://blp-api-vercel-production.up.railway.app/api/withdrawal"
	logrus.Info("Making direct API call to: ", apiURL)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		logrus.Error("Request creation error: ", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Create and execute the request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("API request failed: ", err)
		return
	}
	defer resp.Body.Close()

	// Read and log the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("Failed to read response body: ", err)
		return
	}

	logrus.Infof("API Response Status: %d", resp.StatusCode)
	logrus.Infof("API Response Body: %s", string(respBody))
}
