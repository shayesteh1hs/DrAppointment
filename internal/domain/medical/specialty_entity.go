package medical

import (
	"time"

	"github.com/google/uuid"
)

// Specialty represents a medical specialty
type Specialty struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (s Specialty) GetId() string {
	return s.ID.String()
}
