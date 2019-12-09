package ftxotc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	exc "github.com/yin75620/crypto-berserker/exchange"
)

type PriceType exc.PriceType
type SideType string

const (
	Buy  SideType = exc.Buy
	Sell SideType = exc.Sell
)

type OrderBookResponse struct {
	Result  exc.PriceStatus `json:"result"`
	Success bool            `json:"success"`
}

type FtxotcInit struct {
	ApiKey       string
	ApiSecretKey string
}

type Ftxotc struct {
	client   *http.Client
	initData FtxotcInit
}

func NewFtxotc(c *http.Client, initData FtxotcInit) *Ftxotc {
	Ftxotc := &Ftxotc{}
	Ftxotc.client = c
	Ftxotc.initData = initData
	return Ftxotc
}

var (
	apiURL    = "https://otc.ftx.com/api"
	apiPrefix = "/"
)

// implement exchange
func (Ftxotc *Ftxotc) GetWallet() exc.Wallet {
	w := exc.Wallet{}
	return w
}

func (Ftxotc *Ftxotc) GetFee() exc.Fee {
	fee := exc.Fee{}
	fee.Deposit = 0
	fee.WithDrawl = 0
	fee.Taker = 0.00063175
	fee.Maker = 0.0001805
	return fee
}

func (Ftxotc *Ftxotc) GetName() string {
	return "Ftxotc"
}

func (Ftxotc *Ftxotc) GetMarketInfo(coinPair exc.CoinPair) exc.MarketInfo {
	switch name := coinPair.GetMarketName(); name {
	case "FTT/USD":
		return exc.MarketInfo{VolumeIncrement: 1}
	default:
		return exc.MarketInfo{VolumeIncrement: 0.0001}
	}
}

func (Ftxotc *Ftxotc) GetAccountInfo() []byte {
	return []byte("FTXOTC") //Ftxotc.doGet("account", "")
}

func (ftxotc *Ftxotc) GetBalance() []byte {
	return ftxotc.doGet("balances", "")
}

func (Ftxotc *Ftxotc) GetAllTradingPair() []byte {
	return Ftxotc.doGet("otc/pairs", "")
}

func (Ftxotc *Ftxotc) GetMarket(name string) []byte {
	path := fmt.Sprintf("markets/%s", name)
	return Ftxotc.doGet(path, "")
}

func (ftxotc *Ftxotc) GetQuotes() []byte {
	path := fmt.Sprintf("otc/quotes")
	return ftxotc.doGet(path, "")
}

