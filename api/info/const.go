package infoHandler

const Version string = "BLP API v0.0.5"

var Pairs []string = []string{"ethusd", "btcusd"}

var PairIds []string = []string{"ff61491a931112ddf1bd8147cd1b641375f79f5825126d665480874634fd0ace", "e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43"}

var PairAndIds []Pair = []Pair{
	{Pair: "ethusd", PairId: "ff61491a931112ddf1bd8147cd1b641375f79f5825126d665480874634fd0ace"},
	{Pair: "btcusd", PairId: "e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43"},
}
