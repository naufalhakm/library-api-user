package params

type BookRequest struct {
	AuthorID uint64 `json:"author_id"  validate:"required"`
	Title    string `json:"title"  validate:"required"`
	Stock    int32  `json:"stock"`
}
