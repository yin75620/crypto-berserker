package bybilinear

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
	apiURL = "https://api.Bybit.com/"
)

type JArray exc.JArray

func NewBybilinear(c *http.Client) *Bybilinear {
	bb := Bybilinear{}
	bb.client = c
	bb.apiKey = setting.BYBIT_KEY
	bb.secretKey = setting.BYBIT_SECRET_KEY
	bb.orderBookCenter = ob.NewOrderBookCenter(NewSocket())
	bb.account.MakerFee = -0.00025
	bb.account.TakerFee = 0.00075
	return &bb
}

type Bybilinear struct {
	client          *http.Client
	apiKey          string
	secretKey       string
	orderBookCenter *ob.OrderBookCenter
	account         exc.Account
}

// implement exchange
func (bb *Bybilinear) GetWallet() exc.Wallet {
	var ret GetBalanceResult
	jarray := exc.JArray{}
	coinName := "USDT"
	jarray["coin"] = coinName
	w := exc.NewWallet()

	response, err := bb.doRequest("GET", "v2/private/wallet/balance", jarray)
	if err != nil {
		fmt.Println(err)
		return *w
	}
	fmt.Println(string(response))
	err = json.Unmarshal(response, &ret)
	if err != nil {
		fmt.Println(err)
		return *w
	}

	w.Balances = appendToBalance(w.Balances, ret.Result.USDT, coinName)

	bb.account.WalletInfo = *w

	return *w
}

func appendToBalance(balances []exc.Balance, BybilinearBalance Balance, coinName string) []exc.Balance {
	const TempUSDTToUSDValue = 1

	bal := exc.Balance{
		Coin:         coinName,
		Free:         BybilinearBalance.AvailableBalance,
		FreeUsdValue: BybilinearBalance.AvailableBalance * TempUSDTToUSDValue,
		Total:        BybilinearBalance.Equity,
		UsdValue:     TempUSDTToUSDValue * BybilinearBalance.Equity,
	}
	balances = append(balances, bal)
	return balances
}

func (bb *Bybilinear) GetAccountInfo() []byte {

	jarray := exc.JArray{}
	coinName := "USDT"
	jarray["coin"] = coinName

	response := bb.GetWallet()
	fmt.Println(response)

	//res := bb.SendGetLeverage()
	//fmt.Println(res)
	// float64(res.Result["BTCUSDT"].Leverage)
	bb.account.Leverage = 100 // there has no api, just manual set to 100.

	return []byte(fmt.Sprintf("%v", response))
}

func (bb *Bybilinear) PostOrder(order exc.ExchangeOrder) (string, error) {
	return "", nil // not implement
}

func (bb *Bybilinear) GetAccount() exc.Account {
	return bb.account
}

func (bb *Bybilinear) SendGetLeverage() GetLeverageResult {
	leverageResult := GetLeverageResult{}
	res, err := bb.doRequest("GET", "/user/leverage", exc.JArray{})
	if err != nil {
		fmt.Println(err)
		return leverageResult
	}
	fmt.Println(string(res))
	err = json.Unmarshal(res, &leverageResult)

	if err != nil {
		fmt.Println("SendGetLeverage", err)
		return leverageResult
	}
	//leverageResult.Result["BTCUSDT"] = LeverageItem{Leverage: 100}
	return leverageResult
}

func (bb *Bybilinear) PostFuturesOrder(order exc.FuturesOrder) (string, error) {
	bo := BybilinearOrder{}
	bo.Side = strings.Title(order.CommodityOrder.Side)
	bo.Symbol = strings.ToUpper(order.Futures.GetLinkMarketName())
	bo.OrderType = strings.Title(string(order.CommodityOrder.OrderType))
	bo.Quantity = order.CommodityOrder.Size
	bo.Price = order.CommodityOrder.Price
	bo.TimeInForce = "GoodTillCancel"

	return bb.doPostOrder(bo)
}

type BybilinearOrder struct {
	Side        string  `json:"side"`          //side	true	string	方向, 有效选项:Buy, Sell (Buy Sell )
	Symbol      string  `json:"symbol"`        //symbol	true	string	产品类型, 有效选项:BTCUSD, ETHUSD (BTCUSD ETHUSD )
	OrderType   string  `json:"order_type"`    //order_type	true	string	委托单价格类型, 有效选项:Limit, Market (Limit Market )
	Quantity    float64 `json:"qty"`           //qty	true	integer	委托数量, 单比最大1百万
	Price       float64 `json:"price"`         // price	false	number	委托价格, 在没有仓位时，做多的委托价格需高于市价的10%、低于1百万。如有仓位时则需优于强平价。单笔价格增减最小单位为0.5。如果下限价单，则price为必输字段
	TimeInForce string  `json:"time_in_force"` //time_in_force	true	string	执行策略, 有效选项:GoodTillCancel, ImmediateOrCancel, FillOrKill,PostOnly
	//TakeProfit  string  `json:"take_profit,omitempty"` // take_profit	false	number	止盈价格
	//stop_loss	false	number	止损价格
	ReduceOnly     bool `json:"reduce_only"`      //只减仓
	CloseOnTrigger bool `json:"close_on_trigger"` //false	触发后平仓
	//order_link_id	false	string	机构自定义订单ID, 最大长度36位，且同一机构下自定义ID不可重复
	//trailing_stop	false	number	追踪止损
}

func (bb *Bybilinear) doPostOrder(bo BybilinearOrder) (string, error) {
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

	response, err := bb.doRequest("POST", "private/linear/order/create", jsonMap)
	log.Println(fmt.Sprintf("%s", response))

	return string(response), err
}

func (bb *Bybilinear) GetFee() exc.Fee {
	fee := exc.Fee{}
	fee.Taker = 0.00075
	fee.Maker = -0.00025
	return fee
}
func (bb *Bybilinear) GetName() string {
	return "Bybilinear"
}
func (bb *Bybilinear) GetMarketInfo(coinPair exc.CoinPair) exc.MarketInfo {
	return exc.MarketInfo{}
}

func (bb *Bybilinear) GetVolumeByTotal(total, price float64) float64 {
	return total / price
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

func (bb *Bybilinear) GetFuturesAskBidPair(futures exc.Futures) (exc.PricePair, exc.PricePair) {
	market := futures.GetLinkMarketName()
	return bb.getAskBidPairByMarket(market)
}

func (bb *Bybilinear) GetAskBidPair(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	market := coinPair.GetLinkMakertNameUpper()
	return bb.getAskBidPairByMarket(market)
}

func (bb *Bybilinear) getAskBidPairByMarket(market string) (exc.PricePair, exc.PricePair) {
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

func (bb *Bybilinear) doRequest(method, apiName string, body exc.JArray) ([]byte, error) {
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
		return res, err
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
