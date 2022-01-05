package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func saveUser(user *RegsiterDetail) int {
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

func register(w http.ResponseWriter, r *http.Request) {
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

func selectMentor(w http.ResponseWriter, r *http.Request) {
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

	default:
		serverError(w, r)
		return
	}
}

func getList(sqlStatement string, _id int) ([]Person, error) {
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
	fmt.Println("before query")
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

func viewHandler(w http.ResponseWriter, r *http.Request) {
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

	//  POST @param: name, choice(mentor/student), category_name(string, if mentor) @return: id
	http.HandleFunc("/register", register)

	//  GET @query: requestby,cname[requestby:student],mentorid[requestby:mentor] @return: list of mentor_details/list of student_details
	http.HandleFunc("/view", viewHandler)

	//  POST @param: mentor_id, student_id @return: success/failure
	http.HandleFunc("/select", selectMentor)

	if err := http.ListenAndServe(":5001", nil); err != nil {
		log.Fatal(err)
	}
}
