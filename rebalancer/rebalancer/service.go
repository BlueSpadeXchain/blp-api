package rebalancer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/BlueSpadeXchain/blp-api/rebalancer/pkg/db"
	"github.com/sirupsen/logrus"
	"github.com/supabase-community/supabase-go"
)

func processOrders(supabaseClient *supabase.Client, pairId string, priceMap []float64, minPrice, maxPrice float64) {
	orders, err := db.GetOrdersParsingRange(supabaseClient, pairId, minPrice, maxPrice)
	if err != nil {
		logrus.Error(fmt.Sprintf("could not fetch orders using pair id %v, minPrice %v, maxPrice %v", pairId, minPrice, maxPrice))
	}
	var ordersActivated, ordersFilled, ordersLiquidated, ordersStopped, pnlProfits, pnlLosses, treasuryProfit, vaultProfit, stakeProfit, liquidityProfit float64
	var collateral float64
	for _, order := range *orders {
		LogCreateOrderResponse2(order)
		// since the priceMap is sequenced in decending order, we can figure out which event was triggered first or multiple events were triggers

		// need two degrees of triggers
		// if limit is active -> if tp/stoploss -> if liq/max

		// liq/max/stop -> end loop for order checking against prices
		for _, mark := range priceMap {
			var payout float64 // init to 0
			var newStatus string // status after all prices
			// we set the take profit value to null when it is hit
			// giant if table is fine for now 
			if order.OrderType == "long" {
				// profits
				if order.TakeProfitPrice <= mark {
					payout += order.TakeProfitValue
					order.TakeProfitValue = 0 // should be null
				}
				if order.MaxPrice <= mark {
					payout += order.MaxValue - order.TakeProfitValue
					newStatus = "filled"
					break;
				}
				// losses
				if order.StopLossPrice >= mark {
					// we forgot to add stoploss value to our order table, for now  calculate it (but needs to be change to reduce latency on rebalancer)
					payout -= order.TakeProfitValue
					newStatus = "liquidated"
					break;
				}
				if order.LiquidationPrice >= mark {
					payout += order.MaxValue - order.TakeProfitValue
					break;
				}
			}
			
		}

		if order.
		collateral = order.Collateral 
	}
	// now need to process the orders
// 	ID:                      1d220d57-2864-4d2c-9631-8811d01714ea
// UserID:                  1d2664a39eee6098
// Order Type:              long
// Leverage:                5.00
// Pair ID:                 e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43
// Order Status:            unsigned
// Collateral:              1000.00
// Entry Price:             50000.00
// Liquidation Price:       40000.00
// Limit Order Price:       150000.00
// Max Price:               10000.00
// Max Value:               0.00
// Stop Loss Price:         0.00
// Take Profit Price:       52000.00
// Take Profit Value:       2600.00
// Take Profit Collateral:  500.00
// Created At:              2025-01-26T06:36:44.080027
// Signed At:               
// Started At:              
// Ended At:       




// 2025-01-26 01:47:29 [warning] price map: map[
// e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43:[104655.67 104655.75] 
// ff61491a931112ddf1bd8147cd1b641375f79f5825126d665480874634fd0ace:[3300.2849714999998 3300.2776952599997]]
// ID:                      706b3244-78f1-409a-a0da-c345e4dcbce3
// UserID:                  1d2664a39eee6098
// Order Type:              short
// Leverage:                3.00
// Pair ID:                 e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43
// Order Status:            unsigned
// Collateral:              2000.00
// Entry Price:             50000.00
// Liquidation Price:       66666.67
// Limit Order Price:       -116666.67
// Max Price:               20000.00
// Max Value:               0.00
// Stop Loss Price:         0.00
// Take Profit Price:       0.00
// Take Profit Value:       0.00
// Take Profit Collateral:  0.00
// Created At:              2025-01-26T06:36:44.080027
// Signed At:               
// Started At:              
// Ended At:   

// we now need to form our group merge orders and users
//if orders2.id = 706b3244-78f1-409a-a0da-c345e4dcbce3

}

