package domain

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	UserID      uuid.UUID `json:"user_id"`
	LabelIDs    []uuid.UUID `json:"label_ids"`
}