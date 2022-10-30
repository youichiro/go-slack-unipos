package repository

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func InitDB() (*sql.DB, error) {
	// TODO: リポジトリになってない
	dsn := "host=localhost user=postgres password=postgres dbname=go_slack_unipos_development port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Hour)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
