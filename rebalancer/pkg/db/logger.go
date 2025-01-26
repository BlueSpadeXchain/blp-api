package db

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

func FormatKeyValueLogs(data [][2]string) string {
	var builder strings.Builder
	builder.Grow(len(data) * 10)

	for _, entry := range data {
		builder.WriteString(fmt.Sprintf("  %s: %s\n", entry[0], entry[1]))
	}

	return builder.String()
}

func LogInfo(title string, message string) {
	if logrus.GetLevel() < logrus.InfoLevel {
		return
	}

	logrus.Info(fmt.Sprintf(
		"\033[1m%s\033[0m:\n%s",
		title,
		message,
	))
}

func LogError(message string, errStr string) {
	logrus.Error(fmt.Sprintf(
		"%s: \033[38;5;197m%s\033[0m",
		message,
		errStr,
	))
}

func FormatKeyValueSupabase(data [][2]string) string {
	var builder strings.Builder
	builder.Grow(len(data) * 10)

	for _, entry := range data {
		builder.WriteString(fmt.Sprintf("  %s: %s\n", entry[0], entry[1]))
	}

	return builder.String()
}

func LogSupabaseError(message SupabaseError) {
	logrus.Error(fmt.Sprintf(
		"\033[1m%s\033[0m:\n%s",
		"Supabase response error",
		FormatKeyValueLogs([][2]string{
			{"Code", message.Code},
			{"Details", message.Details},
			{"Hint", message.Hint},
			{"Message", message.Message},
		}),
	))
}
