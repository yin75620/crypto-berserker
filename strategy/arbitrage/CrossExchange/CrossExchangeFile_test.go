package CrossExchange

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/binancef"
	"github.com/yin75620/crypto-berserker/exchange-list/bybit"
)

func TestSaveFile(t *testing.T) {

	bnf := binancef.NewBinancef(http.DefaultClient)
	bybit := bybit.NewBybit(http.DefaultClient)

	crossPairsTable := map[string][]CrossPair{}
	cp := CrossPair{
		Symbol:       "BTCUSDT",
		askExchange:  bnf,
		bidExchange:  bybit,
		askPricePair: exc.PricePair{9600, 100},
		bidPricePair: exc.PricePair{9600, 100},
	}

	cpArray := []CrossPair{
		cp,
	}
	crossPairsTable[cp.GetName()] = cpArray
	savePairMapToFile(crossPairsTable)
}

func TestLoadFile(t *testing.T) {
	exchanges := []exc.Exchange{}
	bnf := binancef.NewBinancef(http.DefaultClient)
	bybit := bybit.NewBybit(http.DefaultClient)
	exchanges = append(exchanges, bnf)
	exchanges = append(exchanges, bybit)

	res, err := loadPairMapFromFile(exchanges)
	fmt.Println(res)
	assert.NoError(t, err)
}
