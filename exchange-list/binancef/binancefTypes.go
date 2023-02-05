package binancef

import (
	"encoding/json"
	"fmt"
	"strconv"

	//"github.com/thrasher-corp/gocryptotrader/currency"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
	"github.com/yin75620/crypto-berserker/exchange/tool"
)

// Response holds basic binance api response data
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// ExchangeInfo holds the full exchange information type
type ExchangeInfo struct {
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	Timezone   string `json:"timezone"`
	Servertime int64  `json:"serverTime"`
	RateLimits []struct {
		RateLimitType string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		Limit         int    `json:"limit"`
	} `json:"rateLimits"`
	ExchangeFilters interface{} `json:"exchangeFilters"`
	Symbols         []struct {
		Symbol             string   `json:"symbol"`
		Status             string   `json:"status"`
		BaseAsset          string   `json:"baseAsset"`
		BaseAssetPrecision int      `json:"baseAssetPrecision"`
		QuoteAsset         string   `json:"quoteAsset"`
		QuotePrecision     int      `json:"quotePrecision"`
		OrderTypes         []string `json:"orderTypes"`
		IcebergAllowed     bool     `json:"icebergAllowed"`
		Filters            []struct {
			FilterType          string  `json:"filterType"` //PRICE_FILTER, MARKET_LOT_SIZE
			MinPrice            float64 `json:"minPrice,string"`
			MaxPrice            float64 `json:"maxPrice,string"`
			TickSize            float64 `json:"tickSize,string"`
			MultiplierUp        float64 `json:"multiplierUp,string"`
			MultiplierDown      float64 `json:"multiplierDown,string"`
			AvgPriceMins        int64   `json:"avgPriceMins"`
			MinQty              float64 `json:"minQty,string"`
			MaxQty              float64 `json:"maxQty,string"`
			StepSize            float64 `json:"stepSize,string"`
			MinNotional         float64 `json:"minNotional,string"`
			ApplyToMarket       bool    `json:"applyToMarket"`
			Limit               int64   `json:"limit"`
			MaxNumAlgoOrders    int64   `json:"maxNumAlgoOrders"`
			MaxNumIcebergOrders int64   `json:"maxNumIcebergOrders"`
		} `json:"filters"`
	} `json:"symbols"`
}

type LeverageBracket struct {
	Symbol   string `json:"symbol"`
	Brackets []struct {
		Bracket          int     `json:"bracket"`
		InitialLeverage  int     `json:"initialLeverage"`
		NotionalCap      int     `json:"notionalCap"`
		NotionalFloor    int     `json:"notionalFloor"`
		MaintMarginRatio float64 `json:"maintMarginRatio,string"`
		Cum              int     `json:"cum"`
	} `json:"brackets"`
}

// OrderBookDataRequestParams represents Klines request data.
type OrderBookDataRequestParams struct {
	Symbol string `json:"symbol"` // Required field; example LTCBTC,BTCUSDT
	Limit  int    `json:"limit"`  // Default 100; max 1000. Valid limits:[5, 10, 20, 50, 100, 500, 1000]
}

// OrderBookData is resp data from orderbook endpoint
type OrderBookData struct {
	Code         int           `json:"code"`
	Msg          string        `json:"msg"`
	LastUpdateID int64         `json:"T"`
	Bids         []interface{} `json:"b,[]string"`
	Asks         []interface{} `json:"a,[]string"`
}

func (obd *OrderBookData) ToOrderBookDetail() ob.OrderBookerResponseDetail {
	res := ob.OrderBookerResponseDetail{}
	res.Time = float64(obd.LastUpdateID)
	res.Action = ob.Partial
	res.Checksum = obd.LastUpdateID
	res.Asks = tool.TransToFloatTwoArray(obd.Asks)
	res.Bids = tool.TransToFloatTwoArray(obd.Bids)

	//transToFloatArray(&res.Asks, obd.Asks)
	//transToFloatArray(&res.Bids, obd.Bids)

	return res
}

func transToFloatArray(floatDoubleArray *[][]float64, data []interface{}) {
	for _, pItem := range data {
		pArray := pItem.([]float64)
		*floatDoubleArray = append(*floatDoubleArray, []float64{pArray[0], pArray[1]})
	}
}

