package types

type SubscriptionRequest struct {
	//orderbook for orderbook market data
	//trades for trade market data
	//ticker for best bid and offer market data
	Channel string `json:"channel,omitempty"` //ticker,orderbook,trade
	Market  string `json:"market,omitempty"`  //ticker,orderbook,trade

	//subscribe to subscribe to a channel
	//unsubscribe to unsubscribe from a channel
	Op string `json:"op,omitempty"` //subscribe,unsubscribe
}

type SubscriptionResponse struct {
	Info    string      `json:"info,omitempty"`
	Channel string      `json:"channel,omitempty"`
	Params  interface{} `json:"params,omitempty"`
}

//{"channel": "orderbook", "market": "BTC/USD", "type": "update", "data": {"time": 1574775055.2251372, "checksum": 3250722053, "bids": [], "asks": [[7089.5, 0.0], [7124.5, 47.3641]], "action": "update"}}

type OrderBookSocketResponseData struct {
	Time     float64     `json:"time,omitempty"`
	Checksum int64       `json:"checksum,omitempty"`
	Asks     [][]float64 `json:"asks,omitempty"`
	Bids     [][]float64 `json:"bids,omitempty"`
	Action   string      `json:"action,omitempty"`
}

type OrderBookSocketResponse struct {
	Channel    string                      `json:"orderbook,omitempty"`
	Market     string                      `json:"market,omitempty"`
	ActionType ActionType                  `json:"type,omitempty"`
	Data       OrderBookSocketResponseData `json:"data,omitempty"`
	Error      error                       // local use
}

func (obsr *OrderBookSocketResponse) IsUpdate() bool {
	return obsr.ActionType == Update
}

type ActionType string

const (
	Add    = "add"
	Remove = "remove"
	Update = "update"
)
