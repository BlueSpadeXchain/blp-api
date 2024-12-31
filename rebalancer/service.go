package rebalancer

import (
	"context"
	"sync"
	"time"

	hermes "github.com/BlueSpadeXchain/blp-api/pkg/hermes"
)

type Config struct {
	// Price Feed Config
	PythBaseURL    string
	PriceFeeds     []string
	UpdateInterval time.Duration

	// Bot Config
	RebalanceThreshold float64
	OrderEndpoint      string
	APIKey             string
}

type Service struct {
	config    Config
	priceServ *hermes.PriceService
	bot       *RebalancerBot
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

type RebalancerBot struct {
	orderEndpoint string
	apiKey        string
	// Add other bot-specific fields
}

// func NewPriceService(cfg PriceServiceConfig) (*PriceService, error) {
// 	ps := &PriceService{
// 		config: cfg,
// 		client: &http.Client{
// 			Timeout: time.Second * 10,
// 		},
// 	}

// 	return ps, nil
// }

// type PriceServiceConfig struct {
// 	BaseURL        string
// 	SSEEndpoint    string
// 	Encoding       string // hex or base64
// 	ParsedEnabled  bool
// 	AllowUnordered bool
// 	BenchmarksOnly bool
// 	IgnoreInvalid  bool
// 	ReconnectDelay time.Duration
// 	CacheTTL       time.Duration
// }

// func NewService(cfg Config) (*Service, error) {
// 	ctx, cancel := context.WithCancel(context.Background())

// 	serviceConfig := hermes.PriceServiceConfig{
// 		BaseURL:        cfg.PythBaseURL,
// 		ReconnectDelay: time.Second * 2,
// 	}

// 	ps, err := hermes.NewPriceService(serviceConfig)
// 	if err != nil {
// 		fmt.Printf("error in NewServce: %v", err.Error())
// 		return &Service{}, err
// 	}

// 	return &Service{
// 		config:    cfg,
// 		priceServ: ps,
// 		bot: &RebalancerBot{
// 			orderEndpoint: cfg.OrderEndpoint,
// 			apiKey:        cfg.APIKey,
// 		},
// 		ctx:    ctx,
// 		cancel: cancel,
// 	}, nil
// }

// func (s *Service) Start() error {
// 	// Setup graceful shutdown
// 	sigChan := make(chan os.Signal, 1)
// 	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

// 	// Start price feed
// 	updates, err := s.priceServ.SubscribeFeeds(s.ctx, s.config.PriceFeeds)
// 	if err != nil {
// 		return err
// 	}

// 	// Start processing updates
// 	s.wg.Add(1)
// 	go func() {
// 		defer s.wg.Done()
// 		s.processUpdates(updates)
// 	}()

// 	// Wait for shutdown signal
// 	<-sigChan
// 	s.Shutdown()

// 	return nil
// }

// func (s *Service) processUpdates(updates <-chan map[string]*hermes.HumanReadablePrice) {
// 	for prices := range updates {
// 		// Process prices and make trading decisions
// 		s.bot.HandlePriceUpdates(prices)
// 	}
// }

// func (s *Service) Shutdown() {
// 	s.cancel()
// 	s.wg.Wait()
// }

// func main() {
// 	cfg := Config{
// 		PythBaseURL:        os.Getenv("PYTH_BASE_URL"),
// 		PriceFeeds:         strings.Split(os.Getenv("PRICE_FEED_IDS"), ","),
// 		UpdateInterval:     time.Second * 2,
// 		RebalanceThreshold: 0.02, // 2%
// 		OrderEndpoint:      os.Getenv("ORDER_ENDPOINT"),
// 		APIKey:             os.Getenv("API_KEY"),
// 	}

// 	service := NewService(cfg)
// 	if err := service.Start(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// we need a main function to execute
