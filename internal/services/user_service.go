package services

import (
	"context"
	"database/sql"
	"fmt"
	"library-api-user/internal/commons/response"
	"library-api-user/internal/grpc/client"
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
}

func NewUserService(db *sql.DB, bookClient *client.BookClient, userRepository repositories.UserRepository, borrowRepository repositories.BorrowRepository, activityRepository repositories.UserActivityRepository) UserService {
	return &UserServiceImpl{
		UserRepository:     userRepository,
		BorrowRepository:   borrowRepository,
		ActivityRepository: activityRepository,
		DB:                 db,
		BookClient:         bookClient,
	}
}

func (service *UserServiceImpl) Detail(ctx context.Context, id uint64) (*params.UserResponse, *response.CustomError) {
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

	user, err := service.UserRepository.FindUserByID(ctx, tx, id)
	if user == nil || err != nil {
		return nil, response.BadRequestError("User is not found!")
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

// Login implements services.AuthSvc
func (service *UserServiceImpl) Update(ctx context.Context, req *params.UserRequest, id uint64) *response.CustomError {
	val := validator.New()
	err := val.Struct(req)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errors []interface{}
		for _, fieldError := range validationErrors {
			error := "error " + fieldError.Field() + " on tag " + fieldError.Tag()
			errors = append(errors, error)
		}
		// service.Logger.Error("[UserService] Failed login request body", map[string]interface{}{
		// 	"error": "Incoming request body that failed to validate.",
		// })
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

	usr, err := service.UserRepository.FindUserByID(ctx, tx, id)
	if usr.Role == "admin" || err != nil {
		return response.BadRequestError("User not compatible to manage.")
	}

	user := models.User{
		ID:        req.ID,
		Email:     req.Email,
		Password:  req.Password,
		Name:      req.Name,
		Role:      req.Role,
		UpdatedAt: time.Now(),
	}

	err = service.UserRepository.UpdateUser(ctx, tx, &user)
	if err != nil {
		// service.Logger.Error("[UserService] Failed login user", map[string]interface{}{
		// 	"error": "Search find user by email.",
		// })
		return response.BadRequestError("Failed update user!")
	}

	return nil
}

func (service *UserServiceImpl) GetAll(ctx context.Context, pagination *models.Pagination) ([]*params.UserResponse, *response.CustomError) {
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

	pagination.Offset = (pagination.Page - 1) * pagination.PageSize

	users, err := service.UserRepository.GetAllUsers(ctx, tx, pagination)
	if err != nil {
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

	fmt.Println(userID)

	err = service.BookClient.DecreaseStock(ctx, bookID)
	if err != nil {
		return response.GeneralError("Failed to decrease book stock: " + err.Error())
	}

	borrowRecord := models.BorrowRecord{
		UserID:     userID,
		BookID:     bookID,
		BorrowedAt: time.Now(),
	}
	err = service.BorrowRepository.CreateBorrow(ctx, tx, &borrowRecord)
	if err != nil {
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
		return response.GeneralError("Failed to create user activity: " + err.Error())
	}

	return nil
}

func (service *UserServiceImpl) ReturnBook(ctx context.Context, userID uint64, bookID uint64) *response.CustomError {
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

	err = service.BookClient.IncreaseStock(ctx, bookID)
	if err != nil {
		return response.GeneralError("Failed to increase book stock: " + err.Error())
	}

	borrowRecord, err := service.BorrowRepository.FindBorrow(ctx, tx, userID, bookID)
	if err != nil {
		return response.GeneralError("Failed to find borrow record: " + err.Error())
	}

	borrowRecord.ReturnedAt = time.Now()
	err = service.BorrowRepository.UpdateBorrow(ctx, tx, borrowRecord)
	if err != nil {
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
		return response.GeneralError("Failed to create user activity: " + err.Error())
	}

	return nil
}
