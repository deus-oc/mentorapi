package dbase

import (
	"database/sql"
	"fmt"
	"strconv"

	e "github.com/deus-oc/mentorapi/env"
	_ "github.com/lib/pq"
)

var db *sql.DB

func GetDB() *sql.DB {
	return db
}


// * DB connections
func Connect() *sql.DB {
	port, err := strconv.Atoi(e.GetEnvVar("PORT"))
	// * connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		e.GetEnvVar("HOST"), port, e.GetEnvVar("USER"), e.GetEnvVar("PASS"), e.GetEnvVar("DBNAME"))

	// * we are now passing the postgres as the sql type and conn string
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
