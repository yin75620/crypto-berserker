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
		fmt.Println(array)
		sArray := array.([]interface{})
		for _, s := range sArray {
			res, _ := strconv.ParseFloat(s.(string), 64)
			askFloatArray = append(askFloatArray, res)
		}
		res = append(res, askFloatArray)
		fmt.Println(askFloatArray)
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
	fmt.Println(path)
	resByte := bm.doNormalRequest("GET", path, "")

	fmt.Println(string(resByte))

	var resJson map[string]interface{}

	err := json.Unmarshal(resByte, &resJson)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resJson)

	quoteResponse := QuoteResponse{}
	quoteResponse.setBy(resJson)

	askPair, _ := quoteResponse.GetPair(1, exc.Ask)
	bidPair, _ := quoteResponse.GetPair(1, exc.Bid)

	return askPair, bidPair
}

func (bm *Bitmax) GetAccountInfo() []byte {
	// 這交易所沒有使用者資料
	return []byte{}
}

type BitmaxOrder struct {
	Coid        string `json:"coid"`        //"xxx...xxx"     a unique identifier of length 32
	Time        int64  `json:"time"`        // 1528988100000   milliseconds since UNIX epoch in UTC
	Symbol      string `json:"symbol"`      //"ETH/BTC"
	OrderPrice  string `json:"orderPrice"`  //"13.5"          optional, limit price of the order. This field is required for limit orders and stop limit orders.
	StopPrice   string `json:"stopPrice"`   //"15.7"          optional, stop price of the order. This field is required for stop market orders and stop limit orders.
	OrderQty    string `json:"orderQty"`    //"3.5"
	OrderType   string `json:"orderType"`   //"limit"         order type, you shall specify one of the following: "limit", "market", "stop_market", "stop_limit".
	Side        string `json:"side"`        //"buy"           "buy" or "sell"
	PostOnly    bool   `json:"postOnly"`    //true            Optional, if true, the order will either be posted to the limit order book or be cancelled, i.e. the order cannot take liquidity; default value is false
	TimeInForce string `json:"timeInForce"` //"GTC"           Optional, default is "GTC". Currently, we support "GTC" (good-till-canceled) and "IOC" (immediate-or-cancel).
}

func (bo *BitmaxOrder) setBy(order exc.ExchangeOrder) {
	bo.Coid = ""
	bo.Time = 0
	bo.Symbol = order.Market
	bo.OrderPrice = order.Price
	//bo.StopPrice = ""
	bo.OrderQty = order.Size
	bo.OrderType = order.OrderType
	bo.Side = order.Side
}

//下訂單
func (bm *Bitmax) PostOrder(order exc.ExchangeOrder) string {

	fo := exc.ExchangeOrder{}
	fo.setBy(order)

	request, err := json.Marshal(fo)
	if err != nil {
		log.Fatal(err)
	}
	body := string(request)
	log.Println(fmt.Sprintf("body:%s", body))
	response := ftx.doPost("orders", body)
	log.Println(fmt.Sprintf("%s", response))
	return string(response)
}

func (bm *Bitmax) doNormalRequest(method, apiName, body string) []byte {
	return bm.doRequest(method, apiName, body, false)
}

func (bm *Bitmax) doAuthRequest(method, apiName, body string) []byte {
	return bm.doRequest(method, apiName, body, true)
}

func (bm *Bitmax) doRequest(method, apiName, body string, needAuth bool) []byte {
	client := bm.client

	var res []byte

	accountGroupStr := ""
	if needAuth {
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
		accountGroupStr = fmt.Sprintf("%d/", bm.accountGroup)
	}

	fullUrl := fmt.Sprintf("%s%s%s%s", apiURL, accountGroupStr, apiPrefix, apiName)
	fmt.Println(fullUrl)
	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Println(err)
		return res
	}

	req.Header.Set("Content-Type", "application/json")
	addHeader(&req.Header, method, apiName, body)

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

func addHeader(header *http.Header, reqMethod, path, body string) {
	ts := exc.GetTimeSpanStr(exc.GetTimeSpan())

	header.Add("x-auth-key", setting.BITMAX_KEY)
	header.Add("x-auth-timestamp", ts)

	if body != "" {
		body = fmt.Sprintf("+%s", body)
	}
	payload := fmt.Sprintf("%s+%s%s", ts, path, body)
	sign, _ := exc.GetParamHmacSHA256Base64Sign(setting.BITMAX_API_SECRET_KEY, payload)
	header.Add("x-auth-signature", sign)
	header.Add("x-auth-coid", setting.BITMAX_COID)
}
