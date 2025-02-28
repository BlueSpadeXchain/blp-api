package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

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

func (e Error) Error() string {
	return fmt.Sprintf("Error (Code: %d, Message: %s)", e.Code, e.Message)
}

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

func ErrInternal(message string) Error {
	origin := GetOrigin()

	return Error{
		Code:    500,
		Message: "Internal server error",
		Details: message,
		Origin:  origin,
	}
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
