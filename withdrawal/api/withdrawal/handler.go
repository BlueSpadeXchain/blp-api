package withdrawHandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BlueSpadeXchain/blp-api/pkg/db"
	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/supabase-community/supabase-go"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("\nRecovered from panic: %v", rec)

			supabaseUrl := os.Getenv("SUPABASE_URL")
			supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
			supabaseClient, err := supabase.NewClient(supabaseUrl, supabaseKey, nil)
			if err == nil {
				logErr := db.LogPanic(supabaseClient, fmt.Sprintf("%v", rec), nil)
				if logErr != nil {
					log.Printf("\nFailed to log panic to Supabase: %v", logErr)
				}
			} else {
				log.Printf("\nFailed to create Supabase client for panic logging: %v", err)
			}

			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	callerUrl := r.Header.Get("Origin")
	if callerUrl == "" {
		callerUrl = r.Header.Get("Referer")
	}

	if callerUrl == "" || !isWhitelistedUrl(callerUrl) {
		logrus.Error("Unauthorized caller", callerUrl)
		logrus.Error("r:", r)
		http.Error(w, "Unauthorized caller", http.StatusForbidden)
		return
	}

	handlerWithCORS := utils.EnableCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		var response interface{}
		var err error
		supabaseUrl := os.Getenv("SUPABASE_URL")
		supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
		supabaseClient, err := supabase.NewClient(supabaseUrl, supabaseKey, nil)
		if err != nil {
			http.Error(w, "Failed to create Supabase client", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		switch query.Get("query") {
		case "withdraw-blu": // used when unstaking blu
			response, err = WithdrawBluRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		case "withdraw-blp": // used when withdrawing balance
			response, err = WithdrawBlpRequest(r, supabaseClient)
			HandleResponse(w, r, supabaseClient, response, err)
			return
		default:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(utils.ErrMalformedRequest("Invalid query parameter"))
			return
		}
	}))
	//userid funded: 1d2664a39eee6098
	handlerWithCORS.ServeHTTP(w, r)
}

func HandleResponse(w http.ResponseWriter, r *http.Request, supabaseClient *supabase.Client, response interface{}, err error) {
	if err != nil {
		if logErr := db.LogError(supabaseClient, err, r.URL.Query().Get("query"), response); logErr != nil {
			utils.LogError("Failed to log error", logErr.Error())
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	utils.LogInfo("response", fmt.Sprint(response))
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}
}
