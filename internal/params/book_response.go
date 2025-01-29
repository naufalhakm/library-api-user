package params

import "time"

type BookResponse struct {
	ID        uint64    `json:"id"`
	AuthorID  uint64    `json:"author_id"`
	Title     string    `json:"title"`
	Stock     int32     `json:"stock"`
	PublishAt time.Time `json:"publish_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
