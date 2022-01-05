package main

import (
	"fmt"
	"log"
	"net/http"

	ctrl "github.com/deus-oc/mentorapi/controllers"
	dbase "github.com/deus-oc/mentorapi/dbase"
)

func main() {

	db := dbase.Connect()
	fmt.Println("Connection Made")
	defer db.Close()

	//  POST @param: name, choice(mentor/student), category_name(string, if mentor) @return: id
	http.HandleFunc("/register", ctrl.Register)

	//  GET @query: requestby,cname[requestby:student],mentorid[requestby:mentor] @return: list of mentor_details/list of student_details
	http.HandleFunc("/view", ctrl.ViewHandler)

	//  POST @param: mentor_id, student_id @return: success/failure
	http.HandleFunc("/select", ctrl.SelectMentor)

	if err := http.ListenAndServe(":5001", nil); err != nil {
		log.Fatal(err)
	}
}
