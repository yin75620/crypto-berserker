package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/yin75620/crypto-berserker/exchange-list/common"
	"github.com/yin75620/crypto-berserker/jtime"
	"github.com/yin75620/crypto-berserker/strategy/assassin/db"
)

func TestMain(t *testing.T) {
	main()
}

func TestFree(t *testing.T) {

	fmt.Println(jtime.UnixToTime(1567965420000))
	now := time.Now()
	fmt.Println(now.Unix())
	fmt.Println(now.UnixMilli())
	fmt.Println(now.UnixMicro())
	fmt.Println(now.UnixNano())

}

func TestEarly(t *testing.T) {
	collector := common.GetCollector("")

	// 連接 MySQL 資料庫
	dbmanager := db.NewDBManager()
	gormDB, err := dbmanager.OpenGormDB()
	if err != nil {
		fmt.Println(err)
	}

	save_bar_into_db(gormDB, collector, 1566965420000, 1577965420000, 1)
}
