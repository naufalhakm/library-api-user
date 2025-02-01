package factory

import (
	"database/sql"
	"library-api-user/internal/config"
	"library-api-user/internal/controllers"
	"library-api-user/internal/grpc/client"
	"library-api-user/internal/logger"
	"library-api-user/internal/repositories"
	"library-api-user/internal/services"
	"log"
)

type Provider struct {
	AuthProvider controllers.AuthController
	UserProvider controllers.UserController
}

func InitFactory(db *sql.DB) *Provider {

	newLog, err := logger.NewLogger("./var/log/user.log")
	if err != nil {
		log.Fatalf("[Logger] Failed to initialize user service logger: %v", err)
	}
	bookClient, err := client.NewBookClient(config.ENV.BookGRCP)
	if err != nil {
		log.Fatalf("Failed to connect to BookService: %v", err)
	}
	userRepo := repositories.NewUserRepository()
	borrowRepo := repositories.NewBorrowRepository()
	activityRepo := repositories.NewUserActivityRepository()

	authSvc := services.NewAuthService(db, userRepo, newLog)
	authController := controllers.NewAuthController(authSvc)

	userSvc := services.NewUserService(db, bookClient, userRepo, borrowRepo, activityRepo, newLog)
	userController := controllers.NewUserController(userSvc)

	return &Provider{
		AuthProvider: authController,
		UserProvider: userController,
	}
}