// OrderBook actual structured data that can be used for orderbook
type OrderBook struct {
	LastUpdateID int64
	Code         int
	Msg          string
	Bids         []struct {
		Price    float64
		Quantity float64
	}
	Asks []struct {
		Price    float64
		Quantity float64
	}
}

// DepthUpdateParams is used as an embedded type for WebsocketDepthStream
type DepthUpdateParams []struct {
	PriceLevel float64
	Quantity   float64
	ingnore    []interface{}
}

// WebsocketDepthStream is the difference for the update depth stream
type WebsocketDepthStream struct {
	Event         string        `json:"e"`
	Timestamp     int64         `json:"E"`
	Pair          string        `json:"s"`
	FirstUpdateID int64         `json:"U"`
	LastUpdateID  int64         `json:"u"`
	UpdateBids    []interface{} `json:"b"`
	UpdateAsks    []interface{} `json:"a"`
}

// RecentTradeRequestParams represents Klines request data.
type RecentTradeRequestParams struct {
	Symbol string `json:"symbol"` // Required field. example LTCBTC, BTCUSDT
	Limit  int    `json:"limit"`  // Default 500; max 500.
}

// RecentTrade holds recent trade data
type RecentTrade struct {
	Code         int     `json:"code"`
	Msg          string  `json:"msg"`
	ID           int64   `json:"id"`
	Price        float64 `json:"price,string"`
	Quantity     float64 `json:"qty,string"`
	Time         float64 `json:"time"`
	IsBuyerMaker bool    `json:"isBuyerMaker"`
	IsBestMatch  bool    `json:"isBestMatch"`
}

// MultiStreamData holds stream data
type MultiStreamData struct {
	Stream string          `json:"stream"`
	Data   json.RawMessage `json:"data"`
}

// TradeStream holds the trade stream data
type TradeStream struct {
	EventType      string `json:"e"`
	EventTime      int64  `json:"E"`
	Symbol         string `json:"s"`
	TradeID        int64  `json:"t"`
	Price          string `json:"p"`
	Quantity       string `json:"q"`
	BuyerOrderID   int64  `json:"b"`
	SellerOrderID  int64  `json:"a"`
	TimeStamp      int64  `json:"T"`
	Maker          bool   `json:"m"`
	BestMatchPrice bool   `json:"M"`
}

// KlineStream holds the kline stream data
type KlineStream struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	Kline     struct {
		StartTime                int64  `json:"t"`
		CloseTime                int64  `json:"T"`
		Symbol                   string `json:"s"`
		Interval                 string `json:"i"`
		FirstTradeID             int64  `json:"f"`
		LastTradeID              int64  `json:"L"`
		OpenPrice                string `json:"o"`
		ClosePrice               string `json:"c"`
		HighPrice                string `json:"h"`
		LowPrice                 string `json:"l"`
		Volume                   string `json:"v"`
		NumberOfTrades           int64  `json:"n"`
		KlineClosed              bool   `json:"x"`
		Quote                    string `json:"q"`
		TakerBuyBaseAssetVolume  string `json:"V"`
		TakerBuyQuoteAssetVolume string `json:"Q"`
	} `json:"k"`
}

// TickerStream holds the ticker stream data
type TickerStream struct {
	EventType              string `json:"e"`
	EventTime              int64  `json:"E"`
	Symbol                 string `json:"s"`
	PriceChange            string `json:"p"`
	PriceChangePercent     string `json:"P"`
	WeightedAvgPrice       string `json:"w"`
	PrevDayClose           string `json:"x"`
	CurrDayClose           string `json:"c"`
	CloseTradeQuantity     string `json:"Q"`
	BestBidPrice           string `json:"b"`
	BestBidQuantity        string `json:"B"`
	BestAskPrice           string `json:"a"`
	BestAskQuantity        string `json:"A"`
	OpenPrice              string `json:"o"`
	HighPrice              string `json:"h"`
	LowPrice               string `json:"l"`
	TotalTradedVolume      string `json:"v"`
	TotalTradedQuoteVolume string `json:"q"`
	OpenTime               int64  `json:"O"`
	CloseTime              int64  `json:"C"`
	FirstTradeID           int64  `json:"F"`
	LastTradeID            int64  `json:"L"`
	NumberOfTrades         int64  `json:"n"`
}

