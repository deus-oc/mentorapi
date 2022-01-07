package controllers

import (
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

	mentorData := dbase.InsertRelation(d.StudentId, d.MentorId)
	switch mentorData {
	case NO_DATA:
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"status": "already present in db"}`))
	case 1:
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"status": "success"}`))
	default:
		serverError(w, r)
	}
}
