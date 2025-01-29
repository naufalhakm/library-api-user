package models

import "time"

type Book struct {
	ID        uint64
	AuthorID  uint64
	Title     string
	Stock     int32
	PublishAt time.Time
	UpdatedAt time.Time
}