// HistoricalTrade holds recent trade data
type HistoricalTrade struct {
	Code         int     `json:"code"`
	Msg          string  `json:"msg"`
	ID           int64   `json:"id"`
	Price        float64 `json:"price,string"`
	Quantity     float64 `json:"qty,string"`
	Time         int64   `json:"time"`
	IsBuyerMaker bool    `json:"isBuyerMaker"`
	IsBestMatch  bool    `json:"isBestMatch"`
}

// AggregatedTrade holds aggregated trade information
type AggregatedTrade struct {
	ATradeID       int64   `json:"a"`
	Price          float64 `json:"p,string"`
	Quantity       float64 `json:"q,string"`
	FirstTradeID   int64   `json:"f"`
	LastTradeID    int64   `json:"l"`
	TimeStamp      int64   `json:"T"`
	Maker          bool    `json:"m"`
	BestMatchPrice bool    `json:"M"`
}

// CandleStick holds kline data
type CandleStick struct {
	OpenTime                 float64
	Open                     float64
	High                     float64
	Low                      float64
	Close                    float64
	Volume                   float64
	CloseTime                float64
	QuoteAssetVolume         float64
	TradeCount               float64
	TakerBuyAssetVolume      float64
	TakerBuyQuoteAssetVolume float64
	//Ignore                   float64
}

func (cs *CandleStick) SetByJArray(array []interface{}) {
	if len(array) < 11 {
		fmt.Println("SetByJArray has error data", array)
	}
	cs.OpenTime = array[0].(float64)

	cs.Open = StringToFloat(array[1].(string))

	cs.High = StringToFloat(array[2].(string))
	cs.Low = StringToFloat(array[3].(string))
	cs.Close = StringToFloat(array[4].(string))
	cs.Volume = StringToFloat(array[5].(string))
	cs.CloseTime = array[6].(float64)
	cs.QuoteAssetVolume = StringToFloat(array[7].(string))
	cs.TradeCount = array[8].(float64)
	cs.TakerBuyAssetVolume = StringToFloat(array[9].(string))
	cs.TakerBuyQuoteAssetVolume = StringToFloat(array[10].(string))
}

func StringToFloat(s string) float64 {
	if fl, err := strconv.ParseFloat(s, 64); err == nil {
		return fl
	}
	return 0
}

// AveragePrice holds current average symbol price
type AveragePrice struct {
	Mins  int64   `json:"mins"`
	Price float64 `json:"price,string"`
}

// PriceChangeStats contains statistics for the last 24 hours trade
type PriceChangeStats struct {
	Symbol             string  `json:"symbol"`
	PriceChange        float64 `json:"priceChange,string"`
	PriceChangePercent float64 `json:"priceChangePercent,string"`
	WeightedAvgPrice   float64 `json:"weightedAvgPrice,string"`
	PrevClosePrice     float64 `json:"prevClosePrice,string"`
	LastPrice          float64 `json:"lastPrice,string"`
	LastQty            float64 `json:"lastQty,string"`
	BidPrice           float64 `json:"bidPrice,string"`
	AskPrice           float64 `json:"askPrice,string"`
	OpenPrice          float64 `json:"openPrice,string"`
	HighPrice          float64 `json:"highPrice,string"`
	LowPrice           float64 `json:"lowPrice,string"`
	Volume             float64 `json:"volume,string"`
	QuoteVolume        float64 `json:"quoteVolume,string"`
	OpenTime           int64   `json:"openTime"`
	CloseTime          int64   `json:"closeTime"`
	FirstID            int64   `json:"fristId"`
	LastID             int64   `json:"lastId"`
	Count              int64   `json:"count"`
}

// SymbolPrice holds basic symbol price
type SymbolPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

// BestPrice holds best price data
type BestPrice struct {
	Symbol   string  `json:"symbol"`
	BidPrice float64 `json:"bidPrice,string"`
	BidQty   float64 `json:"bidQty,string"`
	AskPrice float64 `json:"askPrice,string"`
	AskQty   float64 `json:"askQty,string"`
}

