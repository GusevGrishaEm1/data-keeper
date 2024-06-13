package repo

import (
	"context"

	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/entity"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/repository/postgres"
)

type dataRepo struct {
	*postgres.PostgresDB
}

// NewDataRepo creates new data repository
func NewDataRepo(db *postgres.PostgresDB) *dataRepo {
	return &dataRepo{db}
}

// Insert insert new data for user
func (s *dataRepo) Insert(ctx context.Context, data entity.Data) error {
	query := `
	insert into "data" (uuid, content, content_type, created_at, created_by) 
	values ($1, $2, $3, $4, $5)
	`
	_, err := s.DB.Exec(ctx, query, data.UUID, data.Content, data.ContentType, data.CreatedAt, data.CreatedBy)
	if err != nil {
		return err
	}
	return nil
}

// Delete delete data for user
func (s *dataRepo) Delete(ctx context.Context, user string, uuid string) error {
	query := `delete from "data" where uuid::text = $1 and created_by = $2`
	_, err := s.DB.Exec(ctx, query, uuid, user)
	return err
}

// Update data for user
func (s *dataRepo) Update(ctx context.Context, data entity.Data) error {
	query := `
	update "data" 
	set content = $1
	where uuid::text = $2 and created_by = $3 and content_type = $4`
	_, err := s.DB.Exec(ctx, query, data.Content, data.UUID, data.CreatedBy, data.ContentType)
	return err
}

// Get data by user and content type
func (s *dataRepo) GetByUser(ctx context.Context, user string, contentType entity.ContentType) ([]*entity.Data, error) {
	query := `select uuid, content, content_type from "data" where created_by = $1 and content_type = $2`
	rows, err := s.DB.Query(ctx, query, user, contentType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var datas []*entity.Data
	for rows.Next() {
		var data entity.Data
		err := rows.Scan(&data.UUID, &data.Content, &data.ContentType)
		if err != nil {
			return nil, err
		}
		datas = append(datas, &data)
	}

	return datas, rows.Err()
}

// Get data by user and content type and uuid
func (s *dataRepo) GetByUUID(ctx context.Context, user string, uuid string, contentType entity.ContentType) (*entity.Data, error) {
	query := `select uuid, content, content_type from "data" where created_by = $1 and uuid::text = $2 and content_type = $3`
	row := s.DB.QueryRow(ctx, query, user, uuid, contentType)
	data := &entity.Data{}
	err := row.Scan(&data.UUID, &data.Content, &data.ContentType)
	return data, err
}
