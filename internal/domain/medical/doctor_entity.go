package medical

import (
	"time"

	"github.com/google/uuid"
)

type Doctor struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	SpecialtyID uuid.UUID `json:"specialty_id" db:"specialty_id"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	AvatarURL   string    `json:"avatar_url" db:"avatar_url"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// GetId returns the ID as a string for pagination compatibility
func (d Doctor) GetId() string {
	return d.ID.String()
}
