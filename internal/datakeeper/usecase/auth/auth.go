package auth

import (
	"context"
	"fmt"
	"log/slog"

	customerr "github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/error"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/controller/http/handlers"
	securityservicev1 "github.com/GusevGrishaEm1/protos/gen/go/security_service"
)

type KeyService interface {
	SetKeyForUser(user string, key string) error
	GenerateKey() (string, error)
}

type Service struct {
	authClient securityservicev1.AuthClient
	keyService KeyService
	logger     *slog.Logger
}

// NewAuthService creates new auth service
func NewAuthService(authClient securityservicev1.AuthClient, keyService KeyService, logger *slog.Logger) (*Service, error) {
	return &Service{authClient, keyService, logger}, nil
}

// SignIn sign in user
func (a *Service) SignIn(ctx context.Context, r handlers.LoginRequest) (*handlers.LoginResponse, error) {
	res, err := a.authClient.Login(ctx, &securityservicev1.LoginRequest{Email: r.Email, Password: r.Password})
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	if err = a.keyService.SetKeyForUser(r.Email, r.Key); err != nil {
		return nil, err
	}
	return &handlers.LoginResponse{Token: res.Token}, nil
}

// SignUp sign up user
func (a *Service) SignUp(ctx context.Context, r handlers.RegisterRequest) (*handlers.RegisterResponse, error) {
	res, err := a.authClient.Register(ctx, &securityservicev1.RegisterRequest{Email: r.Email, Password: r.Password})
	if err != nil {
		a.logger.Error(err.Error())
		return nil, err
	}
	key, err := a.keyService.GenerateKey()
	if err != nil {
		return nil, err
	}
	if err = a.keyService.SetKeyForUser(r.Email, key); err != nil {
		return nil, err
	}
	return &handlers.RegisterResponse{Token: res.Token, Key: key}, nil
}

// GetUserFromContext get user from context
func (a *Service) GetUserFromContext(ctx context.Context) (string, error) {
	val, ok := ctx.Value("User").(string)
	if !ok {
		return "", customerr.Error(customerr.NO_USER_IN_CONTEXT)
	}
	return val, nil
}
