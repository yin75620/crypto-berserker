package bybilinear

import (
	"encoding/json"
	"strings"
	"time"
)

type BaseResult struct {
	RetCode         int         `json:"ret_code"`
	RetMsg          string      `json:"ret_msg"`
	ExtCode         string      `json:"ext_code"`
	Result          interface{} `json:"result"`
	TimeNow         string      `json:"time_now"`
	RateLimitStatus int         `json:"rate_limit_status"`
}

type Item struct {
	Price float64 `json:"price,string"`
	Size  float64 `json:"size"`
}

type OrderBook struct {
	Asks []Item    `json:"asks"`
	Bids []Item    `json:"bids"`
	Time time.Time `json:"time"`
}

type RawItem struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
	Size   float64 `json:"size"`
	Side   string  `json:"side"` // Buy/Sell
}

type GetOrderBookResult struct {
	RetCode int       `json:"ret_code"`
	RetMsg  string    `json:"ret_msg"`
	ExtCode string    `json:"ext_code"`
	ExtInfo string    `json:"ext_info"`
	Result  []RawItem `json:"result"`
	TimeNow string    `json:"time_now"`
}

type Balance struct {
	Equity           float64 `json:"equity"`
	AvailableBalance float64 `json:"available_balance"`
	UsedMargin       float64 `json:"used_margin"`
	OrderMargin      float64 `json:"order_margin"`
	PositionMargin   float64 `json:"position_margin"`
	OccClosingFee    float64 `json:"occ_closing_fee"`
	OccFundingFee    float64 `json:"occ_funding_fee"`
	WalletBalance    float64 `json:"wallet_balance"`
	RealisedPnl      float64 `json:"realised_pnl"`
	UnrealisedPnl    float64 `json:"unrealised_pnl"`
	CumRealisedPnl   float64 `json:"cum_realised_pnl"`
	GivenCash        float64 `json:"given_cash"`
	ServiceCash      float64 `json:"service_cash"`
}

type GetBalanceResult struct {
	RetCode          int                  `json:"ret_code"`
	RetMsg           string               `json:"ret_msg"`
	ExtCode          string               `json:"ext_code"`
	ExtInfo          string               `json:"ext_info"`
	Result           GetBalanceResultData `json:"result"`
	TimeNow          string               `json:"time_now"`
	RateLimitStatus  int                  `json:"rate_limit_status"`
	RateLimitResetMs int64                `json:"rate_limit_reset_ms"`
	RateLimit        int                  `json:"rate_limit"`
}

type GetBalanceResultData struct {
	BTC  Balance `json:"BTC"`
	ETH  Balance `json:"ETH"`
	EOS  Balance `json:"EOS"`
	XRP  Balance `json:"XRP"`
	USDT Balance `json:"USDT"`
}

type CreateOrderResult struct {
	RetCode         int    `json:"ret_code"`
	RetMsg          string `json:"ret_msg"`
	ExtCode         string `json:"ext_code"`
	Result          Order  `json:"result"`
	TimeNow         string `json:"time_now"`
	RateLimitStatus int    `json:"rate_limit_status"`
}

type CancelOrderResult struct {
	RetCode         int    `json:"ret_code"`
	RetMsg          string `json:"ret_msg"`
	ExtCode         string `json:"ext_code"`
	Result          Order  `json:"result"`
	TimeNow         string `json:"time_now"`
	RateLimitStatus int    `json:"rate_limit_status"`
}

type OrderListResultData struct {
	Data        []Order `json:"data"`
	CurrentPage int     `json:"current_page"`
	LastPage    int     `json:"last_page"`
}

type OrderListResult struct {
	RetCode         int                 `json:"ret_code"`
	RetMsg          string              `json:"ret_msg"`
	ExtCode         string              `json:"ext_code"`
	Result          OrderListResultData `json:"result"`
	TimeNow         string              `json:"time_now"`
	RateLimitStatus int                 `json:"rate_limit_status"`
}

