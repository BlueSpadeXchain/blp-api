package utils

type VersionResponse struct {
	Version string `json:"version"`
}

type Error struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
	Origin  string `json:"origin"`
}

type PriceUpdate struct {
	EmaPrice Price  `json:"ema_price"`
	ID       string `json:"id"`
	Metadata Meta   `json:"metadata"`
	Price    Price  `json:"price"`
}

type Price struct {
	Conf        string `json:"conf"`
	Expo        int    `json:"expo"`
	Price       string `json:"price"`
	PublishTime int64  `json:"publish_time"`
}

type Meta struct {
	PrevPublishTime    int64 `json:"prev_publish_time"`
	ProofAvailableTime int64 `json:"proof_available_time"`
	Slot               int64 `json:"slot"`
}

type Response struct {
	Binary BinaryData    `json:"binary"`
	Parsed []PriceUpdate `json:"parsed"`
}

type BinaryData struct {
	Encoding string   `json:"encoding"`
	Data     []string `json:"data"`
}
