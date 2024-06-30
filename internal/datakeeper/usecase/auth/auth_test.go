package auth

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/controller/http/handlers"
	security_servicev1 "github.com/GusevGrishaEm1/protos/gen/go/security_service"
)

// mockAuthClient mocks auth client
type mockAuthClient struct {
}

// Login mock
func (m *mockAuthClient) Login(ctx context.Context, in *security_servicev1.LoginRequest, opts ...grpc.CallOption) (*security_servicev1.LoginResponse, error) {
	if in.Email == "existing@example.com" && in.Password == "password" {
		return &security_servicev1.LoginResponse{Token: "some_token"}, nil
	}
	return nil, errors.New("login failed")
}

// Register mock
func (m *mockAuthClient) Register(ctx context.Context, in *security_servicev1.RegisterRequest, opts ...grpc.CallOption) (*security_servicev1.RegisterResponse, error) {
	if in.Email == "new@example.com" && in.Password == "password" {
		return &security_servicev1.RegisterResponse{Token: "some_token"}, nil
	}
	return nil, errors.New("registration failed")
}

// Refresh mock
func (m *mockAuthClient) Refresh(ctx context.Context, in *security_servicev1.RefreshRequest, opts ...grpc.CallOption) (*security_servicev1.RefreshResponse, error) {
	panic("implement me")
}

// Logout mock
func (m *mockAuthClient) Logout(ctx context.Context, in *security_servicev1.LogoutRequest, opts ...grpc.CallOption) (*security_servicev1.LogoutResponse, error) {
	panic("implement me")
}

// KeyService mocks key service
type mockKeyService struct {
}

// SetKeyForUser mock
func (m *mockKeyService) SetKeyForUser(user string, key string) error {
	return nil
}

// GenerateKey mock
func (m *mockKeyService) GenerateKey() (string, error) {
	return "some_key", nil
}

func TestSignIn(t *testing.T) {
	ctx := context.Background()
	mockAuthClient := new(mockAuthClient)
	mockKeyService := new(mockKeyService)

	service, err := NewAuthService(mockAuthClient, mockKeyService, slog.Default())
	assert.NoError(t, err)

	request := handlers.LoginRequest{
		Email:    "existing@example.com",
		Password: "password",
		Key:      "some_key",
	}
	response, err := service.SignIn(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "some_token", response.Token)

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
	mockAuthClient := new(mockAuthClient)
	mockKeyService := new(mockKeyService)

	service, err := NewAuthService(mockAuthClient, mockKeyService, slog.Default())
	assert.NoError(t, err)

	request := handlers.RegisterRequest{
		Email:    "new@example.com",
		Password: "password",
	}
	response, err := service.SignUp(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "some_token", response.Token)

	request = handlers.RegisterRequest{
		Email:    "existing@example.com",
		Password: "password",
	}
	_, err = service.SignUp(ctx, request)
	assert.Error(t, err)
}
