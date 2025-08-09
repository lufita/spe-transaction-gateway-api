package config

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"time"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(ctx context.Context) (*DB, error) {
	dsn := os.Getenv("DB_CONNECTION_STRING")
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	cfg.MaxConnLifetime = time.Hour
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &DB{Pool: pool}, nil
}

func (d *DB) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}
