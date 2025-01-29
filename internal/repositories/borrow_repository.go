package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"library-api-user/internal/models"
)

type BorrowRepository interface {
	CreateBorrow(ctx context.Context, tx *sql.Tx, borrow *models.BorrowRecord) error
	FindBorrow(ctx context.Context, tx *sql.Tx, userID uint64, bookID uint64) (*models.BorrowRecord, error)
	UpdateBorrow(ctx context.Context, tx *sql.Tx, borrow *models.BorrowRecord) error
}

type BorrowRepositoryImpl struct {
}

func NewBorrowRepository() BorrowRepository {
	return &BorrowRepositoryImpl{}
}

func (repository *BorrowRepositoryImpl) CreateBorrow(ctx context.Context, tx *sql.Tx, borrow *models.BorrowRecord) error {

	fmt.Println(borrow)
	query := `
		INSERT INTO borrows (user_id, book_id, borrowed_at) VALUES ($1, $2, $3)`
	_, err := tx.ExecContext(ctx, query,
		borrow.UserID,
		borrow.BookID,
		borrow.BorrowedAt,
	)
	return err
}

func (repository *BorrowRepositoryImpl) FindBorrow(ctx context.Context, tx *sql.Tx, userID uint64, bookID uint64) (*models.BorrowRecord, error) {
	query := `
		SELECT id, user_id, book_id, borrowed_at, returned_at 
		FROM borrows 
		WHERE user_id = $1 AND book_id = $2 AND returned_at IS NULL`
	row := tx.QueryRowContext(ctx, query, userID, bookID)
	var borrow models.BorrowRecord
	err := row.Scan(
		&borrow.ID,
		&borrow.UserID,
		&borrow.BookID,
		&borrow.BorrowedAt,
		&borrow.ReturnedAt,
	)
	if err != nil {
		return nil, err
	}
	return &borrow, nil
}

func (repository *BorrowRepositoryImpl) UpdateBorrow(ctx context.Context, tx *sql.Tx, borrow *models.BorrowRecord) error {
	query := `
		UPDATE borrows SET returned_at = $1 WHERE id = $2`
	_, err := tx.ExecContext(ctx, query,
		borrow.ReturnedAt,
		borrow.ID,
	)
	return err
}