// NewOrderRequest request type
type NewOrderRequest struct {
	// Symbol (currency pair to trade)
	Symbol string
	// Side Buy or Sell
	Side RequestParamsSideType
	// TradeType (market or limit order)
	TradeType RequestParamsOrderType
	// TimeInForce specifies how long the order remains in effect.
	// Examples are (Good Till Cancel (GTC), Immediate or Cancel (IOC) and Fill Or Kill (FOK))
	TimeInForce RequestParamsTimeForceType
	// Quantity
	Quantity         float64
	Price            float64
	NewClientOrderID string
	StopPrice        float64 // Used with STOP_LOSS, STOP_LOSS_LIMIT, TAKE_PROFIT, and TAKE_PROFIT_LIMIT orders.
	IcebergQty       float64 // Used with LIMIT, STOP_LOSS_LIMIT, and TAKE_PROFIT_LIMIT to create an iceberg order.
	NewOrderRespType string
}

// NewOrderResponse is the return structured response from the exchange
type NewOrderResponse struct {
	Code            int     `json:"code"`
	Msg             string  `json:"msg"`
	Symbol          string  `json:"symbol"`
	OrderID         int64   `json:"orderId"`
	ClientOrderID   string  `json:"clientOrderId"`
	TransactionTime int64   `json:"transactTime"`
	Price           float64 `json:"price,string"`
	OrigQty         float64 `json:"origQty,string"`
	ExecutedQty     float64 `json:"executedQty,string"`
	Status          string  `json:"status"`
	TimeInForce     string  `json:"timeInForce"`
	Type            string  `json:"type"`
	Side            string  `json:"side"`
	Fills           []struct {
		Price           float64 `json:"price,string"`
		Qty             float64 `json:"qty,string"`
		Commission      float64 `json:"commission,string"`
		CommissionAsset float64 `json:"commissionAsset,string"`
	} `json:"fills"`
}

// CancelOrderResponse is the return structured response from the exchange
type CancelOrderResponse struct {
	Symbol            string `json:"symbol"`
	OrigClientOrderID string `json:"origClientOrderId"`
	OrderID           int64  `json:"orderId"`
	ClientOrderID     string `json:"clientOrderId"`
}

// QueryOrderData holds query order data
type QueryOrderData struct {
	Code          int     `json:"code"`
	Msg           string  `json:"msg"`
	Symbol        string  `json:"symbol"`
	OrderID       int64   `json:"orderId"`
	ClientOrderID string  `json:"clientOrderId"`
	Price         float64 `json:"price,string"`
	OrigQty       float64 `json:"origQty,string"`
	ExecutedQty   float64 `json:"executedQty,string"`
	Status        string  `json:"status"`
	TimeInForce   string  `json:"timeInForce"`
	Type          string  `json:"type"`
	Side          string  `json:"side"`
	StopPrice     float64 `json:"stopPrice,string"`
	IcebergQty    float64 `json:"icebergQty,string"`
	Time          float64 `json:"time"`
	IsWorking     bool    `json:"isWorking"`
}

// Balance holds query order data
type Asset struct {
	Asset                  string  `json:"asset"`
	InitialMargin          float64 `json:"initialMargin,string"`
	MainMargin             float64 `json:"mainMargin,string"`
	MarginBalance          float64 `json:"marginBalance,string"`
	MaxWithdrawAmount      float64 `json:"maxWithdrawAmount,string"`
	OpenOrderInitialMargin float64 `json:"openOrderInitialMargin,string"`
	PositionInitialMargin  float64 `json:"positionInitialMargin,string"`
	UnrealizedProfit       float64 `json:"unrealizedProfit,string"`
	WalletBalance          float64 `json:"walletBalance,string"`
}

