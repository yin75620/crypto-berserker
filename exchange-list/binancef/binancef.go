package binancef

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	exc "github.com/yin75620/crypto-berserker/exchange"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
	"github.com/yin75620/crypto-berserker/exchange/tool"
	"github.com/yin75620/crypto-berserker/jmath"
)

var (
	apiURL    = "https://fapi.binance.com"
	apiPrefix = "/fapi/"
)

type JArray exc.JArray

func NewBinancef(c *http.Client) *binance {
	bn := binance{}
	bn.orderBookCenter = ob.NewOrderBookCenter(NewSocket())
	bn.client = c
	bn.account.MakerFee = 0.0002
	bn.account.TakerFee = 0.0004
	bn.account.Leverage = 50
	bn.marketInfos = make(map[string]exc.MarketInfo)
	bn.positions = map[string]Position{}
	bn.init = NewBinancefInit()
	bn.init.IniSetting("main.ini")

	return &bn
}

type binance struct {
	client          *http.Client
	orderBookCenter *ob.OrderBookCenter
	account         exc.Account
	init            *BinancefInit
	marketInfos     map[string]exc.MarketInfo
	positions       map[string]Position
}

// implement exchange
func (bn *binance) GetWallet() exc.Wallet {
	res, err := bn.doSignRequest("GET", "v2/account", exc.JArray{})
	if err != nil {
		fmt.Println("Binance", err)
	}

	account := Account{}

	err = json.Unmarshal(res, &account)
	if err != nil {
		log.Fatal("GetWallet json.Unmarshal", err)
	}

	wallet := exc.Wallet{}
	for _, asset := range account.Assets {
		wallet.Balances = append(wallet.Balances, exc.NewBalance(asset.Asset, asset.MaxWithdrawAmount, asset.MarginBalance, asset.MarginBalance))
	}

	for _, position := range account.Positions {
		bn.positions[position.Symbol] = position
	}

	bn.account.WalletInfo = wallet
	bn.account.UnrealizedPnL = account.TotalUnrealizedProfit

	return wallet
}

func (bn *binance) GetBalance() {
	res, _ := bn.doSignRequest("GET", "v2/balance", exc.JArray{})
	fmt.Println(string(res))
}

// for implement
func (bn *binance) GetMaxOrderUSD(symbol string) float64 {
	if value, ok := bn.positions[symbol]; ok {
		return math.Min(value.MaxNotional, float64(value.Leverage)*bn.account.WalletInfo.GetAllBalanceFreeUSDValue())
	} else {
		//default amount
		return bn.account.WalletInfo.GetAllBalanceFreeUSDValue() * bn.account.Leverage
	}
}

func (bn *binance) getUserTrades(symbol string, startTime time.Time, endTime time.Time) []UserTrade {

	jarray := exc.JArray{
		"symbol": symbol,
	}
	zero := time.Time{}
	if startTime != zero {
		jarray.Add(exc.JArray{"startTime": startTime.UnixMilli()})
	}
	if endTime != zero {
		jarray.Add(exc.JArray{"endTime": endTime.UnixMilli()})
	}

	res, _ := bn.doSignRequest("GET", "v1/userTrades", jarray)
	fmt.Println(string(res))

	userTrades := []UserTrade{}
	err := json.Unmarshal(res, &userTrades)
	if err != nil {
		fmt.Println(err)
	}

	return userTrades
}

func (bn *binance) GetTightUserTrades(symbol string) map[exc.UserTradeKey]exc.UserTrade {
	return bn.GetTightUserTradesWithTime(symbol, time.Now().Add(-24*time.Hour), time.Now())
}

func (bn *binance) GetTightUserTradesWithTime(symbol string, startTime time.Time, endTime time.Time) map[exc.UserTradeKey]exc.UserTrade {
	userTrades := bn.getUserTrades(symbol, startTime, endTime)

	euts := map[exc.UserTradeKey]exc.UserTrade{}
	for _, v := range userTrades {
		eut := v.ToExcUserTrade()
		eutKey := v.ToExcUserTradeKey()

		if value, ok := euts[eutKey]; ok {
			value.Combine(eut)
			euts[eutKey] = value
		} else {
			euts[eutKey] = eut
		}
	}
	return euts
}

func (bn *binance) Prepare() []byte {
	bn.prepareMarketInfo()
	response := bn.GetWallet()
	log.Println(response)

	return []byte(fmt.Sprintf("%v", response))
}

func (bn *binance) GetAccountInfo() []byte {

	return bn.Prepare()
}

func (bn *binance) PostOrder(order exc.ExchangeOrder) (string, error) {
	body := exc.JArray{
		"symbol":      order.CoinPair.GetLinkMakertNameUpper(),
		"side":        order.Side,      //"BUY",
		"type":        order.OrderType, //"LIMIT",
		"timeInForce": "GTC",
		"quantity":    order.Size,
		"price":       order.Price,
	}

	resByte, _ := bn.doSignRequest("POST", "v1/order", body)
	return string(resByte), nil
}

