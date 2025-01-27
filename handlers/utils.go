package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("task_tracker")

func GenerateJWT(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ValidateStruct(data interface{}) error {
	validate := validator.New()
	return validate.Struct(data)
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// RespondWithError sends a JSON error response
func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	RespondWithJSON(w, statusCode, map[string]string{"error": message})
}

func CheckPasswordHash(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func nilIfNullString(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}

func nilIfNullTime(ns sql.NullTime) interface{} {
	if ns.Valid {
		return ns.Time
	}
	return nil
}
