package okex

import (
	"time"

	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
	"github.com/yin75620/crypto-berserker/exchange/tool"
)

type SubscriptionRequest struct {
	Op   string   `json:"op,omitempty"`
	Args []string `json:"args,omitempty"`
}

// OrderBookData is resp data from orderbook endpoint
type OrderBookData struct {
	Table int              `json:"table"`
	Data  []InstrumentData `json:"data"`
}

type InstrumentData struct {
	Bids         []interface{} `json:"bids,[]string"`
	Asks         []interface{} `json:"asks,[]string"`
	InstrumentId string        `json:"instrument_id"`
	Timestamp    time.Time     `json:"timestamp"`
}

func (obd *OrderBookData) ToOrderBookDetail() ob.OrderBookerResponseDetail {
	res := ob.OrderBookerResponseDetail{}
	data := obd.Data[0]
	res.Time = float64(data.Timestamp.Unix())
	res.Action = ob.Partial
	res.Checksum = data.Timestamp.Unix()
	res.Asks = tool.TransToFloatTwoArray(data.Asks)
	res.Bids = tool.TransToFloatTwoArray(data.Bids)

	return res
}