type Account struct {
	FeeTier                     int     `json:"feeTier"`
	CanTrade                    bool    `json:"canTrade"`
	CanDeposit                  bool    `json:"canDeposit"`
	CanWithdraw                 bool    `json:"canWithdraw"`
	UpdateTime                  int     `json:"updateTime"`
	MultiAssetsMargin           bool    `json:"multiAssetsMargin"`
	TotalInitialMargin          float64 `json:"totalInitialMargin,string"`
	TotalMaintMargin            float64 `json:"totalMaintMargin,string"`
	TotalWalletBalance          float64 `json:"totalWalletBalance,string"`
	TotalUnrealizedProfit       float64 `json:"totalUnrealizedProfit,string"`
	TotalMarginBalance          float64 `json:"totalMarginBalance,string"`
	TotalPositionInitialMargin  float64 `json:"totalPositionInitialMargin,string"`
	TotalOpenOrderInitialMargin float64 `json:"totalOpenOrderInitialMargin,string"`
	TotalCrossWalletBalance     float64 `json:"totalCrossWalletBalance,string"`
	TotalCrossUnPnl             float64 `json:"totalCrossUnPnl,string"`
	AvailableBalance            float64 `json:"availableBalance,string"`
	MaxWithdrawAmount           float64 `json:"maxWithdrawAmount,string"`
	Assets                      []struct {
		Asset                  string  `json:"asset"`
		WalletBalance          float64 `json:"walletBalance,string"`
		UnrealizedProfit       float64 `json:"unrealizedProfit,string"`
		MarginBalance          float64 `json:"marginBalance,string"`
		MaintMargin            float64 `json:"maintMargin,string"`
		InitialMargin          float64 `json:"initialMargin,string"`
		PositionInitialMargin  float64 `json:"positionInitialMargin,string"`
		OpenOrderInitialMargin float64 `json:"openOrderInitialMargin,string"`
		CrossWalletBalance     float64 `json:"crossWalletBalance,string"`
		CrossUnPnl             float64 `json:"crossUnPnl,string"`
		AvailableBalance       float64 `json:"availableBalance,string"`
		MaxWithdrawAmount      float64 `json:"maxWithdrawAmount,string"`
		MarginAvailable        bool    `json:"marginAvailable"`
		UpdateTime             int64   `json:"updateTime"`
	} `json:"assets"`
	Positions []Position `json:"positions"`
}

type Position struct {
	Symbol                 string  `json:"symbol"`
	InitialMargin          float64 `json:"initialMargin,string"`
	MaintMargin            float64 `json:"maintMargin,string"`
	UnrealizedProfit       float64 `json:"unrealizedProfit,string"`
	PositionInitialMargin  float64 `json:"positionInitialMargin,string"`
	OpenOrderInitialMargin float64 `json:"openOrderInitialMargin,string"`
	Leverage               int     `json:"leverage,string"`
	Isolated               bool    `json:"isolated"`
	EntryPrice             float64 `json:"entryPrice,string"`
	MaxNotional            float64 `json:"maxNotional,string"`
	BidNotional            float64 `json:"bidNotional,string"`
	AskNotional            float64 `json:"askNotional,string"`
	PositionSide           string  `json:"positionSide"`
	PositionAmt            float64 `json:"positionAmt,string"`
	UpdateTime             int64   `json:"updateTime"`
}

type UserTrade struct {
	Buyer           bool    `json:"buyer"`
	Commission      float64 `json:"commission,string"`
	CommissionAsset string  `json:"commissionAsset"`
	ID              int64   `json:"id"`
	Maker           bool    `json:"maker"`
	OrderID         int64   `json:"orderId"`
	Price           float64 `json:"price,string"`
	Qty             float64 `json:"qty,string"`
	QuoteQty        float64 `json:"quoteQty,string"`
	RealizedPnl     float64 `json:"realizedPnl,string"`
	Side            string  `json:"side"`
	PositionSide    string  `json:"positionSide"`
	Symbol          string  `json:"symbol"`
	Time            int64   `json:"time"`
}

