package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
)

const (
	DSN = "postgres://%s:%s@%s:%s/%s?sslmode=disable"
)

func getDSN() string {
	user, password, dbName, port, host := os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_HOST")
	return fmt.Sprintf(DSN, user, password, host, port, dbName)
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func MakeMigrations() error {
	db, err := sql.Open("postgres", getDSN())
	if err != nil {
		return err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "internal/db/migrations"); err != nil {
		panic(err)
	}
	return db.Close()
}

func InitPool() (*pgxpool.Pool, error) {
	dsn := getDSN()
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Errorf("Unable to parse DATABASE_URL=%s with error: %s", dsn, err.Error())
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Errorf("Unable to create connection pool: %s", err.Error())
		return nil, err
	}

	return pool, nil
}
