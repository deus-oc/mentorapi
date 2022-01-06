package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	dbase "github.com/deus-oc/mentorapi/dbase"
)

func SelectMentor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var d DuoIdentity
	if err := decoder.Decode(&d); err != nil {
		log.Print(err)
		serverError(w, r)
		return
	}

	if d.MentorId == 0 || d.StudentId == 0 {
		badRequest(w, r)
		return
	}

	db := dbase.GetDB()
	sqlStatement := `
	INSERT INTO relation (student_id, mentor_id)
	VALUES ($1, $2)
	ON CONFLICT (student_id)
	DO NOTHING
	RETURNING relation_id
	`
	var _id int
	switch err := db.QueryRow(sqlStatement, d.StudentId, d.MentorId).Scan(&_id); err {
	case sql.ErrNoRows:
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"status": "already present in db"}`))
	case nil:
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"status": "success"}`))

	// case missing which is fk_errors on relation i.e. no student_id or no mentor_id present in respective tables

	default:
		serverError(w, r)
		return
	}
}
