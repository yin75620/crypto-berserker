package bybilinear

type BybilinearCancelOrder struct {
	Symbol string `json:"symbol"` //Reuqired
	//OrderID     string `json:"order_id"`
	//OrderLinkID string `json:"order_link_id"`
}

type InstrumentsInfo struct {
	Results struct {
		Category string `json:"category"`
		List     []struct {
			LeverageFilter struct {
				MinLeverage  float64 `json:"minLeverage,string"`
				MaxLeverage  float64 `json:"maxLeverage,string"`
				LeverageStep float64 `json:"leverageStep,string"`
			} `json:"leverageFilter"`
			PriceFilter struct {
				MinPrice float64 `json:"minPrice,string"`
				MaxPrice float64 `json:"maxPrice,string"`
				TickSize float64 `json:"tickSize,string"`
			} `json:"priceFilter"`
			LotSizeFilter struct {
				MaxTradingQty float64 `json:"maxTradingQty,string"`
				QtyStep       float64 `json:"qtyStep,string"`
			} `json:"lotSizeFilter"`
			Symbol string `json:"symbol"`
		} `json:"list"`
	} `json:"result"`
}

// Response represents the top-level structure of the JSON response
type PositionListResponse struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  Result `json:"result"`
	Time    int64  `json:"time"`
}

type Result struct {
	List           []PositionX `json:"list"`
	NextPageCursor string      `json:"nextPageCursor"`
	Category       string      `json:"category"`
}

// Position represents each item within the "list" array of the JSON response
type PositionX struct {
	PositionIdx            int     `json:"positionIdx"`
	RiskId                 int     `json:"riskId"`
	RiskLimitValue         string  `json:"riskLimitValue"`
	Symbol                 string  `json:"symbol"`
	Side                   string  `json:"side"`
	Size                   string  `json:"size"`
	AvgPrice               float64 `json:"avgPrice,string"`
	PositionValue          float64 `json:"positionValue,string"`
	TradeMode              int     `json:"tradeMode"`
	PositionStatus         string  `json:"positionStatus"`
	AutoAddMargin          int     `json:"autoAddMargin"`
	AdlRankIndicator       int     `json:"adlRankIndicator"`
	Leverage               int     `json:"leverage,string"`
	PositionBalance        float64 `json:"positionBalance,string"`
	MarkPrice              float64 `json:"markPrice,string"`
	LiqPrice               string  `json:"liqPrice,omitempty"`
	BustPrice              float64 `json:"bustPrice,string"`
	PositionMM             float64 `json:"positionMM,string"`
	PositionIM             float64 `json:"positionIM,string"`
	TpslMode               string  `json:"tpslMode"`
	TakeProfit             float64 `json:"takeProfit,string"`
	StopLoss               float64 `json:"stopLoss,string"`
	TrailingStop           float64 `json:"trailingStop,string"`
	UnrealisedPnl          float64 `json:"unrealisedPnl,string"`
	CurRealisedPnl         float64 `json:"curRealisedPnl,string"`
	CumRealisedPnl         float64 `json:"cumRealisedPnl,string"`
	Seq                    int64   `json:"seq,string"`
	IsReduceOnly           bool    `json:"isReduceOnly"`
	MmrSysUpdateTime       string  `json:"mmrSysUpdateTime,omitempty"`
	LeverageSysUpdatedTime string  `json:"leverageSysUpdatedTime,omitempty"`
	SessionAvgPrice        string  `json:"sessionAvgPrice,omitempty"`
	CreatedTime            string  `json:"createdTime"`
	UpdatedTime            string  `json:"updatedTime"`
}

type LeverageInfo struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	Result  []struct {
		ID             int      `json:"id"`
		Symbol         string   `json:"symbol"`
		Limit          int      `json:"limit"`
		MaintainMargin float64  `json:"maintain_margin"`
		StartingMargin float64  `json:"starting_margin"`
		Section        []string `json:"section"`
		IsLowestRisk   int      `json:"is_lowest_risk"`
		CreatedAt      string   `json:"created_at"`
		UpdatedAt      string   `json:"updated_at"`
		MaxLeverage    int      `json:"max_leverage"`
	} `json:"result"`
	ExtCode string `json:"ext_code"`
	ExtInfo string `json:"ext_info"`
	TimeNow string `json:"time_now"`
}
