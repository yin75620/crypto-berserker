package coinex

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/setting"
)

var (
	apiURL    = "https://api.coinex.com"
	apiPrefix = "/v1/"
)

type JArray exc.JArray

func NewCoinEx(c *http.Client) *CoinEx {
	ce := CoinEx{}
	ce.client = c
	ce.accessID = setting.COINEX_ACCESS_ID
	return &ce
}

type CoinEx struct {
	client   *http.Client
	accessID string
}

// implement exchange
func (ce *CoinEx) GetWallet() exc.Wallet {
	w := exc.Wallet{}
	return w
}

func (ce *CoinEx) GetAccountInfo() []byte {

	//return ce.doRequest("GET", "balance/info", exc.JArray{})
	return []byte("CoinEx")
}

func (ce *CoinEx) PostOrder(order exc.ExchangeOrder) (string, error) {
	return "", nil
}

func (ce *CoinEx) GetFee() exc.Fee {
	fee := exc.Fee{}
	fee.Taker = 0.0005
	fee.Maker = 0.0005
	return fee
}
func (ce *CoinEx) GetName() string {
	return "CoinEx"
}
func (ce *CoinEx) GetMarketInfo(coinPair exc.CoinPair) exc.MarketInfo {
	return exc.MarketInfo{}
}

type QuoteResponse struct {
	//Timestamp float64         `"json:time"`
	TradeData exc.PriceStatus `"json:data"`
	Code      string          `"json:code"`
	//	Last
}

func (qr *QuoteResponse) setBy(json map[string]interface{}) {
	//qr.Timestamp = json["time"].(float64)

	qr.TradeData.SetByJArray(json["data"].(map[string]interface{}))
}

func (ce *CoinEx) GetAskBidPair(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	reqString := fmt.Sprintf("market/depth?market=%s&limit=%d&merge=0", coinPair.GetLinkMakertName(), depth)
	resByte := ce.doRequest("GET", reqString, exc.JArray{})
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

func (ce *CoinEx) doRequest(method, apiName string, body exc.JArray) []byte {
	client := ce.client

	ts := exc.GetTimeSpan()
	objBody := exc.JArray{
		"access_id": ce.accessID,
		"tonce":     ts,
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

	header.Add("X-MAX-ACCESSKEY", setting.MAICOIN_KEY)
	payload := base64.StdEncoding.EncodeToString([]byte(sendBody))
	sign, _ := exc.GetParamHmacSHA256HexSign(setting.MAICOIN_SECRET_KEY, payload)
	header.Add("X-MAX-PAYLOAD", payload)
	header.Add("X-MAX-SIGNATURE", sign)
}
