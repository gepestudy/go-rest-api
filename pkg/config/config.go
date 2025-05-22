package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type app struct {
	Port string
}

type db struct {
	DSN string
}

var (
	App app
	Db  db
)

func Load() error {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	switch {
	case os.Getenv("APP_PORT") == "":
		return fmt.Errorf("APP_PORT is not set")
	case os.Getenv("DB_DSN") == "":
		return fmt.Errorf("DB_DSN is not set")
	}

	App = app{
		Port: os.Getenv("APP_PORT"),
	}
	Db = db{
		DSN: os.Getenv("DB_DSN"),
	}
	return nil
}
