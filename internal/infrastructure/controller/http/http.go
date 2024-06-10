package http

import (
	"context"
	"log/slog"

	"github.com/GusevGrishaEm1/data-keeper/internal/config"
	"github.com/GusevGrishaEm1/data-keeper/internal/infrastructure/controller/http/handlers"
	"github.com/GusevGrishaEm1/data-keeper/internal/infrastructure/controller/http/middlewares"
	"github.com/GusevGrishaEm1/data-keeper/internal/infrastructure/repository/postgres"
	"github.com/GusevGrishaEm1/data-keeper/internal/infrastructure/repository/postgres/repo"
	"github.com/GusevGrishaEm1/data-keeper/internal/usecase/auth"
	"github.com/GusevGrishaEm1/data-keeper/internal/usecase/card"
	"github.com/GusevGrishaEm1/data-keeper/internal/usecase/file"
	"github.com/GusevGrishaEm1/data-keeper/internal/usecase/key"
	security_servicev1 "github.com/GusevGrishaEm1/protos/gen/go/security_service"
	"google.golang.org/grpc"

	"github.com/labstack/echo"
)

func StartServer(context context.Context, config config.Config, logger *slog.Logger, conn *grpc.ClientConn, db *postgres.PostgresDB) error {
	e := echo.New()
	// key service
	keyService := key.NewKeyService()
	// auth service
	authService, err := auth.NewAuthService(security_servicev1.NewAuthClient(conn), keyService)
	if err != nil {
		return err
	}
	// auth handler
	authHandler := handlers.NewAuthHandler(authService)

	// use middlewares logging
	loggerMiddlewarer := middlewares.NewLoggerMiddleware(logger)
	e.Use(echo.MiddlewareFunc(loggerMiddlewarer.LoggerMiddleware))

	// mapping auth handlers
	e.POST("api/auth/login", authHandler.Login)
	e.POST("api/auth/register", authHandler.Register)

	// auth middlewarer
	authMiddlewarer := middlewares.NewAuthMiddleware(config)

	// data repo
	dataRepo := repo.NewDataRepo(db)
	// card service
	cardService := card.NewCardService(dataRepo, authService, keyService)
	// card handler
	cardHandler := handlers.NewCardHandler(cardService)

	// mapping handlers
	e.POST("api/cards", cardHandler.CreateCard, echo.MiddlewareFunc(authMiddlewarer.AuthMiddleware))
	e.PATCH("api/cards", cardHandler.UpdateCard, echo.MiddlewareFunc(authMiddlewarer.AuthMiddleware))
	e.GET("api/cards", cardHandler.GetCardsByUser, echo.MiddlewareFunc(authMiddlewarer.AuthMiddleware))
	e.DELETE("api/cards/:uuid", cardHandler.DeleteCard, echo.MiddlewareFunc(authMiddlewarer.AuthMiddleware))

	// user's files
	userFileRepo := repo.NewUserFileRepo(db)
	// file service
	fileService := file.NewFileService(dataRepo, userFileRepo, authService, keyService)
	// file handler
	fileHandler := handlers.NewFileHandler(fileService)

	// mapping handlers
	e.POST("api/files", fileHandler.UploadFile, echo.MiddlewareFunc(authMiddlewarer.AuthMiddleware))
	e.DELETE("api/files/:uuid", fileHandler.DeleteFile, echo.MiddlewareFunc(authMiddlewarer.AuthMiddleware))
	e.GET("api/files", fileHandler.GetAllFiles, echo.MiddlewareFunc(authMiddlewarer.AuthMiddleware))
	e.GET("api/files/:uuid", fileHandler.DownloadFile, echo.MiddlewareFunc(authMiddlewarer.AuthMiddleware))

	logger.Info("server started")
	e.Start(config.URL)

	return nil
}
