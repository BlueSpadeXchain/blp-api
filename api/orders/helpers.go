package orderHandler

import (
	"math"
	"time"

	"github.com/BlueSpadeXchain/blp-api/pkg/db"
)

func validateOrderRequest() error {
	// _, relayAddress, err := utils.EnvKey2Ecdsa()
	// if err != nil {
	// 	return err
	// }

	// chainId := new(big.Int)
	// markPrice := new(big.Int)
	// entryPrice := new(big.Int)
	// liqPrice := new(big.Int)

	// if len(params.Signer) != 40 {
	// 	return fmt.Errorf("invalid signer length: must be 40 characters")
	// }

	// if len(params.Signature) != 130 {
	// 	return fmt.Errorf("invalid signature length: must be 130 characters")
	// }

	// if _, ok := chainId.SetString(params.PerpId, 10); !ok {
	// 	return fmt.Errorf("invalid perpId: %s", params.PerpId)
	// }
	// if _, ok := markPrice.SetString(params.MarkPrice, 10); !ok {
	// 	return fmt.Errorf("invalid markPrice: %s", params.MarkPrice)
	// }
	// if _, ok := entryPrice.SetString(params.EntryPrice, 10); !ok {
	// 	return fmt.Errorf("invalid entryPrice: %s", params.EntryPrice)
	// }
	// if _, ok := liqPrice.SetString(params.LiquidationPrice, 10); !ok {
	// 	return fmt.Errorf("invalid liquidationPrice: %s", params.LiquidationPrice)
	// }

	// signer, err := utils.HexToBytes(params.Signer)
	// if err != nil {
	// 	return err
	// }

	// signature, err := utils.HexToBytes(params.Signature)
	// if err != nil {
	// 	return err
	// }

	// var toHash []byte
	// toHash = append(toHash, chainId.Bytes()...)
	// toHash = append(toHash, markPrice.Bytes()...)
	// toHash = append(toHash, entryPrice.Bytes()...)
	// toHash = append(toHash, liqPrice.Bytes()...)

	// hash := crypto.Keccak256Hash(toHash)

	// sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	// if err != nil {
	// 	return err
	// }

	// if !bytes.Equal(sigPublicKey, relayAddress.Bytes()) || !bytes.Equal(sigPublicKey, signer) {
	// 	return fmt.Errorf("invalid signature or hash")
	// }

	return nil
}

func getFeeScalingFactor() float64 {
	return 0.3
}

func getBaseFee() float64 {
	return 0.001
}

func dynamicLeverageFee(leverage float64) float64 {
	//fee percent = 1/ (1+ scaling factor * log(leverage)) * base fee / 100
	return 1 / (1 + getFeeScalingFactor()*math.Log(leverage)) * getBaseFee()
}

func getPerHourFee() float64 {
	return 0.0001
}

func dynamicUtilizationFee(startTimestamp db.CustomTime, globalBorrowed, globalLiquidity float64) float64 {
	elapsedTime := time.Since(startTimestamp.UTC()).Seconds()

	return getPerHourFee() * (elapsedTime / 3600) * globalBorrowed / globalLiquidity
}
