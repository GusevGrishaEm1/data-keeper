package postgres

import (
	"context"
	"fmt"

	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	DB *pgxpool.Pool
}

func NewPostgresDB(ctx context.Context, c config.Config) (*DB, error) {
	postgresURL := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		c.Postgres.User, c.Postgres.Password, c.Postgres.Host, c.Postgres.Port, c.Postgres.DB,
	)
	configPool, err := pgxpool.ParseConfig(postgresURL)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, configPool)
	if err != nil {
		return nil, err
	}

	return &DB{DB: pool}, nil
}
