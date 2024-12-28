package orderHandler

func validateSignature() bool {
	// 	hash := sha3.NewLegacyKeccak256()
	// 	rlp.Encode(hash, []interface{}{order.Signer, order.ChainId, order.MessageId, order.Order})
	// 	expectedHash := hash.Sum(nil)
	// 	sigPublicKey, err := crypto.SigToPub(expectedHash, []byte(order.Signature))
	// 	if err != nil {
	// 		return false
	// 	}
	// 	return crypto.PubkeyToAddress(*sigPublicKey).Hex() == order.Signer
	// }

	// func Handler2(w http.ResponseWriter, r *http.Request) {
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

	//	orderRes := OrderResponse{
	//		Success: true,
	//		Data:    &order,
	//	}
	//
	// json.NewEncoder(w).Encode(orderRes)
	return true
}

// database needs the fields
// data
// values
// hash
// signer
// and when a submission is made we need to have a valid signature to go with it
// and when an edit is made we need a valid signature, and then delete and replace the position
//		the latter might be weird since we are actually a perp (not limit orders)