func (ftxotc *Ftxotc) GetAskBidPair(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {

	buyRes := ftxotc.GetQuoteResponse(coinPair, Buy, 1)
	askPair := buyRes.GetPricePair()
	sellRes := ftxotc.GetQuoteResponse(coinPair, Sell, 1)
	bidPair := sellRes.GetPricePair()
	return askPair, bidPair
}

type QuoteRequest struct {
	BaseCurrency     string  `json:"baseCurrency"`
	QuoteCurrency    string  `json:"quoteCurrency"`
	Side             string  `json:"side"`
	BaseCurrencySize float64 `json:"baseCurrencySize"`
	//QuoteCurrencySize       float64 `json:"quoteCurrencySize"`
	APIOnly                 bool    `json:"apiOnly"`
	SecondsUntilSettlement  float64 `json:"secondsUntilSettlement"`
	CounterpartyAutoSettles bool    `json:"counterpartyAutoSettles"`
	WaitForPrice            bool    `json:"waitForPrice"`
}

type QuoteResponse struct {
	Success bool      `json:"success"`
	Result  QuoteItem `json:"result"`
}

func (qr *QuoteResponse) GetPricePair() exc.PricePair {
	return exc.PricePair{Price: qr.Result.Price, Volume: qr.Result.BaseCurrencySize}
}

func (ftxotc *Ftxotc) GetQuote(coinPair exc.CoinPair, side SideType, size float64) []byte {
	q := QuoteRequest{}
	q.BaseCurrency = coinPair.BaseCoin
	q.QuoteCurrency = coinPair.QuotedCoin
	q.Side = string(side)
	q.BaseCurrencySize = size
	//q.QuoteCurrencySize = 0.253
	q.APIOnly = true
	//q.SecondsUntilSettlement = 20
	//q.CounterpartyAutoSettles = true
	q.WaitForPrice = true

	req, _ := json.Marshal(q)

	fmt.Println(string(req))

	path := fmt.Sprintf("otc/quotes")
	responseByte := ftxotc.doPost(path, string(req))
	return responseByte
}

func (ftxotc *Ftxotc) GetQuoteResponse(coinPair exc.CoinPair, side SideType, size float64) QuoteResponse {
	responseByte := ftxotc.GetQuote(coinPair, side, size)
	response := QuoteResponse{}
	json.Unmarshal(responseByte, &response)
	return response
}

func (ftxotc *Ftxotc) GetQuoteByID(id int64) []byte {
	path := fmt.Sprintf("otc/quotes/%d", id)
	return ftxotc.doGet(path, "")
}

type EOrderType string

const (
	LIMIT  EOrderType = "limit"
	MARKET EOrderType = "market"
)

type FtxotcOrder struct {
	Market    string     `json:"market"`
	Side      string     `json:"side"`
	Price     float64    `json:"price"`
	Size      float64    `json:"size"`
	OrderType EOrderType `json:"order_type"`
	//ReduceOnly bool       `json:"reduceOnly"`
}

func (fo *FtxotcOrder) setBy(order exc.ExchangeOrder) {
	fo.Market = order.CoinPair.GetMarketName()
	fo.Side = order.Side
	fo.Price = order.Price
	fo.Size = order.Size
	fo.OrderType = EOrderType(order.OrderType)
}

//下訂單
func (Ftxotc *Ftxotc) PostOrder(order exc.ExchangeOrder) (string, error) {

	fo := FtxotcOrder{}
	fo.setBy(order)

	request, err := json.Marshal(fo)
	if err != nil {
		log.Fatal(err)
	}
	body := string(request)
	log.Println(fmt.Sprintf("body:%s", body))
	response := Ftxotc.doPost("orders", body)
	log.Println(fmt.Sprintf("%s", response))

	//{"error":"Not enough balances","success":false}
	type OrderResponse struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}
	orderResponse := OrderResponse{}

	json.Unmarshal(response, &orderResponse)

	var resErr error = nil
	if !orderResponse.Success {
		resErr = errors.New(orderResponse.Error)
	}

	return string(response), resErr
}

func (Ftxotc *Ftxotc) doGet(apiName, body string) []byte {
	return Ftxotc.doRequest("GET", apiName, body)
}

func (Ftxotc *Ftxotc) doPost(apiName, body string) []byte {
	return Ftxotc.doRequest("POST", apiName, body)
}

func (Ftxotc *Ftxotc) doRequest(method, apiName, body string) []byte {
	client := Ftxotc.client

	var res []byte

	fullUrl := fmt.Sprintf("%s%s%s", apiURL, apiPrefix, apiName)
	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Println(err)
		return res
	}

	req.Header.Set("Content-Type", "application/json")
	Ftxotc.addHeader(&req.Header, method, apiName, body)

	return exc.SendRequest(client, req)
}

func (Ftxotc *Ftxotc) addHeader(header *http.Header, reqMethod, path, body string) {
	ts := exc.GetTimeSpanStr(exc.GetTimeSpan())

	initData := Ftxotc.initData

	header.Add("FTX-APIKEY", initData.ApiKey)
	header.Add("FTX-TIMESTAMP", ts)
	payload := fmt.Sprintf("%s%s%s%s", ts, reqMethod, apiPrefix+path, body)
	sign, _ := exc.GetParamHmacSHA256HexSign(initData.ApiSecretKey, payload)
	header.Add("FTX-SIGNATURE", sign)
}
