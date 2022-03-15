package DatabaseManager

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
)

func OpenDatabase() (db *sql.DB) {
	sqltrace.Register("mysql", mysql.MySQLDriver{})
	db, err := sqltrace.Open("mysql", "tracker:ddog@tcp(mysql:3306)/tracker_db")

	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(time.Second * 4)

	return db
}

func DBInsert(ctx context.Context, db *sql.DB, table string, columns string, values string) (err error) {
	query_string := "INSERT INTO " + table + "(" + columns + ") VALUES (" + values + ")"
	fmt.Println(query_string)

	insert, err := db.QueryContext(ctx, query_string)
	if err != nil {
		return err
	}

	defer insert.Close()
	return nil
}

func DBDelete(ctx context.Context, db *sql.DB, table string, column string, value string) (err error) {
	query_string := "DELETE FROM " + table + " WHERE " + column + " = \"" + value + "\""

	fmt.Println(query_string)

	insert, err := db.QueryContext(ctx, query_string)
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

func DBSelectTicker(ctx context.Context, db *sql.DB, ticker string) (stock Stock, err error) {
	query_string := "SELECT * FROM stocks_tab WHERE ticker =\"" + ticker + "\";"

	err = db.QueryRowContext(ctx, query_string).Scan(&stock.Ticker, &stock.Company, &stock.Description)

	if err != nil {
		print(err.Error())
		emptyStock := Stock{}
		return emptyStock, errors.New("Stock does not exist in database")
	}

	return stock, err
}
