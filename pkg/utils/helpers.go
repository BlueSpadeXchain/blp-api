package utils

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func WriteJSONResponse(w http.ResponseWriter, r *http.Request, message string) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ParseAndValidateParams(r *http.Request, params interface{}) error {
	val := reflect.ValueOf(params).Elem() // Dereference the pointer to access the underlying struct
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}
	typ := val.Type()

	missingFields := []string{}
	allowedFields := make(map[string]struct{})

	LogInfo("query", fmt.Sprint(typ))
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		queryTag := fieldType.Tag.Get("query")
		optionalTag := fieldType.Tag.Get("optional")

		if queryTag != "" {
			allowedFields[queryTag] = struct{}{}
		}

		if _, exists := typ.FieldByName(fieldType.Name); exists {
			if field.Kind() == reflect.Struct {
				// Recursively parse nested struct fields
				nestedParams := reflect.New(fieldType.Type).Interface()
				if err := ParseAndValidateParams(r, nestedParams); err != nil {
					return err
				}
				// After recursion, set the original struct's field value
				field.Set(reflect.ValueOf(nestedParams).Elem())
			} else if queryTag != "" {
				queryValue := r.URL.Query().Get(queryTag)

				// If the field is required (i.e., optional is not set to "true")
				if queryValue == "" && optionalTag != "true" {
					missingFields = append(missingFields, queryTag)
				} else if queryValue != "" {
					field.SetString(queryValue)
				}
			}
		}
	}

	// If there are missing fields, return an error response
	if len(missingFields) > 0 {
		return ErrMalformedRequest(fmt.Sprint("Missing fields: " + strings.Join(missingFields, ", ")))
	}

	return nil
}

func StringifyStructFields(params interface{}, indent string) string {
	var result strings.Builder
	val := reflect.ValueOf(params)

	// Handle nil
	if !val.IsValid() {
		return "nil"
	}

	// Dereference pointer if needed
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := typ.Field(i)
			fieldName := fieldType.Name

			// Handle nested types
			switch field.Kind() {
			case reflect.Struct, reflect.Map:
				result.WriteString(fmt.Sprintf("\n%s\033[1m%s\033[0m:\n", indent, fieldName))
				nestedResult := StringifyStructFields(field.Interface(), indent+"  ")
				result.WriteString(nestedResult)
			default:
				result.WriteString(fmt.Sprintf("\n%s\033[1m%s\033[0m: %v", indent, fieldName, field.Interface()))
			}
		}

	case reflect.Map:
		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()

			// Convert the key to string (most map keys will be strings anyway)
			keyStr := fmt.Sprintf("%v", k.Interface())

			// Handle nested types in map values
			switch v.Kind() {
			case reflect.Struct, reflect.Map:
				result.WriteString(fmt.Sprintf("\n%s\033[1m%s\033[0m:\n", indent, keyStr))
				nestedResult := StringifyStructFields(v.Interface(), indent+"  ")
				result.WriteString(nestedResult)
			default:
				result.WriteString(fmt.Sprintf("\n%s\033[1m%s\033[0m: %v", indent, keyStr, v.Interface()))
			}
		}

	default:
		return fmt.Sprintf("%v", val.Interface())
	}

	return result.String()
}

func PrintStructFields(params interface{}) {
	val := reflect.ValueOf(params)

	// Ensure the value is a struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		fmt.Println("Expected a struct")
		return
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name

		// Check if it's a nested struct
		if field.Kind() == reflect.Struct {
			fmt.Printf("\n%s:\n", fieldName)
			PrintStructFields(field.Interface()) // Recursively print nested struct fields
			fmt.Println()
		} else {
			fmt.Printf("\n%s: %v", fieldName, field.Interface()) // Print field value
		}
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("Error (Code: %d, Message: %s)", e.Code, e.Message)
}

// func ErrMalformedRequest(w http.ResponseWriter, message string) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusBadRequest)

// 	json.NewEncoder(w).Encode(&Error{
// 		Code:    400,
// 		Message: "Malformed request",
// 		Details: message,
// 	})
// }

func GetOrigin() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}
	funcName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(funcName, ".")
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], ".")
	}
	return "unknown"
}

func ErrMalformedRequest(message string) error {
	origin := GetOrigin()

	return Error{
		Code:    400,
		Message: "Malformed request",
		Details: message,
		Origin:  origin,
	}
}