type OrderResponse struct {
	Code          int     `json:"code"`
	Msg           string  `json:"msg"`
	ClientOrderID string  `json:"clientOrderId"`
	CumQty        float64 `json:"cumQty,string"`
	CumQuote      float64 `json:"cumQuote,string"`
	ExecutedQty   float64 `json:"executedQty,string"`
	OrderID       int     `json:"orderId"`
	AvgPrice      float64 `json:"avgPrice,string"`
	OrigQty       float64 `json:"origQty,string"`
	Price         float64 `json:"price,string"`
	ReduceOnly    bool    `json:"reduceOnly"`
	Side          string  `json:"side"`
	PositionSide  string  `json:"positionSide"`
	Status        string  `json:"status"`
	StopPrice     float64 `json:"stopPrice,string"`
	ClosePosition bool    `json:"closePosition"`
	Symbol        string  `json:"symbol"`
	TimeInForce   string  `json:"timeInForce"`
	Type          string  `json:"type"`
	OrigType      string  `json:"origType"`
	ActivatePrice float64 `json:"activatePrice,string"`
	PriceRate     float64 `json:"priceRate,string"`
	UpdateTime    int64   `json:"updateTime"`
	WorkingType   string  `json:"workingType"`
	PriceProtect  bool    `json:"priceProtect"`
}

/*
// Account holds the account data
type Account struct {
	Assets                []Asset `json:"assets"`
	TotalUnrealizedProfit float64 `json:"totalUnrealizedProfit,string"`
	/*"canDeposit": true,
	  "canTrade": true,
	  "canWithdraw": true,
	  "feeTier": 2,
	  "maxWithdrawAmount": "8.41264592",
	  "positions": [
	      {
	         "isolated": false,
	         "leverage": "20",
	         "initialMargin": "0.33683",
	         "maintMargin": "0.02695",
	         "openOrderInitialMargin": "0.00000",
	         "positionInitialMargin": "0.33683",
	         "symbol": "BTCUSDT",
	         "unrealizedProfit": "-0.44537584",
	         "positionSide": "BOTH", // BOTH means that it is the position of One-way Mode
	      },
	      {
	         "isolated": false,
	         "leverage": "20",
	         "initialMargin": "0.00000",
	         "maintMargin": "0.00000",
	         "openOrderInitialMargin": "0.00000",
	         "positionInitialMargin": "0.00000",
	         "symbol": "BTCUSDT",
	         "unrealizedProfit": "0.00000000",
	         "positionSide": "LONG", // LONG or SHORT means that it is the position of Hedge Mode
	      },
	      {
	         "isolated": false,
	         "leverage": "20",
	         "initialMargin": "0.00000",
	         "maintMargin": "0.00000",
	         "openOrderInitialMargin": "0.00000",
	         "positionInitialMargin": "0.00000",
	         "symbol": "BTCUSDT",
	         "unrealizedProfit": "0.00000000",
	         "positionSide": "SHORT", // LONG or SHORT means that it is the position of One-way Mode
	      }
	  ],
	  "totalInitialMargin": "0.33683000",
	  "totalMaintMargin": "0.02695000",
	  "totalMarginBalance": "8.74947592",
	  "totalOpenOrderInitialMargin": "0.00000000",
	  "totalPositionInitialMargin": "0.33683000",
	  "totalUnrealizedProfit": "-0.44537584",
	  "totalWalletBalance": "9.19485176",
	  "updateTime": 0

}*/

// RequestParamsSideType trade order side (buy or sell)
type RequestParamsSideType string

var (
	// BinanceRequestParamsSideBuy buy order type
	BinanceRequestParamsSideBuy = RequestParamsSideType("BUY")

	// BinanceRequestParamsSideSell sell order type
	BinanceRequestParamsSideSell = RequestParamsSideType("SELL")
)

// RequestParamsTimeForceType Time in force
type RequestParamsTimeForceType string

var (
	// BinanceRequestParamsTimeGTC GTC
	BinanceRequestParamsTimeGTC = RequestParamsTimeForceType("GTC")

	// BinanceRequestParamsTimeIOC IOC
	BinanceRequestParamsTimeIOC = RequestParamsTimeForceType("IOC")

	// BinanceRequestParamsTimeFOK FOK
	BinanceRequestParamsTimeFOK = RequestParamsTimeForceType("FOK")
)

// RequestParamsOrderType trade order type
type RequestParamsOrderType string

