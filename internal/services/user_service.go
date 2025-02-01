package services

import (
	"context"
	"database/sql"
	"fmt"
	"library-api-user/internal/commons/response"
	"library-api-user/internal/grpc/client"
	"library-api-user/internal/logger"
	"library-api-user/internal/models"
	"library-api-user/internal/params"
	"library-api-user/internal/repositories"
	"time"

	"github.com/go-playground/validator/v10"
)

type UserService interface {
	Detail(ctx context.Context, id uint64) (*params.UserResponse, *response.CustomError)
	Update(ctx context.Context, req *params.UserRequest, id uint64) *response.CustomError
	GetAll(ctx context.Context, pagination *models.Pagination) ([]*params.UserResponse, *response.CustomError)
	BorrowBook(ctx context.Context, userID uint64, bookID uint64) *response.CustomError
	ReturnBook(ctx context.Context, userID uint64, bookID uint64) *response.CustomError
}

type UserServiceImpl struct {
	UserRepository     repositories.UserRepository
	BorrowRepository   repositories.BorrowRepository
	ActivityRepository repositories.UserActivityRepository
	DB                 *sql.DB
	BookClient         *client.BookClient
	Logger             logger.Logger
}

func NewUserService(db *sql.DB, bookClient *client.BookClient, userRepository repositories.UserRepository, borrowRepository repositories.BorrowRepository, activityRepository repositories.UserActivityRepository, log logger.Logger) UserService {
	return &UserServiceImpl{
		UserRepository:     userRepository,
		BorrowRepository:   borrowRepository,
		ActivityRepository: activityRepository,
		DB:                 db,
		BookClient:         bookClient,
		Logger:             log,
	}
}

