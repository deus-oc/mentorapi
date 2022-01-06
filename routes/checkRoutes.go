package routes

import (
	"net/http"

	ctrl "github.com/deus-oc/mentorapi/controllers"
)

func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ctrl.Register(w, r)
	default:
		ctrl.BadMethod(w, r)
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ctrl.ViewHandler(w, r)
	default:
		ctrl.BadMethod(w, r)
	}
}

func SelectMentor(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ctrl.SelectMentor(w, r)
	default:
		ctrl.BadMethod(w, r)
	}
}
