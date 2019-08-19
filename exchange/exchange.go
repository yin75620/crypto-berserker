package exchange

type Exchange interface {
	GetAskBidPair(marketName string, depth int) (PricePair, PricePair)
	GetAccountInfo() []byte
	PostOrder(order ExchangeOrder) string
}

type PriceType int

const (
	Ask PriceType = iota
	Bid
)

type PricePair struct {
	Price  float64
	Volume float64
}

type EOrderType string

const (
	LIMIT  EOrderType = "limit"
	MARKET EOrderType = "market"
)

type ExchangeOrder struct {
	Market    string     `json:"market"`
	Side      string     `json:"side"`
	Price     float64    `json:"price"`
	Size      float64    `json:"size"`
	OrderType EOrderType `json:"order_type"`
	//ReduceOnly bool       `json:"reduceOnly"`
}

func GetAskBidPair() {

}
