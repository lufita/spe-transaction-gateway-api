package config

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"spe-trx-gateway/lookup"
	"time"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(ctx context.Context) (*DB, error) {
	url := lookup.DbConnString
	cfg, err := pgxpool.ParseConfig(url)
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
