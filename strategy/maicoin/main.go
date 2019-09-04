package main

import (
	"fmt"
	"net/http"

	"github.com/yin75620/crypto-berserker/maicoin"
	Tri "github.com/yin75620/crypto-berserker/strategy/arbitrage/Triangular"
)

var mMai = maicoin.NewMaicoin(http.DefaultClient)
var mTri = Tri.NewTriangular(mMai)

func main() {
	fmt.Println("TEST")

	mTri.SetCoinArrays([][]string{[]string{"max", "twd"}})
	mTri.Start()

	//var mm = maicoin.NewMaicoin(http.DefaultClient)

	//mm.GetAccountInfo()

	//mm.GetFill("usdttwd")

}
