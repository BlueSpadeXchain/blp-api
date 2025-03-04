package withdrawHandler

type WithdrawBluRequestParams struct {
	PendingWithdrawalId string `query:"pending-withdrawal-id"`
	Amount              string `query:"amount"`
	WalletAddress       string `query:"wallet-address"`
	ApiKey              string `query:"api-key"`
}

type WithdrawBalanceRequestParams struct {
	PendingWithdrawalId string `query:"pending-withdrawal-id"`
	Amount              string `query:"amount"`
	WalletAddress       string `query:"wallet-address"`
	ApiKey              string `query:"api-key"`
}
