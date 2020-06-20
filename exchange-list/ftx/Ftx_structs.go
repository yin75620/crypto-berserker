package ftx

type FillResponse struct {
	Fee           float64 `json:"fee"`
	FeeRate       float64 `json:"feeRate"`
	Future        string  `json:"future"`
	Id            int64   `json:"id"`
	Liquidity     string  `json:"liquidity"` //taker
	Market        string  `json:"market"`
	BaseCurrency  string  `json:"baseCurrency"`
	QuoteCurrency string  `json:"quoteCurrency"`
	OrderId       int64   `json:"orderId"`
	Price         float64 `json:"price"`
	Side          string  `json:"side"`
	Size          float64 `json:"size"`
	Time          string  `json:"time"`
	OrderType     string  `json:"type"` // order
}

type AccountResponse struct {
	Result  Account `json:"result"`
	Success bool    `json:"success"`
}

type Account struct {
	BackstopProvider             bool    `json:"backstopProvider"`
	Collateral                   float64 `json:"collateral"`
	FreeCollateral               float64 `json:"freeCollateral"`
	InitialMarginRequirement     float64 `json:"initialMarginRequirement"`
	Liquidating                  bool    `json:"liquidating"`
	MaintenanceMarginRequirement float64 `json:"maintenanceMarginRequirement"`
	MakerFee                     float64 `json:"makerFee"`
	MarginFraction               float64 `json:"marginFraction"`
	OpenMarginFraction           float64 `json:"openMarginFraction"`
	TakerFee                     float64 `json:"takerFee"`
	TotalAccountValue            float64 `json:"totalAccountValue"`
	TotalPositionSize            float64 `json:"totalPositionSize"`
	Username                     string  `json:"username"`
	Leverage                     float64 `json:"leverage"`
	Positions                    []struct {
		Cost                         float64 `json:"cost"`
		EntryPrice                   float64 `json:"entryPrice"`
		Future                       string  `json:"future"`
		InitialMarginRequirement     float64 `json:"initialMarginRequirement"`
		LongOrderSize                float64 `json:"longOrderSize"`
		MaintenanceMarginRequirement float64 `json:"maintenanceMarginRequirement"`
		NetSize                      float64 `json:"netSize"`
		OpenSize                     float64 `json:"openSize"`
		RealizedPnl                  float64 `json:"realizedPnl"`
		ShortOrderSize               float64 `json:"shortOrderSize"`
		Side                         string  `json:"side"`
		Size                         float64 `json:"size"`
		UnrealizedPnl                float64 `json:"unrealizedPnl"`
	}
}

func (ac *Account) GetTotalUnrealizedPnl() float64 {
	sum := 0.0
	for _, v := range ac.Positions {
		sum += v.UnrealizedPnl
	}
	return sum
}
