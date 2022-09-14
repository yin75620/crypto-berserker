package FundingRate

import (
	"log"
	"os"
	"time"

	exc "github.com/yin75620/crypto-berserker/exchange"

	simpleLog "github.com/yin75620/crypto-berserker/log"
	"github.com/yin75620/crypto-berserker/message_tool"
)

type FundingRate struct {
	exchangeClient exc.Exchange
}

func NewFundingRate(exchange exc.Exchange) *FundingRate {
	fr := FundingRate{}
	fr.exchangeClient = exchange
	return &fr
}

func (fr *FundingRate) Start() {
	message_tool.StartTelegram()
	var logFile *os.File = simpleLog.StartLog()
	defer logFile.Close()
	fr.stratStrategy()

	infoStr := string(fr.exchangeClient.GetAccountInfo())
	message_tool.SendTelegram(infoStr)

	//mPreWallet = fr.exchangeClient.GetWallet()
	//mPreTime = time.Now()
	//mPreLoseCheckWallet = fr.exchangeClient.GetWallet()

	liveTime := time.Second * time.Duration(2 /*fr.Init.ShowLiveSecond*/)
	liveTimer := time.NewTimer(liveTime)

	var delay_time_ms int = 1000.0 //fr.Init.DelayTime
	d := time.Duration(time.Millisecond * time.Duration(delay_time_ms))

	t := time.NewTimer(d)
	defer t.Stop()

	for {
		select {
		case <-t.C:

			plusMilliSecond := fr.stratStrategy()
			t.Reset(time.Millisecond * time.Duration(delay_time_ms+plusMilliSecond))
		case <-liveTimer.C:
			log.Println("live")
			liveTimer.Reset(liveTime)
		}

	}

}

func (fr *FundingRate) stratStrategy() int {
	futures := exc.Futures{
		TargetName: "ETH",
		QuoteCoin:  "USD",
	}
	askPair, bidPair := fr.exchangeClient.GetFuturesAskBidPair(futures)
	log.Println(askPair)
	log.Println(bidPair)

	coinPair := exc.CoinPair{
		BaseCoin:   "ETH",
		QuotedCoin: "USD",
	}

	askP, bidP := fr.exchangeClient.GetAskBidPair(coinPair, 1)
	log.Println(askP)
	log.Println(bidP)

	return 0
}

// 策略執行
// 取得期貨5檔
// 取得現貨5檔 (可換月結算期貨)
// 取得資金費率
// 判斷資金費率與期現貨收斂是否同方向
// 1. 資金費率 > 0 則 期貨價格 - 現貨 > 0
// 2. 資金費率 < 0 則 期貨價格 - 現貨 < 0
// 開始建倉: 期貨數量 = 現貨數量
// 1. 當所持有部位 > 本金9倍即停止建倉 (開10倍槓桿)
// 關艙條件: 資金費率反轉持續N次
