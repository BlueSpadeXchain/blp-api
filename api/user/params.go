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

type UnsignedStakeRequestParams struct {
	UserId    string `query:"user-id"`
	StakeType string `query:"stake-type"`
	Amount    string `query:"amount"`
}
type StakeRequestParams struct {
	StakeId     string `query:"user-id"`
	SignatureId string `query:"signature-id"`
	R           string `query:"r" optional:"true"`
	S           string `query:"s" optional:"true"`
	V           string `query:"v" optional:"true"`
}

type EoaStakeRequestParams struct {
	ChainId      string `query:"chain-id"`
	Block        string `query:"block"`
	BlockHash    string `query:"block-hash"`
	TxHash       string `query:"tx-hash"`
	Sender       string `query:"sender"`
	Receiver     string `query:"receiver"`
	DepositNonce string `query:"nonce"`
	Asset        string `query:"asset"`
	Amount       string `query:"amount"`
	Signature    string `query:"signature"`  // must be 65 byte length
	StakeType    string `query:"stake-type"` // BLU or BLP (liquidity)
}

type GetStakesByUserIdRequestParams struct {
	UserId    string `query:"user-id"`
	StakeType string `query:"stake-type" optional:"true"`
	Limit     string `query:"limit" optional:"true"`
}

type GetStakesByUserAddressRequestParams struct {
	WalletAddress string `query:"wallet-address"`
	WalletType    string `query:"wallet-type"`
	StakeType     string `query:"stake-type" optional:"true"`
	Limit         string `query:"limit" optional:"true"`
}

type UnsignedWithdrawRequestParams struct {
	UserId   string `query:"user-id"`
	Amount   string `query:"amount"`
	ChainId  string `query:"chain-id" optional:"true"`
	Receiver string `query:"receiver" optional:"true"`
}
type SignedWithdrawRequestParams struct {
	WithdrawId  string `query:"withdraw-id"`
	SignatureId string `query:"signature-id"`
	R           string `query:"r" optional:"true"`
	S           string `query:"s" optional:"true"`
	V           string `query:"v" optional:"true"`
}

type UnstakeRequestParams struct {
	UserId    string `query:"user-id"`
	Amount    string `query:"amount"`
	StakeType string `query:"stake-type"` // BLU or BLP
	ChainId   string `query:"chain-id" optional:"true"`
	Receiver  string `query:"receiver" optional:"true"`
}
