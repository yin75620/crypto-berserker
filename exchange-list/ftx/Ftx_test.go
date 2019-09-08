package ftx

import (
	"net/http"
	"testing"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/setting"
)

var mftx = NewFtx(http.DefaultClient,
	FtxInit{
		setting.FTX_KEY,
		setting.FTX_API_SECRET_KEY,
		"Saber"})

func TestInfo(t *testing.T) {

	mftx.GetAccountInfo()

	//ftx.GetMarkets()
}

func TestOrder(t *testing.T) {

	var myOrder exc.ExchangeOrder = exc.ExchangeOrder{
		Market:    "FTT/USD",
		Side:      exc.Buy,
		Price:     1.1,
		Size:      1.1255,
		OrderType: exc.LIMIT,
	}
	mftx.PostOrder(myOrder)

	//ftx.GetMarkets()
}
