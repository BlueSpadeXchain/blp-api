package withdrawHandler

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"reflect"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func ConvertStructToQuery(params interface{}) (string, error) {
	fmt.Printf("\n params: %v", params)
	v := reflect.ValueOf(params)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("expected a pointer to a struct")
	}

	query := url.Values{}

	t := v.Elem().Type()
	for i := 0; i < v.Elem().NumField(); i++ {
		field := v.Elem().Field(i)
		fieldName := t.Field(i).Tag.Get("query")

		if field.IsZero() {
			continue
		}

		query.Set(fieldName, fmt.Sprintf("%v", field.Interface()))
	}
	fmt.Printf("\n params end: %v", query.Encode())
	return query.Encode(), nil
}

func SignData(hash []byte, pk *ecdsa.PrivateKey) (string, error) {
	signature, err := crypto.Sign(hash, pk)
	return hex.EncodeToString(signature), err
}

func sendRequest(api, query, body string) {
	url := fmt.Sprintf("%v?query=%v&%v", api, query, body)
	logrus.Info("Request forwarded: ", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error("Request creation error: ", err)
		return
	}

	go func() {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logrus.Error("Request error: ", err)
			return
		}
		defer resp.Body.Close()
	}()
}

// ExecuteResponse represents the structured response from an on-chain function execution
type ExecuteResponse struct {
	TxHash         common.Hash     `json:"txHash"`
	BlockNumber    *big.Int        `json:"blockNumber,omitempty"`
	Status         uint64          `json:"status"`
	GasUsed        uint64          `json:"gasUsed"`
	From           common.Address  `json:"from"`
	To             common.Address  `json:"to"`
	ContractReturn []byte          `json:"contractReturn,omitempty"`
	Logs           []LogData       `json:"logs,omitempty"`
	RawReceipt     json.RawMessage `json:"rawReceipt,omitempty"`
}

// LogData represents a simplified event log
type LogData struct {
	Address common.Address `json:"address"`
	Topics  []common.Hash  `json:"topics"`
	Data    []byte         `json:"data"`
}