var (
	// BinanceRequestParamsOrderLimit Limit order
	BinanceRequestParamsOrderLimit = RequestParamsOrderType("LIMIT")

	// BinanceRequestParamsOrderMarket Market order
	BinanceRequestParamsOrderMarket = RequestParamsOrderType("MARKET")

	// BinanceRequestParamsOrderStopLoss STOP_LOSS
	BinanceRequestParamsOrderStopLoss = RequestParamsOrderType("STOP_LOSS")

	// BinanceRequestParamsOrderStopLossLimit STOP_LOSS_LIMIT
	BinanceRequestParamsOrderStopLossLimit = RequestParamsOrderType("STOP_LOSS_LIMIT")

	// BinanceRequestParamsOrderTakeProfit TAKE_PROFIT
	BinanceRequestParamsOrderTakeProfit = RequestParamsOrderType("TAKE_PROFIT")

	// BinanceRequestParamsOrderTakeProfitLimit TAKE_PROFIT_LIMIT
	BinanceRequestParamsOrderTakeProfitLimit = RequestParamsOrderType("TAKE_PROFIT_LIMIT")

	// BinanceRequestParamsOrderLimitMarker LIMIT_MAKER
	BinanceRequestParamsOrderLimitMarker = RequestParamsOrderType("LIMIT_MAKER")
)

// KlinesRequestParams represents Klines request data.
type KlinesRequestParams struct {
	Symbol    string       // Required field; example LTCBTC, BTCUSDT
	Interval  TimeInterval // Time interval period
	Limit     int          // Default 500; max 500.
	StartTime int64
	EndTime   int64
}

// TimeInterval represents interval enum.
type TimeInterval string

// Vars related to time intervals
var (
	TimeIntervalMinute         = TimeInterval("1m")
	TimeIntervalThreeMinutes   = TimeInterval("3m")
	TimeIntervalFiveMinutes    = TimeInterval("5m")
	TimeIntervalFifteenMinutes = TimeInterval("15m")
	TimeIntervalThirtyMinutes  = TimeInterval("30m")
	TimeIntervalHour           = TimeInterval("1h")
	TimeIntervalTwoHours       = TimeInterval("2h")
	TimeIntervalFourHours      = TimeInterval("4h")
	TimeIntervalSixHours       = TimeInterval("6h")
	TimeIntervalEightHours     = TimeInterval("8h")
	TimeIntervalTwelveHours    = TimeInterval("12h")
	TimeIntervalDay            = TimeInterval("1d")
	TimeIntervalThreeDays      = TimeInterval("3d")
	TimeIntervalWeek           = TimeInterval("1w")
	TimeIntervalMonth          = TimeInterval("1M")
)

