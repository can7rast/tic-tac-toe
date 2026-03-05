package datasource

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func (db *DB) Close() {
	db.Pool.Close()
}
