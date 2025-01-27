package main

import (
	"fmt"
	"net/http"
	"task-tracker/database"
	"task-tracker/routes"

	"github.com/gorilla/mux"
)

func main() {
	db := database.InitDB()
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to the Task Tracker")
	})

	routes.AuthRoutes(router, db)
	routes.TaskRoutes(router, db)
	http.ListenAndServe(":8000", router)

}
