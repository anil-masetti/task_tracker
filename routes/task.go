package routes

import (
	"database/sql"
	"net/http"
	"task-tracker/handlers"

	"github.com/gorilla/mux"
)

func TaskRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/task", handlers.ValidateJWT(func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateTask(w, r, db)
	}, db)).Methods("POST")

	router.HandleFunc("/task/{task_id}/status", handlers.ValidateJWT(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateTaskStatus(w, r, db)
	}, db)).Methods("PUT")

	router.HandleFunc("/tasks", handlers.ValidateJWT(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAllTasks(w, r, db)
	}, db)).Methods("GET")

	router.HandleFunc("/tasks/{user_id}", handlers.ValidateJWT(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetTaskOfAnUser(w, r, db)
	}, db)).Methods("GET")

	router.HandleFunc("/task/{task_id}", handlers.ValidateJWT(func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteTask(w, r, db)
	}, db)).Methods("DELETE")
}
