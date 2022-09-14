package ftx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	exc "github.com/yin75620/crypto-berserker/exchange"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
)

type PriceType exc.PriceType

type OrderBookResponse struct {
	Result  exc.PriceStatus `json:"result"`
	Success bool            `json:"success"`
}

type Ftx struct {
	client          *http.Client
	initData        FtxInit
	orderBookCenter *ob.OrderBookCenter

	account exc.Account
}

func NewFtx(c *http.Client, initData FtxInit) *Ftx {
	ftx := &Ftx{}
	ftx.client = c
	ftx.initData = initData
	ftx.orderBookCenter = ob.NewOrderBookCenter(NewSocket())
	return ftx
}

func (ftx *Ftx) SetInitByIni(filename string) {
	ftx.initData.IniSetting(filename)
}

var (
	apiURL    = "https://ftx.com"
	apiPrefix = "/api/"
)

// implement exchange
func (ftx *Ftx) GetFee() exc.Fee {
	fee := exc.Fee{}
	fee.Deposit = 0
	fee.WithDrawl = 0
	fee.Taker = 0.00063175
	fee.Maker = 0.0001805
	return fee
}

func (ftx *Ftx) GetName() string {
	return "FTX"
}

func (ftx *Ftx) GetMarketInfo(coinPair exc.CoinPair) exc.MarketInfo {
	switch name := coinPair.GetMarketName(); name {
	case "BTC/USD":
		return exc.MarketInfo{VolumeIncrement: 0.0001}
	case "FTT/USD":
		return exc.MarketInfo{VolumeIncrement: 1}
	case "FTT/BTC":
		return exc.MarketInfo{VolumeIncrement: 1}
	default:
		return exc.MarketInfo{}
	}
}

func (ftx *Ftx) GetVolumeByTotal(total, price float64) float64 {
	return total / price
}

func (ftx *Ftx) GetAccountInfo() []byte {
	account := ftx.doGet("account", "")
	accountResult := AccountResponse{}
	err := json.Unmarshal(account, &accountResult)
	if err != nil {
		fmt.Println("GetAccountInfo", err)
	}

	ftx.account.TakerFee = accountResult.Result.TakerFee
	ftx.account.Leverage = accountResult.Result.Leverage
	ftx.account.MakerFee = accountResult.Result.MakerFee
	ftx.account.UnrealizedPnL = accountResult.Result.GetTotalUnrealizedPnl()

	wallet := ftx.GetWallet()

	return []byte(fmt.Sprintf("%v%s", wallet, string(account)))
}

func (ftx *Ftx) GetAccount() exc.Account {
	return ftx.account

}

func (ftx *Ftx) GetWallet() exc.Wallet {
	type BalanceResponse struct {
		Result []exc.Balance `json:"result,omitempty"`
		Sucess string        `json:"success,omitempty"`
	}
	response := BalanceResponse{}

	res := ftx.doGet("wallet/balances", "")
	json.Unmarshal(res, &response)

	w := exc.Wallet{}
	w.Balances = response.Result
	w.CalculerFreeUsdValue()

	ftx.account.WalletInfo = w
	return w
}

// Markets

// resulution: window length in seconds. options: 15, 60, 300, 900, 3600, 14400, 86400, or any multiple of 86400 up to 30*86400
func (ftx *Ftx) GetHistoricalPrices(market string, resolution int64,
	limit int64, startTime int64, endTime int64) (HistoricalPrices, error) {
	var historicalPrices HistoricalPrices
	resp := ftx.doGet(
		"markets/"+market+
			"/candles?resolution="+strconv.FormatInt(resolution, 10)+
			"&limit="+strconv.FormatInt(limit, 10)+
			"&start_time="+strconv.FormatInt(startTime, 10)+
			"&end_time="+strconv.FormatInt(endTime, 10),
		"")
	err := json.Unmarshal(resp, &historicalPrices)
	if err != nil {
		log.Println("Error GetHistoricalPrices:", err)
		return historicalPrices, err
	}
	return historicalPrices, err
}

