package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// 取得當前工作目錄
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	loadenvPath, _ := findProjectRoot(cwd, ".env")

	// 讀取 .env 文件
	err = godotenv.Load(fmt.Sprintf("%s/.env", loadenvPath))
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// 從環境變數讀取 MySQL 連接資訊
	username := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PWD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DB")

	// 組合 DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", username, password, host, port)

	// 設定 MySQL 連接資訊
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// 創建資料庫
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName)
	_, err = db.Exec(createDBQuery)
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}
}

func findProjectRoot(startPath string, rootFileName string) (string, error) {
	currentPath := startPath

	for {
		if _, err := os.Stat(filepath.Join(currentPath, rootFileName)); err == nil {
			return currentPath, nil
		}

		// 上移一層
		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			// 已經到達根目錄，停止搜索
			return "", os.ErrNotExist
		}

		currentPath = parentPath
	}
}
