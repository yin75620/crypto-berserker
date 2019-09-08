package Triangular

import (
	"net/http"
	"testing"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/maicoin"
)

var mMai = maicoin.NewMaicoin(http.DefaultClient)
var mTri *Triangular = NewTriangular(mMai)

func TestExecuteOrder(t *testing.T) {
	dealFlow := mTri.NewDealFlow("ETH", []string{"TWD"})
	mTri.Init.RangePremium = 0.02
	mTri.executeOrder(dealFlow, exc.Ask, 0.793545)
}