func processPrices(supabaseClient *supabase.Client, priceMap map[string][]float64) {

	// Process the collected prices every 5 seconds
	var minPrice, maxPrice float64
	minPrice = 2000000 // set to a large enough value to cover all subvalues
	for id, prices := range priceMap {
		if len(prices) == 0 {
			continue
		}

		for _, price := range prices {
			if price < minPrice {
				minPrice = price
			}
			if price > maxPrice {
				maxPrice = price
			}
		}

		logrus.Warning(fmt.Sprintf("max price: %v, min price: %v", maxPrice, minPrice))
		logrus.Warning(fmt.Sprintf("price map: %v", priceMap))
		processOrders(supabaseClient, id, prices, minPrice, maxPrice)
	}
}

func SubscribeToPriceStream(supabaseClient *supabase.Client, url string, ids []string) {
	var markPriceMap = make(map[string][]float64)
	var mu sync.Mutex

	go func() {
		for {
			time.Sleep(1 * time.Second) // Periodically process prices

			mu.Lock()
			processPrices(supabaseClient, markPriceMap)
			for k := range markPriceMap {
				delete(markPriceMap, k)
			}
			mu.Unlock()
		}
	}()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Fatal("Error creating request:", err)
		return
	}

	q := req.URL.Query()
	for _, id := range ids {
		q.Add("ids[]", id)
	}
	req.URL.RawQuery = q.Encode()

	// Open a connection to the SSE stream
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Fatal("Error connecting to stream:", err)
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		// logrus.Info("Received line:", line)

		response, err := processSSELine(line)
		if err != nil {
			logrus.Error(err.Error())
		}

		parsedResponse, ok := response.(Response)
		if ok && err == nil {
			for _, priceUpdate := range parsedResponse.Parsed {
				// assume exponent always negative
				markPrice, _ := strconv.ParseFloat(priceUpdate.Price.Price, 64)
				exponent := priceUpdate.Price.Expo
				for i := int64(0); i < -int64(exponent); i++ {
					markPrice /= 10
				}

				mu.Lock()
				markPriceMap[priceUpdate.ID] = append(markPriceMap[priceUpdate.ID], markPrice)
				mu.Unlock()

				//logrus.Infof("Received Price Update - ID: %s, MarkPrice: %.6f", priceUpdate.ID, markPrice)

			}
		}
	}

	if err := scanner.Err(); err != nil {
		logrus.Error("Error reading from SSE stream:", err)
		// this should restart and log panic but fix later
	}

	select {}

}

func processSSELine(line string) (interface{}, error) {
	if strings.HasPrefix(line, "data:") {
		// Extract the JSON part of the line (strip "data:")
		jsonData := strings.TrimPrefix(line, "data:")
		jsonData = strings.TrimSpace(jsonData) // Remove any extra spaces or newlines

		var payload Response
		err := json.Unmarshal([]byte(jsonData), &payload)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON payload: %v", err)
		}

		return payload, nil
	} else {
		logrus.Info(fmt.Sprintf("Ignoring unexpected line: %s\n", line))
		return nil, nil
	}
}

