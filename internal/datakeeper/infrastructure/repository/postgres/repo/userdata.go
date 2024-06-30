package repo

import (
	"context"

	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/entity"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/repository/postgres"
)

type DataRepo struct {
	db *postgres.DB
}

// NewDataRepo creates new data repository
func NewDataRepo(db *postgres.DB) *DataRepo {
	return &DataRepo{db}
}

// Insert insert new data for user
func (s *DataRepo) Insert(ctx context.Context, data entity.Data) error {
	query := `
	insert into user_data (uuid, content, content_type, created_at, created_by) 
	values ($1, $2, $3, $4, $5)
	`
	_, err := s.db.DB.Exec(ctx, query, data.UUID, data.Content, data.ContentType, data.CreatedAt, data.CreatedBy)
	if err != nil {
		return err
	}
	return nil
}

// Delete delete data for user
func (s *DataRepo) Delete(ctx context.Context, user string, uuid string) error {
	query := `delete from user_data where uuid::text = $1 and created_by = $2`
	_, err := s.db.DB.Exec(ctx, query, uuid, user)
	return err
}

// Update data for user
func (s *DataRepo) Update(ctx context.Context, data entity.Data) error {
	query := `
	update user_data 
	set content = $1
	where uuid::text = $2 and created_by = $3 and content_type = $4
	`
	_, err := s.db.DB.Exec(ctx, query, data.Content, data.UUID, data.CreatedBy, data.ContentType)
	return err
}

// GetByUser Get data by user and content type
func (s *DataRepo) GetByUser(ctx context.Context, user string, contentType entity.ContentType) ([]*entity.Data, error) {
	query := `
	select uuid, content, content_type
	from user_data
	where created_by = $1 and content_type = $2`
	rows, err := s.db.DB.Query(ctx, query, user, contentType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.Data
	for rows.Next() {
		var data entity.Data
		err := rows.Scan(&data.UUID, &data.Content, &data.ContentType)
		if err != nil {
			return nil, err
		}
		result = append(result, &data)
	}

	return result, rows.Err()
}

// GetByUUID Get data by user and content type and uuid
func (s *DataRepo) GetByUUID(ctx context.Context, user string, uuid string) (*entity.Data, error) {
	query := `
	select uuid, content, content_type, created_by
	from user_data
    where created_by = $1 and uuid::text = $2`
	row := s.db.DB.QueryRow(ctx, query, user, uuid)
	data := &entity.Data{}
	err := row.Scan(&data.UUID, &data.Content, &data.ContentType, &data.CreatedBy)
	return data, err
}
