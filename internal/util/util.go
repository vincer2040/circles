package util

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/vincer2040/circles/internal/db"
)

type PageData struct {
	Route string
}

type CirclesContext struct {
	echo.Context
    DB db.CirclesDB
}

func DbURL() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	db := os.Getenv("DBURL")
	return db, nil
}
