package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	WithdrawHandler "github.com/BlueSpadeXchain/blp-api/withdrawal/api/withdrawal"
	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

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

func main() {
	serverEnv := flag.String("server", "production", "Specify the server environment (local/production)")
	flag.Parse()

	var envFile string
	if *serverEnv == "local" {
		envFile = ".env.local"
	} else {
		envFile = ".env"
	}

	err := godotenv.Load(envFile)
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	debugMode := os.Getenv("DEBUG_MODE_ENABLED")
	if debugMode == "true" || debugMode == "1" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.FatalLevel)
	}

	logrus.SetFormatter(&CustomLogFormatter{})

	logrus.Warning("program starting in debug mode...")

	http.HandleFunc("/api/withdrawal", WithdrawHandler.Handler)

	log.Println("Starting server on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// package main

// import (
// 	"flag"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/BlueSpadeXchain/blp-api/withdrawal/pkg/db"
// 	"github.com/BlueSpadeXchain/blp-api/withdrawal/withdrawal"
// 	"github.com/joho/godotenv"
// 	"github.com/sirupsen/logrus"
// 	"github.com/supabase-community/supabase-go"
// )

// // CustomLogFormatter formats logs in the same way as your API example
// type CustomLogFormatter struct{}

// func (f *CustomLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
// 	level := entry.Level.String()

// 	// Set log colors, using 256-color ANSI escape code
// 	var levelColor string
// 	switch entry.Level {
// 	case logrus.InfoLevel:
// 		levelColor = "\033[38;5;45m" // Green
// 	case logrus.DebugLevel:
// 		levelColor = "\033[34m" // Blue
// 	case logrus.WarnLevel:
// 		levelColor = "\033[33m" // Yellow
// 	case logrus.ErrorLevel:
// 		levelColor = "\033[31m" // Red
// 	case logrus.FatalLevel:
// 		levelColor = "\033[35m" // Magenta
// 	case logrus.PanicLevel:
// 		levelColor = "\033[36m" // Cyan
// 	default:
// 		levelColor = "\033[0m" // Reset color (for unknown levels)
// 	}

// 	logMessage := fmt.Sprintf("\n\033[38;5;180m%s\033[0m [%s%s\033[0m] %s\n",
// 		entry.Time.Format("2006-01-02 15:04:05"), // Timestamp in cream color
// 		levelColor,                               // Colorized log level
// 		level,                                    // Log level name
// 		entry.Message)                            // Log message

// 	return []byte(logMessage), nil
// }

// func fetchData(url string) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		logrus.Error("Failed to fetch data: ", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		logrus.Warn("Received non-OK response: ", resp.Status)
// 		return
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		logrus.Error("Failed to read response body: ", err)
// 		return
// 	}

// 	logrus.Info("Fetched data: ", string(body))
// }

// func main() {
// 	serverEnv := flag.String("server", "production", "Specify the server environment (local/production)")
// 	chainIdFlag := flag.String("chainid", "", "Specify the chain id (default anvil)")
// 	flag.Parse()

// 	var envFile string
// 	if *serverEnv == "local" {
// 		envFile = ".env.local"
// 	} else {
// 		envFile = ".env"
// 	}

// 	err := godotenv.Load(envFile)
// 	if err != nil {
// 		fmt.Println("Error loading .env file")
// 	}

// 	debugMode := os.Getenv("DEBUG_MODE_ENABLED")
// 	if debugMode == "true" || debugMode == "1" {
// 		logrus.SetLevel(logrus.DebugLevel)
// 	} else {
// 		logrus.SetLevel(logrus.FatalLevel)
// 	}

// 	logrus.SetFormatter(&CustomLogFormatter{})

// 	for {
// 		run(chainIdFlag)
// 		logrus.Error("Bot encountered an error. Restarting in 1 seconds...")
// 		time.Sleep(time.Second)
// 	}
// }

// func run(chainIdFlag *string) {
// 	defer func() {
// 		if rec := recover(); rec != nil {
// 			log.Printf("\nRecovered from panic: %v", rec)

// 			supabaseUrl := os.Getenv("SUPABASE_URL")
// 			supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
// 			supabaseClient, err := supabase.NewClient(supabaseUrl, supabaseKey, nil)
// 			if err == nil {
// 				logErr := db.LogPanic(supabaseClient, fmt.Sprintf("%v", rec), nil)
// 				if logErr != nil {
// 					log.Printf("\nFailed to log panic to Supabase: %v", logErr)
// 				}
// 			} else {
// 				log.Printf("\nFailed to create Supabase client for panic logging: %v", err)
// 			}
// 		}
// 	}()

// 	logrus.Info("Supabase listener starting...")

// 	jsonRpc, err := withdrawal.GetChainRpc(*chainIdFlag)
// 	if err != nil {
// 		logrus.Error(err)
// 		return
// 	}

// 	withdrawal.StartListener(jsonRpc, *chainIdFlag)

// }
