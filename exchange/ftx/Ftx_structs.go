package ftx

type FillResponse struct {
	Fee           float64 `json:"fee"`
	FeeRate       float64 `json:"feeRate"`
	Future        string  `json:"future"`
	Id            int64   `json:"id"`
	Liquidity     string  `json:"liquidity"` //taker
	Market        string  `json:"market"`
	BaseCurrency  string  `json:"baseCurrency"`
	QuoteCurrency string  `json:"quoteCurrency"`
	OrderId       int64   `json:"orderId"`
	Price         float64 `json:"price"`
	Side          string  `json:"side"`
	Size          float64 `json:"size"`
	Time          string  `json:"time"`
	OrderType     string  `json:"type"` // order
}
