package ftxotc

import "time"

//QuoteItem 回傳的內容
type QuoteItem struct {
	ID                         int64     `json:"id"`                         //"id": 12345,
	BaseCurrency               string    `json:"baseCurrency"`               //"baseCurrency": "BTC",
	QuoteCurrency              string    `json:"quoteCurrency"`              //	"quoteCurrency": "USDT",
	Side                       string    `json:"side"`                       //"side": "buy",
	BaseCurrencySize           float64   `json:"baseCurrencySize"`           //"baseCurrencySize": 2.1,
	QuoteCurrencySize          float64   `json:"quoteCurrencySize"`          //		"quoteCurrencySize": 7560.42,
	Price                      float64   `json:"price"`                      //"price": 3600.2,
	RquestedAt                 time.Time `json:"requestedAt"`                //"requestedAt": "2018-11-30T17:01:22.137536+00:00",
	QuotedAt                   time.Time `json:"quotedAt"`                   //"quotedAt": "2018-11-30T17:01:23.137536+00:00",
	Expiry                     time.Time `json:"expiry"`                     //"expiry": "2018-11-30T17:01:35.137536+00:00",
	Filled                     bool      `json:"filled"`                     //"filled": false,
	OrderID                    int64     `json:"orderId"`                    //"orderId": 123456,
	CounterpartyAutoSettles    bool      `json:"counterpartyAutoSettles"`    //"counterpartyAutoSettles": false,
	SettledImmediately         bool      `json:"settledImmediately"`         //"settledImmediately": false,
	SettlementTime             time.Time `json:"settlementTime"`             //"settlementTime": "2018-12-02T17:01:35.137536+00:00",
	DeferCostRate              float64   `json:"deferCostRate"`              //"deferCostRate": 0.0004,
	DeferProceedsRate          float64   `json:"deferProceedsRate"`          //"deferProceedsRate": 0.0004,
	SettlementPriority         int64     `json:"settlementPriority"`         //"settlementPriority": 2,
	CostCurrency               string    `json:"costCurrency"`               //"costCurrency": "USDT",
	Cost                       float64   `json:"cost"`                       //"cost": 7560.42,
	ProceedsCurrency           string    `json:"proceedsCurrency"`           //"proceedsCurrency": "BTC",
	Proceeds                   float64   `json:"proceeds"`                   //"proceeds": 2.1,
	TotalDeferCostPaid         float64   `json:"totalDeferCostPaid"`         //"totalDeferCostPaid": 0,
	TotalDeferProceedsPaid     float64   `json:"totalDeferProceedsPaid"`     //"totalDeferProceedsPaid": 0,
	UnsettledCost              float64   `json:"unsettledCost"`              //"unsettledCost": 2000.32,
	UnsettledProceeds          float64   `json:"unsettledProceeds"`          //"unsettledProceeds": 9193,
	UserFullySettledAt         time.Time `json:"userFullySettledAt"`         //"userFullySettledAt": "2018-12-03T17:01:35.137536+00:00",
	CounterpartyFullySettledAt time.Time `json:"counterpartyFullySettledAt"` //"counterpartyFullySettledAt": "2018-12-04T17:01:35.137536+00:00"

}
