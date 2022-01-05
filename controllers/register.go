package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	dbase "github.com/deus-oc/mentorapi/dbase"
)

func saveUser(user *RegsiterDetail) int {
	db := dbase.GetDB()
	var id int

	if user.Choice == "student" { // * add to student
		sqlStatement := `
		INSERT INTO student (student_name)
		VALUES ($1) 
		RETURNING student_id`

		if err := db.QueryRow(sqlStatement, user.Name).Scan(&id); err != nil {
			log.Print(err)
			return DB_ERROR
		}

	} else { // * add to mentor
		if len(user.Category) == 0 {
			return WRONG_DATA
		}
		categoryId := getCategoryId(user.Category, true)
		if categoryId == DB_ERROR {
			return DB_ERROR
		}
		sqlStatement := `
		INSERT INTO mentor(mentor_name, category_id)
		VALUES ($1, $2)
		RETURNING mentor_id
		`
		if err := db.QueryRow(sqlStatement, user.Name, categoryId).Scan(&id); err != nil {
			return DB_ERROR
		}
	}
	return id
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var user RegsiterDetail
	if err := decoder.Decode(&user); err != nil {
		badRequest(w, r)
		return
	} else {
		if (user.Choice != "student" && user.Choice != "mentor") || len(user.Name) == 0 {
			badRequest(w, r)
			return
		} else {
			_id := saveUser(&user)
			if _id == DB_ERROR {
				serverError(w, r)
				return
			} else if _id == WRONG_DATA {
				// * means that no category name inputted after selection of mentor
				badRequest(w, r)
				return
			} else {
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte(`{"status": "success"}`))
			}
		}
	}
}
