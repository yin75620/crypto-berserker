package binancef

import (
	"github.com/go-ini/ini"
	"github.com/yin75620/crypto-berserker/setting"
)

type BinancefInit struct {
	Key       string
	SecretKey string
}

func NewBinancefInit() *BinancefInit {
	return &BinancefInit{
		Key:       setting.BINANCE_KEY,
		SecretKey: setting.BINANCE_SECRET_KEY,
	}
}

func (bfi *BinancefInit) IniSetting(filename string) error {
	cfg, err := ini.Load(filename)
	if err != nil {
		//log.Println("BinancefInit Fail to read file: ", err)
		return err
	}

	const section = "Binancef"

	bfi.Key = cfg.Section(section).Key("Key").MustString(setting.BINANCE_KEY)
	bfi.SecretKey = cfg.Section(section).Key("SecretKey").MustString(setting.BINANCE_SECRET_KEY)
	return nil
}
