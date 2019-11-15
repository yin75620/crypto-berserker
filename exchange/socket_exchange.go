package exchange

type SocketExchange interface {
}

type OrderBookSocketResponse struct {
	//Info      string  `json:"info,omitempty"`
	TimeStamp uint64 `json:"timestamp,string,omitempty"`
	Action    string `json:"action,omitempty"`
	//Market    string  `json:"market,omitempty"`
	//ID        uint64  `json:"id,omitempty"`
	Side   string  `json:"side,omitempty"`
	Volume float64 `json:"volume,string,omitempty"`
	Price  float64 `json:"price,string,omitempty"`
	//OrderType string  `json:"ord_type,omitempty"`
	CoinPair CoinPair
}

func (obsr *OrderBookSocketResponse) IsAdd() bool {
	return obsr.Action == Add
}

func (obsr *OrderBookSocketResponse) IsRemove() bool {
	return obsr.Action == Remove
}

func (obsr *OrderBookSocketResponse) IsUpdate() bool {
	return obsr.Action == Update
}

func (obsr *OrderBookSocketResponse) IsAsk() bool {
	return obsr.Side == "buy"
}

func (obsr *OrderBookSocketResponse) IsBid() bool {
	return obsr.Side == "sell"
}

func (obsr *OrderBookSocketResponse) PricePair() PricePair {
	pp := PricePair{}
	pp.Price = obsr.Price
	pp.Volume = obsr.Volume
	return pp
}

type ActionType string

const (
	Add    = "add"
	Remove = "remove"
	Update = "update"
)
