package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/config"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/controller/http"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/repository/postgres"

	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// load config
	c, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	// init ctx
	ctx := context.Background()

	// auth service client
	authServer := c.AuthService.Host + ":" + strconv.Itoa(c.AuthService.Port)
	authn, err := grpc.NewClient(
		authServer,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithIdleTimeout(c.AuthService.Timeout*time.Second),
	)
	if err != nil {
		panic(err)
	}
	logger := logger()

	postgresURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		c.Postgres.User, c.Postgres.Password, c.Postgres.Host, c.Postgres.Port, c.Postgres.DB,
	)
	logger.Info("postgres URL: " + postgresURL)

	// postgres db
	db, err := postgres.NewPostgresDB(ctx, *c)
	if err != nil {
		panic(err)
	}

	// run migration
	err = migration(err, c)
	if err != nil {
		panic(err)
	}

	// start server
	err = http.StartServer(*c, logger, authn, db)
	if err != nil {
		panic(err)
	}
}

// migration migrate data
func migration(err error, c *config.Config) error {
	postgresURL := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		c.Postgres.User, c.Postgres.Password, c.Postgres.Host, c.Postgres.Port, c.Postgres.DB,
	)
	connToMigrate, err := sql.Open("pgx", postgresURL)
	if err != nil {
		return err
	}
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(connToMigrate, "migrations"); err != nil {
		return err
	}
	err = connToMigrate.Close()
	if err != nil {
		return err
	}
	return nil
}

// logger to log in console
func logger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}
