package main

import (
	"fmt"
	"testing"

	"github.com/yin75620/crypto-berserker/exchange-list/common"
	"github.com/yin75620/crypto-berserker/jtime"
	"github.com/yin75620/crypto-berserker/strategy/assassin/db"
)

func TestMain(t *testing.T) {
	main()
}

func TestFree(t *testing.T) {
	fmt.Println(jtime.UnixToTime(1567965420000))
}

func TestEarly(t *testing.T) {
	collector := common.GetCollector("")

	// 連接 MySQL 資料庫
	dbmanager := db.NewDBManager()
	gormDB, err := dbmanager.OpenGormDB()
	if err != nil {
		fmt.Println(err)
	}

	save_bar_into_db(gormDB, collector, jtime.UnixToTime(1566965420000), jtime.UnixToTime(1577965420000), 1)
}
