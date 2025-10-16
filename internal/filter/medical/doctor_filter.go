package medical

import (
	"github.com/go-playground/validator/v10"
	"github.com/huandu/go-sqlbuilder"
)

type DoctorQueryParam struct {
	Name      string `form:"name"`
	Specialty string `form:"specialty"`
}

func (f DoctorQueryParam) Validate() error {
	validate := validator.New()
	return validate.Struct(f)
}
func (f DoctorQueryParam) Apply(sb *sqlbuilder.SelectBuilder) *sqlbuilder.SelectBuilder {
	if f.Name != "" {
		sb.Where(sb.Equal("name", f.Name))
	}
	return sb
}
