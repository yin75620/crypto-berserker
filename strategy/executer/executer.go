package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-ini/ini"
	"github.com/yin75620/crypto-berserker/exchange-list/bitmax"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx"
	"github.com/yin75620/crypto-berserker/exchange-list/ftxotc"
	"github.com/yin75620/crypto-berserker/exchange-list/maicoin"
	"github.com/yin75620/crypto-berserker/exchange-list/okex"
	"github.com/yin75620/crypto-berserker/setting"
	Tri "github.com/yin75620/crypto-berserker/strategy/arbitrage/Triangular"
)

var mBitmax = bitmax.NewBitmax(http.DefaultClient)

var m_tri = Tri.NewTriangular(mBitmax)

const (
	FTX    = "FTX"
	BITMAX = "BITMAX"
	MAX    = "MAX"
	OKEX   = "OKEX"
	FTXOTC = "FTXOTC"
)

var mSwitchExchange = BITMAX
var mSubAccount = setting.FTX_SUBACCOUNT

const (
	version = "1.0.1-0004"
)

func main() {

	log.Println(version)

	tInit, bunch := iniSetting()
	s, _ := json.Marshal(tInit)
	log.Println(string(s))

	s2, _ := json.Marshal(bunch)
	log.Println(string(s2))

	log.Println(mSwitchExchange)

	tri := m_tri

	if mSwitchExchange == FTX {
		var mFtx = ftx.NewFtx(http.DefaultClient,
			ftx.FtxInit{
				setting.FTX_KEY,
				setting.FTX_API_SECRET_KEY,
				mSubAccount})
		var mTriFtx = Tri.NewTriangular(mFtx)

		tri = mTriFtx
	} else if mSwitchExchange == MAX {
		var mMai = maicoin.NewMaicoin(http.DefaultClient)
		var mTriMai = Tri.NewTriangular(mMai)
		tri = mTriMai
	} else if mSwitchExchange == OKEX {
		var mMai = okex.NewOkex(http.DefaultClient)
		var mTriMai = Tri.NewTriangular(mMai)
		tri = mTriMai
	} else if mSwitchExchange == FTXOTC {
		var mFtxOtc = ftxotc.NewFtxotc(http.DefaultClient,
			ftxotc.FtxotcInit{
				setting.FTXOTC_KEY,
				setting.FTXOTC_SECRET_KEY})
		var mTriFtx = Tri.NewTriangular(mFtxOtc)

		tri = mTriFtx
	} else {
		tri = m_tri
	}

	tri.SetInit(tInit)
	tri.SetCoinBunch(bunch)
	tri.Start()

}

func iniSetting() (Tri.TriangularInit, []Tri.CoinStrip) {
	cfg, err := ini.Load("main.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	mSwitchExchange = cfg.Section("").Key("ExchangeName").String()
	mSubAccount = cfg.Section("FTX").Key("SubAccount").String()

	resInit := Tri.TriangularInit{
		RangePremium:    cfg.Section("").Key("RangePremium").MustFloat64(),
		LeastTotalValue: cfg.Section("").Key("LeastTotalValue").MustFloat64(),
		DelayTime:       cfg.Section("").Key("DelayTime").MustInt(),
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
