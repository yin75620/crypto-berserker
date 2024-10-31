package exchange

type Collector interface {
	GetKlines(symbol string, interval string, startTime, endTime int64, limit int) ([]CandleStick, error)
}
type CandleStick struct {
	TimeInterval             TimeInterval `gorm:"column:time_interval"`
	OpenTime                 float64      `gorm:"column:open_time"`
	Open                     float64      `gorm:"column:open"`
	High                     float64      `gorm:"column:high"`
	Low                      float64      `gorm:"column:low"`
	Close                    float64      `gorm:"column:close"`
	Volume                   float64      `gorm:"column:volume"`
	CloseTime                float64      `gorm:"column:close_time"`
	QuoteAssetVolume         float64      `gorm:"column:quote_asset_volume"`
	TradeCount               float64      `gorm:"column:trade_count"`
	TakerBuyAssetVolume      float64      `gorm:"column:taker_buy_asset_volume"`
	TakerBuyQuoteAssetVolume float64      `gorm:"column:taker_buy_quote_asset_volume"`
}

// TimeInterval represents interval enum.
type TimeInterval string

// Vars related to time intervals
var (
	TimeIntervalMinute         = TimeInterval("1m")
	TimeIntervalThreeMinutes   = TimeInterval("3m")
	TimeIntervalFiveMinutes    = TimeInterval("5m")
	TimeIntervalFifteenMinutes = TimeInterval("15m")
	TimeIntervalThirtyMinutes  = TimeInterval("30m")
	TimeIntervalHour           = TimeInterval("1h")
	TimeIntervalTwoHours       = TimeInterval("2h")
	TimeIntervalFourHours      = TimeInterval("4h")
	TimeIntervalSixHours       = TimeInterval("6h")
	TimeIntervalEightHours     = TimeInterval("8h")
	TimeIntervalTwelveHours    = TimeInterval("12h")
	TimeIntervalDay            = TimeInterval("1d")
	TimeIntervalThreeDays      = TimeInterval("3d")
	TimeIntervalWeek           = TimeInterval("1w")
	TimeIntervalMonth          = TimeInterval("1M")
)
