package Triangular

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx"
	"github.com/yin75620/crypto-berserker/exchange-list/maicoin"
	"github.com/yin75620/crypto-berserker/setting"
)

var mMai = maicoin.NewMaicoin(http.DefaultClient)
var mFtx = ftx.NewFtx(http.DefaultClient, ftx.FtxInit{
	setting.FTX_KEY,
	setting.FTX_API_SECRET_KEY,
	"apicaller"})
var mTri *Triangular = NewTriangular(mFtx)

func TestExecuteOrder(t *testing.T) {
	dealFlow := mTri.NewDealFlow("FTT", []string{"USD"})
	mTri.Init.RangePremium = 0.02
	mTri.executeOrder(dealFlow, exc.Ask, 1.8546)
}

func TestDealFlow(t *testing.T) {
	dealFlow := mTri.NewDealFlow("FTT", []string{"BTC", "USD"})
	laName := dealFlow.getName()

	laPrice := dealFlow.getFinalPair(exc.Ask).Price

	laVolume := dealFlow.getFinalPair(exc.Ask).Volume
	laValue := laPrice * laVolume

	log.Println(fmt.Sprintf("resAsk:%f, laValue:%f, AskCoin:%s", laPrice, laValue, laName))
}
