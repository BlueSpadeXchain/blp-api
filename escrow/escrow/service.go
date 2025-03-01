package escrow

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

type EventData struct {
	From   string `json:"from"`
	Value  string `json:"value"`
	TxHash string `json:"tx_hash"`
}

func StartListener(rpcURL string, chainId string, lastProcessedBlock *uint64) {
	if chainId == "" {
		chainId = "31337"
	}
	pkhex := os.Getenv("EVM_PRIVATE_KEY")
	if pkhex == "" {
		logrus.Fatal("EVM_PRIVATE_KEY is not set")
	}
	pk, _ := crypto.HexToECDSA(pkhex)

	userApi := os.Getenv("USER_API")
	if pkhex == "" {
		logrus.Fatal("USER_API is not set")
	}

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
		logrus.Error("Failed to subscribe to logs:", err)
		time.Sleep(time.Second)
		go StartListener(rpcURL, chainId, lastProcessedBlock) // Restart listener
		return
	}

	DepositEventSig := escrowABI.Events["DepositEvent"].ID
	StakingDepositEventSig := escrowABI.Events["StakingDepositEvent"].ID
	BurnRequestEventSig := escrowABI.Events["BurnRequestEvent"].ID
	// fmt.Printf("\n DepositEventSig: %v", DepositEventSig)
	// fmt.Printf("\n StakingDepositEventSig: %v", StakingDepositEventSig)
	// fmt.Printf("\n BurnRequestEventSig: %v", BurnRequestEventSig)

	logrus.Info("Onchain listener began...")

	for {
		select {
		case err := <-sub.Err():
			logrus.Errorf("Subscription error: %v. Retrying in 5 seconds...", err)
			time.Sleep(time.Second)
			go StartListener(rpcURL, chainId, lastProcessedBlock) // Restart listener
			return
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
					logrus.Fatal(err)
				}
				logrus.Info("DepositEvent:", event)

				signature, err := hashToSignECDSA(crypto.Keccak256(vLog.TxHash.Bytes()), pk)
				if err != nil {
					logrus.Error(err.Error())
				}

				request := &DespositRequestParams{
					ChainId:      chainId,
					Block:        strconv.FormatUint(vLog.BlockNumber, 10),
					BlockHash:    vLog.BlockHash.Hex(),
					TxHash:       vLog.TxHash.Hex(),
					Sender:       event.Account.Hex(),
					Receiver:     event.Account.Hex(),
					DepositNonce: event.Nonce.String(),
					Asset:        event.AssetAddress.Hex(),
					Amount:       event.AssetAmount.String(),
					Signature:    signature,
				}
				body, _ := ConvertStructToQuery(request)
				logrus.Info("body: ", body)
				logrus.Warning("deposit was triggered")
				sendRequest(userApi, "deposit", body)

			case StakingDepositEventSig:
				// need to check if asset is blu, address(0), or
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
				logrus.Info("StakingDepositEvent:", event)
				signature, err := hashToSignECDSA(crypto.Keccak256(vLog.TxHash.Bytes()), pk)
				if err != nil {
					logrus.Error(err.Error())
				}

				request := &StakeRequestParams{
					ChainId:      chainId,
					Block:        strconv.FormatUint(vLog.BlockNumber, 10),
					BlockHash:    vLog.BlockHash.Hex(),
					TxHash:       vLog.TxHash.Hex(),
					Sender:       event.Account.Hex(),
					Receiver:     event.Account.Hex(),
					DepositNonce: event.Nonce.String(),
					Asset:        event.AssetAddress.Hex(),
					Amount:       event.AssetAmount.String(),
					Signature:    signature,
				}
				body, _ := ConvertStructToQuery(request)
				logrus.Warning("stake was triggered")
				sendRequest(userApi, "eoa-stake", body)

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
				logrus.Info("BurnRequestEvent:", event)

			default:
				logrus.Error("Unknown event signature:", vLog.Topics[0].Hex())
			}

			*lastProcessedBlock = vLog.BlockNumber
		}
	}
}

// func StartListener(rpcURL string, chainId string, lastProcessedBlock *uint64) {
// 	// Your existing setup code
// 	if chainId == "" {
// 		chainId = "31337"
// 	}
// 	pkhex := os.Getenv("EVM_PRIVATE_KEY")
// 	if pkhex == "" {
// 		logrus.Fatal("EVM_PRIVATE_KEY is not set")
// 	}
// 	pk, _ := crypto.HexToECDSA(pkhex)

// 	userApi := os.Getenv("USER_API")
// 	if pkhex == "" {
// 		logrus.Fatal("USER_API is not set")
// 	}

// 	client, err := ethclient.Dial(rpcURL)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to Ethereum client: %v", err)
// 	}
// 	defer client.Close()

// 	escrowAddress, err := getAddress(chainId)
// 	if err != nil {
// 		logrus.Fatal("Failed to parse contract address: ", err.Error())
// 	}

// 	// query := ethereum.FilterQuery{
// 	// 	Addresses: []common.Address{escrowAddress},
// 	// }

// 	escrowABI, _ := abi.JSON(strings.NewReader(escrowContractABI))
// 	if err != nil {
// 		logrus.Fatal("Failed to parse contract ABI: ", err.Error())
// 	}

// 	// First, catch up on any missed blocks
// 	currentBlock, err := client.BlockNumber(context.Background())
// 	if err != nil {
// 		logrus.Error("Failed to get current block number:", err)
// 		time.Sleep(time.Second)
// 		go StartListener(rpcURL, chainId, lastProcessedBlock) // Restart listener
// 		return
// 	}

// 	logrus.Infof("Current block: %d, Last processed block: %d", currentBlock, *lastProcessedBlock)

