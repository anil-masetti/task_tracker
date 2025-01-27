package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"task-tracker/database"
	"task-tracker/models"
	"time"
)

func RegistrationHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid user data")
		return
	}

	// Validate user data
	if err := ValidateStruct(user); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Hash password before storing
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	user.Password = hashedPassword

	// Create user in database
	err = database.CreateUser(&user)
	if err != nil {
		fmt.Println(err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate JWT token
	token, err := GenerateJWT(user.ID)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	expirationTime := time.Now().Add(72 * time.Hour)

	// Store the token in the logins table
	err = database.StoreLoginToken(user.ID, token, expirationTime)

	RespondWithJSON(w, http.StatusCreated, map[string]string{"token": token})

}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid user data")
		return
	}

	// Validate user data
	if err := ValidateStruct(user); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var storedUser models.User

	err := db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", user.Username).Scan(&storedUser.ID, &storedUser.Username, &storedUser.Password)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}
	if err := CheckPasswordHash(user.Password, storedUser.Password); err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}
	token, err := GenerateJWT(user.ID)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	expirationTime := time.Now().Add(72 * time.Hour)

	// Store the token in the logins table
	err = database.StoreLoginToken(storedUser.ID, token, expirationTime)
	if err != nil {
		fmt.Println(err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to store token")
		return
	}
	RespondWithJSON(w, http.StatusCreated, map[string]string{"token": token})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	token := r.Header.Get("token")
	if token == "" {
		RespondWithError(w, http.StatusBadRequest, "Token is missing")
		return
	}
	var userID int64
	err := db.QueryRow("SELECT user_id FROM logins WHERE jwt_token=$1", token).Scan(&userID)
	if err == sql.ErrNoRows {
		RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	} else if err != nil {
		fmt.Println(err)
		RespondWithError(w, http.StatusInternalServerError, "Database query failed")
		return
	}

	// Update the token's expiry timestamp to the current time
	_, err = db.Exec("UPDATE logins SET expires_at = $1 WHERE jwt_token = $2", time.Now(), token)
	if err != nil {
		fmt.Println(err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to update token expiry")
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully logged out"})
}
