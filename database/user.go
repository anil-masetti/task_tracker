package database

import (
	"fmt"
	"task-tracker/models"
	"time"
)

func CreateUser(user *models.User) error {
	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
	err := db.QueryRow(query, user.Username, user.Password).Scan(&user.ID)
	return err
}

func StoreLoginToken(userID int64, token string, expiresAt time.Time) error {
	// Insert the token along with the user ID and expiration time
	_, err := db.Exec(`
		INSERT INTO logins (user_id, jwt_token, created_at,expires_at) 
		VALUES ($1, $2, $3, $4)`,
		userID, token, time.Now(), expiresAt)

	if err != nil {
		return fmt.Errorf("failed to store login token: %w", err)
	}
	return nil
}
