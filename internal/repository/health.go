package repository

import "context"

func (db *DB) Check(ctx context.Context) string {
	err := db.pool.Ping(ctx)
	if err != nil {
		return err.Error()
	}
	return "UP"
}
