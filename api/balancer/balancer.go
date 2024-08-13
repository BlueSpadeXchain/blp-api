package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	baseURL       = "https://hermes.pyth.network/v2/updates/price/latest"
	fetchInterval = 10 * time.Second // Adjust the interval as needed
)

// Define structs for parsed data
type Metadata struct {
	Slot               int `json:"slot"`
	ProofAvailableTime int `json:"proof_available_time"`
	PrevPublishTime    int `json:"prev_publish_time"`
}

type PriceDetail struct {
	Price       string `json:"price"`
	Conf        string `json:"conf"`
	Expo        int    `json:"expo"`
	PublishTime int    `json:"publish_time"`
}

type PriceUpdate struct {
	ID       string      `json:"id"`
	Price    PriceDetail `json:"price"`
	EMAPrice PriceDetail `json:"ema_price"`
	Metadata Metadata    `json:"metadata"`
}

type Data struct {
	Binary struct {
		Encoding string   `json:"encoding"`
		Data     []string `json:"data"`
	} `json:"binary"`
	Parsed []PriceUpdate `json:"parsed"`
}

// Fetch and process price updates
func fetchPriceUpdates(ids []string) (Data, error) {
	var queryParams []string
	for _, id := range ids {
		queryParams = append(queryParams, "ids[]="+id)
	}
	queryString := strings.Join(queryParams, "&")
	url := baseURL + "?" + queryString

	resp, err := http.Get(url)
	if err != nil {
		return Data{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Data{}, err
	}

	var data Data
	if err := json.Unmarshal(body, &data); err != nil {
		return Data{}, err
	}

	return data, nil
}

// Monitor price updates
func monitorPriceUpdates(ids []string) {
	for {
		data, err := fetchPriceUpdates(ids)
		if err != nil {
			log.Printf("Error fetching price updates: %v", err)
		} else {
			// Process the parsed data
			for _, update := range data.Parsed {
				fmt.Printf("ID: %s\n", update.ID)
				fmt.Printf("Price: %s (Confidence: %s, Expo: %d, Publish Time: %d)\n",
					update.Price.Price, update.Price.Conf, update.Price.Expo, update.Price.PublishTime)
				fmt.Printf("EMA Price: %s (Confidence: %s, Expo: %d, Publish Time: %d)\n",
					update.EMAPrice.Price, update.EMAPrice.Conf, update.EMAPrice.Expo, update.EMAPrice.PublishTime)
				fmt.Printf("Metadata: %+v\n", update.Metadata)
				fmt.Println()
			}
			fmt.Print("---\n")
		}

		time.Sleep(fetchInterval)
	}
}

func main() {
	ids := []string{
		"72b021217ca3fe68922a19aaf990109cb9d84e9ad004b4d2025ad6f529314419",
		"e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43",
		"4a8e42861cabc5ecb50996f92e7cfa2bce3fd0a2423b0c44c9b423fb2bd25478",
	}

	monitorPriceUpdates(ids)
}
