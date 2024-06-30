package logpass

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/entity"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/controller/http/handlers"
	"github.com/GusevGrishaEm1/data-keeper/internal/lib"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogPassService_Create(t *testing.T) {
	mockRepo := new(MockLogPassRepo)
	mockKeyService := new(MockKeyService)
	mockAuthService := new(MockAuthService)
	service := Service{
		repo:        mockRepo,
		keyService:  mockKeyService,
		authService: mockAuthService,
	}

	ctx := context.Background()
	user := "test_user"
	key := "1234567890123456"
	request := handlers.CreateLogPassRequest{Name: "test_name", Login: "test_login", Password: "test_password"}

	mockAuthService.On("GetUserFromContext", mock.Anything).Return(user, nil)
	mockKeyService.On("GetKeyForUser", user).Return(key, nil)
	mockRepo.On("Insert", mock.Anything, mock.AnythingOfType("entity.Data")).Return(nil)

	response, err := service.Create(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.UUID)

	mockAuthService.AssertCalled(t, "GetUserFromContext", mock.Anything)
	mockKeyService.AssertCalled(t, "GetKeyForUser", user)
	mockRepo.AssertCalled(t, "Insert", mock.Anything, mock.Anything)
}

func TestLogPassService_Update(t *testing.T) {
	mockRepo := new(MockLogPassRepo)
	mockKeyService := new(MockKeyService)
	mockAuthService := new(MockAuthService)
	service := Service{
		repo:        mockRepo,
		keyService:  mockKeyService,
		authService: mockAuthService,
	}

	ctx := context.Background()
	user := "test_user"
	key := "1234567890123456"
	uuidStr := uuid.New().String()
	content := logPassContent{Name: "old_name", Login: "old_login", Password: "old_password"}
	jsonContent, _ := json.Marshal(&content)
	encryptedContent, _ := lib.Encrypt(key, jsonContent)
	data := entity.Data{UUID: uuidStr, Content: encryptedContent, ContentType: entity.LogPass}

	mockAuthService.On("GetUserFromContext", mock.Anything).Return(user, nil)
	mockKeyService.On("GetKeyForUser", user).Return(key, nil)
	mockRepo.On("GetByUUID", mock.Anything, user, uuidStr).Return(&data, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("entity.Data")).Return(nil)

	updateRequest := handlers.UpdateLogPassRequest{UUID: uuidStr, Name: ptrString("new_name")}
	response, err := service.Update(ctx, updateRequest)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, uuidStr, response.UUID)
	mockAuthService.AssertCalled(t, "GetUserFromContext", mock.Anything)
	mockKeyService.AssertCalled(t, "GetKeyForUser", mock.Anything)
	mockRepo.AssertCalled(t, "GetByUUID", mock.Anything, user, uuidStr)
	mockRepo.AssertCalled(t, "Update", mock.Anything, mock.Anything)
}

func ptrString(s string) *string {
	return &s
}

func TestLogPassService_Delete(t *testing.T) {
	mockRepo := new(MockLogPassRepo)
	mockKeyService := new(MockKeyService)
	mockAuthService := new(MockAuthService)
	service := Service{
		repo:        mockRepo,
		keyService:  mockKeyService,
		authService: mockAuthService,
	}

	ctx := context.Background()
	user := "test_user"
	uuidStr := uuid.New().String()

	mockAuthService.On("GetUserFromContext", ctx).Return(user, nil)
	mockRepo.On("Delete", ctx, user, uuidStr).Return(nil)

	deleteRequest := handlers.DeleteLogPassRequest{UUID: uuidStr}
	response, err := service.Delete(ctx, deleteRequest)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, uuidStr, response.UUID)
	mockAuthService.AssertCalled(t, "GetUserFromContext", ctx)
	mockRepo.AssertCalled(t, "Delete", ctx, user, uuidStr)
}

func TestLogPassService_GetAll(t *testing.T) {
	mockRepo := new(MockLogPassRepo)
	mockKeyService := new(MockKeyService)
	mockAuthService := new(MockAuthService)
	service := Service{
		repo:        mockRepo,
		keyService:  mockKeyService,
		authService: mockAuthService,
	}

	ctx := context.Background()
	user := "test_user"
	key := "1234567890123456"
	content := logPassContent{Name: "test_name", Login: "test_login", Password: "test_password"}
	jsonContent, _ := json.Marshal(&content)
	encryptedContent, _ := lib.Encrypt(key, jsonContent)
	data := []*entity.Data{
		{UUID: uuid.New().String(), Content: encryptedContent, ContentType: entity.LogPass},
	}

	mockAuthService.On("GetUserFromContext", mock.Anything).Return(user, nil)
	mockKeyService.On("GetKeyForUser", user).Return(key, nil)
	mockRepo.On("GetByUser", mock.Anything, user, entity.LogPass).Return(data, nil)

	response, err := service.GetAll(ctx, handlers.GetAllLogPassesRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Items, len(data))
	mockAuthService.AssertCalled(t, "GetUserFromContext", mock.Anything)
	mockKeyService.AssertCalled(t, "GetKeyForUser", user)
	mockRepo.AssertCalled(t, "GetByUser", mock.Anything, user, entity.LogPass)
}
