package infoHandler

type GetPairsResponse struct {
	Pairs []string `json:"pairs"`
}

type Pair struct {
	Pair   string `json:"pair"`
	PairId string `json:"pair-id"`
}
