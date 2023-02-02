package binance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	exc "github.com/yin75620/crypto-berserker/exchange"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
)

var (
	apiURL    = "https://api.binance.com"
	apiPrefix = "/api/"
)

type JArray exc.JArray

func NewBinance(c *http.Client) *Binance {
	bn := Binance{}
	bn.orderBookCenter = ob.NewOrderBookCenter(NewSocket())
	bn.client = c
	bn.initData = NewBinanceInit()

	return &bn
}

type Binance struct {
	client          *http.Client
	initData        *BinanceInit
	orderBookCenter *ob.OrderBookCenter
}

func (bn *Binance) SetInitByIni(filename string) {
	bn.initData.IniSetting(filename)
}

// implement exchange
func (bn *Binance) GetWallet() exc.Wallet {
	w := exc.Wallet{}
	return w
}

func (bn *Binance) GetAccountInfo() []byte {

	reqArray := exc.JArray{}

	resByte := bn.doSignRequest("GET", "v3/account", reqArray)

	return resByte
}

func (bn *Binance) PostOrder(order exc.ExchangeOrder) (string, error) {
	body := exc.JArray{
		"symbol":      order.CoinPair.GetLinkMakertNameUpper(),
		"side":        order.Side,      //"BUY",
		"type":        order.OrderType, //"LIMIT",
		"timeInForce": "GTC",
		"quantity":    order.Size,
		"price":       order.Price,
	}

	resByte := bn.doSignRequest("POST", "v3/order", body)

	return string(resByte), nil
}

func (bn *Binance) GetFee() exc.Fee {
	fee := exc.Fee{}
	fee.Taker = 0.00075
	fee.Maker = 0.00075
	return fee
}
func (bn *Binance) GetName() string {
	return "Binance"
}
func (bn *Binance) GetMarketInfo(coinPair exc.CoinPair) exc.MarketInfo {
	return exc.MarketInfo{}
}

// not yet finish implement
func (bn *Binance) GetMaxOrderUSD(symbol string) float64 {
	log.Println("Binance GetMaxOrderUSD not yet finish implement")
	return 100
}
func (bn *Binance) GetTightUserTrades(symbol string) map[exc.UserTradeKey]exc.UserTrade {
	panic("interface not yet implement")
	//return map[exc.UserTradeKey]exc.UserTrade{}
}

func (bn *Binance) GetVolumeByTotal(total, price float64) float64 {
	return total / price
}

func (bn *Binance) GetFuturesAskBidPair(futures exc.Futures) (exc.PricePair, exc.PricePair) {
	return exc.PricePair{}, exc.PricePair{}
}

func (bn *Binance) DeleteAllOrders(symbol string) []byte {
	reqArray := exc.JArray{}
	reqArray["symbol"] = symbol

	resByte := bn.doSignRequest("DELETE", "v3/openOrders", reqArray)

	return resByte
}

// 未實作 UnrealizedPnL
func (bn *Binance) GetAccount() exc.Account {
	return exc.Account{}
}

func (bn *Binance) PostFuturesOrder(order exc.FuturesOrder) (string, error) {
	/*bo := BybitOrder{}
	bo.Side = strings.Title(order.CommodityOrder.Side)
	bo.Symbol = strings.ToUpper(order.Futures.GetLinkMarketName())
	bo.OrderType = strings.Title(string(order.CommodityOrder.OrderType))
	bo.Quantity = int64(order.CommodityOrder.Size)
	bo.Price = order.CommodityOrder.Price
	bo.TimeInForce = "GoodTillCancel"

	return bb.doPostOrder(bo)*/
	return "", nil
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
func (bn *Binance) GetAskBidPair(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	booker := ob.StartOrderBook(bn.orderBookCenter, coinPair.GetLinkMakertName())
	return booker.GetFirstPricePair()
}

func (bn *Binance) GetAskBidPairFromWeb(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	const minLimit = 5
	/*reqString := fmt.Sprintf("depth?symbol=%s&limit=%d",
	strings.ToUpper(coinPair.GetLinkMakertName()),
	minLimit)*/
	reqArray := exc.JArray{}
	reqArray["symbol"] = strings.ToUpper(coinPair.GetLinkMakertName())
	reqArray["limit"] = minLimit
	resByte := bn.doAPIRequest("GET", "v1/depth", reqArray)
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
func (bn *Binance) doSignRequest(method, apiName string, body exc.JArray) []byte {
	return bn.doRequest(method, apiName, true, body)
}

func (bn *Binance) doAPIRequest(method, apiName string, body exc.JArray) []byte {
	return bn.doRequest(method, apiName, false, body)
}

func (bn *Binance) doRequest(method, apiName string, isSign bool, body exc.JArray) []byte {
	client := bn.client

	var res []byte
	fullURL := fmt.Sprintf("%s%s%s", apiURL, apiPrefix, apiName)

	sendContent := body.ToValues().Encode()
	if isSign {
		ts := exc.GetTimeSpan()
		body["timestamp"] = ts
		bodyEncode := body.ToValues().Encode()

		raw := fmt.Sprintf("%s", bodyEncode)
		sign, err := exc.GetParamHmacSHA256HexSign(bn.initData.ApiSecretKey, raw)
		if err != nil {
			fmt.Println(err)
		}
		finalBody := exc.JArray{
			"signature": sign,
		}
		sendContent = fmt.Sprintf("%s&%s", bodyEncode, finalBody.ToValues().Encode())
	}

	finalURL := fullURL + "?" + sendContent
	fmt.Println(finalURL)

	req, err := http.NewRequest(method, finalURL, bytes.NewBuffer([]byte("")))
	if err != nil {
		log.Println(err)
		return res
	}

	req.Header.Set("Content-Type", "content-type application/x-www-form-urlencoded")
	req.Header.Add("X-MBX-APIKEY", bn.initData.ApiKey)

	resByte, err := exc.SendRequest(client, req)
	if err != nil {
		log.Println(err)
	}

	return resByte
}
