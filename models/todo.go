package models

import "time"

type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TodoRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	Date        time.Time `json:"date" validate:"required"`
	Completed   bool      `json:"completed"`
}
