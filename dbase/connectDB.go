package dbase

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// * db constants
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "2997"
	dbname   = "mentorapi"
)

var db *sql.DB

func GetDB() *sql.DB {
	return db
}

func checkTables() error {
	var err error
	queries := []string{studentIQuery, categoryIQuery, mentorIQuery,relationIQuery}
	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			return err
		}

	}
	return nil
}

// * DB connections
func Connect() *sql.DB {
	// * connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// * we are now passing the postgres as the sql type and conn string
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// make sure the tables are present in the database
	err = checkTables()
	if err != nil {
		panic(err)
	}

	return db
}
