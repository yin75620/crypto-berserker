package exchange

type Wallet struct {
	Balances []Balance
}

type Balance struct {
	Coin         string  //coinName
	Free         float64 //can use amount
	FreeUsdValue float64 //can use amount
	Total        float64 //total amount
	UsdValue     float64
}

func NewWallet() *Wallet {
	return &Wallet{}
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

func (w *Wallet) GetAllBalanceProfit(after Wallet) []Balance {
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

func (w *Wallet) GetAllBalanceUSDValue() float64 {
	result := 0.0
	for _, value := range w.Balances {
		result += value.UsdValue
	}
	return result
}

func (w *Wallet) GetAllBalanceFreeUSDValue() float64 {
	result := 0.0
	for _, value := range w.Balances {
		result += value.FreeUsdValue
	}
	return result
}

func (w *Wallet) CalculerFreeUsdValue() {
	for i, value := range w.Balances {
		if value.UsdValue == 0 {
			continue
		}
		value.FreeUsdValue = value.Free * value.Total / value.UsdValue
		w.Balances[i] = value
	}

}