func ErrInternal(message string) Error {
	origin := GetOrigin()

	return Error{
		Code:    500,
		Message: "Internal server error",
		Details: message,
		Origin:  origin,
	}
}

func EnvKey2Ecdsa() (*ecdsa.PrivateKey, common.Address, error) {
	return PrivateKey2Sepc256k1(os.Getenv("RELAY_PRIVATE_KEY"))
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

func Str2Bytes(hexStr string) ([]byte, error) {
	if hexStr == "" {
		return []byte{}, nil // Return empty byte slice for empty input
	}

	if len(hexStr)%2 != 0 {
		return nil, fmt.Errorf("invalid payload: length odd")
	}

	for _, r := range hexStr {
		if _, err := strconv.ParseUint(string(r), 16, 8); err != nil {
			return nil, fmt.Errorf("invalid payload: hex only 0123456789abcdefABCDEF")
		}
	}

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func HasInt(inputArray []int, input int) error {
	for _, value := range inputArray {
		if value == input {
			return nil
		}
	}
	return fmt.Errorf("int not found in input array")
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		LogInfo("API Request", FormatKeyValueLogs([][2]string{
			{"Method", r.Method},
			{"URL", fmt.Sprintf("%v", r.URL)},
		}))

		next.ServeHTTP(w, r)
	})
}

// Helper function to convert []byte to hex string prefixed with "0x".
func ToHexBytes(data []byte) string {
	if len(data) == 0 {
		return "0x"
	}
	return "0x" + hex.EncodeToString(data)
}

func HexToBytes(hexStr string) ([]byte, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Helper function to convert common.Address to hex string prefixed with "0x".
func ToHexAddress(addr common.Address) string {
	return "0x" + hex.EncodeToString(addr[:])
}

// Helper function to convert [4]byte to a uint32 string.
func Uint32ToString(data [4]byte) string {
	value := uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])
	return fmt.Sprintf("%d", value)
}

// Helper function to convert a byte to string.
func Uint8ToString(data byte) string {
	return fmt.Sprintf("%d", data)
}

func Uint256ToBytes(value *big.Int) []byte {
	bytes := value.FillBytes(make([]byte, 32))
	return bytes
}

func HashToEthHash(hash []byte) []byte {
	ethHash := append(EthDomainHeader, hash...)
	return crypto.Keccak256(ethHash)
}

// remember to -= 27 for ethereum signatures
func ValidateEvmEcdsaSignature(hash []byte, signature []byte, address common.Address) (bool, error) {
	if len(signature) != 65 {
		return false, fmt.Errorf("invalid signature length: %d", len(signature))
	}

	ethHash := append(EthDomainHeader, hash...)
	unsignedHash := crypto.Keccak256(ethHash)
	fmt.Printf("\n unsigned hash: %v", hex.EncodeToString(unsignedHash))

	recoveredPubKey, err := crypto.SigToPub(unsignedHash, signature)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %w", err)
	}
	recoveredAddress := crypto.PubkeyToAddress(*recoveredPubKey)

	LogInfo("Recover details", FormatKeyValueLogs([][2]string{
		{"recovered address", recoveredAddress.String()},
		{"expected address ", address.Hex()},
	}))

	return bytes.Equal(recoveredAddress.Bytes(), address.Bytes()), nil
}

func RemoveHex0xPrefix(hex string) string {
	if strings.HasPrefix(hex, "0x") || strings.HasPrefix(hex, "0X") {
		return hex[2:]
	}
	return hex
}

// GetCurrentPriceData queries the API for the current price data for the given pair ID.
func GetCurrentPriceData(pair string) (PriceUpdate, error) {
	baseURL := "https://hermes.pyth.network/v2/updates/price/latest"

	// Create the request with query parameters
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return PriceUpdate{}, fmt.Errorf("error parsing URL: %v", err)
	}

	q := reqURL.Query()
	q.Add("ids[]", pair)
	reqURL.RawQuery = q.Encode()

	resp, err := http.Get(reqURL.String())
	if err != nil {
		return PriceUpdate{}, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PriceUpdate{}, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PriceUpdate{}, fmt.Errorf("error reading response body: %v", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return PriceUpdate{}, fmt.Errorf("error unmarshaling response JSON: %v", err)
	}

	LogResponse(reqURL.String(), response.Parsed[0])

	return response.Parsed[0], nil
}
