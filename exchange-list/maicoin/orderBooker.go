package maicoin

import (
	exc "github.com/yin75620/crypto-berserker/exchange"
)

type OrderBooker struct {
	CoinPair              exc.CoinPair
	TopBidPair            exc.PricePair
	BottomAskPair         exc.PricePair
	Exchange              OrderBookExchange
	SocketResponseChannel chan exc.OrderBookSocketResponse
}

func NewOrderBooker(obex OrderBookExchange, coinPair exc.CoinPair, resChannel chan exc.OrderBookSocketResponse) *OrderBooker {
	ob := &OrderBooker{}
	ob.Exchange = obex
	ob.CoinPair = coinPair
	ob.SocketResponseChannel = resChannel
	return ob
}

type OrderBookExchange interface {
	GetNewsetAskBidPair(coinPair exc.CoinPair) (exc.PricePair, exc.PricePair)
}

func (ob *OrderBooker) Start() {
	ob.ForceUpdatePricePair()
	go ob.Receive()
}

func (ob *OrderBooker) Receive() {
	for {
		res := <-ob.SocketResponseChannel

		// 比較市場
		if res.CoinPair.GetMarketName() != ob.CoinPair.GetMarketName() {
			continue
		}

		// 比較金額是否勝利
		isNeedRefresh := false
		otherPricePair := res.PricePair()
		if res.IsAsk() {
			if ob.TopBidPair.Price > otherPricePair.Price {
				if res.IsUpdate() {
					// 成交最大，直接更新資料
					isNeedRefresh = true
				}
			} else if ob.TopBidPair.Price == otherPricePair.Price {
				if res.IsUpdate() {
					ob.TopBidPair.Volume = otherPricePair.Volume
				} else if res.IsAdd() {
					ob.TopBidPair.Volume += otherPricePair.Volume
				} else if res.IsRemove() {
					ob.TopBidPair.Volume -= otherPricePair.Volume
				}
			} else if ob.TopBidPair.Price < otherPricePair.Price {
				if res.IsAdd() {
					ob.TopBidPair = otherPricePair
				} else if res.IsRemove() || res.IsUpdate() {
					// update current price by rest api
					isNeedRefresh = true
				}
			}
		} else if res.IsBid() {
			if ob.BottomAskPair.Price < otherPricePair.Price {
				if res.IsUpdate() {
					// 成交最大，直接更新資料
					isNeedRefresh = true
				}
			} else if ob.BottomAskPair.Price == otherPricePair.Price {
				if res.IsUpdate() {
					ob.BottomAskPair.Volume = otherPricePair.Volume
				} else if res.IsAdd() {
					ob.BottomAskPair.Volume += otherPricePair.Volume
				} else if res.IsRemove() {
					ob.BottomAskPair.Volume -= otherPricePair.Volume
				}
			} else if ob.BottomAskPair.Price > otherPricePair.Price {
				if res.IsAdd() {
					ob.BottomAskPair = otherPricePair
				} else if res.IsRemove() || res.IsUpdate() {
					// update current price by rest api
					isNeedRefresh = true
				}
			}
		}
		if ob.BottomAskPair.Volume == 0 || ob.TopBidPair.Volume == 0 || isNeedRefresh {
			ob.ForceUpdatePricePair()
		}
	}
}

func (ob *OrderBooker) compareAskPair(otherAskPair exc.PricePair) bool {
	if otherAskPair.Price > ob.TopBidPair.Price {
		return true
	}
	return false
}

func (ob *OrderBooker) compareBidPair(otherBidPair exc.PricePair) bool {
	if otherBidPair.Price < ob.BottomAskPair.Price {
		return true
	}
	return false
}

func (ob *OrderBooker) ForceUpdatePricePair() {
	ask, bid := ob.Exchange.GetNewsetAskBidPair(ob.CoinPair)
	ob.BottomAskPair = ask
	ob.TopBidPair = bid
}

type OrderBookCenter struct {
	OrderBookers map[exc.CoinPair]*OrderBooker
}

func NewOrderBookCenter() *OrderBookCenter {
	obc := &OrderBookCenter{}
	return obc
}

// 給一個幣種，直接開始同步該幣種的最高最低價
func (obc *OrderBookCenter) Register(coinPair exc.CoinPair) {

}

func (obc *OrderBookCenter) AddOrderBooker(orderBooker *OrderBooker) {
	obc.OrderBookers[orderBooker.CoinPair] = orderBooker
}
