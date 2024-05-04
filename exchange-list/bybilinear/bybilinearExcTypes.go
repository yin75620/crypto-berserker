package bybilinear

import (
	"strings"

	exc "github.com/yin75620/crypto-berserker/exchange"
)

func (ut *TradingHistoryOrder) ToExcUserTrade() exc.UserTrade {
	eut := exc.UserTrade{}

	eut.Qty = ut.ExecQty
	eut.Price = ut.ExecPrice
	eut.QuoteQty = ut.ExecValue
	eut.Side = strings.ToUpper(ut.Side)
	eut.Symbol = ut.Symbol
	eut.Time = int64(ut.ExecTime)
	return eut
}

func (ut *TradingHistoryOrder) ToExcUserTradeKey() exc.UserTradeKey {
	eutKey := exc.UserTradeKey{}
	eutKey.Symbol = ut.Symbol
	eutKey.Time = int64(ut.ExecTime) / 1000 * 1000
	return eutKey
}
