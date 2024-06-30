package logpass

import (
	"context"
	"encoding/json"
	"time"

	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/entity"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/controller/http/handlers"
	"github.com/GusevGrishaEm1/data-keeper/internal/lib"
	"github.com/google/uuid"
)

type Repo interface {
	Insert(ctx context.Context, data entity.Data) error
	Update(ctx context.Context, data entity.Data) error
	Delete(ctx context.Context, user string, uuid string) error
	GetByUUID(ctx context.Context, user string, uuid string) (*entity.Data, error)
	GetByUser(ctx context.Context, user string, contentType entity.ContentType) ([]*entity.Data, error)
}

type KeyService interface {
	GetKeyForUser(user string) (string, error)
}

type AuthService interface {
	GetUserFromContext(ctx context.Context) (string, error)
}

type logPassContent struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Service struct {
	repo        Repo
	keyService  KeyService
	authService AuthService
}

func NewLogPassService(repo Repo, keyService KeyService, authService AuthService) *Service {
	return &Service{repo: repo, keyService: keyService, authService: authService}
}

// Create log/pass
func (s *Service) Create(ctx context.Context, r handlers.CreateLogPassRequest) (*handlers.CreateLogPassResponse, error) {
	jsonData, err := json.Marshal(&r)
	if err != nil {
		return nil, err
	}

	user, err := s.authService.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	key, err := s.keyService.GetKeyForUser(user)
	if err != nil {
		return nil, err
	}

	jsonEncrypted, err := lib.Encrypt(key, jsonData)
	if err != nil {
		return nil, err
	}

	newDataToSave := entity.Data{
		UUID:        uuid.New().String(),
		Content:     jsonEncrypted,
		ContentType: entity.LogPass,
		CreatedAt:   time.Now(),
		CreatedBy:   user,
	}

	err = s.repo.Insert(ctx, newDataToSave)
	if err != nil {
		return nil, err
	}

	return &handlers.CreateLogPassResponse{UUID: newDataToSave.UUID}, nil
}

// Update update log/pass
func (s *Service) Update(ctx context.Context, r handlers.UpdateLogPassRequest) (*handlers.UpdateLogPassResponse, error) {
	user, err := s.authService.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	key, err := s.keyService.GetKeyForUser(user)
	if err != nil {
		return nil, err
	}

	fromDB, err := s.repo.GetByUUID(ctx, user, r.UUID)
	if err != nil {
		return nil, err
	}

	jsonDecrypted, err := lib.Decrypt(key, fromDB.Content)
	if err != nil {
		return nil, err
	}

	logPassContent := logPassContent{}
	err = json.Unmarshal(jsonDecrypted, &logPassContent)
	if err != nil {
		return nil, err
	}

	s.setContentToUpdate(r, &logPassContent)

	jsonData, err := json.Marshal(&logPassContent)
	if err != nil {
		return nil, err
	}

	jsonEncrypted, err := lib.Encrypt(key, jsonData)
	if err != nil {
		return nil, err
	}

	fromDB.Content = jsonEncrypted

	err = s.repo.Update(ctx, *fromDB)
	if err != nil {
		return nil, err
	}

	return &handlers.UpdateLogPassResponse{UUID: fromDB.UUID}, nil
}

func (s *Service) setContentToUpdate(r handlers.UpdateLogPassRequest, logpassContent *logPassContent) {
	if r.Name != nil {
		logpassContent.Name = *r.Name
	}
	if r.Login != nil {
		logpassContent.Login = *r.Login
	}
	if r.Password != nil {
		logpassContent.Password = *r.Password
	}
}

// Delete delete log/pass
func (s *Service) Delete(ctx context.Context, r handlers.DeleteLogPassRequest) (*handlers.DeleteLogPassResponse, error) {
	user, err := s.authService.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = s.repo.Delete(ctx, user, r.UUID)
	if err != nil {
		return nil, err
	}

	return &handlers.DeleteLogPassResponse{UUID: r.UUID}, nil
}

// GetAll get all log/pass
func (s *Service) GetAll(ctx context.Context, r handlers.GetAllLogPassesRequest) (*handlers.GetAllLogPassesResponse, error) {
	user, err := s.authService.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	key, err := s.keyService.GetKeyForUser(user)
	if err != nil {
		return nil, err
	}

	data, err := s.repo.GetByUser(ctx, user, entity.LogPass)
	if err != nil {
		return nil, err
	}

	items := make([]handlers.GetAllLogPassResponseItem, 0, len(data))
	for _, v := range data {
		jsonDecrypted, err := lib.Decrypt(key, v.Content)
		if err != nil {
			return nil, err
		}
		item := handlers.GetAllLogPassResponseItem{}
		err = json.Unmarshal(jsonDecrypted, &item)
		item.UUID = v.UUID
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return &handlers.GetAllLogPassesResponse{Items: items}, nil
}
