package handlers

import (
	"context"
	"net/http"
	"time"

	customerr "github.com/GusevGrishaEm1/data-keeper/internal/error"
	"github.com/labstack/echo"
)

// AuthService Auth service interface
type AuthService interface {
	// SignIn Sign in
	SignIn(ctx context.Context, r LoginRequest) (*LoginResponse, error)
	// SignUp Sign up
	SignUp(ctx context.Context, r RegisterRequest) (*RegisterResponse, error)
}

// Login request
type LoginRequest struct {
	// Email user
	Email string `json:"email"`
	// User password
	Password string `json:"password"`
	// Key cypher
	Key string `json:"key"`
}

// Login response
type LoginResponse struct {
	// User token
	Token string `json:"-"`
}

// Register request
type RegisterRequest struct {
	// Email user
	Email string `json:"email"`
	// User password
	Password string `json:"password"`
}

// Register response
type RegisterResponse struct {
	// User token
	Token string `json:"-"`
	// User key
	Key string `json:"key"`
}

// AuthHandler Auth handler
type AuthHandler struct {
	authService AuthService
}

// NewAuthHandler Create new auth handler
func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Authentication
func (h *AuthHandler) Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	res, err := h.authService.SignIn(c.Request().Context(), *req)
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

// Registration
func (h *AuthHandler) Register(c echo.Context) error {
	req := new(RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	res, err := h.authService.SignUp(c.Request().Context(), *req)
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
