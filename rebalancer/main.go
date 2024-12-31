package rebalancer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	// "github.com/BlueSpadeXchain/blp-api/pkg/utils"
)

// CustomLogFormatter formats logs in the same way as your API example
type CustomLogFormatter struct{}

func (f *CustomLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	level := entry.Level.String()

	// Set log colors, using 256-color ANSI escape code
	var levelColor string
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = "\033[38;5;45m" // Green
	case logrus.DebugLevel:
		levelColor = "\033[34m" // Blue
	case logrus.WarnLevel:
		levelColor = "\033[33m" // Yellow
	case logrus.ErrorLevel:
		levelColor = "\033[31m" // Red
	case logrus.FatalLevel:
		levelColor = "\033[35m" // Magenta
	case logrus.PanicLevel:
		levelColor = "\033[36m" // Cyan
	default:
		levelColor = "\033[0m" // Reset color (for unknown levels)
	}

	logMessage := fmt.Sprintf("\n\033[38;5;180m%s\033[0m [%s%s\033[0m] %s\n",
		entry.Time.Format("2006-01-02 15:04:05"), // Timestamp in cream color
		levelColor,                               // Colorized log level
		level,                                    // Log level name
		entry.Message)                            // Log message

	return []byte(logMessage), nil
}

func fetchData(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logrus.Error("Failed to fetch data: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Warn("Received non-OK response: ", resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("Failed to read response body: ", err)
		return
	}

	logrus.Info("Fetched data: ", string(body))
}

func main() {
	// Load environment variables
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Println("Error loading .env.local file")
	}

	// Configure logging
	debugMode := os.Getenv("DEBUG_MODE_ENABLED")
	if debugMode == "true" || debugMode == "1" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&CustomLogFormatter{})

	logrus.Info("Rebalancer bot starting...")

	url := "https://example.com"

	// Periodically fetch data
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logrus.Debug("Fetching data from ", url)
			fetchData(url)
		}
	}
}
