package token

import "time"

type Token struct {
	AuthId  int
	Role    string
	Expired time.Time
}
