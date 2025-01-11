package db

import (
	"fmt"
	"strings"

	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
	"github.com/sirupsen/logrus"
)

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
		utils.FormatKeyValueLogs([][2]string{
			{"Code", message.Code},
			{"Details", message.Details},
			{"Hint", message.Hint},
			{"Message", message.Message},
		}),
	))
}
