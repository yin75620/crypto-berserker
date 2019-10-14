package binance

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/setting"
)

var (
	apiURL    = "https://api.binance.com"
	apiPrefix = "/api/v1/"
)

type JArray exc.JArray

func NewBinance(c *http.Client) *binance {
	bn := binance{}
	bn.client = c

	return &bn
}

type binance struct {
	client *http.Client
}

func (bn *binance) GetAccountInfo() []byte {

	return []byte("Binance")
}

func (bn *binance) PostOrder(order exc.ExchangeOrder) (string, error) {
	return "", nil
}

func (bn *binance) GetFee() exc.Fee {
	fee := exc.Fee{}
	fee.Taker = 0.00075
	fee.Maker = 0.00075
	return fee
}
func (bn *binance) GetName() string {
	return "Binance"
}
func (bn *binance) GetMarketInfo(coinPair exc.CoinPair) exc.MarketInfo {
	return exc.MarketInfo{}
}

type QuoteResponse struct {
	//Timestamp float64         `"json:time"`
	LastUpdateId float64         `"json:lastUpdateId"`
	TradeData    exc.PriceStatus `"json:data"`
}

func (qr *QuoteResponse) setBy(json map[string]interface{}) {
	qr.LastUpdateId = json["lastUpdateId"].(float64)
	qr.TradeData.SetByJArray(json)
}

func (bn *binance) GetAskBidPair(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	const minLimit = 5
	reqString := fmt.Sprintf("depth?symbol=%s&limit=%d",
		strings.ToUpper(coinPair.GetLinkMakertName()),
		minLimit)
	resByte := bn.doRequest("GET", reqString, exc.JArray{})
	//fmt.Println(string(resByte))

	var resJson map[string]interface{}
	err := json.Unmarshal(resByte, &resJson)
	if err != nil {
		log.Fatal(err)
	}

	quoteResponse := QuoteResponse{}
	quoteResponse.setBy(resJson)

	askPair, _ := quoteResponse.TradeData.GetPair(1, exc.Ask)
	bidPair, _ := quoteResponse.TradeData.GetPair(1, exc.Bid)

	return askPair, bidPair
}

func (bn *binance) doRequest(method, apiName string, body exc.JArray) []byte {
	client := bn.client

	ts := exc.GetTimeSpan()
	objBody := exc.JArray{
		"recvWindow": strconv.FormatInt(common.RecvWindow(5*time.Second), 10),
		"timestamp":  ts,
	}

	objBody.Add(body)

	jsonBody, err := json.Marshal(objBody)
	if err != nil {
		log.Fatal(err)
	}
	sendBody := string(jsonBody)

	log.Println(sendBody)

	var res []byte
	fullUrl := fmt.Sprintf("%s%s%s", apiURL, apiPrefix, apiName)
	fmt.Println(fullUrl)

	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer([]byte("")))
	if err != nil {
		log.Println(err)
		return res
	}

	req.Header.Set("Content-Type", "application/json")
	addHeader(&req.Header, method, apiName, ts, sendBody)

	return exc.SendRequest(client, req)
}

func addHeader(header *http.Header, reqMethod, path string, ts int64, sendBody string) {

	header.Add("X-MBX-APIKEY", setting.BINANCE_KEY)
	payload := base64.StdEncoding.EncodeToString([]byte(sendBody))
	sign, _ := exc.GetParamHmacSHA256HexSign(setting.MAICOIN_SECRET_KEY, payload)
	header.Add("X-MAX-PAYLOAD", payload)
	header.Add("X-MAX-SIGNATURE", sign)
}
