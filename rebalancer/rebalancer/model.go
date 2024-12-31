package rebalancer

// Price contains price information along with confidence and exponent
type Price struct {
	Price       string `json:"price"`
	Conf        string `json:"conf"`
	Expo        int    `json:"expo"`
	PublishTime int64  `json:"publish_time"`
}

// Metadata contains metadata related to the price update
type Metadata struct {
	PrevPublishTime    int64 `json:"prev_publish_time"`
	ProofAvailableTime int64 `json:"proof_available_time"`
	Slot               int   `json:"slot"`
}

// PriceUpdate represents the parsed price update with price and metadata
type PriceUpdate struct {
	ID       string   `json:"id"`
	Price    Price    `json:"price"`
	EmaPrice Price    `json:"ema_price"`
	Metadata Metadata `json:"metadata"`
}

// BinaryData contains the raw binary data and encoding format
type BinaryData struct {
	Encoding string   `json:"encoding"`
	Data     []string `json:"data"`
}

// Response represents the overall response structure
type Response struct {
	Binary BinaryData    `json:"binary"`
	Parsed []PriceUpdate `json:"parsed"`
}
