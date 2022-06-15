package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

func InitPool(ctx context.Context) (*pgxpool.Pool, error) {
	dbUser, dbPassword, dbName, dbPort, dbHost :=
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_HOST")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Errorf("Unable to parse DATABASE_URL=%s with error: %s", dsn, err.Error())
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		log.Errorf("Unable to create connection pool: %s", err.Error())
		return nil, err
	}

	return pool, nil
}
