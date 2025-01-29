package utils

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

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

func FormatKeyValueLogs(data [][2]string) string {
	var builder strings.Builder
	builder.Grow(len(data) * 10)

	for _, entry := range data {
		builder.WriteString(fmt.Sprintf("  %s: %s\n", entry[0], entry[1]))
	}

	return builder.String()
}
