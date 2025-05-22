package sqlconnect

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gepestudy/go-rest-api/pkg/config"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", config.Db.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Minute * 5)

	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Connected to msyql!")

	return db, nil
}
