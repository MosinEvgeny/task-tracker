package domain

import (
	"github.com/google/uuid"
)

type Label struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Color  string    `json:"color"`
	UserID uuid.UUID `json:"user_id"`
}
