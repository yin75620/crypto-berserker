package main

import (
	"net/http"

	"github.com/yin75620/crypto-berserker/bitmax"
	"github.com/yin75620/crypto-berserker/strategy/arbitrage/Triangular"
)

var mBitmax = bitmax.NewBitmax(http.DefaultClient)

var m_tri = Triangular.NewTriangular(mBitmax)

func main() {
	m_tri.SetDealCoin([][]string{
		[]string{"FTT", "USDT"},
		[]string{"FTT", "BTC", "USDT"},
	})
	m_tri.Start()
}
