package models

import "time"

type BorrowRecord struct {
	ID         uint64
	UserID     uint64
	BookID     uint64
	BorrowedAt time.Time
	ReturnedAt time.Time
}
