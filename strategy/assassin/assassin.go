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
var BINANCE_BTC_START_TS int64 = 1567965420000

func main() {
	collector := common.GetCollector("")

	// 連接 MySQL 資料庫
	dbmanager := db.NewDBManager()
	gormDB, err := dbmanager.OpenGormDB()
	if err != nil {
		fmt.Println(err)
	}

	barCount := 1000

	// startTime := BINANCE_BTC_START_TIME
	// endTime := startTime.Add(time.Duration(barCount) * time.Second)
	now := time.Now().UnixMilli()
	period := int64(barCount) * 60 * 1000 * 1000

	for ti := BINANCE_BTC_START_TS; ti < now; ti = ti + period {
		startTime := ti
		endTime := startTime + period
		save_bar_into_db(gormDB, collector, startTime, endTime, barCount)
	}
}

func save_bar_into_db(gormDB *gorm.DB, collector exchange.Collector, startTime, endTime int64, barCount int) {

	var count int64
	err := gormDB.Model(&exchange.CandleStick{}).Where("open_time >= ? AND open_time <= ?", startTime, endTime).Count(&count).Error
	if err != nil {
		fmt.Println("Error counting Klines in database:", err)
		return
	}

	if count >= int64(barCount) {
		fmt.Printf("Klines between %v and %v already exist in the database. %d\n", startTime, endTime, count)
		return
	}

	var existingKlines []exchange.CandleStick
	err = gormDB.Where("open_time >= ? AND open_time <= ?", startTime, endTime).Find(&existingKlines).Error
	if err != nil {
		fmt.Println("Error querying Klines from database:", err)
		return
	}

	klines, err := collector.GetKlines("BTCUSDT", "1m", startTime, endTime, barCount)
	if err != nil {
		fmt.Println("Error fetching Klines:", err)
		return
	}

	// 用map存储已有的K线数据的open_time，方便快速查找
	existingMap := make(map[float64]bool)
	for _, kline := range existingKlines {
		existingMap[kline.OpenTime] = true
	}

	// 过滤掉已经存在的K线数据
	var newKlines []exchange.CandleStick
	for _, kline := range klines {
		if !existingMap[kline.OpenTime] {
			newKlines = append(newKlines, kline)
		}
	}

	// 批量插入 K 線資料
	result := gormDB.Create(&newKlines) // 這裡傳入 klines slice
	if result.Error != nil {
		fmt.Println("Error inserting Klines into database:", result.Error)
	} else {
		fmt.Printf("%d Klines inserted successfully\n", result.RowsAffected)
	}

}
