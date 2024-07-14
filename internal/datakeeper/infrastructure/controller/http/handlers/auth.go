package handlers

import (
	"context"
	"net/http"
	"time"

	customerr "github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/error"
	"github.com/labstack/echo/v4"
)

// AuthService Auth service interface
type AuthService interface {
	// SignIn Sign in
	SignIn(ctx context.Context, r LoginRequest) (*LoginResponse, error)
	// SignUp Sign up
	SignUp(ctx context.Context, r RegisterRequest) (*RegisterResponse, error)
}

// LoginRequest Login request
type LoginRequest struct {
	// Login user
	Login string `json:"login"`
	// User password
	Password string `json:"password"`
	// Key cypher
	Key string `json:"key"`
}

// LoginResponse Login response
type LoginResponse struct {
	// Token user's token
	Token string `json:"token"`
}

// RegisterRequest Register request
type RegisterRequest struct {
	// Login user's login
	Login string `json:"login"`
	// Password user's password
	Password string `json:"password"`
}

// RegisterResponse Register response
type RegisterResponse struct {
	// Key user's generated key
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

// Login Authentication
// @Summary Login user
// @Description Authenticate user and get token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/login [post]
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

// Register Registration
// @Summary Register user
// @Description Register new user and get token
// @Tags auth
// @Accept json
// @Produce json
// @Param register body RegisterRequest true "Register request"
// @Success 200 {object} RegisterResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	req := new(RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, customerr.ToJson(err.Error()))
	}

	res, err := h.authService.SignUp(c.Request().Context(), *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerr.ToJson(err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}
