package infoHandler

import "fmt"

const Version string = "BLP API v0.0.5"

var Pairs []string = []string{"ethusd", "btcusd"}

var PairIds []string = []string{
	"ff61491a931112ddf1bd8147cd1b641375f79f5825126d665480874634fd0ace",
	"e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43",
	"ffff73128917a90950cd0473fd2551d7cd274fd5a6cc45641881bbcc6ee73417",
	"2f95862b045670cd22bee3114c39763a4a08beeb663b145d283c31d7d1101c4f",
	"ef0d8b6fda2ceba41da15d4095d1da392a0d2f8ed0c6c7bc0f4cfac8c280b56d",
	"dcef50dd0a4cd2dcc17e45df1676dcb336a11a61c69df7a0299b0150c672d25c",
	"23d7315113f5b1d3ba7a83604c44b94d79f4fd69af77f804fc7f920a6dc65744",
	"879551021853eec7a7dc827578e8e69da7e4fa8148339aa0d3d5296405be4b1a",
	"72b021217ca3fe68922a19aaf990109cb9d84e9ad004b4d2025ad6f529314419",
	"116da895807f81f6b5c5f01b109376e7f6834dc8b51365ab7cdfa66634340e54"}

var PairAndIds []Pair = []Pair{
	{Pair: "ethusd", PairId: "ff61491a931112ddf1bd8147cd1b641375f79f5825126d665480874634fd0ace"},
	{Pair: "btcusd", PairId: "e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43"},
	{Pair: "moodengusd", PairId: "ffff73128917a90950cd0473fd2551d7cd274fd5a6cc45641881bbcc6ee73417"},
	{Pair: "bnbusd", PairId: "2f95862b045670cd22bee3114c39763a4a08beeb663b145d283c31d7d1101c4f"},
	{Pair: "solusd", PairId: "ef0d8b6fda2ceba41da15d4095d1da392a0d2f8ed0c6c7bc0f4cfac8c280b56d"},
	{Pair: "dogeusd", PairId: "dcef50dd0a4cd2dcc17e45df1676dcb336a11a61c69df7a0299b0150c672d25c"},
	{Pair: "suiusd", PairId: "23d7315113f5b1d3ba7a83604c44b94d79f4fd69af77f804fc7f920a6dc65744"},
	{Pair: "trumpusd", PairId: "879551021853eec7a7dc827578e8e69da7e4fa8148339aa0d3d5296405be4b1a"},
	{Pair: "bonkusd", PairId: "72b021217ca3fe68922a19aaf990109cb9d84e9ad004b4d2025ad6f529314419"},
	{Pair: "pnutusd", PairId: "116da895807f81f6b5c5f01b109376e7f6834dc8b51365ab7cdfa66634340e54"},
}

var pairIdMap = map[string]string{
	"btcusd":     "e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43",
	"usdbtc":     "e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43",
	"ethusd":     "ff61491a931112ddf1bd8147cd1b641375f79f5825126d665480874634fd0ace",
	"usdeth":     "ff61491a931112ddf1bd8147cd1b641375f79f5825126d665480874634fd0ace",
	"moodengusd": "ffff73128917a90950cd0473fd2551d7cd274fd5a6cc45641881bbcc6ee73417",
	"usdmoodeng": "ffff73128917a90950cd0473fd2551d7cd274fd5a6cc45641881bbcc6ee73417",
	"bnbusd":     "2f95862b045670cd22bee3114c39763a4a08beeb663b145d283c31d7d1101c4f",
	"usdbnb":     "2f95862b045670cd22bee3114c39763a4a08beeb663b145d283c31d7d1101c4f",
	"solusd":     "ef0d8b6fda2ceba41da15d4095d1da392a0d2f8ed0c6c7bc0f4cfac8c280b56d",
	"usdsol":     "ef0d8b6fda2ceba41da15d4095d1da392a0d2f8ed0c6c7bc0f4cfac8c280b56d",
	"dogeusd":    "dcef50dd0a4cd2dcc17e45df1676dcb336a11a61c69df7a0299b0150c672d25c",
	"usddoge":    "dcef50dd0a4cd2dcc17e45df1676dcb336a11a61c69df7a0299b0150c672d25c",
	"suiusd":     "23d7315113f5b1d3ba7a83604c44b94d79f4fd69af77f804fc7f920a6dc65744",
	"usdsui":     "23d7315113f5b1d3ba7a83604c44b94d79f4fd69af77f804fc7f920a6dc65744",
	"trumpusd":   "879551021853eec7a7dc827578e8e69da7e4fa8148339aa0d3d5296405be4b1a",
	"usdtrump":   "879551021853eec7a7dc827578e8e69da7e4fa8148339aa0d3d5296405be4b1a",
	"bonkusd":    "72b021217ca3fe68922a19aaf990109cb9d84e9ad004b4d2025ad6f529314419",
	"usdbonk":    "72b021217ca3fe68922a19aaf990109cb9d84e9ad004b4d2025ad6f529314419",
	"pnutusd":    "116da895807f81f6b5c5f01b109376e7f6834dc8b51365ab7cdfa66634340e54",
	"usdpnut":    "116da895807f81f6b5c5f01b109376e7f6834dc8b51365ab7cdfa66634340e54",
}

func getPairId(pairString string) (string, error) {
	if pairHex, found := pairIdMap[pairString]; found {
		return pairHex, nil
	}

	return "", fmt.Errorf("unsupported pair name: %s", pairString)
}
