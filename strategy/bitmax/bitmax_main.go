package main

import (
	"net/http"

	"github.com/yin75620/crypto-berserker/bitmax"
	"github.com/yin75620/crypto-berserker/ftx"
	"github.com/yin75620/crypto-berserker/setting"
	"github.com/yin75620/crypto-berserker/strategy/arbitrage/Triangular"
)

var mBitmax = bitmax.NewBitmax(http.DefaultClient)

var m_tri = Triangular.NewTriangular(mBitmax)

var mFtx = ftx.NewFtx(http.DefaultClient,
	ftx.FtxInit{
		setting.FTX_KEY,
		setting.FTX_API_SECRET_KEY,
		setting.FTX_SUBACCOUNT})
var mTriFtx = Triangular.NewTriangular(mFtx)

func main() {

	m_tri.SetDealCoin([]Triangular.FlowString{
		Triangular.FlowString{Coins: []string{"FTT", "USDT"}},
		Triangular.FlowString{Coins: []string{"FTT", "BTC", "USDT"}},
	})
	m_tri.Start()

	/*mTriFtx.SetDealCoin([][]string{
		[]string{"FTT", "USD"},
		[]string{"FTT", "BTC", "USD"},
	})
	mTriFtx.Start()*/
}
