package binancef

import (
	"fmt"
	"net/http"
	"testing"

	exc "github.com/yin75620/crypto-berserker/exchange"
)

func TestAccount(t *testing.T) {
	ce := NewBinancef(http.DefaultClient)
	res := ce.GetAccountInfo()

	fmt.Println(string(res))
}

func TestOrder(t *testing.T) {
	ce := NewBinancef(http.DefaultClient)

	var myOrder exc.ExchangeOrder = exc.ExchangeOrder{
		CoinPair:  exc.CoinPair{"FTT", "USDT"},
		Side:      exc.Buy,
		Price:     2.3,
		Size:      5,
		OrderType: exc.LIMIT,
	}
	res, _ := ce.PostOrder(myOrder)

	fmt.Println(string(res))
}

/*
func TestSign(t *testing.T) {
	key := "vmPUZE6mv9SD5VNHk4HlWFsOr6aKE2zvsw0MuIgwCIPy6utIco14y7Ju91duEh8A"
	secKey := "NhqPtmdSJYdKjVHjA7PZj4Mge3R5YNiP1e3UZjInClVN65XAbvqqM6A7H5fATj0j"
	ce := NewBinance1(key, secKey)
	body := exc.JArray{
		"symbol":      "LTCBTC",
		"side":        "BUY",
		"type":        "LIMIT",
		"timeInForce": "GTC",
		"quantity":    1,
		"price":       0.1,
		"recvWindow":  5000,
		"timestamp":   1499827319559,
	}
	ce.doRequest("GET", "Account", body)

	query := exc.JArray{
		"symbol":      "LTCBTC",
		"side":        "BUY",
		"type":        "LIMIT",
		"timeInForce": "GTC",
	}

	str := fmt.Sprintf("%s%s", query.ToValues().Encode(), body.ToValues().Encode())
	fmt.Println(str)

	//res, err := exc.GetParamHmacSHA256HexSign(secKey, "symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=5000&timestamp=1499827319559")
	res, err := exc.GetParamHmacSHA256HexSign(secKey, str)
	if err != nil {
		fmt.Println(err)
	}
	//symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=5000&timestamp=1499827319559

	//res1 := "c8db56825ae71d6d79447849e617115f4a920fa2acdcab2b053c4b2838bd6b71"
	fmt.Println(res)
}*/

func TestPricePair(t *testing.T) {
	ce := NewBinancef(http.DefaultClient)
	res1, res2 := ce.GetAskBidPair(exc.CoinPair{"BTC", "USDT"}, 1)
	fmt.Println(fmt.Sprintf("res1:%g", res1.Price))
	fmt.Println(fmt.Sprintf("res1:%g", res1.Volume))
	fmt.Println(fmt.Sprintf("res2:%g", res2.Price))
	fmt.Println(fmt.Sprintf("res2:%g", res2.Volume))

	//fmt.Println(string(res))
}
