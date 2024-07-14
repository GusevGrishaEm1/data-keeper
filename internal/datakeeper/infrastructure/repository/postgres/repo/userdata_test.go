package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/config"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/repository/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var repo *DataRepo

func TestMain(m *testing.M) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}

	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %v", err)
		}
	}()

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get container host: %v", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("failed to get mapped port: %v", err)
	}

	dbURI := fmt.Sprintf("postgres://user:password@%s:%s/testdb?sslmode=disable", host, port.Port())
	configPool, err := pgxpool.ParseConfig(dbURI)
	if err != nil {
		log.Fatalf("failed to get config pool: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, configPool)
	if err != nil {
		log.Fatalf("failed to get new pool: %v", err)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping postgres: %v", err)
	}
	portInt, err := strconv.Atoi(port.Port())
	if err != nil {
		log.Fatalf("failed to convert port to int: %v", err)
	}
	postgresC := &config.Postgres{
		Host:     host,
		User:     "user",
		Password: "password",
		Port:     portInt,
		DB:       "testdb",
	}
	err = migration(config.Config{Postgres: *postgresC})
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	dbpostgres, err := postgres.NewPostgresDB(context.TODO(), config.Config{Postgres: *postgresC})
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	repo = NewDataRepo(dbpostgres)

	m.Run()
}

func migration(c config.Config) error {
	postgresURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		c.Postgres.User, c.Postgres.Password, c.Postgres.Host, c.Postgres.Port, c.Postgres.DB,
	)
	connToMigrate, err := sql.Open("pgx", postgresURL)
	if err != nil {
		return err
	}
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(connToMigrate, "../../../../../../migrations"); err != nil {
		return err
	}
	err = connToMigrate.Close()
	if err != nil {
		return err
	}
	return nil
}

func clearTable(ctx context.Context) {
	_, err := repo.db.DB.Exec(ctx, `DELETE FROM "user_data"`)
	if err != nil {
		log.Fatalf("failed to clear table: %v", err)
	}
}

func TestInsert(t *testing.T) {
	ctx := context.Background()
	defer clearTable(ctx)

	data := entity.Data{
		UUID:        uuid.New().String(),
		Content:     []byte("test-content"),
		ContentType: "text",
		CreatedAt:   time.Now(),
		CreatedBy:   "test-user",
	}
	err := repo.Insert(ctx, data)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	defer clearTable(ctx)

	newUUID := uuid.New().String()
	data := entity.Data{
		UUID:        newUUID,
		Content:     []byte("test-content"),
		ContentType: "text",
		CreatedAt:   time.Now(),
		CreatedBy:   "test-user",
	}
	err := repo.Insert(ctx, data)
	assert.NoError(t, err)

	err = repo.Delete(ctx, "test-user", newUUID)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	defer clearTable(ctx)

	newUUID := uuid.New().String()
	data := entity.Data{
		UUID:        newUUID,
		Content:     []byte("test-content"),
		ContentType: "text",
		CreatedAt:   time.Now(),
		CreatedBy:   "test-user",
	}

	err := repo.Insert(ctx, data)
	assert.NoError(t, err)

	updatedData := entity.Data{
		UUID:        newUUID,
		Content:     []byte("updated-content"),
		ContentType: "text",
		CreatedAt:   time.Now(),
		CreatedBy:   "test-user",
	}
	err = repo.Update(ctx, updatedData)
	assert.NoError(t, err)
}

func TestGetAllByUser(t *testing.T) {
	ctx := context.Background()
	defer clearTable(ctx)

	newUUID := uuid.New().String()
	data := entity.Data{
		UUID:        newUUID,
		Content:     []byte("test-content"),
		ContentType: "text",
		CreatedAt:   time.Now(),
		CreatedBy:   "test-user",
	}
	err := repo.Insert(ctx, data)
	assert.NoError(t, err)

	userData, err := repo.GetByUser(ctx, "test-user", "text")
	assert.NoError(t, err)

	assert.NotNil(t, userData)
	assert.Equal(t, 1, len(userData))
	assert.Equal(t, newUUID, userData[0].UUID)
	assert.Equal(t, "test-content", string(userData[0].Content))
	assert.Equal(t, entity.ContentType("text"), userData[0].ContentType)
}
