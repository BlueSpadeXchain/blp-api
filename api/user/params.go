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

type DespositRequestParams struct {
	ChainId       string       `query:"chain-id" optional:"true"` // chain to payout to
	SignatureType string       `query:"sig-type" optional:"true"` // assumes ecdsa secp256k1 only, for now
	Amount        string       `query:"amount"`
	Receiver      string       `query:"receiver"` // needs to be backend whitelisted address, used by block listeners
	Signature     SignatureRaw `query:"signature"`
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

type GetUserByIdRequestParams struct {
	UserId string `query:"user-id"`
}

type GetUserByAddressRequestParams struct {
	Address     string `query:"address"`
	AddressType string `query:"type"` // referance to signature format (ecdsa/secp/edd/etc used by sig validation)
}