func (ftx *Ftx) GetTrades(market string, limit int64, startTime int64, endTime int64) (Trades, error) {
	var trades Trades
	resp := ftx.doGet(
		"markets/"+market+"/trades?"+
			"&limit="+strconv.FormatInt(limit, 10)+
			"&start_time="+strconv.FormatInt(startTime, 10)+
			"&end_time="+strconv.FormatInt(endTime, 10),
		"")
	err := json.Unmarshal(resp, &trades)
	if err != nil {
		log.Println("Error GetTrades:", err)
		return trades, err
	}
	return trades, err
}

///////////

func (ftx *Ftx) GetStructFills() []FillResponse {
	res := ftx.GetFills()
	type tempResponse struct {
		Result []FillResponse `json:"result"`
	}

	tempRes := tempResponse{}

	json.Unmarshal(res, &tempRes)
	return tempRes.Result
}

func (ftx *Ftx) GetFills() []byte {
	return ftx.doGet("fills", "")
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

func (ftx *Ftx) GetAskBidPair(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	market := coinPair.GetMarketName()
	return ftx.getAskBidPairByMarket(market)
}

func (ftx *Ftx) GetFuturesAskBidPair(futures exc.Futures) (exc.PricePair, exc.PricePair) {
	market := futures.GetMarketName()
	ask, bid := ftx.getAskBidPairByMarket(market)
	return ask, bid
}

func (ftx *Ftx) getAskBidPairByMarket(market string) (exc.PricePair, exc.PricePair) {
	if !ftx.orderBookCenter.IsExist(market) {
		channel, _ := ftx.orderBookCenter.Register(market)
		<-channel
		go func() {
			for {
				<-channel
			}
		}()

		//return ftx.getOrderBookFromWeb(coinPair, depth)
	}

	booker := ftx.orderBookCenter.GetBooker(market)
	return booker.GetFirstPricePair()
}

func (ftx *Ftx) getOrderBookFromWeb(coinPair exc.CoinPair, depth int) (exc.PricePair, exc.PricePair) {
	resb := ftx.GetOrderBookResponse(coinPair.GetMarketName(), depth)
	askPair, _ := resb.Result.GetPair(1, exc.Ask)
	bidPair, _ := resb.Result.GetPair(1, exc.Bid)
	return askPair, bidPair
}

func (ftx *Ftx) GetOrderBookResponse(marketName string, depth int) OrderBookResponse {
	response := ftx.GetOrderBook(marketName, depth)
	var bookResponse OrderBookResponse
	json.Unmarshal(response, &bookResponse)
	return bookResponse
}

func (ftx *Ftx) GetOrderBook(marketName string, depth int) []byte {
	path := fmt.Sprintf("/markets/%s/orderbook?depth=%d", marketName, depth)
	response := ftx.doGet(path, "")
	return response
}

type EOrderType string

const (
	LIMIT  EOrderType = "limit"
	MARKET EOrderType = "market"
)

type FtxOrder struct {
	Market    string     `json:"market"`
	Side      string     `json:"side"`
	Price     float64    `json:"price"`
	Size      float64    `json:"size"`
	OrderType EOrderType `json:"order_type"`
	//ReduceOnly bool       `json:"reduceOnly"`
}

func (fo *FtxOrder) setBy(order exc.ExchangeOrder) {
	fo.Market = order.CoinPair.GetMarketName()
	fo.Side = order.Side
	fo.Price = order.Price
	fo.Size = order.Size
	fo.OrderType = EOrderType(order.OrderType)
}

func (fo *FtxOrder) setByFutures(order exc.FuturesOrder) {
	fo.Market = order.Futures.GetMarketName()
	fo.Side = order.Side
	fo.Price = order.Price
	fo.Size = order.Size
	fo.OrderType = EOrderType(order.OrderType)
}

//下訂單
func (ftx *Ftx) PostOrder(order exc.ExchangeOrder) (string, error) {

	fo := FtxOrder{}
	fo.setBy(order)

	return ftx.doPostOrder(fo)
}

func (ftx *Ftx) DeleteAllOrders() (string, error) {
	res, err := ftx.doRequest("DELETE", "orders", "")
	return string(res), err
}

func (ftx *Ftx) PostFuturesOrder(order exc.FuturesOrder) (string, error) {

	fo := FtxOrder{}
	fo.setByFutures(order)

	return ftx.doPostOrder(fo)
}

func (ftx *Ftx) doPostOrder(fo FtxOrder) (string, error) {
	request, err := json.Marshal(fo)
	if err != nil {
		log.Fatal(err)
	}
	body := string(request)
	log.Println(fmt.Sprintf("body:%s", body))
	response, err := ftx.doPost("orders", body)
	if err != nil {
		log.Println("doPostOrder fail.", err)
		return "", err
	}
	log.Println(fmt.Sprintf("%s", response))

	//{"error":"Not enough balances","success":false}
	type OrderResponse struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}
	orderResponse := OrderResponse{}

	json.Unmarshal(response, &orderResponse)

	var resErr error = nil
	if !orderResponse.Success {
		resErr = errors.New(orderResponse.Error)
	}

	return string(response), resErr
}

