package main

import (
	"log"
	"net/http"

	"github.com/yin75620/crypto-berserker/ftx"
)

var m_ftxClient = ftx.NewFtx(http.DefaultClient)

func main() {
	res := m_ftxClient.GetStructFills()
	log.Println(res)
}
