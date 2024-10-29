package main

import (
	"fmt"
	"time"

	"github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/common"
	"github.com/yin75620/crypto-berserker/jtime"
	"github.com/yin75620/crypto-berserker/strategy/assassin/db"
	"gorm.io/gorm"
)

var BINANCE_BTC_START_TIME = jtime.UnixToTime(1567965420000)

func main() {
	collector := common.GetCollector("")

	// 連接 MySQL 資料庫
	dbmanager := db.NewDBManager()
	gormDB, err := dbmanager.OpenGormDB()
	if err != nil {
		fmt.Println(err)
	}

	// startTime, _ := time.Parse(time.RFC3339, "1980-01-01T00:00:00Z")
	// endTime, _ := time.Parse(time.RFC3339, "2024-01-02T00:00:00Z")
	barCount := 1000

	// startTime := BINANCE_BTC_START_TIME
	// endTime := startTime.Add(time.Duration(barCount) * time.Second)
	now := time.Now()

	for ti := BINANCE_BTC_START_TIME; ti.Compare(now) < 1; ti = ti.Add(time.Duration(barCount) * time.Second) {
		startTime := ti
		endTime := startTime.Add(time.Duration(barCount) * time.Second)
		save_bar_into_db(gormDB, collector, startTime, endTime, barCount)
	}
}

func save_bar_into_db(gormDB *gorm.DB, collector exchange.Collector, startTime, endTime time.Time, barCount int) {

	klines, err := collector.GetKlines("BTCUSDT", "1m", startTime, endTime, barCount)
	if err != nil {
		fmt.Println("Error fetching Klines:", err)
		return
	}

	// 顯示取得的 K 線資料
	for _, kline := range klines {
		// 將 K 線資料插入到資料庫
		result := gormDB.Create(&kline)
		if result.Error != nil {
			fmt.Println("Error inserting Kline into database:", result.Error)
		} else {
			fmt.Println("Kline inserted successfully")
		}
		fmt.Printf("Open: %f, Close: %f, Volume: %f\n", kline.Open, kline.Close, kline.Volume)
	}
}
