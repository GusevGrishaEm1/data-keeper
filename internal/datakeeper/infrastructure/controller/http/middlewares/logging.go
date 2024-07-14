package middlewares

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

type LoggerMiddleware struct {
	*slog.Logger
}

func NewLoggerMiddleware(logger *slog.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{logger}
}

func (logger *LoggerMiddleware) LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		if err != nil {
			logger.Error(fmt.Sprintf("Error: %s [%s] %s %s",
				err,
				c.Request().Method,
				c.Path(),
				c.RealIP(),
			))
			return err
		}

		logger.Info(fmt.Sprintf("[%s] %s %s %s",
			c.Request().Method,
			c.Path(),
			c.RealIP(),
			time.Since(start),
		))
		return nil
	}
}
