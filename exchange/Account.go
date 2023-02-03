package exchange

type Account struct {
	TakerFee      float64
	MakerFee      float64
	Leverage      float64
	WalletInfo    Wallet
	UnrealizedPnL float64 //control by some api
}

type UserTrade struct {
	Price       float64
	Qty         float64
	QuoteQty    float64
	Side        string
	Symbol      string
	Time        int64
	RealizedPnl float64
}

type UserTradeKey struct {
	Symbol string
	Time   int64
}

func (ut *UserTrade) Combine(v UserTrade) {

	if ut.Symbol != v.Symbol {
		return
	}
	if ut.Side != v.Side {
		return
	}
	ut.QuoteQty = v.QuoteQty + ut.QuoteQty
	ut.Qty = v.Qty + ut.Qty
	ut.Price = ut.QuoteQty / ut.Qty
}

type UserTradeMap map[UserTradeKey]UserTrade

func (utMap *UserTradeMap) Near(utk UserTradeKey, nearMili int64) (UserTrade, bool) {
	for k, v := range *utMap {
		gap := k.Time - utk.Time
		if gap <= nearMili && gap >= -nearMili {
			return v, true
		}
	}
	return UserTrade{}, false
}
