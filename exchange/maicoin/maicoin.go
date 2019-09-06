package maicoin

import (
	"bytes"
	"encoding/base64"
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

type JArray map[string]interface{}

func NewMaicoin(c *http.Client) *maicoin {
	maicoin := &maicoin{}
	maicoin.client = c
	return maicoin
}

type maicoin struct {
	client *http.Client
}

var (
	apiURL    = "https://max-api.maicoin.com"
	apiPrefix = "/api/v2/"
)

// implement exchange
func (bm *maicoin) GetFee() exc.Fee {
	fee := exc.Fee{}
	fee.Deposit = 0
	fee.WithDrawl = 0
	fee.Taker = 0.00075
	fee.Deposit = 0.00075
	return fee
}

type PriceType exc.PriceType

const (
	Ask PriceType = iota
	Bid
)

type PriceStatus struct {
	Asks [][]float64 `"json:asks"`
	Bids [][]float64 `"json:bids"`
}

func (ps *PriceStatus) GetPair(depth int, pType PriceType) (exc.PricePair, error) {
	switch pType {
	case Ask:
		return ps.getAskPricePair(depth)
	case Bid:
		return ps.getBidPricePair(depth)
	}
	return exc.PricePair{}, errors.New("has no match PriceType")
}

func (ps *PriceStatus) getAskPricePair(depth int) (exc.PricePair, error) {
	return exc.GetPricePair(depth, ps.Asks)
}
func (ps *PriceStatus) getBidPricePair(depth int) (exc.PricePair, error) {
	return exc.GetPricePair(depth, ps.Bids)
}

type QuoteResponse struct {
	Timestamp float64     `"json:timestamp"`
	Asks      [][]float64 `json:"asks,string"`
	Bids      [][]float64 `json:"asks,string"`
}

func (qr *QuoteResponse) setBy(json map[string]interface{}) {
	qr.Timestamp = json["timestamp"].(float64)

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

func (bm *maicoin) GetAskBidPair(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	market := coinPair.GetLinkMakertName()
	fmt.Println(market)
	resByte := bm.doBodyRequest("GET", "depth",
		JArray{
			"market": market,
			"limit":  depth,
		})

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

func (bm *maicoin) GetAccountInfo() []byte {
	// 這交易所沒有使用者資料
	res := bm.doNobodyRequest("GET", "members/profile")
	return res
}

func (bm *maicoin) GetCurrency() []byte {
	res := bm.doNobodyRequest("GET", "members/accounts")
	return res
}

func (bm *maicoin) GetFill(marketName string) []byte {
	res := bm.doBodyRequest("GET", "trades/my",
		JArray{
			"market": "usdttwd",
			"limit":  1,
		})
	return res
}

func (bm *maicoin) GetProducts() []byte {
	res := bm.doNobodyRequest("GET", "products")
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
func (bm *maicoin) PostOrder(order exc.ExchangeOrder) (string, error) {

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

	response := bm.doBodyRequest("POST", "orders", jarray)

	log.Println(fmt.Sprintf("%s", response))

	//{"error":{"code":1001,"message":"market does not have a valid value"}}
	/*type OrderResponse struct {
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
	}*/

	var resErr error = nil

	return string(response), resErr
}

func (bm *maicoin) doNobodyRequest(method, apiName string) []byte {
	return bm.doRequest(method, apiName, JArray{})
}

func (bm *maicoin) doBodyRequest(method, apiName string, body JArray) []byte {
	return bm.doRequest(method, apiName, body)
}

func (bm *maicoin) doRequest(method, apiName string, body JArray) []byte {
	client := bm.client

	ts := exc.GetTimeSpan()
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

	fmt.Printf("%s", sitemap)
	return sitemap
}

func addHeader(header *http.Header, reqMethod, path string, ts int64, sendBody string) {

	header.Add("X-MAX-ACCESSKEY", setting.MAICOIN_KEY)
	payload := base64.StdEncoding.EncodeToString([]byte(sendBody))
	fmt.Println(payload)
	sign, _ := exc.GetParamHmacSHA256HexSign(setting.MAICOIN_SECRET_KEY, payload)
	header.Add("X-MAX-PAYLOAD", payload)
	header.Add("X-MAX-SIGNATURE", sign)
}
