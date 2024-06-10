package postgres

import (
	"context"

	"github.com/GusevGrishaEm1/data-keeper/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	DB *pgxpool.Pool
}

func NewPostgresDB(ctx context.Context, config config.Config) (*PostgresDB, error) {
	configPool, err := pgxpool.ParseConfig(config.PostgresDB.URL)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, configPool)
	if err != nil {
		return nil, err
	}
	return &PostgresDB{DB: pool}, nil
}