func (ftx *Ftx) doGet(apiName, body string) []byte {
	//數量有點多，err先不往上丟
	res, _ := ftx.doRequest("GET", apiName, body)
	return res
}

func (ftx *Ftx) doPost(apiName, body string) ([]byte, error) {
	return ftx.doRequest("POST", apiName, body)
}

func (ftx *Ftx) doRequest(method, apiName, body string) ([]byte, error) {
	client := ftx.client

	var res []byte

	fullUrl := fmt.Sprintf("%s%s%s", apiURL, apiPrefix, apiName)
	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Println(err)
		return res, err
	}

	req.Header.Set("Content-Type", "application/json")
	ftx.addHeader(&req.Header, method, apiName, body)

	return exc.SendRequest(client, req)
}

func (ftx *Ftx) addHeader(header *http.Header, reqMethod, path, body string) {
	ts := exc.GetTimeSpanStr(exc.GetTimeSpan())

	initData := ftx.initData

	header.Add("FTX-KEY", initData.ApiKey)
	header.Add("FTX-TS", ts)
	payload := fmt.Sprintf("%s%s%s%s", ts, reqMethod, apiPrefix+path, body)
	sign, _ := exc.GetParamHmacSHA256HexSign(initData.ApiSecretKey, payload)
	header.Add("FTX-SIGN", sign)
	header.Add("FTX-SUBACCOUNT", initData.SubAccount)
}

// Lending
func (ftx *Ftx) GetLendingRate() []byte {
	res := ftx.doGet("spot_margin/lending_rates", "")
	return res
}

func (ftx *Ftx) GetLendingInfo() LendInfoResponse {
	byteRes := ftx.doGet("spot_margin/lending_info", "")
	response := LendInfoResponse{}
	err := json.Unmarshal(byteRes, &response)
	if err != nil {
		log.Println("GetLendingInfo json.Unmashal error:", err)
		return response
	}

	return response
}

func (ftx *Ftx) PostLendingOffer(lo LendOrder) error {
	body, err := json.Marshal(lo)
	if err != nil {
		log.Println("json.Marshal", err)
		return err
	}

	byteRes, err := ftx.doPost("spot_margin/offers", string(body))
	if err != nil {
		log.Println("ftx.doPost", err)
		return err
	}

	lr := LendOrderResponse{}
	err = json.Unmarshal(byteRes, &lr)
	if err != nil {
		log.Println("json.Unmarshal", err)
		return err
	}

	if !lr.Success {
		log.Println(lr.Result)
		return errors.New(lr.Result)
	}

	return nil
}
