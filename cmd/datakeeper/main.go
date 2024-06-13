package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/config"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/controller/http"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/repository/postgres"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var path string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()
	config, err := config.LoadConfig(path)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	// auth service client
	authconn, err := grpc.NewClient(
		config.URLAUTH,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithIdleTimeout(time.Duration(config.AuthService.Timeout)*time.Second),
	)
	if err != nil {
		panic(err)
	}
	// postgres db
	db, err := postgres.NewPostgresDB(ctx, *config)
	if err != nil {
		panic(err)
	}
	dbMig, err := sql.Open("pgx", config.URLDB)
	if err != nil {
		panic(err)
	}
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	if err := goose.Up(dbMig, "migrations"); err != nil {
		panic(err)
	}
	// start server
	err = http.StartServer(ctx, *config, logger(), authconn, db)
	if err != nil {
		panic(err)
	}
}

// info logger to log in console
func logger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}
