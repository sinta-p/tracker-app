package DatabaseManager

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func OpenDatabase() (db *sql.DB) {
	db, err := sql.Open("mysql", "tracker:@tcp(localhost:3306)/tracker_db")

	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(time.Second * 4)

	return db
}

func DBInsert(db *sql.DB, table string, columns string, values string) (err error) {
	query_string := "INSERT INTO " + table + "(" + columns + ") VALUES (" + values + ")"
	fmt.Println(query_string)

	insert, err := db.Query(query_string)
	if err != nil {
		return err
	}

	defer insert.Close()
	return nil
}

type Stock struct {
	Ticker      string `json:"ticker"`
	Company     string `json:"company"`
	Description string `json:"description"`
}

func DBSelectTicker(db *sql.DB, ticker string) (stock Stock, err error) {
	query_string := "SELECT * FROM stocks_tab WHERE ticker =\"" + ticker + "\";"

	err = db.QueryRow(query_string).Scan(&stock.Ticker, &stock.Company, &stock.Description)

	if err != nil {
		print(err.Error())
		emptyStock := Stock{}
		return emptyStock, errors.New("Stock does not exist in database")
	}

	return stock, err
}
