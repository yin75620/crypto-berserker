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
