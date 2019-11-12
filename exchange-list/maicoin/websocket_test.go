package maicoin

import (
	"fmt"
	"testing"
)

func TestWsStart(t *testing.T) {

	mws := NewSocket()
	mws.Strat()
	ch := mws.SubScribeOrderBook("ethtwd")

	mws2 := NewSocket()
	mws2.Strat()
	ch2 := mws.SubScribeOrderBook("ethusdt")

	go func() {
		for {
			fmt.Println("ch2:", <-ch2)
		}
	}()

	for {
		fmt.Println("ch:", <-ch)
	}

	//time.Sleep(time.Duration(100) * time.Second)
}
