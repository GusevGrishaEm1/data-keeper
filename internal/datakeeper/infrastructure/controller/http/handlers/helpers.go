package handlers

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
)

// ContextConverter convert echo.Context -> context.Context
type ctxConverter interface {
	// ConvertEchoCtxToCtx convert echo.Context -> context.Context
	ConvertEchoCtxToCtx(ctx echo.Context) (context.Context, error)
}

type CtxConverterImpl struct{}

// NewCtxConverter return new service
func NewCtxConverter() *CtxConverterImpl {
	return &CtxConverterImpl{}
}

// ConvertEchoCtxToCtx echo.Context -> context.Context
func (c *CtxConverterImpl) ConvertEchoCtxToCtx(ctx echo.Context) (context.Context, error) {
	user, ok := ctx.Get("User").(string)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return context.WithValue(context.Background(), "User", user), nil
}
