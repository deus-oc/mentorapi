package dbase

import (
	"database/sql"
	"log"
)

type Person struct {
	ID   int
	Name string
}

func GetCategoryId(categoryName string, createNewCategory bool) int {
	row := db.QueryRow(getCategoryIdQuery, categoryName)
	var _id int
	switch err := row.Scan(&_id); err {
	case sql.ErrNoRows:
		// no category_name in records
		if createNewCategory {
			_id = InsertCategory(categoryName)
		} else {
			return WRONG_DATA
		}
	case nil:
		return _id
	default:
		log.Print(err)
		return DB_ERROR
	}
	return _id // * if no rows were found sent new made id or error
}

func GetList(option string, _id int) ([]Person, error) {
	var sqlStatement string
	if option == "mentor" {
		sqlStatement = getMentorQuery
	} else {
		sqlStatement = getStudentQuery
	}

	var persons []Person
	rows, err := db.Query(sqlStatement, _id)
	if err != nil {
		log.Print(err)
		return persons, err
	}
	defer rows.Close()
	for rows.Next() {
		var student Person
		err = rows.Scan(&student.ID, &student.Name)
		if err != nil {
			return persons, err
		}
		persons = append(persons, student)
	}
	err = rows.Err()
	if err != nil {
		return persons, err
	}

	return persons, nil
}
