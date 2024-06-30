package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const secretKey = "your-256-bit-secret"

func generateTestJWT(t *testing.T, email string) string {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatal(err)
	}
	return tokenStr
}

func TestAuthMiddleware(t *testing.T) {
	e := echo.New()
	authMiddleware := NewAuthMiddleware(config.Config{AuthService: config.AuthServiceConfig{JWTKey: secretKey}})

	t.Run("Missing Cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		var handler echo.HandlerFunc = func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		}

		hand := authMiddleware.AuthMiddleware(handler)

		if assert.NoError(t, hand(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		cookie := &http.Cookie{
			Name:  "User",
			Value: "invalid-token",
		}
		req.AddCookie(cookie)
		c := e.NewContext(req, rec)

		handler := authMiddleware.AuthMiddleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})

		if assert.NoError(t, handler(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	})

	t.Run("Valid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		tokenStr := generateTestJWT(t, "test@example.com")
		cookie := &http.Cookie{
			Name:  "User",
			Value: tokenStr,
		}
		req.AddCookie(cookie)
		c := e.NewContext(req, rec)

		handler := authMiddleware.AuthMiddleware(func(c echo.Context) error {
			user, ok := c.Get("User").(string)
			if !ok {
				return c.String(http.StatusInternalServerError, "user not found")
			}
			assert.Equal(t, "test@example.com", user)
			return c.String(http.StatusOK, "test")
		})

		if assert.NoError(t, handler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
}
