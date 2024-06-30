package http

import (
	"context"
	"log/slog"

	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/config"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/controller/http/handlers"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/controller/http/middlewares"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/repository/postgres"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/infrastructure/repository/postgres/repo"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/usecase/auth"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/usecase/file"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/usecase/key"
	"github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/usecase/logpass"
	securityservicev1 "github.com/GusevGrishaEm1/protos/gen/go/security_service"
	"google.golang.org/grpc"

	"github.com/labstack/echo"
)

func StartServer(context context.Context, config config.Config, logger *slog.Logger, conn *grpc.ClientConn, db *postgres.DB) error {
	e := echo.New()

	groupAPI := e.Group("/api")

	// use middlewares logging
	loggerMiddleware := middlewares.NewLoggerMiddleware(logger)
	groupAPI.Use(loggerMiddleware.LoggerMiddleware)

	// key service
	keyService := key.NewKeyService()
	// auth service
	authService, err := auth.NewAuthService(securityservicev1.NewAuthClient(conn), keyService, logger)
	if err != nil {
		return err
	}
	// auth handler
	authHandler := handlers.NewAuthHandler(authService)

	// mapping auth handlers
	groupAuth := groupAPI.Group("/auth")
	groupAuth.POST("/login", authHandler.Login)
	groupAuth.POST("/register", authHandler.Register)

	// auth middleware
	authMiddleware := middlewares.NewAuthMiddleware(config)

	// data repo
	dataRepo := repo.NewDataRepo(db)

	// log/pass service
	logPassService := logpass.NewLogPassService(dataRepo, keyService)
	// converter echo.Context -> context.Context
	ctxConverter := handlers.NewCtxConverter()
	// log/pass handler
	logPassHandler := handlers.NewLogPassHandler(logPassService, ctxConverter)

	// mapping log/pass handlers
	groupLogPass := groupAPI.Group("/logpass")
	groupLogPass.Use(authMiddleware.AuthMiddleware)
	groupLogPass.POST("", logPassHandler.CreateLogPass)
	groupLogPass.PATCH("", logPassHandler.UpdateLogPass)
	groupLogPass.GET("", logPassHandler.GetAllLogPasses)
	groupLogPass.DELETE("", logPassHandler.DeleteLogPass)

	// user's files repo
	userFileRepo := repo.NewUserFileRepo(db)
	// file service
	fileService := file.NewFileService(dataRepo, userFileRepo, authService, keyService)
	// file handler
	fileHandler := handlers.NewFileHandler(fileService, ctxConverter)

	// mapping files handlers
	groupFile := groupAPI.Group("/files")
	groupFile.Use(authMiddleware.AuthMiddleware)
	groupFile.POST("", fileHandler.UploadFile)
	groupFile.DELETE("", fileHandler.DeleteFile)
	groupFile.GET("", fileHandler.GetAllFiles)
	groupFile.GET("/:uuid", fileHandler.DownloadFile)

	logger.Info("server started")
	err = e.Start(config.Port)
	if err != nil {
		return err
	}

	return nil
}
