package orderHandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BlueSpadeXchain/blp-api/pkg/db"
	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/supabase-community/supabase-go"
)

// user sends signed request
// request is validated and then added to queue
// we have order requests and close requests only

// if price is A -> B
// / user makes order for 100 USD and liquidation at 150 usd or 50 usd
// then the price gets to 80 usd, how much money does his have, and how much money does he need to add to stablize his position
// remember this is in the context of I can only add/remove via another order request, or close
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
		// case "open-order": // need to migrate
		// 	response, err = OpenOrderRequest(r)
		// 	HandleResponse(w, r, supabaseClient, response, err)
		// 	return
		// case "close-order":
		// 	response, err = CloseOrderRequest(r)
		// 	HandleResponse(w, r, supabaseClient, response, err)
		// 	return
		case "create-order-unsigned": // returns order with uuid + hash to sign
			response, err = UnsignedOrderRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "create-order-signed": // order must include order uuid, and signature
			response, err = SignedOrderRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "get-order-by-id":
			response, err = GetOrderByIdRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		// case "get-orders-by-id":
		// 	response, err = GetOrdersByIdRequest(r, supabaseClient)
		// 	HandleResponse(w, r, supabaseClient, response, err)
		// 	return
		case "get-orders-by-user-id":
			response, err = GetOrdersByUserIdRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "get-orders-by-user-address":
			response, err = GetOrdersByUserAddressRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "get-orders":
			response, err = GetOrderByIdRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "close-order":
			response, err = CloseOrderRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "cancel-order":
			response, err = CancelOrderRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
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
			utils.LogError("Failed to log error", logErr.Error())
		}

		logrus.Error(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}
}
