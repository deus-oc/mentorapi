package controllers

import (
	"encoding/json"
	"net/http"

	dbase "github.com/deus-oc/mentorapi/dbase"
)

func saveUser(user *RegsiterDetail) int {
	var id int

	if user.Choice == "student" { // * add to student
		id = dbase.InsertStudent(user.Name)

	} else { // * add to mentor
		if len(user.Category) == 0 {
			return WRONG_DATA
		}
		categoryId := dbase.GetCategoryId(user.Category, true)
		if categoryId == DB_ERROR {
			return DB_ERROR
		}
		id = dbase.InsertMentor(user.Name, categoryId)
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
				res := struct {
					Status string
					Id     int
				}{
					Status: "success",
					Id:     _id,
				}
				resJson, err := json.Marshal(res)
				if err != nil {
					serverError(w, r)
					return
				}
				w.WriteHeader(http.StatusAccepted)
				w.Write(resJson)
			}
		}
	}
}
