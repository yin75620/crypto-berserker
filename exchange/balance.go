package exchange

type Wallet struct {
	Balances []Balance
}

type Balance struct {
	Coin     string  //coinName
	Free     float64 //can use amount
	Total    float64 //total amount
	UsdValue float64
}

func (w *Wallet) IsAllBalanceReduce(after Wallet) bool {

	zeroCount := 0
	for _, value := range w.Balances {
		for _, afterValue := range after.Balances {
			if value.Coin == afterValue.Coin {
				changedValue := afterValue.Total - value.Total
				if changedValue > 0 {
					return false
				} else if changedValue == 0 {
					zeroCount++
				}
			}
		}
	}
	if zeroCount == len(w.Balances) {
		return false
	}
	return true
}

func (w *Wallet) GetALlBalanceProfit(after Wallet) []Balance {
	result := []Balance{}
	for _, value := range w.Balances {
		for _, afterValue := range after.Balances {
			if value.Coin == afterValue.Coin {
				changedValue := afterValue.Total - value.Total
				result = append(result, Balance{Coin: value.Coin, Total: changedValue})
			}
		}
	}
	return result
}
