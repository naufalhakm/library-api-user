package services

import (
	"context"
	"database/sql"
	"fmt"
	"library-api-user/internal/commons/response"
	"library-api-user/internal/models"
	"library-api-user/internal/params"
	"library-api-user/internal/repositories"
	"library-api-user/pkg/token"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
)

type AuthService interface {
	Register(ctx context.Context, req *params.RegisterRequest) *response.CustomError
	Login(ctx context.Context, req *params.LoginRequest) (*params.LoginResponse, *response.CustomError)
}

type AuthServiceImpl struct {
	UserRepository repositories.UserRepository
	DB             *sql.DB
}

func NewAuthService(db *sql.DB, userRepository repositories.UserRepository) AuthService {
	return &AuthServiceImpl{
		UserRepository: userRepository,
		DB:             db,
	}
}

func (service *AuthServiceImpl) Register(ctx context.Context, req *params.RegisterRequest) *response.CustomError {
	val := validator.New()
	err := val.Struct(req)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make([]string, len(validationErrors))
		for i, fieldError := range validationErrors {
			errors[i] = fmt.Sprintf("Field '%s' failed validation with tag '%s'", fieldError.Field(), fieldError.Tag())
		}
		// service.Logger.Error("[AuthService] Failed login request body", map[string]interface{}{
		// 	"error": "Incoming request body that failed to validate.",
		// })
		log.Printf("[AuthService] Validation failed: %v", errors)
		return response.BadRequestErrorWithAdditionalInfo(errors)
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return response.GeneralError("Failed Connection to database errors: " + err.Error())
	}
	defer func() {
		err := recover()
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	usr, err := service.UserRepository.FindUserByEmail(ctx, tx, req.Email)
	if usr != nil || err == nil {
		return response.BadRequestError("Email already exists!")
	}

	user := models.User{
		Email:     req.Email,
		Password:  req.Password,
		Name:      req.Name,
		Role:      req.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = service.UserRepository.CreateUser(ctx, tx, &user)
	if err != nil {
		return response.BadRequestError(err.Error())
	}

	return nil
}

// Login implements services.AuthSvc
func (service *AuthServiceImpl) Login(ctx context.Context, req *params.LoginRequest) (*params.LoginResponse, *response.CustomError) {
	val := validator.New()
	err := val.Struct(req)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make([]string, len(validationErrors))
		for i, fieldError := range validationErrors {
			errors[i] = fmt.Sprintf("Field '%s' failed validation with tag '%s'", fieldError.Field(), fieldError.Tag())
		}
		// service.Logger.Error("[AuthService] Failed login request body", map[string]interface{}{
		// 	"error": "Incoming request body that failed to validate.",
		// })
		return nil, response.BadRequestErrorWithAdditionalInfo(errors)
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, response.GeneralError("Failed Connection to database errors: " + err.Error())
	}
	defer func() {
		err := recover()
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	user, err := service.UserRepository.FindUserByEmail(ctx, tx, req.Email)
	if user == nil || err != nil {
		// service.Logger.Error("[AuthService] Failed login user", map[string]interface{}{
		// 	"error": "Search find user by email.",
		// })
		return nil, response.BadRequestError("Invalid email or password")
	}

	token, err := token.GenerateToken(int(user.ID), user.Role)
	if err != nil {
		return nil, response.GeneralErrorWithAdditionalInfo(err.Error())
	}

	return &params.LoginResponse{
		Token: token,
	}, nil
}
