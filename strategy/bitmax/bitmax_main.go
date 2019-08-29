package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-ini/ini"
	"github.com/yin75620/crypto-berserker/bitmax"
	"github.com/yin75620/crypto-berserker/ftx"
	"github.com/yin75620/crypto-berserker/setting"
	Tri "github.com/yin75620/crypto-berserker/strategy/arbitrage/Triangular"
)

var mBitmax = bitmax.NewBitmax(http.DefaultClient)

var m_tri = Tri.NewTriangular(mBitmax)

var mFtx = ftx.NewFtx(http.DefaultClient,
	ftx.FtxInit{
		setting.FTX_KEY,
		setting.FTX_API_SECRET_KEY,
		setting.FTX_SUBACCOUNT})
var mTriFtx = Tri.NewTriangular(mFtx)

func main() {

	tInit, bunch := iniSetting()
	s, _ := json.Marshal(tInit)
	log.Println(string(s))
	m_tri.SetInit(tInit)

	s2, _ := json.Marshal(bunch)
	log.Println(string(s2))
	m_tri.SetCoinBunch(bunch)

	/*m_tri.SetCoinBunch([]Tri.CoinStrip{
		Tri.CoinStrip{Coins: []string{"FTT", "USDT"}},
		Tri.CoinStrip{Coins: []string{"FTT", "BTC", "USDT"}},
	})
	m_tri.SetInit(Tri.TriangularInit{
		RangePremium:    0.1,
		LeastTotalValue: 10,
		RANK_S:          []float64{0.006, -3.0, 1000},
		RANK_N:          []float64{0.001, -2.0, 300},
	})*/
	m_tri.Start()

}

func iniSetting() (Tri.TriangularInit, []Tri.CoinStrip) {
	cfg, err := ini.Load("main.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	resInit := Tri.TriangularInit{
		RangePremium:    cfg.Section("").Key("RangePremium").MustFloat64(),
		LeastTotalValue: cfg.Section("").Key("LeastTotalValue").MustFloat64(),
		RANK_S: []float64{
			cfg.Section("RankS").Key("MinProfit").MustFloat64(),
			cfg.Section("RankS").Key("PlusSecond").MustFloat64(),
			cfg.Section("RankS").Key("TotalValuesUS").MustFloat64(),
		},
		RANK_N: []float64{
			cfg.Section("RankN").Key("MinProfit").MustFloat64(),
			cfg.Section("RankN").Key("PlusSecond").MustFloat64(),
			cfg.Section("RankN").Key("TotalValuesUS").MustFloat64(),
		},
	}

	resFee := cfg.Section("").Key("TakerFee").MustFloat64()
	log.Println(resFee)

	resFlow := []Tri.CoinStrip{}
	for _, value := range cfg.Section("CoinStrip").Keys() {
		res := value.String()
		log.Println(res)

		coins := strings.Split(res, ",")
		resFlow = append(resFlow, Tri.CoinStrip{Coins: coins})
	}
	return resInit, resFlow
}
