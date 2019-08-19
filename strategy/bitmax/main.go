package main

import (
	"fmt"
	"net/http"

	bitmax "github.com/yin75620/crypto-berserker/bitmax"
)

var m_bitmaxClient = bitmax.NewBitmax(http.DefaultClient)

func main() {
	res := m_bitmaxClient.GetAccountInfo()

	fmt.Println(string(res))
}
