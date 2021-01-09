package maicoin

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/setting"
)

type JArray exc.JArray

func NewMaicoin(c *http.Client) *maicoin {
	maicoin := &maicoin{}
	maicoin.client = c
	socket := NewSocket()
	//socket.Strat()
	maicoin.websocket = socket
	maicoin.orderBookers = map[exc.CoinPair]*OrderBooker{}
	//maicoin.orderBooker.Start()
	return maicoin
}

type maicoin struct {
	client       *http.Client
	websocket    *MaincoinWebSocket
	orderBookers map[exc.CoinPair]*OrderBooker
	mutex        sync.Mutex
}

var (
	apiURL    = "https://max-api.maicoin.com"
	apiPrefix = "/api/v2/"
)

// implement exchange
func (ma *maicoin) GetWallet() exc.Wallet {
	w := exc.Wallet{}
	return w
}

func (ma *maicoin) GetFee() exc.Fee {
	fee := exc.Fee{}
	fee.Deposit = 0
	fee.WithDrawl = 0
	fee.Taker = 0.00075
	fee.Maker = 0.00075
	return fee
}
func (ma *maicoin) GetName() string {
	return "MAX"
}

func (ma *maicoin) GetMarketInfo(coinPair exc.CoinPair) exc.MarketInfo {
	return exc.MarketInfo{} //not yet implement
}

type QuoteResponse struct {
	Timestamp float64 `"json:timestamp"`
	exc.PriceStatus
	ErrorInfo exc.ErrorMessage `"json:error"`
}

func (qr *QuoteResponse) setBy(json map[string]interface{}) {
	if json["timestamp"] != nil {
		qr.Timestamp = json["timestamp"].(float64)
	}

	qr.PriceStatus.SetByJArray(json)
}

// GetAskBidPair Ask Bid
func (ma *maicoin) GetAskBidPair(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	if _, ok := ma.orderBookers[coinPair]; !ok {
		ma.mutex.Lock()
		resChannel := ma.websocket.SubScribeOrderBook(coinPair)
		ma.orderBookers[coinPair] = NewOrderBooker(ma, coinPair, resChannel)
		ma.orderBookers[coinPair].Start()
		ma.mutex.Unlock()
	}
	orderBooker := ma.orderBookers[coinPair]
	return orderBooker.BottomAskPair, orderBooker.TopBidPair

}

func (ma *maicoin) GetNewsetAskBidPair(coinPair exc.CoinPair) (exc.PricePair, exc.PricePair) {
	return ma.GetAskBidPairFromWeb(coinPair, 1)
}

// GetAskBidPairFromWeb Ask Bid
func (ma *maicoin) GetAskBidPairFromWeb(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	// get price from cache
	// get price from web

	market := coinPair.GetLinkMakertName()
	fmt.Println(market)
	resByte := ma.doBodyRequest("GET", "depth",
		JArray{
			"market": market,
			"limit":  depth,
		})
	fmt.Println("resByte", string(resByte))
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

func (ma *maicoin) GetAccountInfo() []byte {
	// 這交易所沒有使用者資料
	res := ma.doNobodyRequest("GET", "members/profile")
	return res
}

func (ma *maicoin) GetCurrency() []byte {
	res := ma.doNobodyRequest("GET", "members/accounts")
	return res
}

func (ma *maicoin) GetFill(marketName string) []byte {
	res := ma.doBodyRequest("GET", "trades/my",
		JArray{
			"market": "usdttwd",
			"limit":  1,
		})
	return res
}

func (ma *maicoin) GetProducts() []byte {
	res := ma.doNobodyRequest("GET", "products")
	return res
}

type maicoinOrder struct {
	Market    string  `json:"market"` //"maxtwd"
	Side      string  `json:"side"`   //"buy"           "buy" or "sell"
	Volume    float64 `json:"volume"` //"3.5"
	Price     float64 `json:"price"`  //"13.5"
	StopPrice string  `"json":"stop_price"`
	OrderType string  `json:"ord_type"` //"limit"         order type, you shall specify one of the following: "limit", "market", "stop_market", "stop_limit".
}

func (mai *maicoinOrder) setBy(order exc.ExchangeOrder) {
	mai.Market = order.CoinPair.GetLinkMakertName()
	mai.Side = order.Side
	mai.Volume = order.Size
	mai.Price = order.Price
	//mai.StopPrice = 0
	mai.OrderType = string(order.OrderType)
}

//下訂單
func (ma *maicoin) PostOrder(order exc.ExchangeOrder) (string, error) {

	mai := maicoinOrder{}
	mai.setBy(order)

	request, err := json.Marshal(mai)
	if err != nil {
		log.Fatal(err)
	}
	body := string(request)
	log.Println(fmt.Sprintf("body:%s", body))

	jarray := JArray{}
	json.Unmarshal([]byte(request), &jarray)

	response := ma.doBodyRequest("POST", "orders", jarray)

	log.Println(fmt.Sprintf("%s", response))

	//{"error":{"code":2016,"message":"Failed to create order. Reason: btc_amount_too_small"}}
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

func (ma *maicoin) doNobodyRequest(method, apiName string) []byte {
	return ma.doRequest(method, apiName, JArray{})
}

func (ma *maicoin) doBodyRequest(method, apiName string, body JArray) []byte {
	return ma.doRequest(method, apiName, body)
}

var mlastTs int64 = 0

func (ma *maicoin) doRequest(method, apiName string, body JArray) []byte {
	client := ma.client

	// 防止相同 ts 就失敗的問題。
	ts := exc.GetTimeSpan()
	if ts == mlastTs {
		ts = ts + 1
	}
	mlastTs = ts

	objBody := JArray{
		"path":  apiPrefix + apiName,
		"nonce": ts,
	}

	// 塞入 body
	for k, v := range body {
		objBody[k] = v
	}

	jsonBody, err := json.Marshal(objBody)
	if err != nil {
		log.Fatal(err)
	}
	sendBody := string(jsonBody)

	log.Println(sendBody)

	var res []byte
	fullUrl := fmt.Sprintf("%s%s%s", apiURL, apiPrefix, apiName)

	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer([]byte(sendBody)))
	if err != nil {
		log.Println(err)
		return res
	}

	req.Header.Set("Content-Type", "application/json")
	addHeader(&req.Header, method, apiName, ts, sendBody)

	res, err = exc.SendRequest(client, req)
	if err != nil {
		log.Println(err)
		return res
	}

	return res
}

func addHeader(header *http.Header, reqMethod, path string, ts int64, sendBody string) {

	header.Add("X-MAX-ACCESSKEY", setting.MAICOIN_KEY)
	payload := base64.StdEncoding.EncodeToString([]byte(sendBody))
	sign, _ := exc.GetParamHmacSHA256HexSign(setting.MAICOIN_SECRET_KEY, payload)
	header.Add("X-MAX-PAYLOAD", payload)
	header.Add("X-MAX-SIGNATURE", sign)
}