func (bn *binance) GetFee() exc.Fee {
	fee := exc.Fee{}
	fee.Taker = 0.00075
	fee.Maker = 0.00075
	return fee
}
func (bn *binance) GetName() string {
	return "Binancef"
}
func (bn *binance) GetMarketInfo(coinPair exc.CoinPair) exc.MarketInfo {
	return bn.doGetMarketInfo(coinPair.GetLinkMakertNameUpper())
}

func (bn *binance) doGetMarketInfo(name string) exc.MarketInfo {
	return bn.marketInfos[name]
}

func (bn *binance) GetVolumeByTotal(total, price float64) float64 {
	return total / price
}

func (bn *binance) GetFuturesAskBidPair(futures exc.Futures) (exc.PricePair, exc.PricePair) {
	booker := ob.StartOrderBook(bn.orderBookCenter, strings.ToLower(futures.GetLinkMarketName()))
	return booker.GetFirstPricePair()
}

func (bn *binance) GetAccount() exc.Account {
	return bn.account
}

func (bn *binance) PostFuturesOrder(order exc.FuturesOrder) (string, error) {

	symbol := order.Futures.GetLinkMarketName()
	var merketInfo = bn.doGetMarketInfo(symbol)

	body := exc.JArray{
		"symbol": strings.ToUpper(symbol),
		"side":   strings.ToUpper(order.Side),              //"BUY",
		"type":   strings.ToUpper(string(order.OrderType)), //"LIMIT",

		"quantity": jmath.FloatFloorByFloat(order.Size, merketInfo.VolumeIncrement),
	}
	if order.OrderType == exc.LIMIT {
		body["TimeInForce"] = "GTC"
		body["price"] = jmath.FloatFloorByFloat(order.Price, merketInfo.PriceIncrement)
	}

	resByte, err := bn.doSignRequest("POST", "v1/order", body)
	log.Println(fmt.Sprintf("%s", resByte))

	orderResponse := OrderResponse{}
	json.Unmarshal(resByte, &orderResponse)
	if err != nil {
		fmt.Println(err)
		return string(resByte), err
	}

	if orderResponse.Code != 0 {
		s := fmt.Sprint("Code =", orderResponse.Code, orderResponse.Msg)
		fmt.Println(s)
		return string(resByte), errors.New(s)
	}

	return string(resByte), err

	/*bo := BybitOrder{}
	bo.Side = strings.Title(order.CommodityOrder.Side)
	bo.Symbol = strings.ToUpper(order.Futures.GetLinkMarketName())
	bo.OrderType = strings.Title(string(order.CommodityOrder.OrderType))
	bo.Quantity = int64(order.CommodityOrder.Size)
	bo.Price = order.CommodityOrder.Price
	bo.TimeInForce = "GoodTillCancel"

	return bb.doPostOrder(bo)*/
	//return "", nil
}

func (bn *binance) PostCancelAllOrder(fu exc.Futures) {
	symbol := fu.GetLinkMarketName()

	body := exc.JArray{
		"symbol": strings.ToUpper(symbol),
	}

	resByte, _ := bn.doSignRequest("DELETE", "v1/allOpenOrders", body)
	log.Println(fmt.Sprintf("%s", resByte))
}

func (bn *binance) PostLeverage(symbol string, leverage int) {

	body := exc.JArray{
		"symbol":   symbol,
		"leverage": leverage,
	}

	resByte, _ := bn.doSignRequest("POST", "v1/leverage", body)
	fmt.Println(fmt.Printf("%s", resByte))

}

func (bn *binance) getExchangeInfo() ExchangeInfo {

	exchangeInfo := ExchangeInfo{}
	body := exc.JArray{}
	resByte, _ := bn.doRequest("GET", "v1/exchangeInfo", false, body)

	json.Unmarshal(resByte, &exchangeInfo)
	return exchangeInfo
	//log.Println(bn.exchangeInfo)

}

// call it before order
func (bn *binance) prepareMarketInfo() {
	exInfo := bn.getExchangeInfo()

	if symbols := exInfo.Symbols; len(symbols) != 0 {
		for _, value := range symbols {
			marketInfo := exc.MarketInfo{}
			marketInfo.Name = value.Symbol
			for _, fr := range value.Filters {
				if fr.FilterType == "PRICE_FILTER" {
					marketInfo.PriceIncrement = fr.TickSize
				}
				if fr.FilterType == "MARKET_LOT_SIZE" {
					marketInfo.VolumeIncrement = fr.StepSize
				}
			}
			bn.marketInfos[marketInfo.Name] = marketInfo
		}
	}
}