// 	// If we're behind, query for missed events
// 	if *lastProcessedBlock < currentBlock {
// 		logrus.Infof("Catching up on %d blocks", currentBlock-*lastProcessedBlock)

// 		// Create a query for missed events
// 		catchUpQuery := ethereum.FilterQuery{
// 			FromBlock: new(big.Int).SetUint64(*lastProcessedBlock + 1),
// 			ToBlock:   new(big.Int).SetUint64(currentBlock),
// 			Addresses: []common.Address{escrowAddress},
// 		}

// 		fmt.Printf("\n catchUpQuery result: %v", catchUpQuery)

// 		// Get historical logs
// 		pastLogs, err := client.FilterLogs(context.Background(), catchUpQuery)
// 		if err != nil {
// 			logrus.Error("Failed to get past logs:", err)
// 			time.Sleep(time.Second)
// 			go StartListener(rpcURL, chainId, lastProcessedBlock) // Restart listener
// 			return
// 		}

// 		logrus.Infof("Retrieved %d historical events", len(pastLogs))

// 		if len(pastLogs) != 0 {
// 			// Process historical events in a separate goroutine
// 			go func() {
// 				var n int
// 				for _, vLog := range pastLogs {
// 					processEvent(vLog, &escrowABI, pk, userApi, chainId)

// 					// Update the last processed block
// 					if vLog.BlockNumber > *lastProcessedBlock {
// 						*lastProcessedBlock = vLog.BlockNumber
// 					}
// 				}
// 				n = n + 1
// 				fmt.Printf("\n n: %v", n)
// 				logrus.Infof("Historical event processing complete. Last processed block: %d", *lastProcessedBlock)
// 			}()
// 		} else {
// 			*lastProcessedBlock = currentBlock
// 		}
// 	}

// 	// Set up subscription for new events
// 	logs := make(chan types.Log)
// 	query := ethereum.FilterQuery{
// 		Addresses: []common.Address{escrowAddress},
// 	}

// 	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
// 	if err != nil {
// 		logrus.Error("Failed to subscribe to logs:", err)
// 		time.Sleep(time.Second)
// 		go StartListener(rpcURL, chainId, lastProcessedBlock) // Restart listener
// 		return
// 	}

// 	logrus.Info("Onchain listener began...")

// 	// Start processing new events
// 	go func() {
// 		for {
// 			select {
// 			case err := <-sub.Err():
// 				logrus.Errorf("Subscription error: %v. Retrying in 5 seconds...", err)
// 				time.Sleep(time.Second)
// 				go StartListener(rpcURL, chainId, lastProcessedBlock) // Restart listener
// 				return
// 			case vLog := <-logs:
// 				processEvent(vLog, &escrowABI, pk, userApi, chainId)

// 				// Update the last processed block
// 				if vLog.BlockNumber > *lastProcessedBlock {
// 					*lastProcessedBlock = vLog.BlockNumber
// 				}
// 			}
// 		}
// 	}()
// }

// // Separate function to process a single event
// func processEvent(vLog types.Log, escrowABI *abi.ABI, pk *ecdsa.PrivateKey, userApi string, chainId string) {
// 	fmt.Println("BlockHash:", vLog.BlockHash.Hex())
// 	fmt.Println("BlockNumber:", vLog.BlockNumber)
// 	fmt.Println("TxHash:", vLog.TxHash.Hex())

// 	// Handle events based on signature hash
// 	switch vLog.Topics[0] {
// 	case escrowABI.Events["DepositEvent"].ID:
// 		var event struct {
// 			Sender       common.Address
// 			Account      common.Address
// 			Nonce        *big.Int
// 			AssetAddress common.Address
// 			AssetAmount  *big.Int
// 		}
// 		err := escrowABI.UnpackIntoInterface(&event, "DepositEvent", vLog.Data)
// 		if err != nil {
// 			logrus.Error(err)
// 			return
// 		}
// 		logrus.Info("DepositEvent:", event)
// 		signature, err := hashToSignECDSA(crypto.Keccak256(vLog.TxHash.Bytes()), pk)
// 		if err != nil {
// 			logrus.Error(err.Error())
// 			return
// 		}
// 		request := &DespositRequestParams{
// 			ChainId:      chainId,
// 			Block:        strconv.FormatUint(vLog.BlockNumber, 10),
// 			BlockHash:    vLog.BlockHash.Hex(),
// 			TxHash:       vLog.TxHash.Hex(),
// 			Sender:       event.Account.Hex(),
// 			Receiver:     event.Account.Hex(),
// 			DepositNonce: event.Nonce.String(),
// 			Asset:        event.AssetAddress.Hex(),
// 			Amount:       event.AssetAmount.String(),
// 			Signature:    signature,
// 		}
// 		body, _ := ConvertStructToQuery(request)
// 		logrus.Info("body: ", body)
// 		logrus.Warning("deposit was triggered")
// 		sendRequest(userApi, "deposit", body)

// 	case escrowABI.Events["StakingDepositEvent"].ID:
// 		// Your existing StakingDepositEvent handling code
// 		// ...

// 	case escrowABI.Events["BurnRequestEvent"].ID:
// 		// Your existing BurnRequestEvent handling code
// 		// ...

// 	default:
// 		logrus.Error("Unknown event signature:", vLog.Topics[0].Hex())
// 	}
// }

func hashToSignECDSA(hash []byte, pk *ecdsa.PrivateKey) (string, error) {
	header, _ := hex.DecodeString("19457468657265756d205369676e6564204d6573736167653a0a3332")
	ethHash := append(header, hash...)
	unsignedHash := crypto.Keccak256(ethHash)

	return SignData(unsignedHash, pk)
}
