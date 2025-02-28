package models

import "time"

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAT time.Time `json:"created_at" db:"created_at"`
	UpdatedAT time.Time `json:"updated_at" db:"updated_at"`
}
