package bitmax

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/setting"
)

func NewBitmax(c *http.Client) *Bitmax {
	bitmax := &Bitmax{}
	bitmax.client = c
	return bitmax
}

type Bitmax struct {
	client       *http.Client
	accountGroup int
}

var (
	apiURL    = "https://bitmax.io/"
	apiPrefix = "api/v1/"
)

// implement exchange
func (bm *Bitmax) GetFee() exc.Fee {
	fee := exc.Fee{}
	fee.Deposit = 0
	fee.WithDrawl = 0
	fee.Taker = 0.0004
	fee.Deposit = 0.0004
	return fee
}

type QuoteResponse struct {
	MarketName string      `json:"s"`
	Asks       [][]float64 `json:"asks,string"`
	Bids       [][]float64 `json:"asks,string"`
}

func (qr *QuoteResponse) setBy(json map[string]interface{}) {
	qr.MarketName = json["s"].(string)

	askStrArrays := json["asks"].([]interface{})
	qr.Asks = transToFloatTwoArray(askStrArrays)
	qr.Bids = transToFloatTwoArray(json["bids"].([]interface{}))
}

func transToFloatTwoArray(askStrArrays []interface{}) [][]float64 {
	res := [][]float64{}
	for _, array := range askStrArrays {
		askFloatArray := []float64{}
		sArray := array.([]interface{})
		for _, s := range sArray {
			res, _ := strconv.ParseFloat(s.(string), 64)
			askFloatArray = append(askFloatArray, res)
		}
		res = append(res, askFloatArray)
	}
	return res
}

func (qr *QuoteResponse) GetPair(depth int, pType exc.PriceType) (exc.PricePair, error) {
	switch pType {
	case exc.Ask:
		return qr.getAskPricePair(depth)
	case exc.Bid:
		return qr.getBidPricePair(depth)
	}
	return exc.PricePair{}, errors.New("has no match PriceType")
}

func (qr *QuoteResponse) getAskPricePair(depth int) (exc.PricePair, error) {
	return exc.GetPricePair(depth, qr.Asks)
}
func (qr *QuoteResponse) getBidPricePair(depth int) (exc.PricePair, error) {
	return exc.GetPricePair(depth, qr.Bids)
}

func (bm *Bitmax) GetAskBidPair(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	path := fmt.Sprintf("depth?symbol=%s&n=%d", coinPair.GetSymbal(), depth)
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

func (bm *Bitmax) GetAccountInfo() []byte {
	// 這交易所沒有使用者資料
	// 用 balance 代替
	//res := bm.doAuthRequest("GET", "balance", "")
	res := []byte("Bitmax")
	return res
}

func (bm *Bitmax) GetProducts() []byte {
	res := bm.doNormalRequest("GET", "products", "")
	return res
}

type BitmaxOrder struct {
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

func (bo *BitmaxOrder) setBy(order exc.ExchangeOrder) {
	bo.Coid = exc.Uuid(32)
	bo.Time = 0
	bo.Symbol = order.Market
	bo.OrderPrice = fmt.Sprintf("%g", order.Price)
	//bo.StopPrice = "0"
	bo.OrderQty = fmt.Sprintf("%g", order.Size)
	bo.OrderType = string(order.OrderType)
	bo.Side = order.Side
}

//下訂單
func (bm *Bitmax) PostOrder(order exc.ExchangeOrder) (string, error) {

	ts := exc.GetTimeSpan()
	coid := exc.Uuid(32)
	bo := BitmaxOrder{}
	bo.setBy(order)
	bo.Time = ts
	bo.Coid = coid

	request, err := json.Marshal(bo)
	if err != nil {
		log.Fatal(err)
	}
	body := string(request)
	log.Println(fmt.Sprintf("body:%s", body))

	response := bm.doOrderRequest("order", body, ts, coid)

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

func (bm *Bitmax) doNormalRequest(method, apiName, body string) []byte {
	ts := exc.GetTimeSpan()
	return bm.doRequest(method, apiName, body, false, ts, "")
}

func (bm *Bitmax) doAuthRequest(method, apiName, body string) []byte {
	ts := exc.GetTimeSpan()
	return bm.doRequest(method, apiName, body, true, ts, "")
}

func (bm *Bitmax) doOrderRequest(apiName, body string, ts int64, coid string) []byte {
	return bm.doRequest("POST", apiName, body, true, ts, coid)
}

func (bm *Bitmax) doRequest(method, apiName, body string, needAuth bool, ts int64, coid string) []byte {
	client := bm.client

	var res []byte

	accountGroupStr := ""
	if needAuth {
		accountGroupStr = bm.auth()
	}

	fullUrl := fmt.Sprintf("%s%s%s%s", apiURL, accountGroupStr, apiPrefix, apiName)

	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Println(err)
		return res
	}

	req.Header.Set("Content-Type", "application/json")
	addHeader(&req.Header, method, apiName, ts, coid)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return res
	}

	defer resp.Body.Close()
	sitemap, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("%s", err)
		return res
	}

	//fmt.Printf("%s", sitemap)
	return sitemap
}

func addHeader(header *http.Header, reqMethod, path string, ts int64, sCoid string) {
	strTs := exc.GetTimeSpanStr(ts)

	header.Add("x-auth-key", setting.BITMAX_KEY)
	header.Add("x-auth-timestamp", strTs)

	coid := ""
	if sCoid != "" {
		coid = fmt.Sprintf("+%s", sCoid)
	}
	payload := fmt.Sprintf("%s+%s%s", strTs, path, coid)
	//fmt.Println(payload)
	sign, _ := exc.GetParamHmacSHA256Base64Sign(setting.BITMAX_API_SECRET_KEY, payload)
	header.Add("x-auth-signature", sign)
	header.Add("x-auth-coid", sCoid)
}

func (bm *Bitmax) auth() string {
	if bm.accountGroup == 0 {
		type TempUserInfo struct {
			AccountGroup int `json:"accountGroup"`
		}

		byt := bm.doNormalRequest("GET", "user/info", "")
		userInfo := TempUserInfo{}
		if err := json.Unmarshal(byt, &userInfo); err != nil {
			panic(err)
		}

		bm.accountGroup = userInfo.AccountGroup
	}
	accountGroupStr := fmt.Sprintf("%d/", bm.accountGroup)
	return accountGroupStr
}
