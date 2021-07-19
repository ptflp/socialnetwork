package infoblog

import "time"

type HashTag struct {
	ID        int64     `json:"-" db:"id"`
	Tag       string    `json:"tag" db:"tag"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
