package ftx

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

type OrderBookResponse struct {
	Result  exc.PriceStatus `json:"result"`
	Success bool            `json:"success"`
}

type FtxInit struct {
	ApiKey       string
	ApiSecretKey string
	SubAccount   string
}

type Ftx struct {
	client   *http.Client
	initData FtxInit
}

func NewFtx(c *http.Client, initData FtxInit) *Ftx {
	ftx := &Ftx{}
	ftx.client = c
	ftx.initData = initData
	return ftx
}

var (
	apiURL    = "https://ftx.com/api"
	apiPrefix = "/api/"
)

// implement exchange
func (ftx *Ftx) GetFee() exc.Fee {
	fee := exc.Fee{}
	fee.Deposit = 0
	fee.WithDrawl = 0
	fee.Taker = 0.00063175
	fee.Deposit = 0.0001805
	return fee
}

func (ftx *Ftx) GetAccountInfo() []byte {
	return ftx.doGet("account", "")
}

func (ftx *Ftx) GetStructFills() []FillResponse {
	res := ftx.GetFills()
	type tempResponse struct {
		Result []FillResponse `json:"result"`
	}

	tempRes := tempResponse{}

	json.Unmarshal(res, &tempRes)
	return tempRes.Result
}

func (ftx *Ftx) GetFills() []byte {
	return ftx.doGet("fills", "")
}

func (ftx *Ftx) GetCoins() []byte {
	return ftx.doGet("coins", "")
}

func (ftx *Ftx) GetMarkets() []byte {
	return ftx.doGet("markets", "")
}

func (ftx *Ftx) GetMarket(name string) []byte {
	path := fmt.Sprintf("markets/%s", name)
	return ftx.doGet(path, "")
}

func (ftx *Ftx) GetAskBidPair(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	resb := ftx.GetOrderBookResponse(coinPair.GetMarketName(), depth)
	askPair, _ := resb.Result.GetPair(1, exc.Ask)
	bidPair, _ := resb.Result.GetPair(1, exc.Bid)
	return askPair, bidPair
}

func (ftx *Ftx) GetOrderBookResponse(marketName string, depth int) OrderBookResponse {
	response := ftx.GetOrderBook(marketName, depth)
	var bookResponse OrderBookResponse
	json.Unmarshal(response, &bookResponse)
	return bookResponse
}

func (ftx *Ftx) GetOrderBook(marketName string, depth int) []byte {
	path := fmt.Sprintf("/markets/%s/orderbook?depth=%d", marketName, depth)
	response := ftx.doGet(path, "")
	return response
}

type EOrderType string

const (
	LIMIT  EOrderType = "limit"
	MARKET EOrderType = "market"
)

type FtxOrder struct {
	Market    string     `json:"market"`
	Side      string     `json:"side"`
	Price     float64    `json:"price"`
	Size      float64    `json:"size"`
	OrderType EOrderType `json:"order_type"`
	//ReduceOnly bool       `json:"reduceOnly"`
}

func (fo *FtxOrder) setBy(order exc.ExchangeOrder) {
	fo.Market = order.Market
	fo.Side = order.Side
	fo.Price = order.Price
	fo.Size = order.Size
	fo.OrderType = EOrderType(order.OrderType)
}

//下訂單
func (ftx *Ftx) PostOrder(order exc.ExchangeOrder) (string, error) {

	fo := FtxOrder{}
	fo.setBy(order)

	request, err := json.Marshal(fo)
	if err != nil {
		log.Fatal(err)
	}
	body := string(request)
	log.Println(fmt.Sprintf("body:%s", body))
	response := ftx.doPost("orders", body)
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

func (ftx *Ftx) doGet(apiName, body string) []byte {
	return ftx.doRequest("GET", apiName, body)
}

func (ftx *Ftx) doPost(apiName, body string) []byte {
	return ftx.doRequest("POST", apiName, body)
}

func (ftx *Ftx) doRequest(method, apiName, body string) []byte {
	client := ftx.client

	var res []byte

	fullUrl := fmt.Sprintf("%s/%s", apiURL, apiName)
	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Println(err)
		return res
	}

	req.Header.Set("Content-Type", "application/json")
	ftx.addHeader(&req.Header, method, apiName, body)

	return exc.SendRequest(client, req)
}

func (ftx *Ftx) addHeader(header *http.Header, reqMethod, path, body string) {
	ts := exc.GetTimeSpanStr(exc.GetTimeSpan())

	initData := ftx.initData

	header.Add("FTX-KEY", initData.ApiKey)
	header.Add("FTX-TS", ts)
	payload := fmt.Sprintf("%s%s%s%s", ts, reqMethod, apiPrefix+path, body)
	sign, _ := exc.GetParamHmacSHA256HexSign(initData.ApiSecretKey, payload)
	header.Add("FTX-SIGN", sign)
	header.Add("FTX-SUBACCOUNT", initData.SubAccount)
}
