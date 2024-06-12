package repo

import (
	"context"

	"github.com/GusevGrishaEm1/data-keeper/internal/entity"
	"github.com/GusevGrishaEm1/data-keeper/internal/infrastructure/repository/postgres"
)

type userFileRepo struct {
	*postgres.PostgresDB
}

func NewUserFileRepo(db *postgres.PostgresDB) *userFileRepo {
	return &userFileRepo{db}
}

func (s *userFileRepo) Insert(ctx context.Context, data entity.UserFile) error {
	query := `insert into "user_file" (uuid, content, created_at, created_by) values ($1, $2, $3, $4)`
	_, err := s.DB.Exec(ctx, query, data.UUID, data.Content, data.CreatedAt, data.CreatedBy)
	return err
}

func (s *userFileRepo) Delete(ctx context.Context, user string, uuid string) error {
	query := `delete from "user_file" where uuid::text = $1 and created_by = $2`
	_, err := s.DB.Exec(ctx, query, uuid, user)
	return err
}

func (s *userFileRepo) GetByUUID(ctx context.Context, user string, uuid string) (*entity.UserFile, error) {
	query := `select uuid, content, created_at, created_by from "user_file" where uuid::text = $1 and created_by = $2`
	row := s.DB.QueryRow(ctx, query, uuid, user)
	data := &entity.UserFile{}
	err := row.Scan(&data.UUID, &data.Content, &data.CreatedAt, &data.CreatedBy)
	return data, err
}
