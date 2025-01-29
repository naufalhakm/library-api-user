package models

import "time"

type User struct {
	ID        uint64
	Email     string
	Password  string
	Name      string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
