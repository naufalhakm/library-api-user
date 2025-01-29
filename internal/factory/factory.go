package factory

import (
	"database/sql"
	"library-api-user/internal/config"
	"library-api-user/internal/controllers"
	"library-api-user/internal/grpc/client"
	"library-api-user/internal/repositories"
	"library-api-user/internal/services"
	"log"
)

type Provider struct {
	AuthProvider controllers.AuthController
	UserProvider controllers.UserController
}

func InitFactory(db *sql.DB) *Provider {

	bookClient, err := client.NewBookClient(config.ENV.BookGRCP)
	if err != nil {
		log.Fatalf("Failed to connect to BookService: %v", err)
	}
	userRepo := repositories.NewUserRepository()
	borrowRepo := repositories.NewBorrowRepository()
	activityRepo := repositories.NewUserActivityRepository()

	authSvc := services.NewAuthService(db, userRepo)
	authController := controllers.NewAuthController(authSvc)

	userSvc := services.NewUserService(db, bookClient, userRepo, borrowRepo, activityRepo)
	userController := controllers.NewUserController(userSvc)

	return &Provider{
		AuthProvider: authController,
		UserProvider: userController,
	}
}
