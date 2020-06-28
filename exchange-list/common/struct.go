package common

import (
	"log"
	"net/http"

	"github.com/yin75620/crypto-berserker/exchange-list/binancef"
	"github.com/yin75620/crypto-berserker/exchange-list/bybilinear"
	"github.com/yin75620/crypto-berserker/exchange-list/bybit"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx"
	"github.com/yin75620/crypto-berserker/exchange-list/okex"

	exc "github.com/yin75620/crypto-berserker/exchange"
)

const (
	FTX        = "FTX"
	BINANCE    = "Binance"
	BINANCEF   = "Binancef"
	BYBIT      = "Bybit"
	BYBILINEAR = "Bybilinear"
	OKEX       = "Okex"
)

func GetExchange(exchangeName string) exc.Exchange {
	const initFileName = "Main.ini"
	var res exc.Exchange
	switch exchangeName {
	case FTX:
		fi := ftx.NewFtxInit()
		fe := ftx.NewFtx(http.DefaultClient, *fi)
		fe.SetInitByIni(initFileName)
		res = fe
		break
	case BYBIT:
		res = bybit.NewBybit(http.DefaultClient)
		break
	case BYBILINEAR:
		res = bybilinear.NewBybilinear(http.DefaultClient)
		break
	case BINANCEF:
		res = binancef.NewBinancef(http.DefaultClient)
		break
	case OKEX:
		res = okex.NewOkex(http.DefaultClient)
		break
	default:
		log.Println("error, not define exchangeName:", exchangeName)
		break
	}

	return res
}
