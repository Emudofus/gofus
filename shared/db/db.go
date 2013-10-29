package db

import (
	"database/sql"
	"log"
)

type Configuration struct {
	Driver         string
	DataSourceName string
}

func Open(config *Configuration) (db *sql.DB) {
	db, err := sql.Open(config.Driver, config.DataSourceName)
	if err != nil {
		panic(err.Error())
	}

	var dbname string
	if err = db.QueryRow("select current_database()").Scan(&dbname); err != nil {
		panic(err)
	}

	log.Print("[database] `", dbname, "` successfully opened")
	return
}
