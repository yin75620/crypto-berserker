package bitmax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	exc "github.com/yin75620/crypto-berserker/exchange"
)

type Bitmax struct {
	client *http.Client
}

var (
	apiURL     = "https://bitmax.io/api/v1/products"
	apiVersion = "v1"
	apiPrefix  = "/api/"
)

func (bm *Bitmax) GetAskBidPair(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	coinPair.get
	resb := ftx.GetOrderBookResponse(marketName, depth)
	askPair, _ := resb.Result.GetPair(1, Ask)
	bidPair, _ := resb.Result.GetPair(1, Bid)
	return askPair, bidPair
}

func (bm *Bitmax) GetAccountInfo() []byte {
	return ftx.doGet("account", "")
}

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
}

func (bm *Bitmax) doRequest(method, apiName, body string) []byte {
	client := bm.client

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
