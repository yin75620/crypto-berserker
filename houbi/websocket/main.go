package main

import (
	"fmt"
	"net/http"

	goex "github.com/nntaoli-project/GoEx"
	"github.com/nntaoli-project/GoEx/huobi"

	//"github.com/nntaoli-project/GoEx/okcoin"
	"log"
)

func main() {

	//ws := okcoin.NewOKExFutureWs() //ok期货
	ws := huobi.NewHbdmWs() //huobi期货
	//设置回调函数
	ws.SetCallbacks(func(ticker *goex.FutureTicker) {
		log.Println(ticker)
	}, func(depth *goex.Depth) {
		log.Println(depth)
	}, func(trade *goex.Trade, contract string) {
		log.Println(contract, trade)
	})
	//订阅行情
	ws.SubscribeTrade(goex.BTC_USDT, goex.NEXT_WEEK_CONTRACT)
	ws.SubscribeDepth(goex.BTC_USDT, goex.QUARTER_CONTRACT, 5)
	ws.SubscribeTicker(goex.BTC_USDT, goex.QUARTER_CONTRACT)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
