package dbase

import (
	"database/sql"
	"log"
)

func InsertStudent(name string) int {
	var id int
	if err := db.QueryRow(studentIQuery, name).Scan(&id); err != nil {
		log.Print(err)
		return DB_ERROR
	}
	return id
}

func InsertCategory(categoryName string) int {
	var id int
	if err := db.QueryRow(categoryIQuery, categoryName).Scan(&id); err != nil {
		log.Print(err)
		return DB_ERROR
	}
	return id
}

func InsertMentor(name string, categoryId int) int {
	var id int
	if err := db.QueryRow(mentorIQuery, name, categoryId).Scan(&id); err != nil {
		log.Print(err)
		return DB_ERROR
	}
	return id
}

func InsertRelation(studentId, mentorId int) int {
	var id int
	switch err := db.QueryRow(relationIQuery, studentId, mentorId).Scan(&id); err {
	case sql.ErrNoRows:
		return NO_DATA
	case nil:
		return 1
	default:
		log.Print(err)
		return DB_ERROR
	}
}
