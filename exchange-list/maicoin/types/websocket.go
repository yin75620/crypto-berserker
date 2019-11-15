package types

import (
	exc "github.com/yin75620/crypto-berserker/exchange"
)

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

type OrderBookSocketResponse struct {
	exc.OrderBookSocketResponse
	Market string `json:"market,omitempty"`
	ID     uint64 `json:"id,omitempty"`
	Info   string `json:"info,omitempty"`
	//TimeStamp uint64  `json:"timestamp,string,omitempty"`
	//Action    string  `json:"action,omitempty"`
	//Market    string  `json:"market,omitempty"`
	//ID        uint64  `json:"id,omitempty"`
	//Side      string  `json:"side,omitempty"`
	//Volume    float64 `json:"volume,string,omitempty"`
	//Price     float64 `json:"price,string,omitempty"`
	OrderType string `json:"ord_type,omitempty"`

	//error
	Message string `json:"msg,omitempty"`

	//Success
	Channel string `json:"orderbook,omitempty"`
	//Market
}
