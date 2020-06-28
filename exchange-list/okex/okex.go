package okex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	exc "github.com/yin75620/crypto-berserker/exchange"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
	"github.com/yin75620/crypto-berserker/setting"
)

func NewOkex(c *http.Client) *Okex {
	Okex := &Okex{}
	Okex.client = c
	Okex.orderBookCenter = ob.NewOrderBookCenter(NewSocket())
	return Okex
}

type Okex struct {
	client          *http.Client
	account         exc.Account
	orderBookCenter *ob.OrderBookCenter
}

var (
	apiURL = "https://www.okex.com"
	//apiPrefix = "/api/spot/v3/" //幣幣交易
	//apiPrefix = "/api/futures/v3/" //期貨交易
	apiPrefix = "/api/swap/v3/" //永續合約
)

// implement exchange
func (bm *Okex) GetWallet() exc.Wallet {
	w := exc.Wallet{}
	type WalletResponse struct {
		Info SwapAccountInfo `json:"info"`
	}
	walletRes := WalletResponse{}

	response, err := bm.doRequest("GET", "BTC-USDT-SWAP/accounts", "")
	if err != nil {
		log.Println("GetWallet err:", err)
	}

	json.Unmarshal(response, &walletRes)
	if err != nil {
		log.Fatal("Okex GetWallet", err)
	}
	fmt.Println(string(response))

	fmt.Println(walletRes.Info)

	wInfo := walletRes.Info

	balance := exc.Balance{
		Coin:         wInfo.Currency,
		Free:         wInfo.MaxWithdraw,
		FreeUsdValue: wInfo.MaxWithdraw,
		Total:        wInfo.Equity,
		UsdValue:     wInfo.Equity,
	}
	w.Balances = append(w.Balances, balance)

	bm.account.UnrealizedPnL = wInfo.UnrealizedPnl
	bm.account.WalletInfo = w

	return w
}

func (bm *Okex) GetLeverage() float64 {
	res, err := bm.doRequest("GET", "accounts/BTC-USD-SWAP/settings", "")
	if err != nil {
		log.Println("GetLeverage err:", err)
	}
	setting := SwapAccountsSetting{}
	json.Unmarshal(res, &setting)
	log.Println("leverage", setting.LongLeverage)

	bm.account.Leverage = setting.LongLeverage
	return setting.LongLeverage
}
func (bm *Okex) GetFee() exc.Fee {
	fee := exc.Fee{}
	fee.Deposit = 0
	fee.WithDrawl = 0
	fee.Taker = bm.account.TakerFee
	fee.Maker = bm.account.MakerFee
	return fee
}
func (bm *Okex) doGetFee() exc.Fee {

	res, err := bm.doRequest("GET", "trade_fee", "")
	if err != nil {
		log.Println("GetFee err:", err)
	}
	log.Println("tradeFee", string(res))

	tradeFee := TradeFee{}
	json.Unmarshal(res, &tradeFee)

	fee := exc.Fee{}
	fee.Deposit = 0
	fee.WithDrawl = 0
	fee.Taker = tradeFee.Taker
	fee.Maker = tradeFee.Maker

	bm.account.TakerFee = tradeFee.Taker
	bm.account.MakerFee = tradeFee.Maker

	return fee
}

func (bm *Okex) GetName() string {
	return "Okex"
}

func (bm *Okex) GetMarketInfo(coinPair exc.CoinPair) exc.MarketInfo {
	return exc.MarketInfo{} //not yet implement
}

func (bm *Okex) GetVolumeByTotal(total, price float64) float64 {
	return total / price * 100
}

func (bm *Okex) GetFuturesAskBidPair(futures exc.Futures) (exc.PricePair, exc.PricePair) {
	booker := ob.StartOrderBook(bm.orderBookCenter, futures.GetSwapNameUpper())
	return booker.GetFirstPricePair()
}

func (bm *Okex) GetAccount() exc.Account {
	return bm.account
}