// Order ...
type Order struct {
	OrderID       string     `json:"order_id"`
	UserID        int        `json:"user_id"`
	Symbol        string     `json:"symbol"`
	Side          string     `json:"side"`
	OrderType     string     `json:"order_type"`
	Price         float64    `json:"price"`
	Qty           float64    `json:"qty"`
	TimeInForce   string     `json:"time_in_force"`
	OrderStatus   string     `json:"order_status"`
	LastExecTime  string     `json:"last_exec_time"`
	LastExecPrice float64    `json:"last_exec_price"`
	LeavesQty     float64    `json:"leaves_qty"`
	CumExecQty    float64    `json:"cum_exec_qty"`
	CumExecValue  float64    `json:"cum_exec_value"`
	CumExecFee    float64    `json:"cum_exec_fee"`
	RejectReason  string     `json:"reject_reason"`
	OrderLinkID   string     `json:"order_link_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	ExtFields     *ExtFields `json:"ext_fields,omitempty"`
}

type OrderResponse struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	ExtInfo string `json:"ext_info"`
	Result  struct {
		OrderID        string    `json:"order_id"`
		UserID         int       `json:"user_id"`
		Symbol         string    `json:"symbol"`
		Side           string    `json:"side"`
		OrderType      string    `json:"order_type"`
		Price          float64   `json:"price"`
		Qty            int       `json:"qty"`
		TimeInForce    string    `json:"time_in_force"`
		OrderStatus    string    `json:"order_status"`
		LastExecPrice  int       `json:"last_exec_price"`
		CumExecQty     int       `json:"cum_exec_qty"`
		CumExecValue   int       `json:"cum_exec_value"`
		CumExecFee     int       `json:"cum_exec_fee"`
		ReduceOnly     bool      `json:"reduce_only"`
		CloseOnTrigger bool      `json:"close_on_trigger"`
		OrderLinkID    string    `json:"order_link_id"`
		CreatedTime    time.Time `json:"created_time"`
		UpdatedTime    time.Time `json:"updated_time"`
		TakeProfit     float64   `json:"take_profit"`
		StopLoss       float64   `json:"stop_loss"`
		TpTriggerBy    string    `json:"tp_trigger_by"`
		SlTriggerBy    string    `json:"sl_trigger_by"`
		PositionIdx    int       `json:"position_idx"`
	} `json:"result"`
	TimeNow          string `json:"time_now"`
	RateLimitStatus  int    `json:"rate_limit_status"`
	RateLimitResetMs int64  `json:"rate_limit_reset_ms"`
	RateLimit        int    `json:"rate_limit"`
}

type ExtFields struct {
	OpFrom      string `json:"op_from"`
	Remark      string `json:"remark"`
	OReqNum     int64  `json:"o_req_num"`
	XreqType    string `json:"xreq_type"`
	CrossStatus string `json:"cross_status,omitempty"`
}

type InExtFields struct {
	OpFrom      string `json:"op_from"`
	Remark      string `json:"remark"`
	OReqNum     int64  `json:"o_req_num"`
	XreqType    string `json:"xreq_type"`
	CrossStatus string `json:"cross_status,omitempty"`
}

func (e *ExtFields) MarshalJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e *ExtFields) UnmarshalJSON(b []byte) error {
	s := string(b)
	if strings.HasPrefix(s, "[") {
		return nil
	}
	o := InExtFields{}
	if err := json.Unmarshal(b, &o); err == nil {
		e.OpFrom = o.OpFrom
		e.Remark = o.Remark
		e.OReqNum = o.OReqNum
		e.XreqType = o.XreqType
		e.CrossStatus = o.CrossStatus
		return nil
	} else {
		return err
	}
}

type GetLeverageResult struct {
	RetCode         int                     `json:"ret_code"`
	RetMsg          string                  `json:"ret_msg"`
	ExtCode         string                  `json:"ext_code"`
	Result          map[string]LeverageItem `json:"result"`
	TimeNow         string                  `json:"time_now"`
	RateLimitStatus int                     `json:"rate_limit_status"`
}

type LeverageItem struct {
	Leverage float64 `json:"leverage"`
}

type PositionV1 struct {
	ID                  int                `json:"id"`
	UserID              int                `json:"user_id"`
	RiskID              int                `json:"risk_id"`
	Symbol              string             `json:"symbol"`
	Size                float64            `json:"size"`
	Side                string             `json:"side"`
	EntryPrice          float64            `json:"entry_price"`
	LiqPrice            float64            `json:"liq_price"`
	BustPrice           float64            `json:"bust_price"`
	TakeProfit          float64            `json:"take_profit"`
	StopLoss            float64            `json:"stop_loss"`
	TrailingStop        float64            `json:"trailing_stop"`
	PositionValue       float64            `json:"position_value"`
	Leverage            float64            `json:"leverage"`
	PositionStatus      string             `json:"position_status"`
	AutoAddMargin       float64            `json:"auto_add_margin"`
	OrderMargin         float64            `json:"order_margin"`
	PositionMargin      float64            `json:"position_margin"`
	OccClosingFee       float64            `json:"occ_closing_fee"`
	OccFundingFee       float64            `json:"occ_funding_fee"`
	ExtFields           *PositionExtFields `json:"ext_fields"`
	WalletBalance       float64            `json:"wallet_balance"`
	CumRealisedPnl      float64            `json:"cum_realised_pnl"`
	CumCommission       float64            `json:"cum_commission"`
	RealisedPnl         float64            `json:"realised_pnl"`
	DeleverageIndicator float64            `json:"deleverage_indicator"`
	OcCalcData          string             `json:"oc_calc_data"`
	CrossSeq            float64            `json:"cross_seq"`
	PositionSeq         float64            `json:"position_seq"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
	UnrealisedPnl       float64            `json:"unrealised_pnl"`
}

