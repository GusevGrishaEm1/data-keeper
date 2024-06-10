package handlers

import (
	"context"
	"net/http"
	"time"

	customerr "github.com/GusevGrishaEm1/data-keeper/internal/error"
	"github.com/labstack/echo"
)

// Сервис Аутентификации
type AuthService interface {
	// Аутентификация
	SignIn(ctx context.Context, r LoginRequest) (*LoginResponse, error)
	// Регистрация
	SignUp(ctx context.Context, r RegisterRequest) (*RegisterResponse, error)
}

// Запрос Аутентификации
type LoginRequest struct {
	// Email user
	Email string `json:"email"`
	// User password
	Password string `json:"password"`
	// Key cypher
	Key string `json:"key"`
}

// Ответ Аутентификации
type LoginResponse struct {
	Token string `json:"-"`
}

// Запрос Регистрации
type RegisterRequest struct {
	// Email user
	Email string `json:"email"`
	// User password
	Password string `json:"password"`
}

// Ответ Регистрации
type RegisterResponse struct {
	Token string `json:"-"`
	Key   string `json:"key"`
}

// Обработчик Аутентификации
type AuthHandler struct {
	authService AuthService
}

// Создание Обработчика Аутентификации
func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Аутентификация
func (h *AuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}
	res, err := h.authService.SignIn(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}
	c.SetCookie(&http.Cookie{
		Name:    "User",
		Value:   res.Token,
		Expires: time.Now().Add(24 * time.Hour),
	})
	return c.JSON(http.StatusOK, res)
}

// Регистрация
func (h *AuthHandler) Register(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}
	res, err := h.authService.SignUp(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}
	c.SetCookie(&http.Cookie{
		Name:    "User",
		Value:   res.Token,
		Expires: time.Now().Add(24 * time.Hour),
	})
	return c.JSON(http.StatusOK, res)
}