// WithdrawalFees the large list of predefined withdrawal fees
// Prone to change
/*
var WithdrawalFees = map[currency.Code]float64{
	currency.BNB:     0.13,
	currency.BTC:     0.0005,
	currency.NEO:     0,
	currency.ETH:     0.01,
	currency.LTC:     0.001,
	currency.QTUM:    0.01,
	currency.EOS:     0.1,
	currency.SNT:     35,
	currency.BNT:     1,
	currency.GAS:     0,
	currency.BCC:     0.001,
	currency.BTM:     5,
	currency.USDT:    3.4,
	currency.HCC:     0.0005,
	currency.OAX:     6.5,
	currency.DNT:     54,
	currency.MCO:     0.31,
	currency.ICN:     3.5,
	currency.ZRX:     1.9,
	currency.OMG:     0.4,
	currency.WTC:     0.5,
	currency.LRC:     12.3,
	currency.LLT:     67.8,
	currency.YOYO:    1,
	currency.TRX:     1,
	currency.STRAT:   0.1,
	currency.SNGLS:   54,
	currency.BQX:     3.9,
	currency.KNC:     3.5,
	currency.SNM:     25,
	currency.FUN:     86,
	currency.LINK:    4,
	currency.XVG:     0.1,
	currency.CTR:     35,
	currency.SALT:    2.3,
	currency.MDA:     2.3,
	currency.IOTA:    0.5,
	currency.SUB:     11.4,
	currency.ETC:     0.01,
	currency.MTL:     2,
	currency.MTH:     45,
	currency.ENG:     2.2,
	currency.AST:     14.4,
	currency.DASH:    0.002,
	currency.BTG:     0.001,
	currency.EVX:     2.8,
	currency.REQ:     29.9,
	currency.VIB:     30,
	currency.POWR:    8.2,
	currency.ARK:     0.2,
	currency.XRP:     0.25,
	currency.MOD:     2,
	currency.ENJ:     26,
	currency.STORJ:   5.1,
	currency.KMD:     0.002,
	currency.RCN:     47,
	currency.NULS:    0.01,
	currency.RDN:     2.5,
	currency.XMR:     0.04,
	currency.DLT:     19.8,
	currency.AMB:     8.9,
	currency.BAT:     8,
	currency.ZEC:     0.005,
	currency.BCPT:    14.5,
	currency.ARN:     3,
	currency.GVT:     0.13,
	currency.CDT:     81,
	currency.GXS:     0.3,
	currency.POE:     134,
	currency.QSP:     36,
	currency.BTS:     1,
	currency.XZC:     0.02,
	currency.LSK:     0.1,
	currency.TNT:     47,
	currency.FUEL:    79,
	currency.MANA:    18,
	currency.BCD:     0.01,
	currency.DGD:     0.04,
	currency.ADX:     6.3,
	currency.ADA:     1,
	currency.PPT:     0.41,
	currency.CMT:     12,
	currency.XLM:     0.01,
	currency.CND:     58,
	currency.LEND:    84,
	currency.WABI:    6.6,
	currency.SBTC:    0.0005,
	currency.BCX:     0.5,
	currency.WAVES:   0.002,
	currency.TNB:     139,
	currency.GTO:     20,
	currency.ICX:     0.02,
	currency.OST:     32,
	currency.ELF:     3.9,
	currency.AION:    3.2,
	currency.CVC:     10.9,
	currency.REP:     0.2,
	currency.GNT:     8.9,
	currency.DATA:    37,
	currency.ETF:     1,
	currency.BRD:     3.8,
	currency.NEBL:    0.01,
	currency.VIBE:    17.3,
	currency.LUN:     0.36,
	currency.CHAT:    60.7,
	currency.RLC:     3.4,
	currency.INS:     3.5,
	currency.IOST:    105.6,
	currency.STEEM:   0.01,
	currency.NANO:    0.01,
	currency.AE:      1.3,
	currency.VIA:     0.01,
	currency.BLZ:     10.3,
	currency.SYS:     1,
	currency.NCASH:   247.6,
	currency.POA:     0.01,
	currency.ONT:     1,
	currency.ZIL:     37.2,
	currency.STORM:   152,
	currency.XEM:     4,
	currency.WAN:     0.1,
	currency.WPR:     43.4,
	currency.QLC:     1,
	currency.GRS:     0.2,
	currency.CLOAK:   0.02,
	currency.LOOM:    11.9,
	currency.BCN:     1,
	currency.TUSD:    1.35,
	currency.ZEN:     0.002,
	currency.SKY:     0.01,
	currency.THETA:   24,
	currency.IOTX:    90.5,
	currency.QKC:     24.6,
	currency.AGI:     29.81,
	currency.NXS:     0.02,
	currency.SC:      0.1,
	currency.EON:     10,
	currency.NPXS:    897,
	currency.KEY:     223,
	currency.NAS:     0.1,
	currency.ADD:     100,
	currency.MEETONE: 300,
	currency.ATD:     100,
	currency.MFT:     175,
	currency.EOP:     5,
	currency.DENT:    596,
	currency.IQ:      50,
	currency.ARDR:    2,
	currency.HOT:     1210,
	currency.VET:     100,
	currency.DOCK:    68,
	currency.POLY:    7,
	currency.VTHO:    21,
	currency.ONG:     0.1,
	currency.PHX:     1,
	currency.HC:      0.005,
	currency.GO:      0.01,
	currency.PAX:     1.4,
	currency.EDO:     1.3,
	currency.WINGS:   8.9,
	currency.NAV:     0.2,
	currency.TRIG:    49.1,
	currency.APPC:    12.4,
	currency.PIVX:    0.02,
}
*/
// WithdrawResponse contains status of withdrawal request
type WithdrawResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	ID      int64  `json:"id"`
}
