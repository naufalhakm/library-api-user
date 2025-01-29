package models

import "time"

type UserActivity struct {
	ID                uint64
	UserID            uint64
	BookID            uint64
	ActivityType      string
	ActivityTimestamp time.Time
}
