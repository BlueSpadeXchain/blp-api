package hermes

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// PriceData represents the structured price feed data
type PriceData struct {
	ID           string    `json:"id"`
	Price        float64   `json:"price"`
	Conf         float64   `json:"conf"`
	Timestamp    time.Time `json:"timestamp"`
	PrevPubTime  time.Time `json:"prevPubTime"`
	PubTime      time.Time `json:"pubTime"`
	EncodedPrice string    `json:"encodedPrice,omitempty"`
	Parsed       bool      `json:"parsed"`
}

// BatchPriceUpdate represents multiple price updates
type BatchPriceUpdate struct {
	Updates map[string]PriceData `json:"updates"`
	Error   string               `json:"error,omitempty"`
}

// PriceServiceConfig holds configuration options
type PriceServiceConfig struct {
	BaseURL        string
	SSEEndpoint    string
	Encoding       string // hex or base64
	ParsedEnabled  bool
	AllowUnordered bool
	BenchmarksOnly bool
	IgnoreInvalid  bool
	ReconnectDelay time.Duration
	CacheTTL       time.Duration
}

// PriceService manages price feed subscriptions and caching
type PriceService struct {
	config      PriceServiceConfig
	prices      sync.Map // map[string]*PriceData
	subscribers sync.Map // map[string][]chan BatchPriceUpdate
	client      *http.Client
	cache       sync.Map
}

// NewPriceService creates a new PriceService instance
func NewPriceService(cfg PriceServiceConfig) (*PriceService, error) {
	ps := &PriceService{
		config: cfg,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	return ps, nil
}

// SubscribeMultiple subscribes to multiple price feeds
func (ps *PriceService) SubscribeMultiple(ctx context.Context, priceIDs []string) (<-chan BatchPriceUpdate, error) {
	updates := make(chan BatchPriceUpdate, 100)

	queryParams := ps.buildQueryParams(priceIDs)
	sseURL := fmt.Sprintf("%s%s?%s", ps.config.BaseURL, ps.config.SSEEndpoint, queryParams)

	sseResp, err := ps.client.Get(sseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSE: %w", err)
	}

	defer sseResp.Body.Close()

	go ps.processSSEEvents(ctx, bufio.NewReader(sseResp.Body), updates)

	go func() {
		<-ctx.Done()
		close(updates)
	}()

	return updates, nil
}

// processSSEEvents reads and processes SSE events
func (ps *PriceService) processSSEEvents(ctx context.Context, reader *bufio.Reader, updates chan BatchPriceUpdate) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error reading SSE: %v", err)
				return
			}

			var update BatchPriceUpdate
			if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &update); err != nil {
				log.Printf("Error parsing SSE: %v", err)
				continue
			}

			for id, price := range update.Updates {
				ps.cache.Store(id, price)
			}

			select {
			case updates <- update:
			default:
				log.Println("Update channel full, dropping update")
			}
		}
	}
}

// buildQueryParams constructs query parameters
func (ps *PriceService) buildQueryParams(priceIDs []string) string {
	params := []string{}

	// Add price feed IDs
	for _, id := range priceIDs {
		params = append(params, fmt.Sprintf("ids[]=%s", id))
	}

	// Add optional parameters
	if ps.config.Encoding != "" {
		params = append(params, fmt.Sprintf("encoding=%s", ps.config.Encoding))
	}
	params = append(params, fmt.Sprintf("parsed=%t", ps.config.ParsedEnabled))
	params = append(params, fmt.Sprintf("allow_unordered=%t", ps.config.AllowUnordered))
	params = append(params, fmt.Sprintf("benchmarks_only=%t", ps.config.BenchmarksOnly))
	params = append(params, fmt.Sprintf("ignore_invalid_price_ids=%t", ps.config.IgnoreInvalid))

	return strings.Join(params, "&")
}

// GetBatchPrices retrieves latest prices for multiple feeds
func (ps *PriceService) GetBatchPrices(ctx context.Context, priceIDs []string) (*BatchPriceUpdate, error) {
	update := &BatchPriceUpdate{
		Updates: make(map[string]PriceData),
	}

	// Check cache first
	allCached := true
	for _, id := range priceIDs {
		if cached, ok := ps.cache.Load(id); ok {
			update.Updates[id] = cached.(PriceData)
		} else {
			allCached = false
			break
		}
	}

	// If all data is cached, return immediately
	if allCached {
		return update, nil
	}

	// Fetch fresh data for IDs not in cache
	// queryParams := ps.buildQueryParams(priceIDs)
	// resp, err := ps.fetchPrices(ctx, queryParams)
	// if err != nil {
	// 	return nil, err
	// }

	// // Update cache and build response
	// for id, price := range resp.Updates {
	// 	ps.cache.Store(id, price)
	// 	update.Updates[id] = price
	// }

	return update, nil
}
