package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"task-tracker/models"
	"time"

	"github.com/gorilla/mux"
)

func CreateTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid task data")
		return
	}

	// Validate the task data
	if task.Title == "" {
		RespondWithError(w, http.StatusBadRequest, "Task title is required")
		return
	}
	task.Status = "todo"
	if task.AssigneeID == 0 {
		RespondWithError(w, http.StatusBadRequest, "Assignee ID is required")
		return
	}

	// Verify if the assignee exists
	var assigneeExists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)", task.AssigneeID).Scan(&assigneeExists)
	if err != nil || !assigneeExists {
		RespondWithError(w, http.StatusBadRequest, "Invalid assignee ID")
		return
	}

	// Insert the task into the database
	query := `INSERT INTO tasks (title, description, status, assignee_id,created_at) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var taskID int64
	err = db.QueryRow(query, task.Title, task.Description, task.Status, task.AssigneeID, time.Now()).Scan(&taskID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	// Return the created task's details
	task.ID = taskID
	RespondWithJSON(w, http.StatusCreated, task)
}

func GetAllTasks(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	query := `
		SELECT t.id, t.title, t.description, t.status, t.created_at, t.updated_at, u.id AS assignee_id, u.username AS assigned_user
		FROM tasks t
		LEFT JOIN users u ON t.assignee_id = u.id
		ORDER BY t.created_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve tasks")
		return
	}
	defer rows.Close()

	var tasks []map[string]interface{}

	// Iterate through the rows and construct the task list
	for rows.Next() {
		var (
			id               int64
			title            string
			description      sql.NullString
			status           string
			createdAt        string
			updatedAt        sql.NullTime
			assignedUserID   int64
			assignedUsername string
		)

		// Scan each row into variables
		err := rows.Scan(&id, &title, &description, &status, &createdAt, &updatedAt, &assignedUserID, &assignedUsername)
		if err != nil {
			fmt.Println(err)
			RespondWithError(w, http.StatusInternalServerError, "Failed to parse tasks")
			return
		}

		// Create a task map
		task := map[string]interface{}{
			"id":          id,
			"title":       title,
			"description": nilIfNullString(description),
			"status":      status,
			"created_at":  createdAt,
			"updated_at":  nilIfNullTime(updatedAt),
		}

		// Add assigned user details if available
		if assignedUserID != 0 {
			task["assignee_id"] = assignedUserID
			task["assigned_user"] = assignedUsername
		}

		tasks = append(tasks, task)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve tasks")
		return
	}

	// Respond with the list of tasks
	RespondWithJSON(w, http.StatusOK, tasks)
}

func GetTaskOfAnUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Extract the user_id from the path parameter
	userID := mux.Vars(r)["user_id"]

	// Query to get tasks assigned to a specific user
	query := `
		SELECT t.id, t.title, t.description, t.status, t.created_at, t.updated_at
		FROM tasks t
		WHERE t.assignee_id = $1
		ORDER BY t.created_at DESC
	`

	// Execute the query with userID as parameter
	rows, err := db.Query(query, userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve tasks")
		return
	}
	defer rows.Close()

	// Slice to hold all tasks for the user
	var tasks []map[string]interface{}

	// Iterate through the rows and construct the task list
	for rows.Next() {
		var (
			id          int64
			title       string
			description sql.NullString
			status      string
			createdAt   string
			updatedAt   sql.NullTime
		)

		// Scan each row into variables
		err := rows.Scan(&id, &title, &description, &status, &createdAt, &updatedAt)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to parse tasks")
			return
		}

		// Create a task map
		task := map[string]interface{}{
			"id":          id,
			"title":       title,
			"description": nilIfNullString(description),
			"status":      status,
			"created_at":  createdAt,
			"updated_at":  nilIfNullTime(updatedAt),
		}

		// Add the task to the list
		tasks = append(tasks, task)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve tasks")
		return
	}

	// Respond with the list of tasks
	RespondWithJSON(w, http.StatusOK, tasks)
}

func UpdateTaskStatus(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	taskID := mux.Vars(r)["task_id"]

	// Define a struct to hold the status from the request body
	var request struct {
		Status string `json:"status"`
	}

	// Decode the request body to get the new status
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	// Validate that the status is one of the valid options
	validStatuses := map[string]bool{
		"todo":        true,
		"in_progress": true,
		"completed":   true,
	}

	if !validStatuses[request.Status] {
		RespondWithError(w, http.StatusBadRequest, "Invalid status")
		return
	}

	// Update the task's status in the database
	query := `
		UPDATE tasks
		SET status = $1
		WHERE id = $2
		RETURNING id, title, status
	`

	var id int64
	var title, status string

	err := db.QueryRow(query, request.Status, taskID).Scan(&id, &title, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusNotFound, "Task not found")
		} else {
			RespondWithError(w, http.StatusInternalServerError, "Failed to update task status")
		}
		return
	}

	// Respond with the updated task details
	task := map[string]interface{}{
		"id":     id,
		"title":  title,
		"status": status,
	}
	RespondWithJSON(w, http.StatusOK, task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Extract the task_id from the URL path parameter
	taskID := mux.Vars(r)["task_id"]

	// Check if the task exists by querying it
	var title string
	err := db.QueryRow("SELECT title FROM tasks WHERE id = $1", taskID).Scan(&title)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusNotFound, "Task not found")
		} else {
			RespondWithError(w, http.StatusInternalServerError, "Failed to check task existence")
		}
		return
	}

	// Task exists, proceed with deletion
	_, err = db.Exec("DELETE FROM tasks WHERE id = $1", taskID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	// Respond with success
	RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Task deleted successfully",
	})
}
