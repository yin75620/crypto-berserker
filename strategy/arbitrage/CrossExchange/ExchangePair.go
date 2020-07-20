package CrossExchange

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	exc "github.com/yin75620/crypto-berserker/exchange"
	common "github.com/yin75620/crypto-berserker/exchange-list/common"
)

type ExchangePair struct {
	exchange exc.Exchange
	ask      exc.PricePair
	bid      exc.PricePair
}

func NewExchangePair(exchange exc.Exchange, ask, bid exc.PricePair) *ExchangePair {
	ep := ExchangePair{}
	ep.exchange = exchange
	ep.ask = ask
	ep.bid = bid
	return &ep
}

func (ep *ExchangePair) FillPair() {

}

type CrossPairJson struct {
	AskExchangeName string
	BidExchangeName string
	AskPair         exc.PricePair
	BidPair         exc.PricePair
	OrderVolume     float64
	Symbol          string
}

type CrossPair struct {
	askExchange  exc.Exchange
	bidExchange  exc.Exchange
	askPricePair exc.PricePair
	bidPricePair exc.PricePair
	orderVolume  float64
	Symbol       string
}

func NewCrossPair(ask, bid exc.Exchange, askPircePair, bidPricePair exc.PricePair, symbol string) *CrossPair {
	cp := CrossPair{}
	cp.askExchange = ask
	cp.bidExchange = bid
	cp.askPricePair = askPircePair
	cp.bidPricePair = bidPricePair
	cp.Symbol = symbol
	return &cp
}

func (cp *CrossPair) IsDefault() bool {
	a, _ := json.Marshal(cp)
	b, _ := json.Marshal(CrossPair{})
	return bytes.Compare(a, b) == 0

	//return cp.askExchange == nil && cp.bidExchange == nil && cp.askPricePair.Price == 0 && cp.bidPricePair.Price == 0
}

func (cp *CrossPair) toJson() CrossPairJson {
	j := CrossPairJson{}
	j.AskExchangeName = cp.askExchange.GetName()
	j.BidExchangeName = cp.bidExchange.GetName()
	j.AskPair = cp.askPricePair
	j.BidPair = cp.bidPricePair
	j.OrderVolume = cp.orderVolume
	j.Symbol = cp.Symbol
	return j
}

func (cp *CrossPair) GetName() string {
	return fmt.Sprintf("A%sB%s", cp.askExchange.GetName(), cp.bidExchange.GetName())
}

func (cp *CrossPair) GetMatchName() string {
	return fmt.Sprintf("A%sB%s", cp.bidExchange.GetName(), cp.askExchange.GetName())
}

func (cp *CrossPair) GetAskInfo() (exc.Exchange, exc.PricePair) {
	return cp.askExchange, cp.askPricePair
}

func (cp *CrossPair) GetBidInfo() (exc.Exchange, exc.PricePair) {
	return cp.bidExchange, cp.bidPricePair
}

func (cp *CrossPair) GetAskVolumeByTotal(total float64) float64 {
	askExchange, askPair := cp.GetAskInfo()
	askVolume := askExchange.GetVolumeByTotal(total, askPair.Price)
	return askVolume
}

func (cp *CrossPair) GetBidVolumeByTotal(total float64) float64 {
	bidExchange, bidPair := cp.GetBidInfo()
	bidVolume := bidExchange.GetVolumeByTotal(total, bidPair.Price)
	return bidVolume
}

//profit
func (cp *CrossPair) GetProfit() float64 {

	aPrice := cp.GetAskPriceWithFee()
	profitNumber := cp.GetProfitNumber()

	profit := profitNumber / aPrice

	return profit
}

func (cp *CrossPair) GetProfitNumber() float64 {
	aPrice := cp.GetAskPriceWithFee()
	bPrice := cp.GetBidPriceWithFee()

	// 出現錯誤，放慢速度
	if aPrice <= 0 {
		log.Println("laPrice <= 0")
		return 0
	}
	return bPrice - aPrice
}

func (cp *CrossPair) GetAskPriceWithFee() float64 {
	return cp.askPricePair.Price * (1.0 + cp.askExchange.GetFee().Taker)
}

