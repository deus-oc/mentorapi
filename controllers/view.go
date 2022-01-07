package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	dbase "github.com/deus-oc/mentorapi/dbase"
)

func getMentors(w http.ResponseWriter, r *http.Request, categoryName string) {
	w.Header().Set("content-type", "application/json")
	// * fetch the category id
	categoryId := dbase.GetCategoryId(categoryName, false)
	if categoryId == WRONG_DATA {
		w.WriteHeader(http.StatusNoContent)
		return

	} else if categoryId == DB_ERROR {
		serverError(w, r)
		return
	} else {
		// * fetch all the data via categoryId from Mentor table and send it
		mentors, err := dbase.GetList("mentor", categoryId)
		if err != nil {
			serverError(w, r)
			return
		}

		// make json and send
		dataJson, e := json.Marshal(mentors)
		if e != nil {
			log.Print(e)
			serverError(w, r)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write(dataJson)
	}
}

func getStudents(w http.ResponseWriter, r *http.Request, mentorId int) {
	w.Header().Set("content-type", "application/json")

	students, err := dbase.GetList("student", mentorId)
	if err != nil {
		serverError(w, r)
		return
	}

	if len(students) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// make json and send
	dataJson, e := json.Marshal(students)
	if e != nil {
		log.Print(e)
		serverError(w, r)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(dataJson)
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	userReq := r.URL.Query().Get("requestby")
	if userReq == "student" {
		categoryName := r.URL.Query().Get("cname")
		if len(categoryName) == 0 {
			badRequest(w, r)
			return
		} else {
			getMentors(w, r, categoryName)
		}
	} else if userReq == "mentor" {
		mentorId := r.URL.Query().Get("mentorid")
		_id, err := strconv.Atoi(mentorId)
		if err != nil {
			serverError(w, r)
			return
		}
		getStudents(w, r, _id)
	} else {
		badRequest(w, r)
		return
	}
}
