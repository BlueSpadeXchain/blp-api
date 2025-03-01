package escrow

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

var chainRpcMap = map[string]string{
	"":         "ws://127.0.0.1:8545",
	"31337":    "ws://127.0.0.1:8545",
	"0x03106A": "ws://testnet-rpc.bitlayer.org",
	"200810":   "ws://testnet-rpc.bitlayer.org",
	"0x4268":   "ws://ethereum-holesky-rpc.publicnode.com",
	//"17000":    "wss://ethereum-holesky-rpc.publicnode.com",
	"17000":    "wss://holy-hidden-cherry.ethereum-holesky.quiknode.pro/05f81ada01537c3719d152120293816b3835f642",
	"0xAA36A7": "ws://ethereum-sepolia.publicnode.com",
	"11155111": "ws://ethereum-sepolia.publicnode.com",
	"0xF35A":   "ws://rpc.devnet.citrea.xyz",
	"62298":    "ws://rpc.devnet.citrea.xyz",
	"998":      "ws://api.hyperliquid-testnet.xyz/evm",
	"8453":     "wss://holy-hidden-cherry.base-mainnet.quiknode.pro/05f81ada01537c3719d152120293816b3835f642",
}

func GetChainRpc(chainId string) (string, error) {
	if jsonrpc, found := chainRpcMap[chainId]; found {
		return jsonrpc, nil
	}

	return "", fmt.Errorf("unsupporting chain id: %s", chainId)
}

// "chain": {"escrowAddress"}
var addressMap = map[string]string{
	"":      "0x5FbDB2315678afecb367f032d93F642f64180aa3",
	"31337": "0x5FbDB2315678afecb367f032d93F642f64180aa3", // anvil default
	"17000": "0x73Ae6bC869286f0b0D67483538593adB15c7f66f",
	"8453":  "0x88a03091ea64a90938d0f0906FeD15B57a36F5C8", //base
}

func getAddress(chainId string) (common.Address, error) {
	if escrowAddress, found := addressMap[chainId]; found {
		return common.HexToAddress(escrowAddress), nil
	}

	return common.Address{}, fmt.Errorf("unsupporting chain id: %s", chainId)
}

// escrow, blu, usdc
var addresessMap = map[string][3]string{
	"":      {"0x5FbDB2315678afecb367f032d93F642f64180aa3", "", ""},
	"31337": {"0x5FbDB2315678afecb367f032d93F642f64180aa3", "", ""},
	"17000": {"0x73Ae6bC869286f0b0D67483538593adB15c7f66f", "0x7711C2219a436B48cA03f0740fB7EbA87C4a439e", "0x31ab43583dD532FE8E00a521322338a8E2bB0C4B"},
	"8453":  {"0x88a03091ea64a90938d0f0906FeD15B57a36F5C8", "0x7711C2219a436B48cA03f0740fB7EbA87C4a439e", "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"},
}

func getAddresses(chainId string) (common.Address, common.Address, common.Address, error) {
	if addresses, found := addresessMap[chainId]; found {

		return common.HexToAddress(addresses[0]), common.HexToAddress(addresses[1]), common.HexToAddress(addresses[2]), nil
	}
	return common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("unsupporting chain id: %s", chainId)
}
