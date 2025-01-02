package escrow

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

type EventData struct {
	From   string `json:"from"`
	Value  string `json:"value"`
	TxHash string `json:"tx_hash"`
}

func StartListener(rpcURL string, chainId string) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	defer client.Close()

	escrowAddress, err := getAddress(chainId)
	if err != nil {
		logrus.Fatal("Failed to parse contract address: ", err.Error())
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{escrowAddress},
	}

	escrowABI, _ := abi.JSON(strings.NewReader(escrowContractABI))
	if err != nil {
		logrus.Fatal("Failed to parse contract ABI: ", err.Error())
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal("Failed to subscribe to logs:", err)
	}

	DepositEventSig := escrowABI.Events["DepositEvent"].ID
	StakingDepositEventSig := escrowABI.Events["StakingDepositEvent"].ID
	BurnRequestEventSig := escrowABI.Events["BurnRequestEvent"].ID
	// fmt.Printf("\n DepositEventSig: %v", DepositEventSig)
	// fmt.Printf("\n StakingDepositEventSig: %v", StakingDepositEventSig)
	// fmt.Printf("\n BurnRequestEventSig: %v", BurnRequestEventSig)

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println("BlockHash:", vLog.BlockHash.Hex())
			fmt.Println("BlockNumber:", vLog.BlockNumber)
			fmt.Println("TxHash:", vLog.TxHash.Hex())

			// Handle events based on signature hash
			switch vLog.Topics[0] {
			case DepositEventSig:
				var event struct {
					Sender       common.Address
					Account      common.Address
					Nonce        *big.Int
					AssetAddress common.Address
					AssetAmount  *big.Int
				}
				err := escrowABI.UnpackIntoInterface(&event, "DepositEvent", vLog.Data)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("DepositEvent:", event)

			case StakingDepositEventSig:
				var event struct {
					Sender       common.Address
					Account      common.Address
					Nonce        *big.Int
					AssetAddress common.Address
					AssetAmount  *big.Int
				}
				err := escrowABI.UnpackIntoInterface(&event, "StakingDepositEvent", vLog.Data)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("StakingDepositEvent:", event)

			case BurnRequestEventSig:
				var event struct {
					Nonce        *big.Int
					AssetAddress common.Address
					AssetAmount  *big.Int
				}
				err := escrowABI.UnpackIntoInterface(&event, "BurnRequestEvent", vLog.Data)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("BurnRequestEvent:", event)

			default:
				fmt.Println("Unknown event signature:", vLog.Topics[0].Hex())
			}
		}
	}
}

// now need to make request to backend to add deposits
// BlockNumber: 14
// TxHash: 0xe9f1fe395e55ca3037a5d248b87de7f5c124a2f558a9f8493ce6fa6fe9c8e9fd
// DepositEvent: {0x70997970C51812dc3A010C7d01b50e0d17dc79C8 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 4 0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0 100000000000000000000}

// func handleEvent(contractABI abi.ABI, vLog ethereum.Log) {
// 	// Parse the event log
// 	event := struct {
// 		From  common.Address
// 		Value *big.Int
// 	}{}

// 	err := contractABI.UnpackIntoInterface(&event, "DepositEvent", vLog.Data)
// 	if err != nil {
// 		log.Printf("Failed to unpack event data: %v", err)
// 		return
// 	}

// 	// Format and send the event to the backend API
// 	data := EventData{
// 		From:   event.From.Hex(),
// 		Value:  event.Value.String(),
// 		TxHash: vLog.TxHash.Hex(),
// 	}

// 	sendToBackendAPI(data)
// }

// func sendToBackendAPI(event EventData) {
// 	backendURL := os.Getenv("BACKEND_API_URL")
// 	if backendURL == "" {
// 		log.Println("BACKEND_API_URL is not set")
// 		return
// 	}

// 	payload, err := json.Marshal(event)
// 	if err != nil {
// 		log.Printf("Failed to marshal event data: %v", err)
// 		return
// 	}

// 	resp, err := http.Post(backendURL, "application/json", bytes.NewReader(payload))
// 	if err != nil {
// 		log.Printf("Failed to send event to backend: %v", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	fmt.Printf("Event sent to backend successfully. Status code: %d\n", resp.StatusCode)
// }
