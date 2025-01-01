package escrow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EventData struct {
	From   string `json:"from"`
	Value  string `json:"value"`
	TxHash string `json:"tx_hash"`
}

func StartListener(rpcURL string) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	defer client.Close()

	// Load ABI
	abiFile, err := os.ReadFile("./escrow/abi/contract.json")
	if err != nil {
		log.Fatalf("Failed to load ABI file: %v", err)
	}

	contractABI, err := abi.JSON(string(abiFile))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	// Define contract address and topics
	contractAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	// Subscribe to events
	logs := make(chan ethereum.LogFilterer)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to events: %v", err)
	}
	defer sub.Unsubscribe()

	fmt.Println("Listening for events...")

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Subscription error: %v", err)
		case vLog := <-logs:
			fmt.Println("Event received!")
			handleEvent(contractABI, vLog)
		}
	}
}

func handleEvent(contractABI abi.ABI, vLog ethereum.Log) {
	// Parse the event log
	event := struct {
		From  common.Address
		Value *big.Int
	}{}

	err := contractABI.UnpackIntoInterface(&event, "DepositEvent", vLog.Data)
	if err != nil {
		log.Printf("Failed to unpack event data: %v", err)
		return
	}

	// Format and send the event to the backend API
	data := EventData{
		From:   event.From.Hex(),
		Value:  event.Value.String(),
		TxHash: vLog.TxHash.Hex(),
	}

	sendToBackendAPI(data)
}

func sendToBackendAPI(event EventData) {
	backendURL := os.Getenv("BACKEND_API_URL")
	if backendURL == "" {
		log.Println("BACKEND_API_URL is not set")
		return
	}

	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event data: %v", err)
		return
	}

	resp, err := http.Post(backendURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Printf("Failed to send event to backend: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Event sent to backend successfully. Status code: %d\n", resp.StatusCode)
}
