package repository

import "github.com/jackc/pgx/v4/pgxpool"

type DB struct {
	pool *pgxpool.Pool
}

func NewDBService(pool *pgxpool.Pool) *DB {
	return &DB{
		pool: pool,
	}
}
