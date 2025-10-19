package medical

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SearchDoctorsRequest struct {
	SpecialtyID uuid.UUID `form:"specialty_id" validate:"omitempty"`
	Name        string    `form:"name" validate:"omitempty,max=100"`
}

// Validate validates the search doctors request
func (r *SearchDoctorsRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
