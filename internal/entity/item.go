package entity

import (
	"errors"
	"time"
)

var (
	ErrValidationFailed      = errors.New("validation failed")
	ErrDuplicateCode         = errors.New("item code already exists in this company")
	ErrItemNotFound          = errors.New("item not found")
	ErrAlreadyArchived       = errors.New("item is already archived")
	ErrInternalServerError   = errors.New("internal server error")
)

const (
	TypeProduct = "PRODUCT"
	TypeService = "SERVICE"

	StatusActive   = "ACTIVE"
	StatusInactive = "INACTIVE"
	StatusArchived = "ARCHIVED"
)

type Item struct {
	ID           string    `json:"item_id"`
	CompanyID    string    `json:"company_id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Price        float64   `json:"price"`
	CategoryName string    `json:"category_name,omitempty"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ItemFilter struct {
	Type    string
	Status  string
	Keyword string
}
