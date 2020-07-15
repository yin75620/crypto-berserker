package CrossExchange

import (
	"fmt"

	"github.com/go-ini/ini"
)

type CrossExchangeInit struct {
	DelayMilliSecond        int64
	OverPrice               float64 //= 0.02 // 交易時，要溢價多少。 Ex:目前價位 9000 => 會用9180買進
	MinSellProfit           float64 // = -0.0007
	MinSumProfit            float64 //= 0.0001
	MaxHoldVolume           float64 // 1000.0
	MaxHoldeExchangePercent float64 // 1 = 100%
	MaxHoldBuffer           float64 //0.01
	MinCreateProfit         float64 // 0.001
	MinVolume               float64 //1鎂
	StopLosePercent         float64 // -0.02 多少%要停損
	IsEnableDBLog           bool    // 是否要啟用db記錄 log
}

func NewCrossExchangeInit() *CrossExchangeInit {
	return &CrossExchangeInit{
		DelayMilliSecond:        500,
		OverPrice:               0.02,
		MinSellProfit:           -0.05,
		MinSumProfit:            0.0001,
		MaxHoldVolume:           1000.0,
		MaxHoldeExchangePercent: 1,
		MaxHoldBuffer:           0.01,
		MinCreateProfit:         0.001,
		MinVolume:               1,
		StopLosePercent:         -0.02,
		IsEnableDBLog:           true,
	}
}

func (cei *CrossExchangeInit) IniSetting(filename string) error {
	cfg, err := ini.Load(filename)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		return err
	}

	const section = "CrossExchange"

	cei.DelayMilliSecond = cfg.Section(section).Key("DelayMilliSecond").MustInt64()
	cei.OverPrice = cfg.Section(section).Key("OverPrice").MustFloat64()
	cei.MinSellProfit = cfg.Section(section).Key("MinSellProfit").MustFloat64()
	cei.MinSumProfit = cfg.Section(section).Key("MinSumProfit").MustFloat64()
	cei.MaxHoldVolume = cfg.Section(section).Key("MaxHoldVolume").MustFloat64()
	cei.MaxHoldeExchangePercent = cfg.Section(section).Key("MaxHoldeExchangePercent").MustFloat64()
	cei.MaxHoldBuffer = cfg.Section(section).Key("MaxHoldBuffer").MustFloat64()
	cei.MinCreateProfit = cfg.Section(section).Key("MinCreateProfit").MustFloat64()
	cei.MinVolume = cfg.Section(section).Key("MinVolume").MustFloat64()
	cei.StopLosePercent = cfg.Section(section).Key("StopLosePercent").MustFloat64()
	cei.IsEnableDBLog = cfg.Section(section).Key("EnableDBLog").MustBool()

	return nil
}
