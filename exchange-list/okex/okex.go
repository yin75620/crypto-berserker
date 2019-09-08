package okex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/setting"
)

func NewOkex(c *http.Client) *Okex {
	Okex := &Okex{}
	Okex.client = c
	return Okex
}

type Okex struct {
	client       *http.Client
	accountGroup int
}

var (
	apiURL    = "https://www.okex.com"
	apiPrefix = "/api/spot/v3/" //幣幣交易
)

// implement exchange
func (bm *Okex) GetFee() exc.Fee {
	fee := exc.Fee{}
	fee.Deposit = 0
	fee.WithDrawl = 0
	fee.Taker = 0.0015
	fee.Maker = 0.0015
	return fee
}

func (bm *Okex) GetName() string {
	return "Okex"
}

func (bm *Okex) GetMarketInfo(coinPair exc.CoinPair) exc.MarketInfo {
	return exc.MarketInfo{} //not yet implement
}

type QuoteResponse struct {
	exc.PriceStatus
	timeStamp string `json:"timeStamp"`
}

func (qr *QuoteResponse) setBy(json map[string]interface{}) {
	qr.timeStamp = json["timestamp"].(string)

	qr.PriceStatus.SetByJArray(json)
}

func (bm *Okex) GetAskBidPair(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	path := fmt.Sprintf("instruments/%s/book?size=%d",
		coinPair.GetSymbal(), depth)
	resByte := bm.doNormalRequest("GET", path, "")

	var resJson map[string]interface{}

	err := json.Unmarshal(resByte, &resJson)
	if err != nil {
		log.Fatal(err)
	}

	quoteResponse := QuoteResponse{}
	quoteResponse.setBy(resJson)

	askPair, _ := quoteResponse.GetPair(1, exc.Ask)
	bidPair, _ := quoteResponse.GetPair(1, exc.Bid)

	return askPair, bidPair
}

func (bm *Okex) GetAccountInfo() []byte {
	res := bm.doNormalRequest("GET", "accounts", "")
	return res
}

func (bm *Okex) GetProducts() []byte {
	res := bm.doNormalRequest("GET", "products", "")
	return res
}

type OkexOrder struct {
	Coid       string `json:"coid"`       //"xxx...xxx"     a unique identifier of length 32
	Time       int64  `json:"time"`       // 1528988100000   milliseconds since UNIX epoch in UTC
	Symbol     string `json:"symbol"`     //"ETH/BTC"
	OrderPrice string `json:"orderPrice"` //"13.5"          optional, limit price of the order. This field is required for limit orders and stop limit orders.
	//StopPrice  string `json:"stopPrice"`  //"15.7"          optional, stop price of the order. This field is required for stop market orders and stop limit orders.
	OrderQty  string `json:"orderQty"`  //"3.5"
	OrderType string `json:"orderType"` //"limit"         order type, you shall specify one of the following: "limit", "market", "stop_market", "stop_limit".
	Side      string `json:"side"`      //"buy"           "buy" or "sell"
	//PostOnly    bool   `json:"postOnly"`    //true            Optional, if true, the order will either be posted to the limit order book or be cancelled, i.e. the order cannot take liquidity; default value is false
	//TimeInForce string `json:"timeInForce"` //"GTC"           Optional, default is "GTC". Currently, we support "GTC" (good-till-canceled) and "IOC" (immediate-or-cancel).
}

func (bo *OkexOrder) setBy(order exc.ExchangeOrder) {
	bo.Coid = exc.Uuid(32)
	bo.Time = exc.GetTimeSpan()
	bo.Symbol = order.Market
	bo.OrderPrice = fmt.Sprintf("%g", order.Price)
	//bo.StopPrice = "0"
	bo.OrderQty = fmt.Sprintf("%g", order.Size)
	bo.OrderType = string(order.OrderType)
	bo.Side = order.Side
}

//下訂單
func (bm *Okex) PostOrder(order exc.ExchangeOrder) (string, error) {

	bo := OkexOrder{}
	bo.setBy(order)

	request, err := json.Marshal(bo)
	if err != nil {
		log.Fatal(err)
	}
	body := string(request)
	log.Println(fmt.Sprintf("body:%s", body))

	response := bm.doOrderRequest("order", body, bo.Time, bo.Coid)

	log.Println(fmt.Sprintf("%s", response))

	//{"code":6010,"message":"Not enough balance."}
	type OrderResponse struct {
		Code    float64 `json:"code"`
		Message string  `json:"message"`
	}
	orderResponse := OrderResponse{}

	if err := json.Unmarshal(response, &orderResponse); err != nil {
		panic(err)
	}

	var resErr error = nil
	if orderResponse.Code != 0 {
		resErr = errors.New(orderResponse.Message)
	}

	return string(response), resErr
}

func (bm *Okex) doNormalRequest(method, apiName, body string) []byte {
	ts := exc.GetTimeSpan()
	return bm.doRequest(method, apiName, body, false, ts, "")
}

func (bm *Okex) doOrderRequest(apiName, body string, ts int64, coid string) []byte {
	return bm.doRequest("POST", apiName, body, true, ts, coid)
}

func (bm *Okex) doRequest(method, apiName, body string, needAuth bool, ts int64, coid string) []byte {
	client := bm.client

	var res []byte

	fullUrl := fmt.Sprintf("%s%s%s", apiURL, apiPrefix, apiName)

	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Println(err)
		return res
	}

	req.Header.Set("Content-Type", "application/json")
	addHeader(&req.Header, method, apiName, ts, body)

	sendRes := exc.SendRequest(client, req)
	//{"code":6010,"message":"Not enough balance."}
	type OrderResponse struct {
		Code    float64 `json:"code"`
		Message string  `json:"message"`
	}
	orderResponse := OrderResponse{}

	json.Unmarshal(sendRes, &orderResponse)

	var resErr error = nil
	if orderResponse.Code != 0 {
		resErr = errors.New(orderResponse.Message)
		panic(resErr)
	}

	return exc.SendRequest(client, req)
}

func addHeader(header *http.Header, reqMethod, path string, ts int64, body string) {
	signPath := fmt.Sprintf("%v%v", apiPrefix, path)
	//exc.GetTimeSpan()

	tsStr := exc.GetUTC() // exc.GetTimeSpanStr(ts)
	total := tsStr + reqMethod + signPath + string(body)
	sign, _ := exc.GetParamHmacSHA256Base64Sign(setting.OKEX_SECRET_KEY, total)

	header.Add("OK-ACCESS-KEY", setting.OKEX_KEY)
	header.Add("OK-ACCESS-SIGN", sign)
	header.Add("OK-ACCESS-TIMESTAMP", tsStr)
	header.Add("OK-ACCESS-PASSPHRASE", setting.OKEX_PASSPHASE)

}
