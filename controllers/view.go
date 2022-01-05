package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	dbase "github.com/deus-oc/mentorapi/dbase"
)

func getList(sqlStatement string, _id int) ([]Person, error) {
	db := dbase.GetDB()
	rows, err := db.Query(sqlStatement, _id)

	var persons []Person

	if err != nil {
		log.Print(err)
		return persons, err
	}
	// * called because in case of error in for rows.Next(), Close() is not called automatically
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

	// return the list of students via json
	return persons, nil
}

func getMentors(w http.ResponseWriter, r *http.Request, categoryName string) {
	w.Header().Set("content-type", "application/json")
	// * fetch the category id
	_id := getCategoryId(categoryName, false)
	if _id == WRONG_DATA {
		badRequest(w, r)
		return

	} else if _id == DB_ERROR {
		serverError(w, r)
		return
	} else {
		// * fetch all the data via categoryId from Mentor table and send it
		sqlStatement := `SELECT mentor_id,mentor_name
		FROM mentor
		WHERE category_id=$1`

		mentors, err := getList(sqlStatement, _id)
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
	sqlStatement := `SELECT s.student_id, s.student_name
	FROM relation r
	INNER JOIN student s
	ON s.student_id = r.student_id
	WHERE mentor_id=$1`

	students, err := getList(sqlStatement, mentorId)
	if err != nil {
		serverError(w, r)
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