func ExecuteFunction(client ethclient.Client, contractAddress common.Address, parsedABI abi.ABI, methodName string, value *big.Int, args ...interface{}) (response *ExecuteResponse, err error) {
	log.Printf("Starting ExecuteFunction for method: %s", methodName)

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}
	log.Printf("Chain ID: %s", chainId.String())

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}
	log.Printf("Suggested gas price: %s", gasPrice.String())

	privateKey, relayAddress, err := Key2Ecdsa(os.Getenv("EVM_PRIVATE_KEY"))
	if err != nil {
		return nil, fmt.Errorf("failed to get private key: %w", err)
	}
	log.Printf("Relay address: %s", relayAddress.Hex())

	// Get account balance
	balance, err := client.BalanceAt(context.Background(), relayAddress, nil)
	if err != nil {
		log.Printf("WARNING: Could not get balance: %v", err)
	} else {
		log.Printf("Account balance: %s wei", balance.String())
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}
	auth.Value = value

	data, err := parsedABI.Pack(methodName, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %w", err)
	}
	log.Printf("Packed data length: %d bytes", len(data))
	log.Printf("Method name: %s, Value: %s", methodName, value.String())

	initialGasEstimate := uint64(100000)

	callMsg := ethereum.CallMsg{
		From:     relayAddress,
		To:       &contractAddress,
		Gas:      initialGasEstimate,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	}

	// Try simulating the contract call
	log.Printf("Simulating contract call...")
	_, err = client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return nil, fmt.Errorf("contract call simulation failed: %w", err)
	}
	log.Printf("Contract call simulation successful")

	nonce, err := client.PendingNonceAt(context.Background(), relayAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}
	log.Printf("Nonce: %d", nonce)

	log.Printf("Estimating gas...")
	estimatedGas, err := client.EstimateGas(context.Background(), callMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}
	log.Printf("Estimated gas: %d", estimatedGas)

	// Calculate total cost
	gasLimit := 120 * estimatedGas / 100
	log.Printf("Gas limit (120%% of estimate): %d", gasLimit)

	totalGasCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	log.Printf("Total gas cost: %s wei", totalGasCost.String())

	totalCost := new(big.Int).Add(totalGasCost, value)
	log.Printf("Total transaction cost (gas + value): %s wei", totalCost.String())

	if balance != nil {
		if balance.Cmp(totalCost) < 0 {
			log.Printf("WARNING: Account balance (%s) is less than total cost (%s)", balance.String(), totalCost.String())
		}
	}

	tx := types.NewTransaction(nonce, contractAddress, value, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainId), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	log.Printf("Sending transaction...")
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		gasFeeCap := new(big.Int).Mul(gasPrice, big.NewInt(2)) // Example: double the gas price as a fee cap
		log.Printf("Transaction failed. You might try with gasFeeCap of %s or lower gasLimit", gasFeeCap.String())
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	log.Printf("Transaction sent successfully! Hash: %s", signedTx.Hash().Hex())

	// Create response with transaction hash immediately
	response = &ExecuteResponse{
		TxHash: signedTx.Hash(),
		From:   relayAddress,
		To:     contractAddress,
	}

	// Wait for transaction receipt
	receipt, err := bind.WaitMined(context.Background(), &client, signedTx)
	if err != nil {
		return response, fmt.Errorf("failed to get transaction receipt: %w", err)
	}

	// Fill in the response with receipt data
	response.BlockNumber = receipt.BlockNumber
	response.Status = receipt.Status
	response.GasUsed = receipt.GasUsed

	// Extract logs
	if len(receipt.Logs) > 0 {
		logs := make([]LogData, 0, len(receipt.Logs))
		for _, log := range receipt.Logs {
			if len(log.Data) > 0 {
				response.ContractReturn = log.Data // Take the first non-empty log data as return value
			}
			logs = append(logs, LogData{
				Address: log.Address,
				Topics:  log.Topics,
				Data:    log.Data,
			})
		}
		response.Logs = logs
	}

	// Add raw receipt JSON for full details
	rawJSON, err := json.Marshal(receipt)
	if err != nil {
		// Log but don't fail
		fmt.Printf("Warning: Failed to marshal receipt to JSON: %v\n", err)
	} else {
		response.RawReceipt = rawJSON
	}

	fmt.Printf("Transaction successful - hash: %s, block: %s\n",
		response.TxHash.Hex(),
		response.BlockNumber.String())

	return response, nil
}

func Key2Ecdsa(key string) (*ecdsa.PrivateKey, common.Address, error) {
	return PrivateKey2Sepc256k1(key)
}

func PrivateKey2Sepc256k1(privateKeyString string) (privateKey *ecdsa.PrivateKey, publicAddress common.Address, err error) {
	privateKey, err = crypto.HexToECDSA(privateKeyString)
	if err != nil {
		err = fmt.Errorf("error converting private key: %v", err)
		return
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		err = fmt.Errorf("error casting public key to ECDSA")
		return
	}

	publicAddress = crypto.PubkeyToAddress(*publicKeyECDSA)
	return
}

func isWhitelistedUrl(callerUrl string) bool {
	whitelistUrl := os.Getenv("WHITELISTED_URL") // https://blp-api-vercel.vercel.app/ || http://localhost:8080

	if whitelistUrl == "" {
		logrus.Error("No whitelist URL configured in environment variables")
		return false
	}

	whitelistParsed, err1 := url.Parse(whitelistUrl)
	callerParsed, err2 := url.Parse(callerUrl)

	if err1 != nil || err2 != nil {
		logrus.Error(fmt.Sprintf("Error parsing URLs: whitelist=%v, caller=%v", err1, err2))
		return false
	}

	// Compare origin (scheme + host + port)
	return whitelistParsed.Scheme+"://"+whitelistParsed.Host == callerParsed.Scheme+"://"+callerParsed.Host
}