type QuoteResponse struct {
	exc.PriceStatus
	timeStamp string `json:"timeStamp"`
}

func (qr *QuoteResponse) setBy(json map[string]interface{}) {
	qr.timeStamp = json["timestamp"].(string)

	qr.PriceStatus.SetByJArray(json)
}

func (bm *Okex) GetAskBidPair(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	path := fmt.Sprintf("instruments/%s/book?size=%d",
		coinPair.GetSymbal(), depth)
	resByte, err := bm.doNormalRequest("GET", path, "")
	if err != nil {
		log.Println("GetAskBidPair err:", err)
	}

	var resJson map[string]interface{}

	err = json.Unmarshal(resByte, &resJson)
	if err != nil {
		log.Fatal(err)
	}

	quoteResponse := QuoteResponse{}
	quoteResponse.setBy(resJson)

	askPair, _ := quoteResponse.GetPair(1, exc.Ask)
	bidPair, _ := quoteResponse.GetPair(1, exc.Bid)

	return askPair, bidPair
}

func (bm *Okex) GetAccountInfo() []byte {

	str := ""
	out, err := json.Marshal(bm.GetWallet())
	if err != nil {
		panic(err)
	}
	str += string(out)

	out, err = json.Marshal(bm.GetLeverage())
	if err != nil {
		panic(err)
	}
	str += string(out)

	out, err = json.Marshal(bm.doGetFee())
	if err != nil {
		panic(err)
	}
	str += string(out)

	return []byte(str)
}

func (bm *Okex) GetProducts() []byte {
	res, err := bm.doNormalRequest("GET", "products", "")
	if err != nil {
		log.Println("GetAskBidPair err:", err)
	}
	return res
}

type OkexOrder struct {
	//Coid       string  `json:"client_oid,omitempty"`         //"xxx...xxx"     a unique identifier of length 32
	Symbol     string  `json:"instrument_id"`                // 合约名称，如BTC-USD-SWAP,BTC-USDT-SWAP
	Price      float64 `json:"orderPrice,string"`            //"13.5"          optional, limit price of the order. This field is required for limit orders and stop limit orders.
	Size       float64 `json:"size,string"`                  //（以张计数）
	Side       int     `json:"type,string"`                  //	可填参数：1:开多2:开空3:平多4:平空
	MatchPrice int     `json:"match_price,string,omitempty"` //是否以对手价下单。0:不是; 1:是。当以对手价下单，order_type只能选择0（普通委托）
	//OrderTYpe string `json:"order_type"` //0：普通委托（order_type不填或填0都是普通委托）1：只做Maker（Post only）2：全部成交或立即取消（FOK）3：立即成交并取消剩余（IOC）4：市价委托
	//PostOnly    bool   `json:"postOnly"`    //true            Optional, if true, the order will either be posted to the limit order book or be cancelled, i.e. the order cannot take liquidity; default value is false
	//TimeInForce string `json:"timeInForce"` //"GTC"           Optional, default is "GTC". Currently, we support "GTC" (good-till-canceled) and "IOC" (immediate-or-cancel).
}

const (
	OPEN_LONG   = 1
	OPEN_SHORT  = 2
	CLOSE_LONG  = 3
	CLOSE_SHORT = 4
)

func (oo *OkexOrder) setBy(order exc.ExchangeOrder) {

	// 只允許做市價單
	if order.OrderType == exc.LIMIT {
		// 因為無法判斷強制平倉或溢價只允許做市價單
	}

	oo.Symbol = order.Market
	panic("not yet implement ExchangeOrder")
}

