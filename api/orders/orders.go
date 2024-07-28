package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/supabase/postgrest-go"
	"github.com/vercel/go-bridge/go/bridge"
	"golang.org/x/crypto/sha3"
)

type OrderRequest struct {
	Signer    string `json:"signer"`
	ChainId   int    `json:"chainId"`
	MessageId string `json:"messageId"`
	OrderData string `json:"orderData"`
	Signature string `json:"signature"`
}

type OrderResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    *Order `json:"data,omitempty"`
}

type Order struct {
	ID        string `json:"id"`
	Signer    string `json:"signer"`
	ChainId   int    `json:"chainId"`
	MessageId string `json:"messageId"`
	OrderData string `json:"orderData"`
	Status    string `json:"status"`
}

func validateSignature(order OrderRequest) bool {
	hash := sha3.NewLegacyKeccak256()
	rlp.Encode(hash, []interface{}{order.Signer, order.ChainId, order.MessageId, order.OrderData})
	expectedHash := hash.Sum(nil)
	sigPublicKey, err := crypto.SigToPub(expectedHash, []byte(order.Signature))
	if err != nil {
		return false
	}
	return crypto.PubkeyToAddress(*sigPublicKey).Hex() == order.Signer
}

func handler(w http.ResponseWriter, r *http.Request) {
	var orderReq OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !validateSignature(orderReq) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	client := postgrest.NewClient("https://arlgbqlmnvdeglgwtxic.supabase.co", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImFybGdicWxtbnZkZWdsZ3d0eGljIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MjIxMjgwMzcsImV4cCI6MjAzNzcwNDAzN30.0rs1ghN-Nt31Hjx5IbaXwN9c4wX38FO0tvC5b9qWUaA")

	// Check if user has sufficient funds
	resp, err := client.From("accounts").Select("*").Eq("signer", orderReq.Signer).Single().Execute()
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "User account not found", http.StatusNotFound)
		return
	}

	var userFunds map[string]interface{}
	if err := json.Unmarshal(resp.Body, &userFunds); err != nil {
		http.Error(w, "Failed to parse user account data", http.StatusInternalServerError)
		return
	}

	// Add the order to the database
	order := Order{
		Signer:    orderReq.Signer,
		ChainId:   orderReq.ChainId,
		MessageId: orderReq.MessageId,
		OrderData: orderReq.OrderData,
		Status:    "pending",
	}

	_, err = client.From("orders").Insert(order, false, "", "").Execute()
	if err != nil {
		http.Error(w, "Failed to add order", http.StatusInternalServerError)
		return
	}

	orderRes := OrderResponse{
		Success: true,
		Data:    &order,
	}
	json.NewEncoder(w).Encode(orderRes)
}

func main() {
	bridge.Start(handler)
}
