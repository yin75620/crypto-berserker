package main

import (
	"fmt"
	"runtime"
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

func TestMemoryCheck(t *testing.T) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Printf("Allocated memory: %d bytes", ms.Alloc)

	fmt.Printf("Allocated memory: %d bytes", ms.Alloc)

	fmt.Printf("Heap memory usage: %d bytes", ms.HeapAlloc)

}
