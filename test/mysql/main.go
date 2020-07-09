package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("Go MySQL Tutorial")

	//CREATE SCHEMA `crypto` DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ;

	// create table sql
	/*
			CREATE TABLE `crypto`.`log_cross_exchange_tick` (
		  `id` INT NOT NULL AUTO_INCREMENT,
		  `ask_exchange` VARCHAR(45) NULL,
		  `ask_c_price` DOUBLE NULL,
		  `ask_s_price` DOUBLE NULL,
		  `ask_total_volume` DOUBLE NULL,
		  `bid_exchange` VARCHAR(45) NULL,
		  `bid_c_price` DOUBLE NULL,
		  `bid_s_price` DOUBLE NULL,
		  `bid_total_volume` DOUBLE NULL,
		  `profit` DOUBLE NULL,
		  `min_total_volume` DOUBLE NULL,
		  `create_time` TIMESTAMP NULL,
		  PRIMARY KEY (`id`))
		COMMENT = 'save every x ms data';

	*/

	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	db, err := sql.Open("mysql", "root:79125@tcp(127.0.0.1:3306)/crypto")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	// perform a db.Query insert
	insert, err := db.Query(`INSERT INTO crypto.log_cross_exchange_tick 
	(ask_exchange, ask_c_price, ask_s_price, ask_total_volume, bid_exchange, bid_c_price, bid_s_price, bid_total_volume, profit, min_total_volume,
	 create_time) VALUES ('x1', '0.1', '0.2', '0.3', 'x2', '1.1', '1.2', '1.3', '2.1', '3.1', '1594305390');`)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()

}
