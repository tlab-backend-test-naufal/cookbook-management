package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

const defaultRestPort = "8080"

func RestPort() string {
	port := os.Getenv("REST_PORT")
	if port == "" {
		port = defaultRestPort
	}

	return port
}

func BuildPostgres() (*sqlx.DB, error) {
	dataSourceURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_SSLMODE"))

	db, err := sqlx.Connect("postgres", dataSourceURL)
	if err != nil {
		return nil, err
	}

	if maxOpenConns, err := strconv.Atoi(os.Getenv("POSTGRES_MAX_OPEN_CONNS")); err == nil {
		db.SetMaxOpenConns(maxOpenConns)
	}
	if maxIdleConns, err := strconv.Atoi(os.Getenv("POSTGRES_MAX_IDLE_CONNS")); err == nil {
		db.SetMaxIdleConns(maxIdleConns)
	}
	if maxIdleTime, err := strconv.Atoi(os.Getenv("POSTGRES_MAX_IDLE_TIME_MINUTES")); err == nil {
		db.SetConnMaxIdleTime(time.Duration(maxIdleTime) * time.Minute)
	}
	if maxLifeTime, err := strconv.Atoi(os.Getenv("POSTGRES_MAX_LIFE_TIME_HOURS")); err == nil {
		db.SetConnMaxLifetime(time.Duration(maxLifeTime) * time.Hour)
	}

	return db, nil
}
