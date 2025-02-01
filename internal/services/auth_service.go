package services

import (
	"context"
	"database/sql"
	"fmt"
	"library-api-user/internal/commons/response"
	"library-api-user/internal/logger"
	"library-api-user/internal/models"
	"library-api-user/internal/params"
	"library-api-user/internal/repositories"
	"library-api-user/pkg/token"
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
	Logger         logger.Logger
}

func NewAuthService(db *sql.DB, userRepository repositories.UserRepository, log logger.Logger) AuthService {
	return &AuthServiceImpl{
		UserRepository: userRepository,
		DB:             db,
		Logger:         log,
	}
}

func (service *AuthServiceImpl) Register(ctx context.Context, req *params.RegisterRequest) *response.CustomError {
	val := validator.New()
	if err := val.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make([]string, len(validationErrors))
		for i, fieldError := range validationErrors {
			errors[i] = fmt.Sprintf("Field '%s' failed validation with tag '%s'", fieldError.Field(), fieldError.Tag())
		}
		service.Logger.Error("[AuthService] Validation failed - Register", map[string]interface{}{
			"error": errors,
		})
		return response.BadRequestErrorWithAdditionalInfo(errors)
	}

	tx, err := service.DB.Begin()
	if err != nil {
		service.Logger.Error("[AuthService] Failed to begin transaction - Register", map[string]interface{}{
			"error": err.Error(),
		})
		return response.GeneralError("Failed to begin transaction: " + err.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			service.Logger.Error("[AuthService] Transaction rolled back due to panic - Register", map[string]interface{}{
				"error": r,
			})
		} else if err != nil {
			tx.Rollback()
			service.Logger.Error("[AuthService] Transaction rolled back due to error - Register", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			tx.Commit()
		}
	}()

	existingUser, err := service.UserRepository.FindUserByEmail(ctx, tx, req.Email)
	if existingUser != nil || err == nil {
		service.Logger.Error("[AuthService] Failed email already exists! - Register", map[string]interface{}{
			"email": req.Email,
		})
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

	if err := service.UserRepository.CreateUser(ctx, tx, &user); err != nil {
		service.Logger.Error("[AuthService] Failed to create user - Register", map[string]interface{}{
			"email": req.Email,
			"error": err.Error(),
		})
		return response.GeneralError("Failed to create user: " + err.Error())
	}

	return nil
}

func (service *AuthServiceImpl) Login(ctx context.Context, req *params.LoginRequest) (*params.LoginResponse, *response.CustomError) {
	val := validator.New()
	if err := val.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make([]string, len(validationErrors))
		for i, fieldError := range validationErrors {
			errors[i] = fmt.Sprintf("Field '%s' failed validation with tag '%s'", fieldError.Field(), fieldError.Tag())
		}
		service.Logger.Error("[AuthService] Validation failed - Login", map[string]interface{}{
			"error": errors,
		})
		return nil, response.BadRequestErrorWithAdditionalInfo(errors)
	}

	tx, err := service.DB.Begin()
	if err != nil {
		service.Logger.Error("[AuthService] Failed to begin transaction - Login", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, response.GeneralError("Failed to begin transaction: " + err.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			service.Logger.Error("[AuthService] Transaction rolled back due to panic - Login", map[string]interface{}{
				"error": r,
			})
		} else if err != nil {
			tx.Rollback()
			service.Logger.Error("[AuthService] Transaction rolled back due to error - Login", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			tx.Commit()
		}
	}()

	user, err := service.UserRepository.FindUserByEmail(ctx, tx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			service.Logger.Warn("[AuthService] User not found - Login", map[string]interface{}{
				"email": req.Email,
			})
			return nil, response.BadRequestError("Invalid email or password")
		}
		service.Logger.Error("[AuthService] Failed to find user by email - Login", map[string]interface{}{
			"email": req.Email,
			"error": err.Error(),
		})
		return nil, response.GeneralError("Failed to find user: " + err.Error())
	}

	if user.Password != req.Password {
		service.Logger.Warn("[AuthService] Invalid password - Login", map[string]interface{}{
			"email": req.Email,
		})
		return nil, response.BadRequestError("Invalid email or password")
	}

	token, err := token.GenerateToken(int(user.ID), user.Role)
	if err != nil {
		service.Logger.Error("[AuthService] Failed to generate token - Login", map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		})
		return nil, response.GeneralError("Failed to generate token: " + err.Error())
	}

	return &params.LoginResponse{
		Token: token,
	}, nil
}
