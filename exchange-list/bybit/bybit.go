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
	"strings"

	exc "github.com/yin75620/crypto-berserker/exchange"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
	"github.com/yin75620/crypto-berserker/object_tool"
	"github.com/yin75620/crypto-berserker/setting"
)

var (
	apiURL = "https://api.bybit.com/"
)

type JArray exc.JArray

func NewBybit(c *http.Client) *Bybit {
	bb := Bybit{}
	bb.client = c
	bb.apiKey = setting.BYBIT_KEY
	bb.secretKey = setting.BYBIT_SECRET_KEY
	bb.orderBookCenter = ob.NewOrderBookCenter(NewSocket())
	bb.account.MakerFee = -0.00025
	bb.account.TakerFee = 0.00075
	return &bb
}

type Bybit struct {
	client          *http.Client
	apiKey          string
	secretKey       string
	orderBookCenter *ob.OrderBookCenter
	account         exc.Account
}

// implement exchange
func (bb *Bybit) GetWallet() exc.Wallet {
	var ret GetBalanceResult
	jarray := exc.JArray{}
	coinName := "BTC"
	jarray["coin"] = coinName

	response := bb.doRequest("GET", "v2/private/wallet/balance", jarray)
	fmt.Println(string(response))
	err := json.Unmarshal(response, &ret)

	w := exc.NewWallet()
	w.Balances = appendToBalance(w.Balances, ret.Result.BTC, coinName)

	if err != nil {
		fmt.Println(err)
		return *w
	}

	return *w
}

func appendToBalance(balances []exc.Balance, bybitBalance Balance, coinName string) []exc.Balance {
	const TempBTCToUSDValue = 9000.0

	bal := exc.Balance{
		Coin:     coinName,
		Free:     bybitBalance.AvailableBalance,
		Total:    bybitBalance.Equity,
		UsdValue: TempBTCToUSDValue * bybitBalance.Equity,
	}
	balances = append(balances, bal)
	return balances
}

func (bb *Bybit) GetAccountInfo() []byte {

	jarray := exc.JArray{}
	coinName := "BTC"
	jarray["coin"] = coinName

	response := bb.doRequest("GET", "v2/private/wallet/balance", jarray)
	fmt.Println(string(response))

	res := bb.SendGetLeverage()
	fmt.Println(res)
	bb.account.Leverage = float64(res.Result["BTCUSD"].Leverage)

	return response
}

func (bb *Bybit) PostOrder(order exc.ExchangeOrder) (string, error) {
	return "", nil // not implement
}

func (bb *Bybit) GetAccount() exc.Account {
	return bb.account
}

func (bb *Bybit) SendGetLeverage() GetLeverageResult {
	res := bb.doRequest("GET", "/user/leverage", exc.JArray{})
	leverageResult := GetLeverageResult{}
	err := json.Unmarshal(res, &leverageResult)
	if err != nil {
		fmt.Println("SendGetLeverage", err)
	}
	return leverageResult
}

func (bb *Bybit) PostFuturesOrder(order exc.FuturesOrder) (string, error) {
	bo := BybitOrder{}
	bo.Side = strings.Title(order.CommodityOrder.Side)
	bo.Symbol = strings.ToUpper(order.Futures.GetLinkMarketName())
	bo.OrderType = strings.Title(string(order.CommodityOrder.OrderType))
	bo.Quantity = int64(order.CommodityOrder.Size)
	bo.Price = order.CommodityOrder.Price
	bo.TimeInForce = "GoodTillCancel"

	return bb.doPostOrder(bo)
}

type BybitOrder struct {
	Side        string  `json:"side"`          //side	true	string	方向, 有效选项:Buy, Sell (Buy Sell )
	Symbol      string  `json:"symbol"`        //symbol	true	string	产品类型, 有效选项:BTCUSD, ETHUSD (BTCUSD ETHUSD )
	OrderType   string  `json:"order_type"`    //order_type	true	string	委托单价格类型, 有效选项:Limit, Market (Limit Market )
	Quantity    int64   `json:"qty"`           //qty	true	integer	委托数量, 单比最大1百万
	Price       float64 `json:"price"`         // price	false	number	委托价格, 在没有仓位时，做多的委托价格需高于市价的10%、低于1百万。如有仓位时则需优于强平价。单笔价格增减最小单位为0.5。如果下限价单，则price为必输字段
	TimeInForce string  `json:"time_in_force"` //time_in_force	true	string	执行策略, 有效选项:GoodTillCancel, ImmediateOrCancel, FillOrKill,PostOnly
	//TakeProfit  string  `json:"take_profit,omitempty"` // take_profit	false	number	止盈价格
	//stop_loss	false	number	止损价格
	//reduce_only	false	bool	只减仓
	//close_on_trigger	false	bool	触发后平仓
	//order_link_id	false	string	机构自定义订单ID, 最大长度36位，且同一机构下自定义ID不可重复
	//trailing_stop	false	number	追踪止损
}

func (bb *Bybit) doPostOrder(bo BybitOrder) (string, error) {
	request, err := json.Marshal(bo)
	if err != nil {
		log.Fatal(err)
	}
	body := string(request)
	log.Println(fmt.Sprintf("body:%s", body))

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(request, &jsonMap)
	if err != nil {
		panic(err)
	}

	response := bb.doRequest("POST", "v2/private/order/create", jsonMap)
	log.Println(fmt.Sprintf("%s", response))

	var resErr error

	return string(response), resErr
}

func (bb *Bybit) GetFee() exc.Fee {
	fee := exc.Fee{}
	fee.Taker = 0.00075
	fee.Maker = -0.00025
	return fee
}
func (bb *Bybit) GetName() string {
	return "Bybit"
}
func (bb *Bybit) GetMarketInfo(coinPair exc.CoinPair) exc.MarketInfo {
	return exc.MarketInfo{}
}

func (bb *Bybit) GetVolumeByTotal(total, price float64) float64 {
	return total
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

func (bb *Bybit) GetFuturesAskBidPair(futures exc.Futures) (exc.PricePair, exc.PricePair) {
	market := futures.GetLinkMarketName()
	return bb.getAskBidPairByMarket(market)
}

func (bb *Bybit) GetAskBidPair(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	market := coinPair.GetLinkMakertNameUpper()
	return bb.getAskBidPairByMarket(market)
}

func (bb *Bybit) getAskBidPairByMarket(market string) (exc.PricePair, exc.PricePair) {
	if !bb.orderBookCenter.IsExist(market) {
		channel, _ := bb.orderBookCenter.Register(market)
		<-channel
		go func() {
			for {
				<-channel
			}
		}()

		//return ftx.getOrderBookFromWeb(coinPair, depth)
	}

	booker := bb.orderBookCenter.GetBooker(market)
	return booker.GetFirstPricePair()
}

func (bb *Bybit) doRequest(method, apiName string, body exc.JArray) []byte {
	client := bb.client

	ts := exc.GetTimeSpan()
	objBody := exc.JArray{
		"api_key":   bb.apiKey,
		"timestamp": ts,
	}
	objBody.Add(body)

	sign := GetSignature(objBody, bb.secretKey)

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
