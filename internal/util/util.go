package util

import (
	"crypto/rand"
	"encoding/hex"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/vincer2040/circles/internal/db"
)

type PageData struct {
	Route string
}

type CirclesContext struct {
	echo.Context
	DB    db.CirclesDB
	Store *sessions.CookieStore
}

func DbURL() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	db := os.Getenv("DBURL")
	return db, nil
}

func GenerateRandomKey(length int) (string, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}
