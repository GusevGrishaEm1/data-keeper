package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

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
		Login:    "test@example.com",
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
		"login":    loginRequest.Login,
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
		Login:    "test@example.com",
		Password: "password",
	}
	registerResponse := &RegisterResponse{Key: "key"}

	mockAuthService.On("SignUp", mock.Anything, registerRequest).Return(registerResponse, nil)

	e := setupServer(mockAuthService)

	server := httptest.NewServer(e)
	defer server.Close()

	expect := httpexpect.Default(t, server.URL)

	// Prepare request data
	reqData := map[string]interface{}{
		"login":    registerRequest.Login,
		"password": registerRequest.Password,
	}

	// Perform request
	expect.POST("/register").
		WithJSON(reqData).
		Expect().
		Status(http.StatusOK)

	mockAuthService.AssertExpectations(t)
}
