package medical

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type SearchDoctorsRequest struct {
	SpecialtyID int    `form:"specialty_id" validate:"omitempty"`
	Name        string `form:"name" validate:"omitempty"`
}

// Validate validates the search doctors request
func (r *SearchDoctorsRequest) Validate() error {
	return validate.Struct(r)
}
