package models

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"sync"
)

var DB_URL string
var DB_PORT string
var DB_USER string
var DB_PW string
var DB_NAME string

var dbLock *sync.Mutex = new(sync.Mutex)

var connection *sql.DB

func safelyConnect() error {
	dbLock.Lock()
	if connection == nil {
		openFmt := "postgres://%s:%s@%s:5432/%s?sslmode=disable"
		dsn := fmt.Sprintf(openFmt, DB_USER, DB_PW, DB_URL, DB_NAME)
		var err error
		connection, err = sql.Open("postgres", dsn)
		if err != nil {
			return err
		}
	}
	err := connection.Ping()
	dbLock.Unlock()
	if err != nil {
		return err
	}
	return nil
}
