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
	"17000":    "wss://ethereum-holesky-rpc.publicnode.com",
	"0xAA36A7": "ws://ethereum-sepolia.publicnode.com",
	"11155111": "ws://ethereum-sepolia.publicnode.com",
	"0xF35A":   "ws://rpc.devnet.citrea.xyz",
	"62298":    "ws://rpc.devnet.citrea.xyz",
	"998":      "ws://api.hyperliquid-testnet.xyz/evm",
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
	"31337": "0x5FbDB2315678afecb367f032d93F642f64180aa3",
	"17000": "0x73Ae6bC869286f0b0D67483538593adB15c7f66f",
}

func getAddress(chainId string) (common.Address, error) {
	if escrowAddress, found := addressMap[chainId]; found {
		return common.HexToAddress(escrowAddress), nil
	}

	return common.Address{}, fmt.Errorf("unsupporting chain id: %s", chainId)
}
