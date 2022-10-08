package repository

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresRepository struct {
	db *sql.DB
}

func (repo *PostgresRepository) Connect() error {
	dsn := "host=localhost user=postgres password=postgres dbname=go_slack_unipos_development port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Hour)

	err = db.Ping()
	if err != nil {
		return err
	}

	repo.db = db
	return nil
}

func (repo *PostgresRepository) Close() {
	repo.db.Close()
}
