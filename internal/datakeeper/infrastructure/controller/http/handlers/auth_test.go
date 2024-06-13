package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
)

// Mock service
type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) SignIn(ctx context.Context, r LoginRequest) (*LoginResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*LoginResponse), args.Error(1)
}

func (m *mockAuthService) SignUp(ctx context.Context, r RegisterRequest) (*RegisterResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*RegisterResponse), args.Error(1)
}

func setupServer(mockAuthService *mockAuthService) *echo.Echo {
	e := echo.New()
	handler := NewAuthHandler(mockAuthService)

	e.POST("/login", handler.Login)
	e.POST("/register", handler.Register)

	return e
}

func TestAuthHandler_Login(t *testing.T) {
	mockAuthService := new(mockAuthService)
	loginRequest := LoginRequest{
		Email:    "test@example.com",
		Password: "password",
		Key:      "key",
	}
	loginResponse := &LoginResponse{Token: "token"}

	mockAuthService.On("SignIn", mock.Anything, loginRequest).Return(loginResponse, nil)

	e := setupServer(mockAuthService)

	server := httptest.NewServer(e)
	defer server.Close()

	expect := httpexpect.Default(t, server.URL)

	// Prepare request data
	reqData := map[string]interface{}{
		"email":    loginRequest.Email,
		"password": loginRequest.Password,
		"key":      loginRequest.Key,
	}

	// Perform request
	expect.POST("/login").
		WithJSON(reqData).
		Expect().
		Status(http.StatusOK).
		Cookie("User").Value().IsEqual("token")

	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Register(t *testing.T) {
	mockAuthService := new(mockAuthService)
	registerRequest := RegisterRequest{
		Email:    "test@example.com",
		Password: "password",
	}
	registerResponse := &RegisterResponse{Token: "token", Key: "key"}

	mockAuthService.On("SignUp", mock.Anything, registerRequest).Return(registerResponse, nil)

	e := setupServer(mockAuthService)

	server := httptest.NewServer(e)
	defer server.Close()

	expect := httpexpect.Default(t, server.URL)

	// Prepare request data
	reqData := map[string]interface{}{
		"email":    registerRequest.Email,
		"password": registerRequest.Password,
	}

	// Perform request
	expect.POST("/register").
		WithJSON(reqData).
		Expect().
		Status(http.StatusOK).
		Cookie("User").Value().IsEqual("token")

	mockAuthService.AssertExpectations(t)
}
