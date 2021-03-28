package exchange

type FirstBuyExchange interface {
	DeleteAllOrders() (string, error)
	PostOrder(order ExchangeOrder) (string, error)
}
