package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() *sql.DB {
	connStr := "host=172.18.0.2 user=postgres password=postgres dbname=task_tracker sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to verify connection: %v", err)
	}
	fmt.Println("Database connected successfully!")
	return db
}
