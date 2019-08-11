package ftx

import (
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

type PricePair struct {
	price  float64
	volume float64
}

type PriceStatus struct {
	Asks [][]float64 `"json:asks"`
	Bids [][]float64 `"json:bids"`
}

func (ps *PriceStatus) getAskPricePair(depth int) (PricePair, error) {
	var res = PricePair{}
	size := len(ps.Asks)
	if depth > size {
		return res, errors.New("depth can't over size")
	}

	index := depth - 1
	res.price = ps.Asks[index][0] // first prize, second volume
	res.volume = ps.Asks[index][1]
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

// 看看 Api提供的交易對
func (ftx *Ftx) GetAsk(marketName string, depth int) float64 {
	response := ftx.GetOrderBook(marketName, depth)
	var bookResponse OrderBookResponse
	json.Unmarshal(response, &bookResponse)
	res, _ := bookResponse.Result.getAskPricePair(depth)
	return res.price
}

func (ftx *Ftx) GetBidPrice() {

}

func (ftx *Ftx) doGet(apiName, body string) []byte {
	return ftx.doRequest("GET", apiName, body)
}

func (ftx *Ftx) doRequest(method, apiName, body string) []byte {
	client := ftx.client

	var res []byte

	fullUrl := fmt.Sprintf("%s/%s", apiURL, apiName)
	req, err := http.NewRequest(method, fullUrl, nil)
	if err != nil {
		log.Println(err)
		return res
	}

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
	//boyd := "" //之後再用
	payload := fmt.Sprintf("%s%s%s%s", ts, reqMethod, apiPrefix+path, body)
	log.Println(payload)
	sign, _ := GetParamHmacSHA256HexSign(bsk.FTX_API_SECRET_KEY, payload)
	log.Println(sign)
	header.Add("FTX-SIGN", sign)
}
