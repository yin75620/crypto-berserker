package Triangular

import (
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
	"Saber"})
var mTri *Triangular = NewTriangular(mFtx)

func TestExecuteOrder(t *testing.T) {
	dealFlow := mTri.NewDealFlow("FTT", []string{"USD"})
	mTri.Init.RangePremium = 0.02
	mTri.executeOrder(dealFlow, exc.Ask, 1.8546)
}
