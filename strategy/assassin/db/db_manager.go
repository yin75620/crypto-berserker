package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBManager struct {
	gormDB *gorm.DB
	sqlDB  *sql.DB
}

func NewDBManager() *DBManager {
	dbm := &DBManager{}
	return dbm
}

// GetDB 提供外部取得 Gorm 連線的函式
func (manager *DBManager) OpenGormDB() (*gorm.DB, error) {
	dsn := getdsn()
	// 使用 Gorm 連接 MySQL
	sqldb := mysql.Open(dsn)
	gormDb, err := gorm.Open(sqldb, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	manager.gormDB = gormDb
	return manager.gormDB, nil
}

func (manager *DBManager) OpenSqlDB() (*sql.DB, error) {
	if manager.sqlDB != nil {
		manager.sqlDB.Close()
	}
	dsn := getdsn()
	// 設定 MySQL 連接資訊
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %v", err)
	}

	manager.sqlDB = db
	return manager.sqlDB, nil
}

func getdsn() string {
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbName)

	return dsn
}

// findProjectRoot 用來找尋 .env 檔案的根目錄
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
