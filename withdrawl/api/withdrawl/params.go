package withdrawHandler

type WithdrawBluRequestParams struct {
	PendingWithdrawlId string `query:"pending-withdrawl-id"`
	Amount             string `query:"amount"`
	WalletAddress      string `query:"wallet-address"`
}

type WithdrawBlpRequestParams struct{}
