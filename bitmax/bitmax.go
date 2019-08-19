package bitmax

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

/*
func (bm *Bitmax) GetAskBidPair(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {

}*/

func (bm *Bitmax) GetAccountInfo() []byte {
	return bm.doAuthRequest("GET", "balance", "")
}

/*
//下訂單
func (bm *Bitmax) PostOrder(order exc.ExchangeOrder) string {

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
	return string(response)
}

func (bm *Bitmax) GetDepth(marketName string, depth int) OrderBookResponse {
	response := ftx.GetOrderBook(marketName, depth)
	var bookResponse OrderBookResponse
	json.Unmarshal(response, &bookResponse)
	return bookResponse
}*/

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
		if bm.accountGroup != 0 {
			accountGroupStr = fmt.Sprintf("%d/", bm.accountGroup)
		} else {
			bm.doNormalRequest("GET", "user/info", "")
		}
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
