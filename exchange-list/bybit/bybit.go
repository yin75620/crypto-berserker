package bybit

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/object_tool"
	"github.com/yin75620/crypto-berserker/setting"
)

var (
	apiURL = "https://api.bybit.com/"
)

type JArray exc.JArray

func NewBybit(c *http.Client) *Bybit {
	ce := Bybit{}
	ce.client = c
	ce.apiKey = setting.BYBIT_KEY
	ce.secretKey = setting.BYBIT_SECRET_KEY
	return &ce
}

type Bybit struct {
	client          *http.Client
	apiKey          string
	secretKey       string
	orderBookCenter *OrderBookCenter
}

// implement exchange
func (ce *Bybit) GetWallet() exc.Wallet {
	w := exc.Wallet{}
	return w
}

func (ce *Bybit) GetAccountInfo() []byte {

	return ce.doRequest("GET", "open-api/wallet/fund/records", exc.JArray{})
}

func (ce *Bybit) PostOrder(order exc.ExchangeOrder) (string, error) {
	return "", nil
}

func (ce *Bybit) GetFee() exc.Fee {
	fee := exc.Fee{}
	fee.Taker = 0.0007
	fee.Maker = 0.0007
	return fee
}
func (ce *Bybit) GetName() string {
	return "Bybit"
}
func (ce *Bybit) GetMarketInfo(coinPair exc.CoinPair) exc.MarketInfo {
	return exc.MarketInfo{}
}

type QuoteResponse struct {
	//Timestamp float64         `"json:time"`
	TradeData exc.PriceStatus `"json:data"`
	Code      string          `"json:code"`
	//	Last
}

func (qr *QuoteResponse) setBy(json map[string]interface{}) {
	qr.TradeData.SetByJArray(json)
}

func (ce *Bybit) GetAskBidPair(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
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

func (ce *Bybit) doRequest(method, apiName string, body exc.JArray) []byte {
	client := ce.client

	ts := exc.GetTimeSpan()
	objBody := exc.JArray{
		"api_key":   ce.apiKey,
		"timestamp": ts,
	}
	objBody.Add(body)

	sign := GetSignature(objBody, ce.secretKey)

	objBody.Add(exc.JArray{
		"sign": sign,
	})

	jsonBody, err := json.Marshal(objBody)
	if err != nil {
		log.Fatal(err)
	}

	var res []byte
	fullURL := fmt.Sprintf("%s%s", apiURL, apiName)

	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println(err)
		return res
	}

	if method == "GET" {
		q := req.URL.Query()
		for key, obj := range objBody {
			q.Add(key, object_tool.ToString(obj))
		}
		req.URL.RawQuery = q.Encode()
	}

	req.Header.Add("Content-Type", "application/json")

	return exc.SendRequest(client, req)
}

func GetSignature(params map[string]interface{}, key string) string {
	keys := make([]string, len(params))
	i := 0
	_val := ""
	for k, _ := range params {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		value := object_tool.ToString(params[k])
		_val += k + "=" + value + "&"
	}
	_val = _val[0 : len(_val)-1]
	fmt.Println(_val)
	h := hmac.New(sha256.New, []byte(key))
	io.WriteString(h, _val)
	return fmt.Sprintf("%x", h.Sum(nil))
}
