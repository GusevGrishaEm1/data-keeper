package repo

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/GusevGrishaEm1/data-keeper/internal/entity"
	"github.com/GusevGrishaEm1/data-keeper/internal/infrastructure/repository/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var repo *dataRepo

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

	_, err = pool.Exec(ctx, `
		create table if not exists "data" (
		uuid uuid primary key,
		content bytea not null,
		content_type varchar(255) not null,
		created_at timestamp not null,
		created_by varchar(255) not null
	);

	create index if not exists data_idx on "data" (created_by);
	`)

	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	repo = NewDataRepo(&postgres.PostgresDB{DB: pool})

	m.Run()
}

func clearTable() {
	_, err := repo.DB.Exec(context.TODO(), `DELETE FROM "data"`)
	if err != nil {
		log.Fatalf("failed to clear table: %v", err)
	}
}

func TestInsert(t *testing.T) {
	defer clearTable()
	ctx := context.Background()
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
	defer clearTable()
	ctx := context.Background()
	newuuid := uuid.New().String()
	data := entity.Data{
		UUID:        newuuid,
		Content:     []byte("test-content"),
		ContentType: "text",
		CreatedAt:   time.Now(),
		CreatedBy:   "test-user",
	}
	err := repo.Insert(ctx, data)
	assert.NoError(t, err)
	err = repo.Delete(ctx, "test-user", newuuid)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	defer clearTable()
	ctx := context.Background()
	newuuid := uuid.New().String()
	data := entity.Data{
		UUID:        newuuid,
		Content:     []byte("test-content"),
		ContentType: "text",
		CreatedAt:   time.Now(),
		CreatedBy:   "test-user",
	}
	err := repo.Insert(ctx, data)
	assert.NoError(t, err)
	updatedData := entity.Data{
		UUID:        newuuid,
		Content:     []byte("updated-content"),
		ContentType: "text",
		CreatedAt:   time.Now(),
		CreatedBy:   "test-user",
	}
	err = repo.Update(ctx, updatedData)
	assert.NoError(t, err)
}

func TestGetAllByUser(t *testing.T) {
	defer clearTable()
	ctx := context.Background()
	newuuid := uuid.New().String()
	data := entity.Data{
		UUID:        newuuid,
		Content:     []byte("test-content"),
		ContentType: "text",
		CreatedAt:   time.Now(),
		CreatedBy:   "test-user",
	}
	err := repo.Insert(ctx, data)
	assert.NoError(t, err)
	datas, err := repo.GetByUser(ctx, "test-user", "text")
	assert.NoError(t, err)
	assert.NotNil(t, datas)
	assert.Equal(t, 1, len(datas))
	assert.Equal(t, newuuid, datas[0].UUID)
	assert.Equal(t, "test-content", string(datas[0].Content))
	assert.Equal(t, entity.ContentType("text"), datas[0].ContentType)
}
