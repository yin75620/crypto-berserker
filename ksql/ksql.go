package ksql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yin75620/crypto-berserker/setting"
)

type Ksql struct {
	db      *sql.DB
	isStart bool
}

func NewKsql() *Ksql {
	ksql := Ksql{}
	ksql.isStart = false

	return &ksql
}

func (k *Ksql) Start() error {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/crypto", setting.DB_PASSWORD))
	// if there is an error opening the connection, handle it
	if err != nil {
		return err
	}
	k.isStart = true
	k.db = db
	return nil
}

func (k *Ksql) Insert(s string) error {
	if !k.isStart {
		return nil
	}
	// perform a db.Query insert
	insert, err := k.db.Query(s)

	// if there is an error inserting, handle it
	if err != nil {
		return err
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()
	return nil
}

func (k *Ksql) End() {
	k.db.Close()
	k.isStart = false
}