func (oo *OkexOrder) setByFutures(order exc.FuturesOrder) {
	oo.Symbol = order.Futures.GetSwapNameUpper()

	//oo.Coid = exc.Uuid(32)

	if !order.IsClose {
		if order.Side == exc.Buy {
			oo.Side = OPEN_LONG
		} else {
			oo.Side = OPEN_SHORT
		}
	} else {
		if order.Side == exc.Buy {
			oo.Side = CLOSE_LONG
		} else {
			oo.Side = CLOSE_SHORT
		}
	}

	oo.Price = order.Price
	if order.OrderType == exc.MARKET {

		//oo.Price = order.Price
		oo.MatchPrice = 1 //1=yes,0=no
	} else {
		// 因為無法判斷強制平倉或溢價只允許做市價單
		oo.MatchPrice = 1
	}

	oo.Size = order.Size

}

func (ox *Okex) PostFuturesOrder(order exc.FuturesOrder) (string, error) {

	oo := OkexOrder{}
	oo.setByFutures(order)

	return ox.doPostOrder(oo)
}

func (ox *Okex) PostOrder(order exc.ExchangeOrder) (string, error) {

	oo := OkexOrder{}
	oo.setBy(order)

	return ox.doPostOrder(oo)
}

//下訂單
func (bm *Okex) doPostOrder(order OkexOrder) (string, error) {

	request, err := json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}
	body := string(request)
	log.Println(fmt.Sprintf("body:%s", body))

	response, err := bm.doPostRequest("order", body)
	if err != nil {
		log.Println("GetAskBidPair err:", err)
	}

	log.Println(fmt.Sprintf("%s", response))

	//{"error_message":"","result":"true","error_code":"0","client_oid":null,"order_id":"530018249871622144"}
	type OrderResponse struct {
		Code    float64 `json:"error_code,string"`
		Message string  `json:"error_message"`
	}
	orderResponse := OrderResponse{}

	if err := json.Unmarshal(response, &orderResponse); err != nil {
		panic(err)
	}

	var resErr error = nil
	if orderResponse.Code != 0 {
		resErr = errors.New(orderResponse.Message)
	}

	return string(response), resErr
}

func (bm *Okex) doNormalRequest(method, apiName, body string) ([]byte, error) {
	return bm.doRequest(method, apiName, body)
}

func (bm *Okex) doPostRequest(apiName, body string) ([]byte, error) {
	return bm.doRequest("POST", apiName, body)
}

func (bm *Okex) doRequest(method, apiName, body string) ([]byte, error) {
	ts := exc.GetTimeSpan()

	client := bm.client

	var res []byte

	fullUrl := fmt.Sprintf("%s%s%s", apiURL, apiPrefix, apiName)

	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Println(err)
		return res, err
	}

	req.Header.Set("Content-Type", "application/json")
	addHeader(&req.Header, method, apiName, ts, body)

	sendRes, err := exc.SendRequest(client, req)
	if err != nil {
		log.Println("okex SendRequest:", err)
	}
	log.Println("RES:", string(sendRes))

	//{"error_message":"","result":"true","error_code":"0","client_oid":null,"order_id":"530018249871622144"}
	type OrderResponse struct {
		Code    float64 `json:"error_code,string"`
		Message string  `json:"error_message"`
	}
	orderResponse := OrderResponse{}

	json.Unmarshal(sendRes, &orderResponse)

	var resErr error = nil
	if orderResponse.Code != 0 {
		resErr = errors.New(orderResponse.Message)
		panic(resErr)
	}

	return sendRes, err
}

func addHeader(header *http.Header, reqMethod, path string, ts int64, body string) {
	signPath := fmt.Sprintf("%v%v", apiPrefix, path)
	//exc.GetTimeSpan()

	tsStr := exc.GetUTC() // exc.GetTimeSpanStr(ts)
	total := tsStr + reqMethod + signPath + string(body)
	sign, _ := exc.GetParamHmacSHA256Base64Sign(setting.OKEX_SECRET_KEY, total)

	header.Add("OK-ACCESS-KEY", setting.OKEX_KEY)
	header.Add("OK-ACCESS-SIGN", sign)
	header.Add("OK-ACCESS-TIMESTAMP", tsStr)
	header.Add("OK-ACCESS-PASSPHRASE", setting.OKEX_PASSPHASE)

}
