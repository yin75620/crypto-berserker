package bybilinear

import exc "github.com/yin75620/crypto-berserker/exchange"

func (ut *UserTrade) ToExcUserTrade() exc.UserTrade {
	eut := exc.UserTrade{}

	eut.Qty = ut.ExecQty
	eut.Price = ut.ExecPrice
	eut.QuoteQty = ut.ExecValue
	eut.Side = ut.Side
	eut.Symbol = ut.Symbol
	eut.Time = int64(ut.TradeTimeMs) / 1000 * 1000
	return eut
}

func (ut *UserTrade) ToExcUserTradeKey() exc.UserTradeKey {
	eutKey := exc.UserTradeKey{}
	eutKey.Symbol = ut.Symbol
	eutKey.Time = int64(ut.TradeTimeMs) / 1000 * 1000
	return eutKey
}