type Position struct {
	ID                  int                `json:"id"`
	UserID              int                `json:"user_id"`
	RiskID              int                `json:"risk_id"`
	Symbol              string             `json:"symbol"`
	Size                float64            `json:"size"`
	Side                string             `json:"side"`
	EntryPrice          float64            `json:"entry_price,string"`
	LiqPrice            float64            `json:"liq_price,string"`
	BustPrice           float64            `json:"bust_price,string"`
	TakeProfit          float64            `json:"take_profit,string"`
	StopLoss            float64            `json:"stop_loss,string"`
	TrailingStop        float64            `json:"trailing_stop,string"`
	PositionValue       float64            `json:"position_value,string"`
	Leverage            float64            `json:"leverage,string"`
	PositionStatus      string             `json:"position_status"`
	AutoAddMargin       float64            `json:"auto_add_margin"`
	OrderMargin         float64            `json:"order_margin,string"`
	PositionMargin      float64            `json:"position_margin,string"`
	OccClosingFee       float64            `json:"occ_closing_fee,string"`
	OccFundingFee       float64            `json:"occ_funding_fee,string"`
	ExtFields           *PositionExtFields `json:"ext_fields"`
	WalletBalance       float64            `json:"wallet_balance,string"`
	CumRealisedPnl      float64            `json:"cum_realised_pnl,string"`
	CumCommission       float64            `json:"cum_commission,string"`
	RealisedPnl         float64            `json:"realised_pnl,string"`
	DeleverageIndicator float64            `json:"deleverage_indicator"`
	OcCalcData          string             `json:"oc_calc_data"`
	CrossSeq            float64            `json:"cross_seq"`
	PositionSeq         float64            `json:"position_seq"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
	UnrealisedPnl       float64            `json:"unrealised_pnl"`
}

type PositionExtFields struct {
	Remark string `json:"_remark"`
}

type PositionListResult struct {
	BaseResult
	ExtInfo interface{}  `json:"ext_info"`
	Result  []PositionV1 `json:"result"`
}

type GetPositionResult struct {
	BaseResult
	ExtInfo interface{} `json:"ext_info"`
	Result  Position    `json:"result"`
}

type TradingHistoryResponse struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	Result  struct {
		PageToken string      `json:"page_token"`
		Data      []UserTrade `json:"data"`
	} `json:"result"`
	ExtCode string `json:"ext_code"`
	ExtInfo string `json:"ext_info"`
	TimeNow string `json:"time_now"`

	RateLimitStatus  int `json:"rate_limit_status"`
	RateLimitResetMs int `json:"rate_limit_reset_ms"`
	RateLimit        int `json:"rate_limit"`
}

type UserTrade struct {
	OrderID          string  `json:"order_id"`
	OrderLinkID      string  `json:"order_link_id"`
	Side             string  `json:"side"`
	Symbol           string  `json:"symbol"`
	ExecID           string  `json:"exec_id"`
	Price            float64 `json:"price"`
	OrderPrice       float64 `json:"order_price"`
	OrderQty         float64 `json:"order_qty"`
	OrderType        string  `json:"order_type"`
	FeeRate          float64 `json:"fee_rate"`
	ExecPrice        float64 `json:"exec_price"`
	ExecType         string  `json:"exec_type"`
	ExecQty          float64 `json:"exec_qty"`
	ExecFee          float64 `json:"exec_fee"`
	ExecValue        float64 `json:"exec_value"`
	LeavesQty        float64 `json:"leaves_qty"`
	ClosedSize       float64 `json:"closed_size"`
	LastLiquidityInd string  `json:"last_liquidity_ind"`
	TradeTime        int     `json:"trade_time"`
	TradeTimeMs      int     `json:"trade_time_ms"`
}
