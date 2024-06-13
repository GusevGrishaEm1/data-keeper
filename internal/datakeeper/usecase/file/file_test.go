package file

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/entity"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/controller/http/handlers"
	"github.com/GusevGrishaEm1/data-keeper/internal/lib"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the dependencies
type mockDataRepo struct {
	mock.Mock
}

func (m *mockDataRepo) Insert(ctx context.Context, data entity.Data) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *mockDataRepo) Update(ctx context.Context, data entity.Data) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *mockDataRepo) Delete(ctx context.Context, user string, uuid string) error {
	args := m.Called(ctx, user, uuid)
	return args.Error(0)
}

func (m *mockDataRepo) GetByUUID(ctx context.Context, user string, uuid string, contentType entity.ContentType) (*entity.Data, error) {
	args := m.Called(ctx, user, uuid, contentType)
	return args.Get(0).(*entity.Data), args.Error(1)
}

func (m *mockDataRepo) GetByUser(ctx context.Context, user string, contentType entity.ContentType) ([]*entity.Data, error) {
	args := m.Called(ctx, user, contentType)
	return args.Get(0).([]*entity.Data), args.Error(1)
}

type mockUserFileRepo struct {
	mock.Mock
}

func (m *mockUserFileRepo) Insert(ctx context.Context, data entity.UserFile) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *mockUserFileRepo) Delete(ctx context.Context, user string, uuid string) error {
	args := m.Called(ctx, user, uuid)
	return args.Error(0)
}

func (m *mockUserFileRepo) GetByUUID(ctx context.Context, user string, uuid string) (*entity.UserFile, error) {
	args := m.Called(ctx, user, uuid)
	return args.Get(0).(*entity.UserFile), args.Error(1)
}

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) GetUserFromContext(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

type MockKeyService struct {
	mock.Mock
}

func (m *MockKeyService) GetKeyForUser(user string) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

// Unit tests
func TestUploadFile(t *testing.T) {
	mockDataRepo := new(mockDataRepo)
	mockUserFileRepo := new(mockUserFileRepo)
	mockAuthService := new(mockAuthService)
	mockKeyService := new(MockKeyService)

	service := NewFileService(mockDataRepo, mockUserFileRepo, mockAuthService, mockKeyService)

	ctx := context.Background()
	user := "test-user"
	key := "352fa5gdhvdryhwr"
	fileContent := []byte("test file content")

	mockAuthService.On("GetUserFromContext", ctx).Return(user, nil)
	mockKeyService.On("GetKeyForUser", user).Return(key, nil)

	mockDataRepo.On("Insert", ctx, mock.AnythingOfType("entity.Data")).Return(nil)

	mockUserFileRepo.On("Insert", ctx, mock.AnythingOfType("entity.UserFile")).Return(nil)

	req := handlers.UploadFileRequest{
		Name:   "test-file",
		Format: "txt",
		File:   fileContent,
	}

	resp, err := service.UploadFile(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.UUID)
}

func TestDeleteFile(t *testing.T) {
	mockDataRepo := new(mockDataRepo)
	mockUserFileRepo := new(mockUserFileRepo)
	mockAuthService := new(mockAuthService)
	mockKeyService := new(MockKeyService)

	service := NewFileService(mockDataRepo, mockUserFileRepo, mockAuthService, mockKeyService)

	ctx := context.Background()
	user := "test-user"
	fileUUID := uuid.New().String()

	mockAuthService.On("GetUserFromContext", ctx).Return(user, nil)
	mockDataRepo.On("Delete", ctx, user, fileUUID).Return(nil)
	mockUserFileRepo.On("Delete", ctx, user, fileUUID).Return(nil)

	req := handlers.DeleteFileRequest{
		UUID: fileUUID,
	}

	resp, err := service.DeleteFile(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, fileUUID, resp.UUID)
}

func TestGetAllFiles(t *testing.T) {
	mockDataRepo := new(mockDataRepo)
	mockUserFileRepo := new(mockUserFileRepo)
	mockAuthService := new(mockAuthService)
	mockKeyService := new(MockKeyService)

	service := NewFileService(mockDataRepo, mockUserFileRepo, mockAuthService, mockKeyService)

	ctx := context.Background()
	user := "test-user"
	key := "352fa5gdhvdryhwr"

	mockAuthService.On("GetUserFromContext", ctx).Return(user, nil)
	mockKeyService.On("GetKeyForUser", user).Return(key, nil)

	fileContent := FileContent{
		Name:   "test-file",
		Format: "txt",
		Size:   10,
	}

	jsonContent, err := json.Marshal(fileContent)
	assert.NoError(t, err)
	encryptedContent, err := lib.Encrypt(key, jsonContent)
	assert.NoError(t, err)

	data := []*entity.Data{
		{
			UUID:        uuid.New().String(),
			Content:     encryptedContent,
			ContentType: entity.FILE,
			CreatedAt:   time.Now(),
			CreatedBy:   user,
		},
	}

	mockDataRepo.On("GetByUser", ctx, user, entity.FILE).Return(data, nil)

	resp, err := service.GetAllFiles(ctx, handlers.GetAllFilesRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, "test-file", resp.Items[0].Name)
	assert.Equal(t, "txt", resp.Items[0].Format)
	assert.Equal(t, 10, resp.Items[0].Size)
}

func TestDownloadFile(t *testing.T) {
	mockDataRepo := new(mockDataRepo)
	mockUserFileRepo := new(mockUserFileRepo)
	mockAuthService := new(mockAuthService)
	mockKeyService := new(MockKeyService)

	service := NewFileService(mockDataRepo, mockUserFileRepo, mockAuthService, mockKeyService)

	ctx := context.Background()
	user := "test-user"
	key := "352fa5gdhvdryhwr"
	fileUUID := uuid.New().String()

	mockAuthService.On("GetUserFromContext", ctx).Return(user, nil)
	mockKeyService.On("GetKeyForUser", user).Return(key, nil)

	metadata := FileContent{
		Name:   "test-file",
		Format: "txt",
		Size:   10,
	}
	metadataJson, err := json.Marshal(metadata)
	assert.NoError(t, err)
	encryptedMetadata, err := lib.Encrypt(key, metadataJson)
	assert.NoError(t, err)

	data := entity.Data{
		UUID:        fileUUID,
		Content:     encryptedMetadata,
		ContentType: entity.FILE,
		CreatedAt:   time.Now(),
		CreatedBy:   user,
	}

	mockDataRepo.On("GetByUUID", ctx, user, fileUUID, entity.FILE).Return(&data, nil)

	fileContent := []byte("encrypted file content")
	encryptedFileContent, err := lib.Encrypt(key, fileContent)
	assert.NoError(t, err)

	file := entity.UserFile{
		UUID:      fileUUID,
		Content:   encryptedFileContent,
		CreatedAt: time.Now(),
		CreatedBy: user,
	}

	mockUserFileRepo.On("GetByUUID", ctx, user, fileUUID).Return(&file, nil)

	resp, err := service.DownloadFile(ctx, handlers.DownloadFileRequest{UUID: fileUUID})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-file", resp.Name)
	assert.Equal(t, "txt", resp.Format)
	assert.Equal(t, []byte("encrypted file content"), resp.File)
}
