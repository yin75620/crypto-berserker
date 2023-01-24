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

type PositionListResponse struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	ExtInfo string `json:"ext_info"`
	Result  []struct {
		Data struct {
			UserID              int     `json:"user_id"`
			Symbol              string  `json:"symbol"`
			Side                string  `json:"side"`
			Size                int     `json:"size"`
			PositionValue       float64 `json:"position_value"`
			EntryPrice          float64 `json:"entry_price"`
			LiqPrice            float64 `json:"liq_price"`
			BustPrice           float64 `json:"bust_price"`
			Leverage            int     `json:"leverage"`
			AutoAddMargin       int     `json:"auto_add_margin"`
			IsIsolated          bool    `json:"is_isolated"`
			PositionMargin      float64 `json:"position_margin"`
			OccClosingFee       float64 `json:"occ_closing_fee"`
			RealisedPnl         float64 `json:"realised_pnl"`
			CumRealisedPnl      float64 `json:"cum_realised_pnl"`
			FreeQty             int     `json:"free_qty"`
			TpSlMode            string  `json:"tp_sl_mode"`
			UnrealisedPnl       float64 `json:"unrealised_pnl"`
			DeleverageIndicator int     `json:"deleverage_indicator"`
			RiskID              int     `json:"risk_id"`
			StopLoss            float64 `json:"stop"`
			TakeProfit          float64 `json:"take_profit"`
			TrailingStop        int     `json:"trailing_stop"`
			PositionIdx         int     `json:"position_idx"`
			Mode                string  `json:"mode"`
		} `json:"data"`
		IsValid bool `json:"is_valid"`
	} `json:"result"`
	TimeNow          string `json:"time_now"`
	RateLimitStatus  int    `json:"rate_limit_status"`
	RateLimitResetMs int    `json:"rate_limit_reset_ms"`
	RateLimit        int    `json:"rate_limit"`
}
