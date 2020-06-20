package exchange

type Account struct {
	TakerFee      float64
	MakerFee      float64
	Leverage      float64
	WalletInfo    Wallet
	UnrealizedPnL float64 //control by some api
}
