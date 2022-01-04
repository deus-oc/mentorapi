package main

import (
	"database/sql"
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
	dbname   = "demo"
)

// * controllers
func register(w http.ResponseWriter, r *http.Request) {
}

func viewMentor(w http.ResponseWriter, r *http.Request) {

}

func selectMentor(w http.ResponseWriter, r *http.Request) {

}

func viewStudent(w http.ResponseWriter, r *http.Request) {

}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "not found"}`))
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

	db := Connect()
	fmt.Println("Connection Made")
	defer db.Close()

	//  POST @param: mentor or Student (bool), category(string, if mentor)/null @return: success/failure
	http.HandleFunc("/register", register)
	//  GET @param: category_name @return: list of mentor_name, id
	http.HandleFunc("/viewmentor", viewMentor)
	//  POST @param: mentor_name or mentor_id, student_id or student_name @return: success/failure
	http.HandleFunc("/selectmentor", selectMentor)
	//  GET @param: mentor_id or mentor_name @return: list of students
	http.HandleFunc("/viewstudent", viewStudent)

	if err := http.ListenAndServe(":5001", nil); err != nil {
		log.Fatal(err)
	}
}
