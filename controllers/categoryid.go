package controllers

import (
	"database/sql"
	"log"

	dbase "github.com/deus-oc/mentorapi/dbase"
)

func getCategoryId(cName string, canMakeNew bool) int {
	db := dbase.GetDB()
	sqlStatement := `SELECT category_id 
	FROM CATEGORY
	WHERE category_name=$1`
	row := db.QueryRow(sqlStatement, cName)
	var _id int
	switch err := row.Scan(&_id); err {
	case sql.ErrNoRows:
		// no category_name in records
		if canMakeNew {
			sqlStatement := `
			INSERT INTO category (category_name)
			VALUES ($1) 
			RETURNING category_id`
			if err := db.QueryRow(sqlStatement, cName).Scan(&_id); err != nil {
				log.Print(err)
				return DB_ERROR
			}
		} else {
			return WRONG_DATA
		}
	case nil:
		return _id
	default:
		log.Print(err)
		return DB_ERROR
	}
	return _id // * if no rows were found sent new made id
}
