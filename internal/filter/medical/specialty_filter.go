package medical

import (
	"github.com/go-playground/validator/v10"
	"github.com/huandu/go-sqlbuilder"
)

type SpecialtyFilter struct {
	Name string `form:"name" validate:"omitempty"`
}

func (f SpecialtyFilter) Validate() error {
	validate := validator.New()
	return validate.Struct(f)
}

func (f SpecialtyFilter) Apply(sb *sqlbuilder.SelectBuilder) *sqlbuilder.SelectBuilder {
	if f.Name != "" {
		sb.Where(sb.Equal("name", f.Name))
	}
	return sb
}
