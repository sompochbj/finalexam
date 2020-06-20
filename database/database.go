package database

import (
	"database/sql"
	"os"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	createTb := `CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);`

	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table customers", err)
	}
}

// Conn ...
func Conn() *sql.DB {
	return db
}