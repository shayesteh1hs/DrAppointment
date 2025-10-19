package medical

import (
	"github.com/go-playground/validator/v10"
)

type SearchDoctorsRequest struct {
	SpecialtyID int    `form:"specialty_id" validate:"omitempty,min=1"`
	Name        string `form:"name" validate:"omitempty,max=100"`
}

// Validate validates the search doctors request
func (r *SearchDoctorsRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
