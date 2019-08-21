package ftx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	bsk "github.com/yin75620/crypto-berserker/setting"
)

type PriceType int

const (
	Ask PriceType = iota
	Bid
)

type PricePair struct {
	Price  float64
	Volume float64
}

type PriceStatus struct {
	Asks [][]float64 `"json:asks"`
	Bids [][]float64 `"json:bids"`
}

func (ps *PriceStatus) GetPair(depth int, pType PriceType) (PricePair, error) {
	switch pType {
	case Ask:
		return ps.getAskPricePair(depth)
	case Bid:
		return ps.getBidPricePair(depth)
	}
	return PricePair{}, errors.New("has no match PriceType")
}

func (ps *PriceStatus) getAskPricePair(depth int) (PricePair, error) {
	return getPricePair(depth, ps.Asks)
}
func (ps *PriceStatus) getBidPricePair(depth int) (PricePair, error) {
	return getPricePair(depth, ps.Bids)
}

func getPricePair(depth int, prices [][]float64) (PricePair, error) {
	var res = PricePair{}
	size := len(prices)
	if depth > size {
		return res, errors.New("depth can't over size")
	}

	index := depth - 1
	res.Price = prices[index][0] // first prize, second volume
	res.Volume = prices[index][1]
	return res, nil
}

type OrderBookResponse struct {
	Result  PriceStatus `json:"result"`
	Success bool        `json:"success"`
}

type Ftx struct {
	client *http.Client
}

var (
	apiURL    = "https://ftx.com/api"
	apiPrefix = "/api/"
)

func NewFtx(c *http.Client) *Ftx {
	ftx := &Ftx{}
	ftx.client = c
	return ftx
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

func (ftx *Ftx) GetTopOrderBook(marketName string) []byte {
	return ftx.GetOrderBook(marketName, 1)
}

func (ftx *Ftx) GetOrderBook(marketName string, depth int) []byte {
	path := fmt.Sprintf("/markets/%s/orderbook?depth=%d", marketName, depth)
	response := ftx.doGet(path, "")
	return response
}

func (ftx *Ftx) GetOrderBookResponse(marketName string, depth int) OrderBookResponse {
	response := ftx.GetOrderBook(marketName, depth)
	var bookResponse OrderBookResponse
	json.Unmarshal(response, &bookResponse)
	return bookResponse
}

// 看看 Api提供的交易對
func (ftx *Ftx) GetPair(marketName string, depth int, pType PriceType) PricePair {
	var bookResponse OrderBookResponse = ftx.GetOrderBookResponse(marketName, depth)
	res, _ := bookResponse.Result.GetPair(depth, pType)
	return res
}

func (ftx *Ftx) GetAskPair(marketName string, depth int) PricePair {
	var bookResponse OrderBookResponse = ftx.GetOrderBookResponse(marketName, depth)
	res, _ := bookResponse.Result.getAskPricePair(depth)
	return res
}

func (ftx *Ftx) GetAsk(marketName string, depth int) float64 {
	res := ftx.GetAskPair(marketName, depth)
	return res.Price
}

func (ftx *Ftx) GetBidPair(marketName string, depth int) PricePair {
	var bookResponse OrderBookResponse = ftx.GetOrderBookResponse(marketName, depth)
	res, _ := bookResponse.Result.getBidPricePair(depth)
	return res
}

func (ftx *Ftx) GetBid(marketName string, depth int) float64 {
	res := ftx.GetBidPair(marketName, depth)
	return res.Price
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

//下訂單
func (ftx *Ftx) PostOrder(order FtxOrder) string {
	request, err := json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}
	body := string(request)
	log.Println(fmt.Sprintf("body:%s", body))
	response := ftx.doPost("orders", body)
	log.Println(fmt.Sprintf("%s", response))
	return string(response)
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
	nanos := time.Now().UnixNano() / 1000000
	ts := strconv.FormatInt(nanos, 10)

	header.Add("FTX-KEY", bsk.FTX_KEY)
	header.Add("FTX-TS", ts)
	payload := fmt.Sprintf("%s%s%s%s", ts, reqMethod, apiPrefix+path, body)
	sign, _ := GetParamHmacSHA256HexSign(bsk.FTX_API_SECRET_KEY, payload)
	header.Add("FTX-SIGN", sign)
	header.Add("FTX-SUBACCOUNT", bsk.FTX_SUBACCOUNT)
}
