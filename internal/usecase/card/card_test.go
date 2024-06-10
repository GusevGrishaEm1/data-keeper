package card

import (
	"context"
	"testing"

	"github.com/GusevGrishaEm1/data-keeper/internal/entity"
	"github.com/GusevGrishaEm1/data-keeper/internal/infrastructure/controller/http/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepo представляет фейковый репозиторий для тестов
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Insert(ctx context.Context, data entity.Data) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockRepo) Update(ctx context.Context, data entity.Data) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockRepo) Delete(ctx context.Context, user string, uuid string) error {
	args := m.Called(ctx, user, uuid)
	return args.Error(0)
}

func (m *MockRepo) GetByUUID(ctx context.Context, user string, uuid string, contentType entity.ContentType) (*entity.Data, error) {
	args := m.Called(ctx, user, uuid)
	card, _ := args.Get(0).(*entity.Data)
	return card, args.Error(1)
}

func (m *MockRepo) GetByUser(ctx context.Context, user string, contentType entity.ContentType) ([]*entity.Data, error) {
	args := m.Called(ctx, user, contentType)
	cards, _ := args.Get(0).([]*entity.Data)
	return cards, args.Error(1)
}

// MockAuthService представляет фейковый сервис аутентификации для тестов
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GetUserFromContext(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

// MockKeyService представляет фейковый сервис ключей для тестов
type MockKeyService struct {
	mock.Mock
}

func (m *MockKeyService) GetKeyForUser(user string) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func TestCreateCard(t *testing.T) {
	mockRepo := new(MockRepo)
	mockAuthService := new(MockAuthService)
	mockKeyService := new(MockKeyService)

	service := NewCardService(mockRepo, mockAuthService, mockKeyService)

	mockAuthService.On("GetUserFromContext", mock.Anything).Return("user123", nil)
	mockKeyService.On("GetKeyForUser", "user123").Return("1234567890123456", nil)

	request := handlers.CreateCardRequest{
		Key:     "1234",
		Number:  "1234567890123456",
		CVV:     "123",
		Name:    "John Doe",
		Expires: "12/25",
	}

	mockRepo.On("Insert", mock.Anything, mock.AnythingOfType("*entity.Data")).Return(nil)

	response, err := service.CreateCard(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.UUID)

	mockAuthService.AssertExpectations(t)
	mockKeyService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
