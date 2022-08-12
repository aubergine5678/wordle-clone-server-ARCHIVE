package db_client

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DBClient *sql.DB

func InitialiseDBConnection() {
	log.Println("Initialising DB connection...")

	// If you really want to get this running locally for some reason, enter your MySQL Database details below:
	db, err := sql.Open("mysql", "[username]:[password]@[protocol](hostname:port)/[db-name]?[options]")
	// Template for connecting: "<username>:<password>@<protocol>(<hostname>:<port>)/<dbname>?<options>"

	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	DBClient = db
	log.Println("Successfully connected to DB")
}
