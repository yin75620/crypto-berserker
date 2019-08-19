package bifinex

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	exc "github.com/yin75620/crypto-berserker/exchange"
)

type Bitfinex struct {
	client *http.Client
}

var (
	apiURL    = "https://bitmax.io/api/v1/products"
	apiPrefix = "/api/"
)

func (bf *Bitfinex) GetAskBidPair(marketName string, depth int) (exc.PricePair, exc.PricePair) {
	resb := ftx.GetOrderBookResponse(marketName, depth)
	askPair, _ := resb.Result.GetPair(1, Ask)
	bidPair, _ := resb.Result.GetPair(1, Bid)
	return askPair, bidPair
}
func (bf *Bitfinex) GetAccountInfo() []byte {
	return ftx.doGet("account", "")
}

//下訂單
func (bf *Bitfinex) PostOrder(order exc.ExchangeOrder) string {

	fo := FtxOrder{}
	fo.setBy(order)

	request, err := json.Marshal(fo)
	if err != nil {
		log.Fatal(err)
	}
	body := string(request)
	log.Println(fmt.Sprintf("body:%s", body))
	response := ftx.doPost("orders", body)
	log.Println(fmt.Sprintf("%s", response))
	return string(response)
}
