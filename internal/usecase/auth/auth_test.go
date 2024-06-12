package auth

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/GusevGrishaEm1/data-keeper/internal/infrastructure/controller/http/handlers"
	security_servicev1 "github.com/GusevGrishaEm1/protos/gen/go/security_service"
)

// MockAuthClient представляет фейковый клиент для тестов
type MockAuthClient struct {
}

func (m *MockAuthClient) Login(ctx context.Context, in *security_servicev1.LoginRequest, opts ...grpc.CallOption) (*security_servicev1.LoginResponse, error) {
	if in.Email == "existing@example.com" && in.Password == "password" {
		return &security_servicev1.LoginResponse{Token: "some_token"}, nil
	}
	return nil, errors.New("login failed")
}

func (m *MockAuthClient) Register(ctx context.Context, in *security_servicev1.RegisterRequest, opts ...grpc.CallOption) (*security_servicev1.RegisterResponse, error) {
	if in.Email == "new@example.com" && in.Password == "password" {
		return &security_servicev1.RegisterResponse{Token: "some_token"}, nil
	}
	return nil, errors.New("registration failed")
}

func (m *MockAuthClient) Refresh(ctx context.Context, in *security_servicev1.RefreshRequest, opts ...grpc.CallOption) (*security_servicev1.RefreshResponse, error) {
	panic("implement me")
}
func (m *MockAuthClient) Logout(ctx context.Context, in *security_servicev1.LogoutRequest, opts ...grpc.CallOption) (*security_servicev1.LogoutResponse, error) {
	panic("implement me")
}

// mockKeyService представляет фейковый сервис ключей для тестов
type mockKeyService struct {
}

func (m *mockKeyService) SetKeyForUser(user string, key string) error {
	return nil
}

func (m *mockKeyService) GenerateKey() (string, error) {
	return "some_key", nil
}

func TestSignIn(t *testing.T) {
	ctx := context.Background()
	mockAuthClient := new(MockAuthClient)
	mockKeyService := new(mockKeyService)

	service, err := NewAuthService(mockAuthClient, mockKeyService, slog.Default())
	assert.NoError(t, err)

	// Успешный вход
	request := handlers.LoginRequest{
		Email:    "existing@example.com",
		Password: "password",
		Key:      "some_key",
	}
	response, err := service.SignIn(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "some_token", response.Token)

	// Неудачный вход
	request = handlers.LoginRequest{
		Email:    "nonexisting@example.com",
		Password: "password",
		Key:      "some_key",
	}
	_, err = service.SignIn(ctx, request)
	assert.Error(t, err)
}

func TestSignUp(t *testing.T) {
	ctx := context.Background()
	mockAuthClient := new(MockAuthClient)
	mockKeyService := new(mockKeyService)

	service, err := NewAuthService(mockAuthClient, mockKeyService, slog.Default())
	assert.NoError(t, err)

	// Успешная регистрация
	request := handlers.RegisterRequest{
		Email:    "new@example.com",
		Password: "password",
	}
	response, err := service.SignUp(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "some_token", response.Token)

	// Неудачная регистрация
	request = handlers.RegisterRequest{
		Email:    "existing@example.com",
		Password: "password",
	}
	_, err = service.SignUp(ctx, request)
	assert.Error(t, err)
}
