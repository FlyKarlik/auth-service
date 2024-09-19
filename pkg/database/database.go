package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectionPostgresSQLX(databaseURL string) (*sqlx.DB, error) {
	conStr := fmt.Sprint(databaseURL)
	db, err := sqlx.Open("postgres", conStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

var DB *gorm.DB

func ConnectionPostgresGormDB(dbDsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbDsn))
	if err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}

func GetDB() *gorm.DB {
	return DB
}
