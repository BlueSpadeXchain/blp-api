package orderHandler

// func Handler(w http.ResponseWriter, r *http.Request) {
// 	var orderReq OrderRequest
// 	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if !validateSignature(orderReq) {
// 		http.Error(w, "Invalid signature", http.StatusUnauthorized)
// 		return
// 	}

// 	client := postgrest.NewClient("https://arlgbqlmnvdeglgwtxic.supabase.co", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImFybGdicWxtbnZkZWdsZ3d0eGljIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MjIxMjgwMzcsImV4cCI6MjAzNzcwNDAzN30.0rs1ghN-Nt31Hjx5IbaXwN9c4wX38FO0tvC5b9qWUaA")

// 	// Check if user has sufficient funds
// 	resp, err := client.From("accounts").Select("*").Eq("signer", orderReq.Signer).Single().Execute()
// 	if err != nil || resp.StatusCode != http.StatusOK {
// 		http.Error(w, "User account not found", http.StatusNotFound)
// 		return
// 	}

// 	var userFunds map[string]interface{}
// 	if err := json.Unmarshal(resp.Body, &userFunds); err != nil {
// 		http.Error(w, "Failed to parse user account data", http.StatusInternalServerError)
// 		return
// 	}

// 	// Add the order to the database
// 	order := OrderData{
// 		Signer:    orderReq.Signer,
// 		ChainId:   orderReq.ChainId,
// 		MessageId: orderReq.MessageId,
// 		OrderData: orderReq.OrderData,
// 		Status:    "pending",
// 	}

// 	_, err = client.From("orders").Insert(order, false, "", "").Execute()
// 	if err != nil {
// 		http.Error(w, "Failed to add order", http.StatusInternalServerError)
// 		return
// 	}

// 	orderRes := OrderResponse{
// 		Success: true,
// 		Data:    &order,
// 	}
// 	json.NewEncoder(w).Encode(orderRes)
// }

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/BlueSpadeXchain/blp-api/pkg/db"
	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
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
		case "request-order": // similar to data from unsigned-data, but no op data
			response, err = OrderRequest(r)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "get-orders": // similar to data from unsigned-data, but no op data
			response, err = GetOrdersRequest(r)
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
