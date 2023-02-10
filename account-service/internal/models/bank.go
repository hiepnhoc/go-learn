package models

import (
	"time"
)

type Bank struct {
	Id           string     `json:"id" db:"id" validate:"omitempty"`
	BankCode     string     `json:"bank_code" db:"bank_code" validate:"required"`
	BankName     string     `json:"bank_name" db:"bank_name" validate:"required"`
	BankNumber   string     `json:"bank_number" db:"bank_number" validate:"required"`
	CustomerName string     `json:"customer_name" db:"customer_name" validate:"required"`
	Content      *string    `json:"content" db:"content" `
	IsDefault    bool       `json:"is_default" db:"is_default" validate:"required"`
	Owner        string     `json:"owner" db:"owner"  validate:"required"`
	CreatedAt    time.Time  `json:"created_at,omitempty" db:"created_at" goqu:"skipinsert"`
	UpdatedAt    time.Time  `json:"updated_at,omitempty" db:"updated_at" goqu:"skipinsert"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at" goqu:"skipinsert"`
}
