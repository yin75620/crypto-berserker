package common

import (
	"log"
	"net/http"

	"github.com/yin75620/crypto-berserker/exchange-list/binancef"
	"github.com/yin75620/crypto-berserker/exchange-list/bybilinear"
	"github.com/yin75620/crypto-berserker/exchange-list/bybit"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx"

	exc "github.com/yin75620/crypto-berserker/exchange"
)

const (
	FTX        = "FTX"
	BINANCE    = "Binance"
	BINANCEF   = "Binancef"
	BYBIT      = "Bybit"
	BYBILINEAR = "Bybilinear"
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
	default:
		log.Println("error, not define exchangeName:", exchangeName)
		break
	}

	return res
}
