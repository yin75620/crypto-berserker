package binancef

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	exc "github.com/yin75620/crypto-berserker/exchange"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
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
}

// implement exchange
func (bn *binance) GetWallet() exc.Wallet {
	res := bn.doSignRequest("GET", "v1/account", exc.JArray{})

	account := Account{}

	err := json.Unmarshal(res, &account)
	if err != nil {
		log.Fatal("GetWallet json.Unmarshal", err)
	}

	wallet := exc.Wallet{}
	for _, asset := range account.Assets {
		wallet.Balances = append(wallet.Balances, exc.NewBalance(asset.Asset, asset.MaxWithdrawAmount, asset.MarginBalance, asset.MarginBalance))
	}
	bn.account.WalletInfo = wallet
	bn.account.UnrealizedPnL = account.TotalUnrealizedProfit

	return wallet
}

func (bn *binance) Prepare() {
	bn.prepareMarketInfo()
}

func (bn *binance) GetAccountInfo() []byte {

	bn.prepareMarketInfo()
	response := bn.GetWallet()
	log.Println(response)

	return []byte(fmt.Sprintf("%v", response))
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

	resByte := bn.doSignRequest("POST", "v1/order", body)
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

	resByte := bn.doSignRequest("POST", "v1/order", body)
	log.Println(fmt.Sprintf("%s", resByte))
	return string(resByte), nil

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

	resByte := bn.doSignRequest("DELETE", "v1/allOpenOrders", body)
	log.Println(fmt.Sprintf("%s", resByte))
}

func (bn *binance) getExchangeInfo() ExchangeInfo {

	exchangeInfo := ExchangeInfo{}
	body := exc.JArray{}
	resByte := bn.doRequest("GET", "v1/exchangeInfo", false, body)

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
	resByte := bn.doSignRequest("GET", "v1/leverageBracket", body)

	lb := []LeverageBracket{}
	json.Unmarshal(resByte, &lb)

	fmt.Println(string(resByte))

	for _, value := range lb {
		if value.Symbol == "BTCUSDT" {
			for _, bracket := range value.Brackets {
				fmt.Println(bracket)
			}
		}
	}
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
func (bn *binance) doSignRequest(method, apiName string, body exc.JArray) []byte {
	return bn.doRequest(method, apiName, true, body)
}

func (bn *binance) doAPIRequest(method, apiName string, body exc.JArray) []byte {
	return bn.doRequest(method, apiName, false, body)
}

func (bn *binance) doRequest(method, apiName string, isSign bool, body exc.JArray) []byte {
	client := bn.client

	var res []byte
	fullURL := fmt.Sprintf("%s%s%s", apiURL, apiPrefix, apiName)

	sendContent := body.ToValues().Encode()
	if isSign {
		ts := exc.GetTimeSpan()
		body["timestamp"] = ts
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
		return res
	}

	req.Header.Set("Content-Type", "content-type application/x-www-form-urlencoded")
	req.Header.Add("X-MBX-APIKEY", bn.init.Key)

	resByte, err := exc.SendRequest(client, req)
	if err != nil {
		log.Println("binance doRequest err:", err)
	}

	return resByte
}
