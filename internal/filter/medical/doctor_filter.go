package medical

import (
	"github.com/go-playground/validator/v10"
	"github.com/huandu/go-sqlbuilder"
)

type DoctorQueryParam struct {
	Name        string `form:"name"`
	SpecialtyID int    `form:"specialty_id"`
}

func (f DoctorQueryParam) Validate() error {
	validate := validator.New()
	return validate.Struct(f)
}
func (f DoctorQueryParam) Apply(sb *sqlbuilder.SelectBuilder) *sqlbuilder.SelectBuilder {
	if f.Name != "" {
		sb.Where(sb.Like("name", "%"+f.Name+"%"))
	}
	if f.SpecialtyID != uuid.Nil {
		sb.Where(sb.Equal("specialty_id", f.SpecialtyID.String()))
	}
	return sb
}
