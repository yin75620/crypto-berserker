package main

import (
	"fmt"
	"net/http"

	fx "github.com/yin75620/crypto-berserker/ftx"
)

func main() {
	fmt.Println("TEST")

	var ftx = fx.NewFtx(http.DefaultClient)

	ftx.GetAccountInfo()

	//ftx.GetMarkets()
}
