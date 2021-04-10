package binance

import (
	"github.com/go-ini/ini"
	"github.com/yin75620/crypto-berserker/setting"
)

type BinanceInit struct {
	ApiKey       string
	ApiSecretKey string
}

func (bi *BinanceInit) IniSettingByFile(cfg *ini.File) {
	const section = "BINANCE"

	bi.ApiKey = setting.BINANCE_KEY
	bi.ApiSecretKey = setting.BINANCE_SECRET_KEY

	apiKey := cfg.Section(section).Key("ApiKey").String()
	if apiKey != "" {
		bi.ApiKey = apiKey
	}

	apiSecretKey := cfg.Section(section).Key("ApiSecretKey").String()
	if apiSecretKey != "" {
		bi.ApiSecretKey = apiSecretKey
	}
}
