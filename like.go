package infoblog

import "time"

type Like struct {
	ID        int64     `json:"-" db:"id"`
	Type      int64     `json:"type" db:"type"`
	ForeignID int64     `json:"foreign_id" db:"foreign_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
