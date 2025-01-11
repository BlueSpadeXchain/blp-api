package orderHandler

import (
	"fmt"
)

var PriceFeedIds = []string{
	"e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43", // BTC/USD
	"ff61491a931112ddf1bd8147cd1b641375f79f5825126d665480874634fd0ace", // ETH/USD
}

var pairMap = map[string]string{
	"btcusd": "e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43",
	"usdbtc": "e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43",
	"ethusd": "ff61491a931112ddf1bd8147cd1b641375f79f5825126d665480874634fd0ace",
	"usdeth": "ff61491a931112ddf1bd8147cd1b641375f79f5825126d665480874634fd0ace",
}

func getPair(pairString string) (string, error) {
	if pairHex, found := pairMap[pairString]; found {
		return pairHex, nil
	}

	return "", fmt.Errorf("unsupported pair name: %s", pairString)
}