func (cp *CrossPair) GetBidPriceWithFee() float64 {
	return cp.bidPricePair.Price * (1.0 - cp.bidExchange.GetFee().Taker)
}

func (cp *CrossPair) GetProfitString() string {

	aPrice := cp.GetAskPriceWithFee()
	bPrice := cp.GetBidPriceWithFee()

	askStr := fmt.Sprintf("ask Exchange:%s, Cprice:%4.4f, S price:%4.4f, total:%.4f", cp.askExchange.GetName(), aPrice, cp.askPricePair.Price, cp.askPricePair.Total())
	bidStr := fmt.Sprintf("bid Exchange:%s, Cprice:%4.4f, S price:%4.4f, total:%.4f", cp.bidExchange.GetName(), bPrice, cp.bidPricePair.Price, cp.bidPricePair.Total())

	profit := 0.0
	if aPrice > 0 {
		profit = (bPrice - aPrice) / aPrice
	}

	return fmt.Sprintf("[%s] A%s,B%s, Profit:%.5f, %s, %s", cp.Symbol, cp.askExchange.GetName(), cp.bidExchange.GetName(), profit, askStr, bidStr)
}

func (cp *CrossPair) GetMinTotalVolume() float64 {

	aVolume := cp.askPricePair.Total()
	bVolume := cp.bidPricePair.Total()
	// 出現錯誤，放慢速度
	if aVolume <= 0 {
		log.Println("aVolume <= 0")
		return 0
	}

	minTotalVolume := math.Min(aVolume, bVolume)
	//log.Println(fmt.Sprintf("minTotalVolume:%f", minTotalVolume))
	return minTotalVolume
}

func (cp *CrossPair) GetMaxHoldVolume() float64 {
	aAccount := cp.askExchange.GetAccount()
	bAccount := cp.bidExchange.GetAccount()
	aw := aAccount.WalletInfo
	bw := bAccount.WalletInfo
	atw := (aw.GetAllBalanceFreeUSDValue() - aAccount.UnrealizedPnL) * aAccount.Leverage
	btw := (bw.GetAllBalanceFreeUSDValue() - bAccount.UnrealizedPnL) * bAccount.Leverage
	return math.Min(atw, btw)
}

func (cp *CrossPair) GetMinPricePercent(matchPair CrossPair) float64 {
	askGap := matchPair.bidPricePair.Price - cp.askPricePair.Price
	bidGap := cp.bidPricePair.Price - matchPair.askPricePair.Price

	askPercent := askGap / cp.askPricePair.Price
	bidPercent := bidGap / cp.bidPricePair.Price

	minPercent := math.Min(askPercent, bidPercent)
	return minPercent

}

func (cp *CrossPair) GetDbInserString() string {

	return fmt.Sprintf(`INSERT INTO crypto.log_cross_exchange_tick 
	(symbol, ask_exchange, ask_c_price, ask_s_price, ask_total_volume, bid_exchange, bid_c_price, bid_s_price, bid_total_volume, profit, min_total_volume,
	 create_time) VALUES ('%s', '%d', '%f', '%f', '%f', '%d', '%f', '%f', '%f', '%f', '%f', '%d');`,
		cp.Symbol,
		common.GetIndexByName(cp.askExchange.GetName()), cp.GetAskPriceWithFee(), cp.askPricePair.Price, cp.askPricePair.Total(),
		common.GetIndexByName(cp.bidExchange.GetName()), cp.GetBidPriceWithFee(), cp.bidPricePair.Price, cp.bidPricePair.Total(),
		cp.GetProfit(),
		cp.GetMinTotalVolume(),
		time.Now().Unix(),
	)
}

// M0S0 表示 MainAsk, SubBid 有艙位
// MASB 表示 MainBid, SubAsk 有艙位
// MBSA 表示無艙位
type PositionStatus int

const (
	M0S0 PositionStatus = iota
	MASB
	MBSA
)

type CrossPosition struct {
	mainExchangeName string
	subExchangeName  string
	mainPricePair    exc.PricePair
	subPricePair     exc.PricePair
	status           PositionStatus
}
