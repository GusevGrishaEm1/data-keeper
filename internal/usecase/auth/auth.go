package auth

import (
	"context"

	customerr "github.com/GusevGrishaEm1/data-keeper/internal/error"
	"github.com/GusevGrishaEm1/data-keeper/internal/infrastructure/controller/http/handlers"
	security_servicev1 "github.com/GusevGrishaEm1/protos/gen/go/security_service"
)

type KeyService interface {
	SetKeyForUser(user string, key string) error
	GenerateKey() (string, error)
}

type authService struct {
	authClient security_servicev1.AuthClient
	keyService KeyService
}

// NewAuthService creates new auth service
func NewAuthService(authClient security_servicev1.AuthClient, keyService KeyService) (*authService, error) {
	// conn, err := grpc.NewClient(
	// 	config.AuthService.URL,
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// 	grpc.WithIdleTimeout(time.Second*time.Duration(config.AuthService.Timeout)),
	// )
	// if err != nil {
	// 	return nil, err
	// }
	//return &authService{security_servicev1.NewAuthClient(conn), keyService}, nil
	return &authService{authClient, keyService}, nil
}

// SignIn sign in user
func (a *authService) SignIn(ctx context.Context, r handlers.LoginRequest) (*handlers.LoginResponse, error) {
	res, err := a.authClient.Login(ctx, &security_servicev1.LoginRequest{Email: r.Email, Password: r.Password})
	if err != nil {
		return nil, err
	}
	if err = a.keyService.SetKeyForUser(r.Email, r.Key); err != nil {
		return nil, err
	}
	return &handlers.LoginResponse{Token: res.Token}, nil
}

// SignUp sign up user
func (a *authService) SignUp(ctx context.Context, r handlers.RegisterRequest) (*handlers.RegisterResponse, error) {
	res, err := a.authClient.Register(ctx, &security_servicev1.RegisterRequest{Email: r.Email, Password: r.Password})
	if err != nil {
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
func (a *authService) GetUserFromContext(ctx context.Context) (string, error) {
	val, ok := ctx.Value("User").(string)
	if !ok {
		return "", customerr.Error(customerr.NO_USER_IN_CONTEXT)
	}
	return val, nil
}
