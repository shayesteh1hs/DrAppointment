package medical

import (
	"github.com/go-playground/validator/v10"
)

type SearchDoctorsRequest struct {
	Specialty string `form:"specialty" validate:"omitempty"`
	Name      string `form:"name" validate:"omitempty"`
}

// Validate validates the search doctors request
func (r *SearchDoctorsRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
