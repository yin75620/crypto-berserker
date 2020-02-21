package CrossExchange

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/bybit"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx"
	"github.com/yin75620/crypto-berserker/setting"
)

func TestSaveFile(t *testing.T) {

	ft := ftx.NewFtx(http.DefaultClient, ftx.FtxInit{
		setting.FTX_KEY,
		setting.FTX_API_SECRET_KEY,
		"tester"})
	bybit := bybit.NewBybit(http.DefaultClient)

	crossPairsTable := map[string][]CrossPair{}
	cp := CrossPair{
		askExchange:  ft,
		bidExchange:  bybit,
		askPricePair: exc.PricePair{9600, 100},
		bidPricePair: exc.PricePair{9600, 100},
		orderVolume:  200,
	}

	cpArray := []CrossPair{
		cp,
	}
	crossPairsTable[cp.GetName()] = cpArray
	savePairMapToFile(crossPairsTable)
}

func TestLoadFile(t *testing.T) {
	exchanges := []exc.Exchange{}
	ft := ftx.NewFtx(http.DefaultClient, ftx.FtxInit{
		setting.FTX_KEY,
		setting.FTX_API_SECRET_KEY,
		"tester"})
	bybit := bybit.NewBybit(http.DefaultClient)
	exchanges = append(exchanges, ft)
	exchanges = append(exchanges, bybit)

	res, err := loadPairMapFromFile(exchanges)
	fmt.Println(res)
	assert.NoError(t, err)
}
