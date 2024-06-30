package logpass

import (
	"context"

	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/entity"
	"github.com/stretchr/testify/mock"
)

// MockLogPassRepo is a mock implementation of Repo
type MockLogPassRepo struct {
	mock.Mock
}

func (m *MockLogPassRepo) Insert(ctx context.Context, data entity.Data) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockLogPassRepo) Update(ctx context.Context, data entity.Data) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockLogPassRepo) Delete(ctx context.Context, user string, uuid string) error {
	args := m.Called(ctx, user, uuid)
	return args.Error(0)
}

func (m *MockLogPassRepo) GetByUUID(ctx context.Context, user string, uuid string) (*entity.Data, error) {
	args := m.Called(ctx, user, uuid)
	return args.Get(0).(*entity.Data), args.Error(1)
}

func (m *MockLogPassRepo) GetByUser(ctx context.Context, user string, contentType entity.ContentType) ([]*entity.Data, error) {
	args := m.Called(ctx, user, contentType)
	return args.Get(0).([]*entity.Data), args.Error(1)
}

// MockKeyService is a mock implementation of KeyService
type MockKeyService struct {
	mock.Mock
}

func (m *MockKeyService) GetKeyForUser(user string) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

// MockAuthService is a mock implementation of AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GetUserFromContext(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}
