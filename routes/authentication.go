package routes

import (
	"database/sql"
	"net/http"
	"task-tracker/handlers"

	"github.com/gorilla/mux"
)

// create login endpoint using mux

func AuthRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegistrationHandler(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		handlers.LogoutHandler(w, r, db)
	}).Methods("POST")
}
