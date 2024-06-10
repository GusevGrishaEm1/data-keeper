package middlewares

import (
	"net/http"

	"github.com/GusevGrishaEm1/data-keeper/internal/config"
	customerr "github.com/GusevGrishaEm1/data-keeper/internal/error"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo"
)

type AuthMiddleware struct {
	jwtkey string
}

func NewAuthMiddleware(config config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		jwtkey: config.AuthService.JWTKey,
	}
}

func (m *AuthMiddleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Request().Cookie("User")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, customerr.ToJson(err.Error()))
		}

		tokenStr := cookie.Value

		token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.jwtkey), nil
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
