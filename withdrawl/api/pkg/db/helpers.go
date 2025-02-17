package db

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/supabase-community/supabase-go"
)

func LogError(client *supabase.Client, err error, message string, context interface{}) error {
	logData := map[string]interface{}{
		"log_level": "ERROR",
		"error":     err,
		"message":   message,
		"context":   context,
	}

	_, _, dbErr := client.From("debug_logs").Insert(logData, false, "", "minimal", "").Execute()
	if dbErr != nil {
		log.Printf("Failed to insert log: %v", dbErr)
		return dbErr
	}
	return nil
}

func LogPanic(client *supabase.Client, message string, context interface{}) error {

	if err := LogError(client, fmt.Errorf("log panic"), message, context); err != nil {
		log.Printf("Failed to log panic: %v", err)
		return err
	}

	return nil
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
