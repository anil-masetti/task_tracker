package models

type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	AssigneeID  int64  `json:"assignee_id"`
}
