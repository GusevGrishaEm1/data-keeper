package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/config"
	customerr "github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/error"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo"
)

// AuthMiddleware auth middleware
type AuthMiddleware struct {
	jwtKey string
	logger *slog.Logger
}

// NewAuthMiddleware creates new auth middleware
func NewAuthMiddleware(config config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		jwtKey: config.AuthService.JWTKey,
		logger: slog.Default(),
	}
}

// AuthMiddleware auth middleware
// TODO: add refresh token
// TODO: add logout
// Get cookie from request and parse token
// Get email from token and set it in request context
func (m *AuthMiddleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("User")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, customerr.ToJson(err.Error()))
		}

		token, err := jwt.ParseWithClaims(cookie.Value, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.jwtKey), nil
		})
		if err != nil {
			return c.JSON(http.StatusUnauthorized, customerr.ToJson(err.Error()))
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			email := claims["email"].(string)
			c.Set("User", email)
		} else {
			return c.JSON(http.StatusUnauthorized, customerr.INVALID_TOKEN)
		}

		return next(c)
	}
}
