package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yin75620/crypto-berserker/strategy/assassin/db"
)

func main() {
	dbmanager := db.NewDBManager()
	db, err := dbmanager.OpenSqlDB()
	if err != nil {
		fmt.Println(err)
	}

	// // 設定 MySQL 連接資訊
	// db, err := sql.Open("mysql", dsn)
	// if err != nil {
	// 	log.Fatalf("Error opening database: %v", err)
	// }
	// defer db.Close()

	// 創建資料庫
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", "stock_ai_avalon1")
	_, err = db.Exec(createDBQuery)
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}
}
