package file

import (
	"context"
	"encoding/json"
	"time"

	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/entity"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/controller/http/handlers"
	"github.com/GusevGrishaEm1/data-keeper/internal/lib"
	"github.com/google/uuid"
)

type FileContent struct {
	Name   string `json:"name"`
	Format string `json:"format"`
	Size   int    `json:"size"`
}

type UserFileRepo interface {
	Insert(ctx context.Context, data entity.UserFile) error
	Delete(ctx context.Context, user string, uuid string) error
	GetByUUID(ctx context.Context, user string, uuid string) (*entity.UserFile, error)
}

type DataRepo interface {
	Insert(ctx context.Context, data entity.Data) error
	Update(ctx context.Context, data entity.Data) error
	Delete(ctx context.Context, user string, uuid string) error
	GetByUUID(ctx context.Context, user string, uuid string, contentType entity.ContentType) (*entity.Data, error)
	GetByUser(ctx context.Context, user string, contentType entity.ContentType) ([]*entity.Data, error)
}

type AuthService interface {
	GetUserFromContext(ctx context.Context) (string, error)
}

type KeyService interface {
	GetKeyForUser(user string) (string, error)
}

type cardService struct {
	dataRepo     DataRepo
	userFileRepo UserFileRepo
	authService  AuthService
	keyService   KeyService
}

func NewFileService(dataRepo DataRepo, userFileRepo UserFileRepo, authService AuthService, keyService KeyService) *cardService {
	return &cardService{
		dataRepo:     dataRepo,
		userFileRepo: userFileRepo,
		authService:  authService,
		keyService:   keyService,
	}
}

// Upload file
func (s *cardService) UploadFile(ctx context.Context, r handlers.UploadFileRequest) (*handlers.UploadFileResponse, error) {
	user, err := s.authService.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	key, err := s.keyService.GetKeyForUser(user)
	if err != nil {
		return nil, err
	}

	// insert file meta to data
	fileContent := FileContent{
		Name:   r.Name,
		Format: r.Format,
		Size:   len(r.File),
	}
	jsonContent, err := json.Marshal(fileContent)
	if err != nil {
		return nil, err
	}

	encryptedContent, err := lib.Encrypt(key, jsonContent)
	if err != nil {
		return nil, err
	}

	data := entity.Data{
		UUID:        uuid.New().String(),
		Content:     encryptedContent,
		ContentType: entity.FILE,
		CreatedAt:   time.Now(),
		CreatedBy:   user,
	}

	err = s.dataRepo.Insert(ctx, data)
	if err != nil {
		return nil, err
	}

	// insert file data to user_file
	encryptedContent, err = lib.Encrypt(key, r.File)
	if err != nil {
		return nil, err
	}

	err = s.userFileRepo.Insert(ctx, entity.UserFile{
		UUID:      data.UUID,
		Content:   encryptedContent,
		CreatedAt: data.CreatedAt,
		CreatedBy: data.CreatedBy,
	})
	if err != nil {
		return nil, err
	}

	return &handlers.UploadFileResponse{UUID: data.UUID}, nil
}

func (s *cardService) DeleteFile(ctx context.Context, r handlers.DeleteFileRequest) (*handlers.DeleteFileResponse, error) {
	user, err := s.authService.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = s.dataRepo.Delete(ctx, user, r.UUID)
	if err != nil {
		return nil, err
	}

	if err = s.userFileRepo.Delete(ctx, user, r.UUID); err != nil {
		return nil, err
	}

	return &handlers.DeleteFileResponse{UUID: r.UUID}, nil
}

func (s *cardService) GetAllFiles(ctx context.Context, r handlers.GetAllFilesRequest) (*handlers.GetAllFilesResponse, error) {
	user, err := s.authService.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	key, err := s.keyService.GetKeyForUser(user)
	if err != nil {
		return nil, err
	}

	data, err := s.dataRepo.GetByUser(ctx, user, entity.FILE)
	if err != nil {
		return nil, err
	}

	items := make([]handlers.GetAllFilesResponceItem, 0, len(data))
	for _, item := range data {
		decryptedContent, err := lib.Decrypt(key, item.Content)
		if err != nil {
			return nil, err
		}
		fileDB := &FileContent{}
		err = json.Unmarshal(decryptedContent, fileDB)
		if err != nil {
			return nil, err
		}
		items = append(items, handlers.GetAllFilesResponceItem{
			UUID:   item.UUID,
			Name:   fileDB.Name,
			Format: fileDB.Format,
			Size:   fileDB.Size,
		})
	}

	return &handlers.GetAllFilesResponse{Items: items}, nil
}

func (s *cardService) DownloadFile(ctx context.Context, r handlers.DownloadFileRequest) (*handlers.DownloadFileResponse, error) {
	user, err := s.authService.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	key, err := s.keyService.GetKeyForUser(user)
	if err != nil {
		return nil, err
	}

	data, err := s.dataRepo.GetByUUID(ctx, user, r.UUID, entity.FILE)
	if err != nil {
		return nil, err
	}

	decryptedContent, err := lib.Decrypt(key, data.Content)
	if err != nil {
		return nil, err
	}

	fileDB := &FileContent{}
	err = json.Unmarshal(decryptedContent, fileDB)
	if err != nil {
		return nil, err
	}

	fileContent, err := s.userFileRepo.GetByUUID(ctx, user, data.UUID)
	if err != nil {
		return nil, err
	}

	decryptedFileContent, err := lib.Decrypt(key, fileContent.Content)
	if err != nil {
		return nil, err
	}

	return &handlers.DownloadFileResponse{
		Name:   fileDB.Name,
		Format: fileDB.Format,
		File:   decryptedFileContent,
	}, nil
}
