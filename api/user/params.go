package userHandler

type SignatureRaw struct {
	V string `query:"v"` // tvm is often a garbage value at least from ts
	R string `query:"r"`
	S string `query:"s"`
}

type Signature struct {
	V uint64
	R uint64
	S uint64
}

type SignatureBytes struct {
	V byte
	R []byte
	S []byte
}

type WithdrawRequestParams struct {
	ChainId       string       `query:"chain-id" optional:"true"` // chain to payout to
	SignatureType string       `query:"sig-type" optional:"true"` // assumes ecdsa secp256k1 only, for now
	Amount        string       `query:"amount"`
	Receiver      string       `query:"receiver"`
	Signature     SignatureRaw `query:"signature"`
}

// now need to make request to backend to add deposits
// BlockNumber: 14
// TxHash: 0xe9f1fe395e55ca3037a5d248b87de7f5c124a2f558a9f8493ce6fa6fe9c8e9fd
// DepositEvent: {0x70997970C51812dc3A010C7d01b50e0d17dc79C8 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 4 0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0 100000000000000000000}

// nothing optional, only accept txs from the escrow listener
type DespositRequestParams struct {
	ChainId      string `query:"chain-id"`
	Block        string `query:"block"`
	BlockHash    string `query:"block-hash"`
	TxHash       string `query:"tx-hash"`
	Sender       string `query:"sender"`
	Receiver     string `query:"receiver"`
	DepositNonce string `query:"nonce"`
	Asset        string `query:"asset"`
	Amount       string `query:"amount"`
	Signature    string `query:"signature"` // must be 65 byte length
}

type UserDataRequestParams struct {
	UserId string `query:"user-id"`
}

type AddAuthorizedWalletRequestParams struct {
	UserId        string       `query:"user-id"`
	Address       string       `query:"address"`
	AddressFormat string       `query:"format"`
	Signature     SignatureRaw `query:"signature"`
}

type RemoveAuthorizedWalletRequestParams struct {
	UserId        string       `query:"user-id"`
	Address       string       `query:"address"`
	AddressFormat string       `query:"format"`
	Signature     SignatureRaw `query:"signature"`
}

type GetUserByUserIdRequestParams struct {
	UserId string `query:"user-id"`
}

type GetUserByUserAddressRequestParams struct {
	Address     string `query:"address"`
	AddressType string `query:"type"` // referance to signature format (ecdsa/secp/edd/etc used by sig validation)
}

type GetDepositsByUserIdRequestParams struct {
	UserId string `query:"user-id"`
}

type GetDepositsByUserAddressRequestParams struct {
	WalletAddress string `query:"wallet-address"`
	WalletType    string `query:"wallet-type"`
}
