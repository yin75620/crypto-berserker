package main

import (
	"log"
	"time"

	goex "github.com/nntaoli-project/GoEx"
	"github.com/nntaoli-project/GoEx/builder"
)

func main() {
	apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second)
	//apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second).HttpProxy("socks5://127.0.0.1:1080")

	//build spot api
	//api := apiBuilder.APIKey("").APISecretkey("").ClientID("123").Build(goex.BITSTAMP)
	api := apiBuilder.APIKey("").APISecretkey("").Build(goex.HUOBI_PRO)
	log.Println(api.GetExchangeName())
	log.Println(api.GetTicker(goex.BTC_USD))
	log.Println(api.GetDepth(2, goex.BTC_USD))
	//log.Println(api.GetAccount())
	//log.Println(api.GetUnfinishOrders(goex.BTC_USD))

	//build future api
	futureApi := apiBuilder.APIKey("").APISecretkey("").BuildFuture(goex.HBDM)
	log.Println(futureApi.GetExchangeName())
	log.Println(futureApi.GetFutureTicker(goex.BTC_USD, goex.QUARTER_CONTRACT))
	log.Println(futureApi.GetFutureDepth(goex.BTC_USD, goex.QUARTER_CONTRACT, 5))
	//log.Println(futureApi.GetFutureUserinfo()) // account
	//log.Println(futureApi.GetFuturePosition(goex.BTC_USD , goex.QUARTER_CONTRACT))//position info
}
