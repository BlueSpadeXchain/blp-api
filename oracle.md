# Oracle Free Tier Instance Setup for Go Bot

## SSH Access
1. Open a terminal and use your private SSH key to connect:
   ```bash
   ssh -i /path/to/your/private-key opc@129.146.189.121
   ```
2. Ensure your private key file has the correct permissions:
   ```bash
   chmod 600 /path/to/your/private-key
   ```

---

## Update and Install Dependencies
Once logged into the instance, update the system and install necessary software:
```bash
sudo dnf update -y
sudo dnf install -y git golang
```

---

## Configure Go Environment
1. Check the installed Go version:
   ```bash
   go version
   ```
   If Go isn’t installed, follow Oracle’s official guide to install Go manually.
2. Set up a Go workspace:
   ```bash
   mkdir -p ~/go/src
   echo 'export GOPATH=$HOME/go' >> ~/.bashrc
   echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
   source ~/.bashrc
   ```

---

## Deploy Your Bot
1. Clone your bot’s repository:
   ```bash
   git clone https://github.com/your-repo/your-bot.git ~/go/src/your-bot
   cd ~/go/src/your-bot
   ```
2. Build the bot:
   ```bash
   go build -o rebalancer
   ```
3. Run the bot to ensure it works:
   ```bash
   ./rebalancer
   ```

---

## Configure Bot as a Service
Set up your bot to run automatically:
1. Create a systemd service file:
   ```bash
   sudo nano /etc/systemd/system/rebalancer.service
   ```
2. Add the following:
   ```ini
   [Unit]
   Description=Go Rebalancer Bot
   After=network.target

   [Service]
   User=opc
   WorkingDirectory=/home/opc/go/src/your-bot
   ExecStart=/home/opc/go/src/your-bot/rebalancer
   Restart=on-failure

   [Install]
   WantedBy=multi-user.target
   ```
3. Reload systemd and start the service:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl start rebalancer
   sudo systemctl enable rebalancer
   ```

---

## Networking Configuration
1. **Set Up Firewall Rules**:
   Ensure your bot can communicate with your backend:
   ```bash
   sudo firewall-cmd --permanent --add-port=443/tcp
   sudo firewall-cmd --reload
   ```

2. **Validate Connectivity**:
   Check that the bot can access your backend (e.g., via curl):
   ```bash
   curl -H "Authorization: Bearer your-api-key" https://your-backend-endpoint
   ```

---

## Cost Minimization
Since you’re using the Free Tier:
- The **VM.Standard.E2.1.Micro** instance is free, and your boot volume is already the minimum size (1GB = $0.05/month).
- Regularly monitor **storage usage** with:
  ```bash
  df -h
  ```
- Offload logs (e.g., to S3 or object storage) if logs grow too large.

---

## Scaling Considerations
If you need more bots in the future:
- **Multiple Instances**: You can create additional Free Tier instances in other regions or availability domains.
- **Containerization**: Use Docker to containerize the bot, making it easy to replicate across instances.
- **Lightweight Alternatives**: If capacity becomes a challenge, consider lightweight bot frameworks or regions with lower demand.

---

# Updated `hermes` API for Pyth Hermes

```go
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

	for _, id := range priceIDs {
		params = append(params, fmt.Sprintf("ids[]=%s", id))
	}

	if ps.config.Encoding != "" {
		params = append(params, fmt.Sprintf("encoding=%s", ps.config.Encoding))
	}
	params = append(params, fmt.Sprintf("parsed=%v", ps.config.ParsedEnabled))
	params = append(params, fmt.Sprintf("allow_unordered=%v", ps.config.AllowUnordered))
	params = append(params, fmt.Sprintf("benchmarks_only=%v", ps.config.BenchmarksOnly))
	params = append(params, fmt.Sprintf("ignore_invalid_price_ids=%v", ps.config.IgnoreInvalid))

	return strings.Join(params,
