package utils

import (
	"fmt"
	"reflect"
	"strings"
)

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
