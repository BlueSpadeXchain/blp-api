package escrow

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/ethereum/go-ethereum/crypto"
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
