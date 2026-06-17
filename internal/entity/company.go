package entity

import "time"

type Company struct {
	ID        string    `json:"uuid"`
	Name      string    `json:"name"`
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
