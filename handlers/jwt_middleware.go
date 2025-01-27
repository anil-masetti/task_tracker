package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func ValidateJWT(next http.HandlerFunc, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		token := r.Header.Get("token")
		if token == "" {
			RespondWithError(w, http.StatusUnauthorized, "Token is required")
			return

		}

		// Verify the token exists in the database and is valid
		var userID int64

		var expiresAt time.Time

		err := db.QueryRow(
			"SELECT user_id, expires_at FROM logins WHERE jwt_token = $1",
			token,
		).Scan(&userID, &expiresAt)
		fmt.Println(userID, expiresAt.Local().Unix(), time.Now().Local().Unix())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {

				RespondWithError(w, http.StatusBadRequest, "Invalid or unrecognized token")

				return
			}
			RespondWithError(w, http.StatusInternalServerError, "Error Verifying Token")
			return
		}

		// Check if the token is expired
		if time.Now().Unix() > expiresAt.Unix() {
			fmt.Println("Token Expired")
			RespondWithError(w, http.StatusUnauthorized, "Please Login Again...")
			return
		}
		next(w, r)

	}
}
