package repositories

import (
	"context"
	"database/sql"
	"library-api-user/internal/models"
)

type UserActivityRepository interface {
	CreateActivity(ctx context.Context, tx *sql.Tx, activity *models.UserActivity) error
}

type UserActivityRepositoryImpl struct {
}

func NewUserActivityRepository() UserActivityRepository {
	return &UserActivityRepositoryImpl{}
}

func (repository *UserActivityRepositoryImpl) CreateActivity(ctx context.Context, tx *sql.Tx, activity *models.UserActivity) error {
	query := `INSERT INTO user_activities (user_id, book_id, activity_type, activity_timestamp) VALUES ($1, $2, $3, $4)`
	_, err := tx.ExecContext(ctx, query,
		activity.UserID,
		activity.BookID,
		activity.ActivityType,
		activity.ActivityTimestamp,
	)
	return err
}
