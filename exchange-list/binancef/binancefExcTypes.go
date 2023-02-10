package binancef

import exc "github.com/yin75620/crypto-berserker/exchange"

func (ut *UserTrade) ToExcUserTrade() exc.UserTrade {
	eut := exc.UserTrade{}

	eut.Qty = ut.Qty
	eut.Price = ut.Price
	eut.QuoteQty = ut.QuoteQty
	eut.Side = ut.Side
	eut.Symbol = ut.Symbol
	eut.Time = ut.Time
	return eut
}

func (ut *UserTrade) ToExcUserTradeKey() exc.UserTradeKey {
	eutKey := exc.UserTradeKey{}
	eutKey.Symbol = ut.Symbol
	eutKey.Time = ut.Time / 1000 * 1000
	return eutKey
}
