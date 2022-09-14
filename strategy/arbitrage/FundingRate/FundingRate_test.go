package FundingRate

import (
	"log"
	"net/http"
	"testing"

	"github.com/yin75620/crypto-berserker/exchange-list/ftx"
	"github.com/yin75620/crypto-berserker/setting"
)

func TestMain(t *testing.T) {
	ft := ftx.NewFtx(http.DefaultClient, ftx.FtxInit{
		setting.FTX_KEY,
		setting.FTX_API_SECRET_KEY,
		"tester"})

	fr := NewFundingRate(ft)
	fr.Start()
	log.Println("End")
}
