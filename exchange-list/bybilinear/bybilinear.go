package bybilinear

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	exc "github.com/yin75620/crypto-berserker/exchange"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
	"github.com/yin75620/crypto-berserker/jmath"
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
	bb.marketInfos = make(map[string]exc.MarketInfo)
	bb.LeverageInfos = map[string]exc.LeverageInfo{}
	bb.account.MakerFee = -0.00025
	bb.account.TakerFee = 0.00075
	bb.account.Leverage = 50 // there has no api, just manual set to 50.
	return &bb
}

type Bybilinear struct {
	client          *http.Client
	apiKey          string
	secretKey       string
	orderBookCenter *ob.OrderBookCenter
	marketInfos     map[string]exc.MarketInfo
	LeverageInfos   map[string]exc.LeverageInfo
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
	//fmt.Println(string(response))
	err = json.Unmarshal(response, &ret)
	if err != nil {
		fmt.Println(err)
		return *w
	}

	w.Balances = appendToBalance(w.Balances, ret.Result.USDT, coinName)

	bb.account.WalletInfo = *w
	bb.account.UnrealizedPnL = ret.Result.USDT.UnrealisedPnl
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

func (bb *Bybilinear) Prepare() []byte {
	response := bb.GetWallet()

	bb.prepareMarketInfo()
	bb.prepareLeverage()

	return []byte(fmt.Sprintf("%v", response))
}

// for implement
func (bb *Bybilinear) GetMaxOrderUSD(symbol string) float64 {
	if value, ok := bb.LeverageInfos[symbol]; ok {
		return bb.account.WalletInfo.GetAllBalanceFreeUSDValue() * float64(value.Leverage)
	} else {
		return bb.account.WalletInfo.GetAllBalanceFreeUSDValue() * bb.account.Leverage
	}
}

func (bb *Bybilinear) GetAccountInfo() []byte {
	return bb.Prepare()
}

func (bb *Bybilinear) PostOrder(order exc.ExchangeOrder) (string, error) {
	return "", nil // not implement
}

func (bb *Bybilinear) GetAccount() exc.Account {
	return bb.account
}

func (bb *Bybilinear) SendGetAccount() {

}

func (bb *Bybilinear) SendGetLeverage() GetLeverageResult {
	leverageResult := GetLeverageResult{}
	res, err := bb.doRequest("GET", "/user/leverage", exc.JArray{})
	if err != nil {
		fmt.Println(err)
		return leverageResult
	}
	//fmt.Println(string(res))
	err = json.Unmarshal(res, &leverageResult)

	if err != nil {
		fmt.Println("SendGetLeverage", err)
		return leverageResult
	}
	//leverageResult.Result["BTCUSDT"] = LeverageItem{Leverage: 100}
	return leverageResult
}

// Not Yet finish
func (bb *Bybilinear) PostCancelOrder(f exc.Futures) {
	boc := BybilinearCancelOrder{}
	boc.Symbol = f.GetLinkMarketName()
	bb.doPostCancelOrder(boc)
}

func (bb *Bybilinear) PostCancelAllOrder(f exc.Futures) {

	boc := BybilinearCancelOrder{}
	boc.Symbol = f.GetLinkMarketName()
	bb.doPostCancelAllOrder(boc)
}

func (bb *Bybilinear) PostFuturesOrder(order exc.FuturesOrder) (string, error) {
	symbol := order.Futures.GetLinkMarketName()
	var merketInfo = bb.doGetMarketInfo(symbol)

	bo := BybilinearOrder{}
	bo.Side = strings.Title(order.CommodityOrder.Side)
	bo.Symbol = strings.ToUpper(order.Futures.GetLinkMarketName())
	bo.OrderType = strings.Title(string(order.CommodityOrder.OrderType))
	bo.Quantity = jmath.FloatFloorByFloat(order.CommodityOrder.Size, merketInfo.VolumeIncrement)
	bo.Price = jmath.FloatFloorByFloat(order.CommodityOrder.Price, merketInfo.PriceIncrement)
	bo.TimeInForce = "GoodTillCancel"
	bo.CloseOnTrigger = order.IsClose

	res, err := bb.doPostOrder(bo)
	if err != nil {
		fmt.Println(err)
		return res, err
	}

	return res, err
}

func (bb *Bybilinear) getInstrumentInfo() InstrumentsInfo {

	body := exc.JArray{
		"category": "linear",
	}

	instrumentsInfo := InstrumentsInfo{}

	resByte, _ := bb.doRequest("GET", "derivatives/v3/public/instruments-info", body)
	//println(string(resByte))
	json.Unmarshal(resByte, &instrumentsInfo)

	return instrumentsInfo

}

func (bb *Bybilinear) prepareMarketInfo() {

	instrumentsInfo := bb.getInstrumentInfo()
	for _, value := range instrumentsInfo.Results.List {
		marketInfo := exc.MarketInfo{}
		marketInfo.Name = value.Symbol
		marketInfo.PriceIncrement = value.PriceFilter.TickSize
		marketInfo.VolumeIncrement = value.LotSizeFilter.QtyStep
		bb.marketInfos[marketInfo.Name] = marketInfo
	}

}

func (bb *Bybilinear) prepareLeverage() {

	resByte, _ := bb.doRequest("GET", "private/linear/position/list", exc.JArray{})

	pr := PositionListResponse{}
	json.Unmarshal(resByte, &pr)

	for _, value := range pr.Result {
		li := exc.LeverageInfo{}
		li.Name = value.Data.Symbol
		li.Leverage = value.Data.Leverage
		bb.LeverageInfos[li.Name] = li
	}
}

func (bb *Bybilinear) setAllLeverage() {

	resByte, _ := bb.doRequest("GET", "/public/linear/risk-limit", exc.JArray{})

	li := LeverageInfo{}
	json.Unmarshal(resByte, &li)

	fmt.Println(string(resByte))
	fmt.Println(li)

	currentSymbol := ""
	isSymbolSetFinish := false

	for _, v := range li.Result {

		if currentSymbol != v.Symbol {
			currentSymbol = v.Symbol
			isSymbolSetFinish = false
		}

		//一種symbol只要設定最大的那一次就好
		if isSymbolSetFinish {
			continue
		}

		total := bb.account.WalletInfo.GetAllBalanceUSDValue() * float64(v.MaxLeverage)
		if int(total) < v.Limit {
			bb.PostSetLeverage(v.Symbol, v.MaxLeverage)
			isSymbolSetFinish = true
		}
	}

}

func (bb *Bybilinear) PostSetLeverage(symbol string, leverage int) {
	req := exc.JArray{
		"symbol":        symbol,
		"buy_leverage":  leverage,
		"sell_leverage": leverage,
	}

	resByte, err := bb.doRequest("POST", "/private/linear/position/set-leverage", req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(resByte))

}

func (bb *Bybilinear) getUserTrades(symbol string, startTime time.Time, endTime time.Time) []UserTrade {

	jarray := exc.JArray{
		"symbol": symbol,
	}

	zero := time.Time{}
	if startTime != zero {
		jarray.Add(exc.JArray{"start_time": startTime.UnixMilli()})
	}
	if endTime != zero {
		jarray.Add(exc.JArray{"end_time": endTime.UnixMilli()})
	}

	res, err := bb.doRequest("GET", "/private/linear/trade/execution/history-list", jarray)

	if err != nil {
		log.Fatal(err)
	}

	thr := TradingHistoryResponse{}
	err = json.Unmarshal(res, &thr)
	if err != nil {
		fmt.Println(err)
	}

	return thr.Result.Data
}

func (bb *Bybilinear) GetTightUserTrades(symbol string) map[exc.UserTradeKey]exc.UserTrade {
	return bb.GetTightUserTradesWithTime(symbol, time.Now().Add(-24*time.Hour), time.Now())
}

func (bb *Bybilinear) GetTightUserTradesWithTime(symbol string, startTime time.Time, endTime time.Time) map[exc.UserTradeKey]exc.UserTrade {
	userTrades := bb.getUserTrades(symbol, startTime, endTime)

	euts := map[exc.UserTradeKey]exc.UserTrade{}
	for _, v := range userTrades {
		eut := v.ToExcUserTrade()
		eutKey := v.ToExcUserTradeKey()

		if value, ok := euts[eutKey]; ok {
			value.Combine(eut)
			euts[eutKey] = value
		} else {
			euts[eutKey] = eut
		}
	}
	return euts
}

type BybilinearOrder struct {
	Side        string  `json:"side"`          //side	true	string	方向, 有效选项:Buy, Sell (Buy Sell )
	Symbol      string  `json:"symbol"`        //symbol	true	string	产品类型, 有效选项:BTCUSD, ETHUSD (BTCUSD ETHUSD )
	OrderType   string  `json:"order_type"`    //order_type	true	string	委托单价格类型, 有效选项:Limit, Market (Limit Market )
	Quantity    float64 `json:"qty,string"`    //qty	true	integer	委托数量, 单比最大1百万
	Price       float64 `json:"price,string"`  // price	false	number	委托价格, 在没有仓位时，做多的委托价格需高于市价的10%、低于1百万。如有仓位时则需优于强平价。单笔价格增减最小单位为0.5。如果下限价单，则price为必输字段
	TimeInForce string  `json:"time_in_force"` //time_in_force	true	string	执行策略, 有效选项:GoodTillCancel, ImmediateOrCancel, FillOrKill,PostOnly
	//TakeProfit  string  `json:"take_profit,omitempty"` // take_profit	false	number	止盈价格
	//stop_loss	false	number	止损价格
	ReduceOnly     bool `json:"reduce_only"`      //只减仓
	CloseOnTrigger bool `json:"close_on_trigger"` //false	触发后平仓
	//order_link_id	false	string	机构自定义订单ID, 最大长度36位，且同一机构下自定义ID不可重复
	//trailing_stop	false	number	追踪止损
	PositionIdx int `json:"position_idx"`
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

	or := OrderResponse{}
	json.Unmarshal(response, &or)
	if err != nil {
		return string(response), err
	}
	if or.RetCode != 0 {
		err = errors.New(fmt.Sprintf("RectCode !=0, RectCode = %d", or.RetCode))
	}

	return string(response), err
}

func (bb *Bybilinear) doPostCancelOrder(cancelPos BybilinearCancelOrder) (string, error) {
	request, err := json.Marshal(cancelPos)
	if err != nil {
		log.Fatal(err)
	}
	body := string(request)
	//jtest
	log.Println(fmt.Sprintf("body:%s", body))

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(request, &jsonMap)
	if err != nil {
		panic(err)
	}

	response, err := bb.doRequest("POST", "private/linear/order/cancel", jsonMap)
	//jtest
	log.Println(fmt.Sprintf("%s", response))

	return string(response), err
}

func (bb *Bybilinear) doPostCancelAllOrder(cancelPos BybilinearCancelOrder) (string, error) {
	request, err := json.Marshal(cancelPos)
	if err != nil {
		log.Fatal(err)
	}
	body := string(request)
	//jtest
	log.Println(fmt.Sprintf("body:%s", body))

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(request, &jsonMap)
	if err != nil {
		panic(err)
	}

	response, err := bb.doRequest("POST", "private/linear/order/cancel-all", jsonMap)
	//jtest
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
	return bb.doGetMarketInfo(coinPair.GetLinkMakertNameUpper())
}

func (bb *Bybilinear) doGetMarketInfo(name string) exc.MarketInfo {

	return bb.marketInfos[name]
	/*
		switch name {
		case "BTCUSDT":
			return exc.MarketInfo{PriceIncrement: 0.5, VolumeIncrement: 0.001}
		case "ETHUSDT":
			return exc.MarketInfo{PriceIncrement: 0.01, VolumeIncrement: 0.01}
		case "SOLUSDT":
			return exc.MarketInfo{PriceIncrement: 0.005, VolumeIncrement: 0.1}

		default:
			return exc.MarketInfo{}
		}*/
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
	//fmt.Println(_val)
	h := hmac.New(sha256.New, []byte(key))
	io.WriteString(h, _val)
	return fmt.Sprintf("%x", h.Sum(nil))
}