func (service *UserServiceImpl) Detail(ctx context.Context, id uint64) (*params.UserResponse, *response.CustomError) {
	tx, err := service.DB.Begin()
	if err != nil {
		service.Logger.Error("[UserService] Failed to begin transaction - Detail", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, response.GeneralError("Failed to begin transaction: " + err.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			service.Logger.Error("[UserService] Transaction rolled back due to panic - Detail", map[string]interface{}{
				"error": r,
			})
		} else if err != nil {
			tx.Rollback()
			service.Logger.Error("[UserService] Transaction rolled back due to error - Detail", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			tx.Commit()
		}
	}()

	user, err := service.UserRepository.FindUserByID(ctx, tx, id)
	if err != nil {
		service.Logger.Error("[UserService] Failed to retrieve user by ID - Detail", map[string]interface{}{
			"user_id": id,
			"error":   err.Error(),
		})
		return nil, response.NotFoundError("User not found")
	}

	return &params.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (service *UserServiceImpl) Update(ctx context.Context, req *params.UserRequest, id uint64) *response.CustomError {
	val := validator.New()
	err := val.Struct(req)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errors []interface{}
		for _, fieldError := range validationErrors {
			errors = append(errors, fmt.Sprintf("error %s on tag %s", fieldError.Field(), fieldError.Tag()))
		}
		service.Logger.Error("[UserService] Validation failed - Update", map[string]interface{}{
			"error": errors,
		})
		return response.BadRequestErrorWithAdditionalInfo(errors)
	}

	tx, err := service.DB.Begin()
	if err != nil {
		service.Logger.Error("[UserService] Failed to begin transaction - Update", map[string]interface{}{
			"error": err.Error(),
		})
		return response.GeneralError("Failed to begin transaction: " + err.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			service.Logger.Error("[UserService] Transaction rolled back due to panic - Update", map[string]interface{}{
				"error": r,
			})
		} else if err != nil {
			tx.Rollback()
			service.Logger.Error("[UserService] Transaction rolled back due to error - Update", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			tx.Commit()
		}
	}()

	_, err = service.UserRepository.FindUserByID(ctx, tx, id)
	if err != nil {
		service.Logger.Error("[UserService] Failed to find user by ID - Update", map[string]interface{}{
			"user_id": id,
			"error":   err.Error(),
		})
		return response.NotFoundError("User not found")
	}

	user := models.User{
		ID:        id,
		Email:     req.Email,
		Password:  req.Password,
		Name:      req.Name,
		Role:      req.Role,
		UpdatedAt: time.Now(),
	}

	err = service.UserRepository.UpdateUser(ctx, tx, &user)
	if err != nil {
		service.Logger.Error("[UserService] Failed to update user - Update", map[string]interface{}{
			"user_id": id,
			"error":   err.Error(),
		})
		return response.GeneralError("Failed to update user: " + err.Error())
	}

	return nil
}

func (service *UserServiceImpl) GetAll(ctx context.Context, pagination *models.Pagination) ([]*params.UserResponse, *response.CustomError) {
	tx, err := service.DB.Begin()
	if err != nil {
		service.Logger.Error("[UserService] Failed to begin transaction - GetAll", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, response.GeneralError("Failed to begin transaction: " + err.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			service.Logger.Error("[UserService] Transaction rolled back due to panic - GetAll", map[string]interface{}{
				"error": r,
			})
		} else if err != nil {
			tx.Rollback()
			service.Logger.Error("[UserService] Transaction rolled back due to error - GetAll", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			tx.Commit()
		}
	}()

	pagination.Offset = (pagination.Page - 1) * pagination.PageSize

	users, err := service.UserRepository.GetAllUsers(ctx, tx, pagination)
	if err != nil {
		service.Logger.Error("[UserService] Failed to fetch users - GetAll", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, response.GeneralError("Failed to fetch users: " + err.Error())
	}

	userResponses := make([]*params.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = &params.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Password:  user.Password,
			Name:      user.Name,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	pagination.PageCount = (pagination.TotalCount + pagination.PageSize - 1) / pagination.PageSize

	return userResponses, nil
}

func (service *UserServiceImpl) BorrowBook(ctx context.Context, userID uint64, bookID uint64) *response.CustomError {
	tx, err := service.DB.Begin()
	if err != nil {
		service.Logger.Error("[UserService] Failed to begin transaction - BorrowBook", map[string]interface{}{
			"error": err.Error(),
		})
		return response.GeneralError("Failed to begin transaction: " + err.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			service.Logger.Error("[UserService] Transaction rolled back due to panic - BorrowBook", map[string]interface{}{
				"error": r,
			})
		} else if err != nil {
			tx.Rollback()
			service.Logger.Error("[UserService] Transaction rolled back due to error - BorrowBook", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			tx.Commit()
		}
	}()

	err = service.BookClient.DecreaseStock(ctx, bookID)
	if err != nil {
		service.Logger.Error("[UserService] Failed to decrease book stock - BorrowBook", map[string]interface{}{
			"book_id": bookID,
			"error":   err.Error(),
		})
		return response.GeneralError("Failed to decrease book stock: " + err.Error())
	}

	borrowRecord := models.BorrowRecord{
		UserID:     userID,
		BookID:     bookID,
		BorrowedAt: time.Now(),
	}
	err = service.BorrowRepository.CreateBorrow(ctx, tx, &borrowRecord)
	if err != nil {
		service.Logger.Error("[UserService] Failed to create borrow record - BorrowBook", map[string]interface{}{
			"user_id": userID,
			"book_id": bookID,
			"error":   err.Error(),
		})
		return response.GeneralError("Failed to create borrow record: " + err.Error())
	}

	userActivity := models.UserActivity{
		UserID:            userID,
		BookID:            bookID,
		ActivityType:      "borrowed",
		ActivityTimestamp: time.Now(),
	}
	err = service.ActivityRepository.CreateActivity(ctx, tx, &userActivity)
	if err != nil {
		service.Logger.Error("[UserService] Failed to create user activity - BorrowBook", map[string]interface{}{
			"user_id": userID,
			"book_id": bookID,
			"error":   err.Error(),
		})
		return response.GeneralError("Failed to create user activity: " + err.Error())
	}

	return nil
}

func (service *UserServiceImpl) ReturnBook(ctx context.Context, userID uint64, bookID uint64) *response.CustomError {
	tx, err := service.DB.Begin()
	if err != nil {
		service.Logger.Error("[UserService] Failed to begin transaction - ReturnBook", map[string]interface{}{
			"error": err.Error(),
		})
		return response.GeneralError("Failed to begin transaction: " + err.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			service.Logger.Error("[UserService] Transaction rolled back due to panic - ReturnBook", map[string]interface{}{
				"error": r,
			})
		} else if err != nil {
			tx.Rollback()
			service.Logger.Error("[UserService] Transaction rolled back due to error - ReturnBook", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			tx.Commit()
		}
	}()

	err = service.BookClient.IncreaseStock(ctx, bookID)
	if err != nil {
		service.Logger.Error("[UserService] Failed to increase book stock - ReturnBook", map[string]interface{}{
			"book_id": bookID,
			"error":   err.Error(),
		})
		return response.GeneralError("Failed to increase book stock: " + err.Error())
	}

	borrowRecord, err := service.BorrowRepository.FindBorrow(ctx, tx, userID, bookID)
	if err != nil {
		service.Logger.Error("[UserService] Failed to find borrow record - ReturnBook", map[string]interface{}{
			"user_id": userID,
			"book_id": bookID,
			"error":   err.Error(),
		})
		return response.GeneralError("Failed to find borrow record: " + err.Error())
	}

	borrowRecord.ReturnedAt = time.Now()
	err = service.BorrowRepository.UpdateBorrow(ctx, tx, borrowRecord)
	if err != nil {
		service.Logger.Error("[UserService] Failed to update borrow record - ReturnBook", map[string]interface{}{
			"user_id": userID,
			"book_id": bookID,
			"error":   err.Error(),
		})
		return response.GeneralError("Failed to update borrow record: " + err.Error())
	}

	userActivity := models.UserActivity{
		UserID:            userID,
		BookID:            bookID,
		ActivityType:      "returned",
		ActivityTimestamp: time.Now(),
	}
	err = service.ActivityRepository.CreateActivity(ctx, tx, &userActivity)
	if err != nil {
		service.Logger.Error("[UserService] Failed to create user activity - ReturnBook", map[string]interface{}{
			"user_id": userID,
			"book_id": bookID,
			"error":   err.Error(),
		})
		return response.GeneralError("Failed to create user activity: " + err.Error())
	}

	return nil
}
