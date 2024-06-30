package repo

import (
	"context"

	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/entity"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/repository/postgres"
)

type UserFileRepo struct {
	db *postgres.DB
}

func NewUserFileRepo(db *postgres.DB) *UserFileRepo {
	return &UserFileRepo{db}
}

func (s *UserFileRepo) Insert(ctx context.Context, data entity.UserFile) error {
	query := `
	insert into file_repository (uuid, content, created_at, created_by)
	values ($1, $2, $3, $4)`
	_, err := s.db.DB.Exec(ctx, query, data.UUID, data.Content, data.CreatedAt, data.CreatedBy)
	return err
}

func (s *UserFileRepo) Delete(ctx context.Context, user string, uuid string) error {
	query := `
	delete from file_repository
    where uuid::text = $1 and created_by = $2`
	_, err := s.db.DB.Exec(ctx, query, uuid, user)
	return err
}

func (s *UserFileRepo) GetByUUID(ctx context.Context, user string, uuid string) (*entity.UserFile, error) {
	query := `
	select uuid, content, created_at, created_by
	from file_repository
	where uuid::text = $1 and created_by = $2`
	row := s.db.DB.QueryRow(ctx, query, uuid, user)
	data := &entity.UserFile{}
	err := row.Scan(&data.UUID, &data.Content, &data.CreatedAt, &data.CreatedBy)
	return data, err
}
