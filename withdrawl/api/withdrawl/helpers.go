package withdrawHandler

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
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

func ExecuteFunction(client ethclient.Client, contractAddress common.Address, parsedABI abi.ABI, methodName string, value *big.Int, args ...interface{}) (receiptJSON []byte, err error) {
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	// privateKey, relayAddress, err := utils.EnvKey2Ecdsa()
	// if err != nil {
	// 	return nil, err
	// }

	privateKey, relayAddress, err := Key2Ecdsa(os.Getenv("EVM_PRIVATE_KEY"))
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return nil, err
	}
	auth.Value = big.NewInt(1000000000000000000)

	data, err := parsedABI.Pack(methodName, args...)
	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{
		From:     relayAddress,
		To:       &contractAddress,
		Gas:      0,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	}

	_, err = client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return nil, err
	}

	nonce, err := client.PendingNonceAt(context.Background(), relayAddress)
	if err != nil {
		return nil, err
	}

	estimatedGas, err := client.EstimateGas(context.Background(), callMsg)
	if err != nil {
		return nil, err
	}

	gasLimit := 120 * estimatedGas / 100

	tx := types.NewTransaction(nonce, contractAddress, value, gasLimit, gasPrice, data)

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainId), privateKey)
	if err != nil {
		return nil, err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}

	reciept, err := bind.WaitMined(context.Background(), &client, signedTx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("tx receipt: \n%v", reciept)

	var returnedData []byte
	for _, log := range reciept.Logs {
		if len(log.Data) > 0 {
			returnedData = log.Data
			break
		}
	}

	// receiptJSON, err = json.Marshal(receipt)
	// if err != nil {
	// 	log.Fatalf("Failed to JSON marshal receipt: %v", err)
	// 	return nil, err
	// }

	return returnedData, nil
}

func Key2Ecdsa(key string) (*ecdsa.PrivateKey, common.Address, error) {
	return PrivateKey2Sepc256k1(key)
}

func PrivateKey2Sepc256k1(privateKeyString string) (privateKey *ecdsa.PrivateKey, publicAddress common.Address, err error) {
	privateKey, err = crypto.HexToECDSA(privateKeyString)
	if err != nil {
		err = fmt.Errorf("Error converting private key: %v", err)
		return
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		err = fmt.Errorf("Error casting public key to ECDSA")
		return
	}

	publicAddress = crypto.PubkeyToAddress(*publicKeyECDSA)
	return
}
