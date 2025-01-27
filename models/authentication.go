package models

// User represents a user in the system
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Token string `json:"token"`
}
