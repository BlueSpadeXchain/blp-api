package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/supabase-community/postgrest-go"
	"github.com/vercel/go-bridge/go/bridge"
	"golang.org/x/crypto/sha3"
)

type OrderRequest struct {
	Signer    string    `json:"signer"`
	CreatedOn string    `json:"createdOn"`
	ChainId   string    `json:"chainId"`
	Order     OrderData `json:"order"`
	MessageId string    `json:"messageId"`
	Signature string    `json:"signature"`
	Nonce     int64     `json:"nonce"`
}

// actually we don't need to store the hash or signature on the order, only for message validity

type OrderData struct {
	OrderId          string `json:"orderId"`
	NetValue         string `json:"netValue"`
	Amount           string `json:"amount"`
	Collateral       string `json:"collateral"`
	MarkPrice        string `json:"markPrice"`
	EntryPrice       string `json:"EntryPrice"`
	LiquidationPrice string `json:"liquidationPrice"`
}

type OrderResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Data    *OrderData `json:"data,omitempty"` // this will assign a new orderId if not already applied
}

func validateSignature(order OrderRequest) bool {
	hash := sha3.NewLegacyKeccak256()
	rlp.Encode(hash, []interface{}{order.Signer, order.ChainId, order.MessageId, order.Order})
	expectedHash := hash.Sum(nil)
	sigPublicKey, err := crypto.SigToPub(expectedHash, []byte(order.Signature))
	if err != nil {
		return false
	}
	return crypto.PubkeyToAddress(*sigPublicKey).Hex() == order.Signer
}

func Handler(w http.ResponseWriter, r *http.Request) {
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
	order := OrderData{
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
	bridge.Start(Handler)
}

// database needs the fields
// data
// values
// hash
// signer
// and when a submission is made we need to have a valid signature to go with it
// and when an edit is made we need a valid signature, and then delete and replace the position
//		the latter might be weird since we are actually a perp (not limit orders)
