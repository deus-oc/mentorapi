package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// * db constants
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "2997"
	dbname   = "mentorapi"
)

// * enumerations
const (
	DB_ERROR = iota - 2
	WRONG_DATA
)

var db *sql.DB

type RegsiterDetail struct {
	Name     string `json:"name"`
	Choice   string `json:"choice"`
	Category string `json:"category_name"`
}

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DuoIdentity struct {
	StudentId int `json:"student_id"`
	MentorId  int `json:"mentor_id"`
}

// * controllers

func getCategoryId(cName string, canMakeNew bool) int {
	sqlStatement := `SELECT category_id 
	FROM CATEGORY
	WHERE category_name=$1`
	row := db.QueryRow(sqlStatement, cName)
	var _id int
	switch err := row.Scan(&_id); err {
	case sql.ErrNoRows:
		// no category_name in records
		if canMakeNew {
			log.Println("making new category")
			sqlStatement := `
			INSERT INTO category (category_name)
			VALUES ($1) 
			RETURNING category_id`
			if err := db.QueryRow(sqlStatement, cName).Scan(&_id); err != nil {
				return DB_ERROR
			}
			log.Println("Made new category")
		} else {
			return WRONG_DATA
		}
	case nil:
		return _id
	default:
		return DB_ERROR
	}
	return _id // * if no rows were found sent new made id
}

func saveUser(user *RegsiterDetail) int {
	var id int

	// log.Println("user name is ", user.Name)
	if user.Choice == "student" { // * add to student
		sqlStatement := `
		INSERT INTO student (student_name)
		VALUES ($1) 
		RETURNING student_id`

		// log.Println("sql after inserting into student")
		if err := db.QueryRow(sqlStatement, user.Name).Scan(&id); err != nil {
			return DB_ERROR
		}
		// log.Println("sql after scanning id")

	} else { // * add to mentor
		// log.Println("userCategory is ", user.Category)
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

func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var user RegsiterDetail
	if err := decoder.Decode(&user); err != nil {
		// bad request
		badRequest(w, r)
		return
	} else {
		log.Println(user.Name)
		log.Println(user.Choice)
		log.Println(user.Category)
		if (user.Choice != "student" && user.Choice != "mentor") || len(user.Name) == 0 {
			// json not constructed correct, no choice taken
			badRequest(w, r)
			return
		} else {
			_id := saveUser(&user)
			log.Println("id is: ", _id)
			if _id == DB_ERROR {
				// return server error
				log.Println("DB ERROR")
				serverError(w, r)
				return
			} else if _id == WRONG_DATA {
				// * means that no category name inputted after selection of mentor
				log.Println("BAD REQUEST")
				badRequest(w, r)
				return
			} else {
				// * registered
				log.Println("registered")
			}
		}
	}
}

func selectMentor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var d DuoIdentity
	if err := decoder.Decode(&d); err != nil {
		// server error
		serverError(w, r)
		return
	}

	if d.MentorId == 0 || d.StudentId == 0 {
		badRequest(w, r)
		return
	}

	// insert into relation table
	sqlStatement := `
	INSERT INTO relation (student_id, mentor_id)
	VALUES ($1, $2)`

	if err := db.QueryRow(sqlStatement, d.StudentId, d.MentorId); err != nil {
		// send database error
		serverError(w, r)
		return
	}

	// send success message
	fmt.Println("success")

}

func viewMentors(w http.ResponseWriter, r *http.Request, categoryName string) {
	w.Header().Set("content-type", "application/json")
	// * fetch the category id
	_id := getCategoryId(categoryName, false)
	fmt.Println("value of getCategoryID is ", _id)
	if _id == WRONG_DATA {
		// wrong data in query
		badRequest(w, r)
		return

	} else if _id == DB_ERROR {
		// database error
		serverError(w, r)
		return
	} else {
		// * fetch all the data via categoryId from Mentor table and send it
		// var mentors []Person
		sqlStatement := `SELECT mentor_id,mentor_name
		FROM mentor
		WHERE category_id=$1`

		rows, err := db.Query(sqlStatement, _id)
		if err != nil {
			// DB error
			serverError(w, r)
			return
		}

		// * called because in case of error in for rows.Next(), Close() is not called automatically
		defer rows.Close()

		var mentors []Person

		for rows.Next() {
			var mentor Person
			err = rows.Scan(&mentor.ID, &mentor.Name)
			if err != nil {
				serverError(w, r)
				return
			}
			mentors = append(mentors, mentor)
		}
		err = rows.Err()
		if err != nil {
			serverError(w, r)
			return
		}

		// return the list of mentors via json
		for _, mentor := range mentors {
			fmt.Println(mentor.ID, mentor.Name)
		}

		// make json and send
		dataJson, e := json.Marshal(mentors)
		if e != nil {
			// server error
			serverError(w, r)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write(dataJson)
	}
}

func viewStudents(w http.ResponseWriter, r *http.Request, mentorId string) {
	w.Header().Set("content-type", "application/json")
	fmt.Println("before query")
	sqlStatement := `SELECT s.student_id, s.student_name
	FROM relation r
	INNER JOIN student s
	ON s.student_id = r.student_id
	WHERE mentor_id=$1`

	rows, err := db.Query(sqlStatement, mentorId)
	if err != nil {
		// DB error
		serverError(w, r)
		return
	}
	fmt.Println("after query")
	// * called because in case of error in for rows.Next(), Close() is not called automatically
	defer rows.Close()

	var students []Person

	for rows.Next() {
		var student Person
		err = rows.Scan(&student.ID, &student.Name)
		if err != nil {
			serverError(w, r)
			return
		}
		students = append(students, student)
	}

	fmt.Println("all rows parsed")

	err = rows.Err()
	if err != nil {
		serverError(w, r)
		return
	}

	fmt.Println("no error while parsing")

	// return the list of students via json
	for _, student := range students {
		fmt.Println(student.ID, student.Name)
	}

	// make json and send
	dataJson, e := json.Marshal(students)
	if e != nil {
		// server error
		serverError(w, r)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(dataJson)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	userReq := r.URL.Query().Get("choice")
	if userReq == "mentor" {
		categoryName := r.URL.Query().Get("cname")
		if len(categoryName) == 0 {
			// bad query values
			badRequest(w, r)
			return
		} else {
			viewMentors(w, r, categoryName)
		}
	} else if userReq == "student" {
		mentorId := r.URL.Query().Get("mentorid")
		viewStudents(w, r, mentorId)
	} else {
		// bad request/bad query values
		badRequest(w, r)
		return
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "not found"}`))
}

func serverError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"error": "Server Error"}`))
}

func badRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"error": "Wrong Parameter/Query"}`))
}

// * DB connections
func Connect() *sql.DB {
	// * connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// * we are now passing the postgres as the sql type and conn string
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

func main() {

	db = Connect()
	fmt.Println("Connection Made")
	defer db.Close()

	//  POST @param: choice(mentor/student), category(string, if mentor) @return: id
	http.HandleFunc("/register", register)

	//  GET @query: choice,cname[choice:mentor],mentorid[choice:student] @return: list of mentor_details/list of student_details
	http.HandleFunc("/view", viewHandler)

	//  POST @param: mentor_id, student_id @return: success/failure
	http.HandleFunc("/select", selectMentor)

	if err := http.ListenAndServe(":5001", nil); err != nil {
		log.Fatal(err)
	}
}