func (bn *binance) prepareLeverage() {
	body := exc.JArray{}
	resByte, _ := bn.doSignRequest("GET", "v1/leverageBracket", body)

	lb := []LeverageBracket{}
	json.Unmarshal(resByte, &lb)

	//fmt.Println(string(resByte))

	for _, value := range lb {
		if value.Symbol == "BTCUSDT" {
			for _, bracket := range value.Brackets {
				tool.PrintStructNameValue(bracket)
			}
		}
	}
}

func (bn *binance) setAllLeverage() {
	body := exc.JArray{}
	resByte, _ := bn.doSignRequest("GET", "v1/leverageBracket", body)

	lb := []LeverageBracket{}
	json.Unmarshal(resByte, &lb)

	for _, value := range lb {
		for _, bracket := range value.Brackets {
			total := float64(bracket.InitialLeverage) * bn.account.WalletInfo.GetAllBalanceUSDValue()
			if total < float64(bracket.NotionalCap) {
				bn.PostLeverage(value.Symbol, bracket.InitialLeverage)
				break
			}
		}
	}
}

func (bn *binance) getKline() {
	body := exc.JArray{
		"symbol":   "BTCUSDT",
		"interval": "1m",
		//"startTime": 1674724440000,
		//"endTime":LONG
		"limit": 3,
	}

	start := time.Now()
	resByte, _ := bn.doSignRequest("GET", "v1/klines", body)
	d := time.Since(start)
	fmt.Println(d)
	body.ToValues()
	klr := [][]interface{}{}
	json.Unmarshal(resByte, &klr)

	candles := []CandleStick{}
	for _, value := range klr {
		cs := CandleStick{}
		cs.SetByJArray(value)

		candles = append(candles, cs)
	}
	fmt.Println(string(resByte))
	fmt.Println(candles)

}

type QuoteResponse struct {
	//Timestamp float64         `"json:time"`
	LastUpdateId float64         `"json:lastUpdateId"`
	TradeData    exc.PriceStatus `"json:data"`
}

func (qr *QuoteResponse) setBy(json map[string]interface{}) {
	qr.LastUpdateId = json["lastUpdateId"].(float64)
	qr.TradeData.SetByJArray(json)
}
func (bn *binance) GetAskBidPair(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	booker := ob.StartOrderBook(bn.orderBookCenter, coinPair.GetLinkMakertName())
	return booker.GetFirstPricePair()
}

func (bn *binance) GetAskBidPairFromWeb(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	const minLimit = 5
	/*reqString := fmt.Sprintf("depth?symbol=%s&limit=%d",
	strings.ToUpper(coinPair.GetLinkMakertName()),
	minLimit)*/
	reqArray := exc.JArray{}
	reqArray["symbol"] = strings.ToUpper(coinPair.GetLinkMakertName())
	reqArray["limit"] = minLimit
	resByte, _ := bn.doAPIRequest("GET", "v1/depth", reqArray)
	fmt.Println(string(resByte))

	var resJson map[string]interface{}
	err := json.Unmarshal(resByte, &resJson)
	if err != nil {
		log.Fatal(err)
	}

	quoteResponse := QuoteResponse{}
	quoteResponse.setBy(resJson)

	askPair, _ := quoteResponse.TradeData.GetPair(1, exc.Ask)
	bidPair, _ := quoteResponse.TradeData.GetPair(1, exc.Bid)

	return askPair, bidPair
}
func (bn *binance) doSignRequest(method, apiName string, body exc.JArray) ([]byte, error) {
	return bn.doRequest(method, apiName, true, body)
}

func (bn *binance) doAPIRequest(method, apiName string, body exc.JArray) ([]byte, error) {
	return bn.doRequest(method, apiName, false, body)
}

func (bn *binance) doRequest(method, apiName string, isSign bool, body exc.JArray) ([]byte, error) {
	client := bn.client

	var res []byte
	fullURL := fmt.Sprintf("%s%s%s", apiURL, apiPrefix, apiName)

	sendContent := body.ToValues().Encode()
	if isSign {
		ts := exc.GetTimeSpan()
		body["timestamp"] = ts
		body["recvWindow"] = 10000
		bodyEncode := body.ToValues().Encode()

		raw := fmt.Sprintf("%s", bodyEncode)
		sign, err := exc.GetParamHmacSHA256HexSign(bn.init.SecretKey, raw)
		if err != nil {
			fmt.Println(err)
		}
		finalBody := exc.JArray{
			"signature": sign,
		}
		sendContent = fmt.Sprintf("%s&%s", bodyEncode, finalBody.ToValues().Encode())
	}

	finalURL := fullURL + "?" + sendContent
	//fmt.Println(finalURL)

	req, err := http.NewRequest(method, finalURL, bytes.NewBuffer([]byte("")))
	if err != nil {
		log.Println(err)
		return res, err
	}

	req.Header.Set("Content-Type", "content-type application/x-www-form-urlencoded")
	req.Header.Add("X-MBX-APIKEY", bn.init.Key)

	resByte, err := exc.SendRequest(client, req)
	if err != nil {
		log.Println("binance doRequest err:", err)
	}

	return resByte, err
}
