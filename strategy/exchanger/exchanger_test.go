package main

import (
	"testing"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/common"
)

func TestExecuteOrder(t *testing.T) {
	iniSetting()

	exchanges := []exc.Exchange{}

	for _, v := range mExchangeStrings {
		ex := common.GetExchange(v)
		exchanges = append(exchanges, ex)
	}

	sendWallet(exchanges)

}
