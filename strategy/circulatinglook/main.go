package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	simpleLog "github.com/yin75620/crypto-berserker/log"
)

func main() {
	slog := simpleLog.StartLog()
	defer slog.Close()

	LogAmount("https://api.raydium.io/ray/circulating")
}

func LogAmount(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	log.Println(string(body))
	NextRound(LogAmount, url)
}

func NextRound(fc func(string), url string) {
	now := time.Now().UTC()
	nextHour := now.Add(time.Second * 1)
	midnoon := time.Date(nextHour.Year(), nextHour.Month(), nextHour.Day(),
		nextHour.Hour(), 0, 0, 0, now.Location())
	duration := midnoon.Sub(now)
	time.Sleep(duration)
	fc(url)
}
