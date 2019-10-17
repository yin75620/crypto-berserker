package types

type SubscriptionRequest struct {
	Cmd     string      `json:"cmd,omitempty"`     //subscribe,unsubscribe
	Channel string      `json:"channel,omitempty"` //ticker,orderbook,trade
	Params  interface{} `json:"params,omitempty"`
}

type SubscriptionResponse struct {
	Info    string      `json:"info,omitempty"`
	Channel string      `json:"channel,omitempty"`
	Params  interface{} `json:"params,omitempty"`
}
