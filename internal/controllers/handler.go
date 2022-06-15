package controllers

import "github.com/jackc/pgx/v4/pgxpool"

type handler struct {
	pool *pgxpool.Pool
}

func NewHandler(pool *pgxpool.Pool) *handler {
	return &handler{
		pool: pool,
	}
}
