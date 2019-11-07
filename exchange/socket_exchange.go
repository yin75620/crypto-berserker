package exchange

type SocketExchange interface {
}

type OrderBookSocketResponse struct {
	//Info      string  `json:"info,omitempty"`
	TimeStamp uint64  `json:"timestamp,string,omitempty"`
	Action    string  `json:"action,omitempty"`
	Market    string  `json:"market,omitempty"`
	ID        uint64  `json:"id,omitempty"`
	Side      string  `json:"side,omitempty"`
	Volume    float64 `json:"volume,string,omitempty"`
	Price     float64 `json:"price,string,omitempty"`
	//OrderType string  `json:"ord_type,omitempty"`
}

type ActionType string

const (
	Add    = "add"
	Remove = "remove"
	Update = "update"
)
