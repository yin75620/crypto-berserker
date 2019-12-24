package order_booker

type IWebSocket interface {
	SubScribeOrderBook(market string) chan OrderBookerResponseDetail
}

/*
type IOrderBookerResponse interface {
	GetOrderBookerResponseDetail() OrderBookerResponseDetail
}*/

type OrderBookerResponseDetail struct {
	Time     float64     `json:"time,omitempty"`
	Checksum int64       `json:"checksum,omitempty"`
	Asks     [][]float64 `json:"asks,omitempty"`
	Bids     [][]float64 `json:"bids,omitempty"`
	Action   ActionType  `json:"action,omitempty"`
	Market   string      `json:"market,omitempty"`
}

type ActionType string

const (
	Update  ActionType = "update"  //update ,if delete use 0 at amount
	Partial            = "partial" //first
)

func (obrd *OrderBookerResponseDetail) IsUpdate() bool {
	return Update == obrd.Action
}
