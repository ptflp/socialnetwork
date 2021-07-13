package infoblog

import "time"

type HashTag struct {
	ID        int64 `json:"id" db:"id"`
	Tag       string
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