/*
data:{
	"binary":{
		"encoding":"hex",
		"data":[
			"504e41550100000003b801000000040d00e351a260da76aec683d5b6179b6f8e817d8d7d0727e84ceaf6dac0deca05a4fc0ce55d0a5addaf6a871ed4aaacb785752db5d24b4a8bfa90158962a2517cc1f400024cb0a66e9f4adcbab8568d8f6032360dbf93def497ccbd02402a8abb7169e3ed1de3fa827f86ba7b3992679dcd7c6136719121254ce96d78309e1bcfe263d00400035876e6d45911aa1c73028be8d41a5412701732c3f7005a86440500310b3eaa0f635a3dfcd9d3757a31f78ad5b3e7c4c7f1561b04c7e1af3d360dee650030ade100049613ef5f6ddbe9271ffb7bd41ee781cbe8c55180bca45b794318a9d199e229f32bf8f44b8da144430d09826406861d13cbc8b9a53f50bd30a616c74d6849d7c101061cd7547a31bc897fc7e2443e9cbee426ca992048d064712f5596a954fa5689374741e6ec2842532fdd480604d34292336c6744150a144658bd9092bf9e87434b000852205ad844251a1349b6fbe1e7d9f27b8681d0628f47f7769142cc421b0b890433dc5ed49285b2717169edfa5b07544f8d53eec122079b5ceb22ccefbd4b63ef000a231ae39d67ed1f15b199ea34018660bfebb4e8f275037667f594868c74a64b362b548d8649136b230159342a6b2761bbca4c4996a6ce21d3504f1ff935bd670c010bb9e2a7ec75a714a5f77760968e50f32d9c8a70295f265e71aff8dc2c52a2d1d73753f36979cf2b2ec788961743b729976675305e7725a1a4862741bc7f9aea33000c5554b32e7c51e4f5383319af3345ab57a434bb6df4251d126ba6b8667fec55366e83b8acc085bbd5c75bcc1d99338f4bf1c75aaf08fe9d007593042cc183aedd000d717e0a95bc5f4db1a8f5a0795258f030b5190d059d6cff387b1b0783cb1012242c27f23aee7ccd09729a5feb72fb22d4b8bccd9d25bc5d6f4e49892026975031000f18fa460fd4967520c9976682c721d9dd0b532bf6d6e286af38998d41975f1eb4562be5a05911d14def552318667eb0d9a8e0e8b9307be23bf112aa700b446030011069b26caf4a5b2afad642c2c2e6efb9d5ee57dc37f39988f22fe7dc2f8001f6fb3fa03dfca0b7bf03dbbf9b02b2885a4785186d552de621b826ed4ef16a87741d00118e4b2346adcfa23e76e089b74753d49fc27e5c041be6c627959d5a51a5d9b8f706a943fa815ebac68606c9807bfa3fa4daec36a4d5a7aebb4fa02ce3229b8a73016773daea00000000001ae101faedac5851e32b9b23b5f9411a8c2bac4aae3ed4dd7b811dd1a72ea4aa7100000000062a89a7014155575600000000000b376d6500002710d86affe1de75b312f0a6f7383f0b0732330ea7c601005500e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b430000088e17d0360600000002305d11befffffff8000000006773daea000000006773daea0000088ac179e42000000001bc5439500b31be0d91a82eae8b2c080cf02a88e5f96a22a52fa70923bc214a8d0f3b0b4d5f4251d4fd655a2cc50ab6b8017b7ea1b2dea695791776ca14d5331491993a0e8576e8e4e5991983c3be76742f55e851f5d47879b0c8728abb48ff3042a66f2ab688922740cb33732f7b765f1c96c91a99a5a7ddb339c155f10ccd5fe1471fd63c6e8c59271b06cd7027380aee4d1a5605dc190821952944fc9168c58f7a5b04d2271cf84cf207de2550df8db4984bb7e473c269e452cf31930a62a0040e4edb21c6a6740d2619be0f0a92c4d15a3355ecc0050863529ff14422e9fbda"
		]
	},
	"parsed":[{
		"id":"e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43",
		"price":{
			"price":"9406377899526",
			"conf":"9401340350",
			"expo":-8,
			"publish_time":1735645930
		},
		"ema_price":{
			"price":"9392044500000",
			"conf":"7454603600",
			"expo":-8,
			"publish_time":1735645930
		},
		"metadata":{
			"slot":188181861,
			"proof_available_time":1735645932,
			"prev_publish_time":1735645930
		}
	}]
}
*/
