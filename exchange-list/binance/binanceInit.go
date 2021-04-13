package binance

import (
	exchange "github.com/yin75620/crypto-berserker/exchange/struct"
	"github.com/yin75620/crypto-berserker/setting"
)

type BinanceInit struct {
	exchange.ExchangeInit
}

func NewBinanceInit() *BinanceInit {
	bi := BinanceInit{}
	bi.ApiKey = setting.BINANCE_KEY
	bi.ApiSecretKey = setting.BINANCE_SECRET_KEY
	bi.SectionName = "Binance"
	return &bi
}
