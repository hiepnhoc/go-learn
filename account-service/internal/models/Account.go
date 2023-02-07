package models

import (
	"time"
)

// Account model
type Account struct {
	Id        string    `json:"id"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
