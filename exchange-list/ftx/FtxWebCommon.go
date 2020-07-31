package ftx

import (
	"github.com/gorilla/websocket"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx/types"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
	"github.com/yin75620/crypto-berserker/exchange/tool"
)

const WEBSOCKET_URL = "wss://ftx.com/ws/"

func toDetail(res types.OrderBookSocketResponse) ob.OrderBookerResponseDetail {
	obr := ob.OrderBookerResponseDetail{}
	if res.IsUpdate() {
		obr.Action = ob.Update
	} else {
		obr.Action = ob.Partial
	}

	obr.Asks = res.Data.Asks
	obr.Bids = res.Data.Bids
	obr.Checksum = res.Data.Checksum
	obr.Error = res.Error
	obr.Market = res.Market
	obr.Time = res.Data.Time

	return obr
}

func createConn() *websocket.Conn {
	return tool.CreateConn(WEBSOCKET_URL)
}

func sendSubcribe(c *websocket.Conn, channelName string, marketName string) {
	req := types.SubscriptionRequest{
		Op:      "subscribe",
		Channel: channelName,
		Market:  marketName,
	}

	tool.Send(c, req)
}
