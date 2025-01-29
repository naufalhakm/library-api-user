package repositories

import (
	"context"
	"database/sql"
	"errors"
	"library-api-user/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, tx *sql.Tx, user *models.User) error
	FindUserByID(ctx context.Context, tx *sql.Tx, id uint64) (*models.User, error)
	FindUserByEmail(ctx context.Context, tx *sql.Tx, email string) (*models.User, error)
	UpdateUser(ctx context.Context, tx *sql.Tx, user *models.User) error
	DeleteUser(ctx context.Context, tx *sql.Tx, id uint64) error
	GetAllUsers(ctx context.Context, tx *sql.Tx, pagination *models.Pagination) ([]*models.User, error)
}

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) CreateUser(ctx context.Context, tx *sql.Tx, user *models.User) error {
	query := `INSERT INTO users (email, password, name, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	response, err := tx.ExecContext(ctx, query, user.Email, user.Password, user.Name, user.Role, user.CreatedAt, user.UpdatedAt)
	if err != nil || response == nil {
		return errors.New("Failed to create a user, transaction rolled back. Reason: " + err.Error())
	}

	return nil
}

func (repository *UserRepositoryImpl) FindUserByID(ctx context.Context, tx *sql.Tx, id uint64) (*models.User, error) {
	query := "SELECT id, email, password, name, role, created_at, updated_at FROM users WHERE id = $1 LIMIT 1"
	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user = models.User{}
	if rows.Next() {
		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return &user, nil
	} else {
		return nil, errors.New("user is not found")
	}
}
func (repository *UserRepositoryImpl) FindUserByEmail(ctx context.Context, tx *sql.Tx, email string) (*models.User, error) {
	query := "SELECT id, email, password, name, role, created_at, updated_at FROM users WHERE email = $1 LIMIT 1"
	rows, err := tx.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user = models.User{}
	if rows.Next() {
		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return &user, nil
	} else {
		return nil, errors.New("user is not found")
	}
}

func (repository *UserRepositoryImpl) UpdateUser(ctx context.Context, tx *sql.Tx, user *models.User) error {
	query := `UPDATE users SET email = $1, password = $2, name = $3, role = $4, updated_at = $5 WHERE id = $6`

	_, err := tx.ExecContext(ctx, query,
		user.Email,
		user.Password,
		user.Name,
		user.Role,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return errors.New("Failed to update a user, transaction rolled back. Reason: " + err.Error())
	}
	return nil
}

func (repository *UserRepositoryImpl) DeleteUser(ctx context.Context, tx *sql.Tx, id uint64) error {
	SQL := `DELETE FROM users WHERE id = $1`

	_, err := tx.ExecContext(ctx, SQL, id)
	if err != nil {
		return errors.New("Failed to update a user, transaction rolled back. Reason: " + err.Error())
	}
	return nil
}

func (repository *UserRepositoryImpl) GetAllUsers(ctx context.Context, tx *sql.Tx, pagination *models.Pagination) ([]*models.User, error) {
	query := `SELECT id, email, password, name, role, created_at, updated_at, count(*) over() FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := tx.QueryContext(ctx, query, pagination.PageSize, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt, &pagination.TotalCount)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)

	}
	return users, nil
}
