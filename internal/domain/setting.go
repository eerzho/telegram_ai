package domain

import "time"

type Setting struct {
	UserID    int64     `db:"user_id" json:"user_id"`
	ChatID    int64     `db:"chat_id" json:"chat_id"`
	Style     string    `db:"style" json:"style"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
