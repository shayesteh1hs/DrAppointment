package medical

import (
	"strings"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
)

type DoctorQueryParam struct {
	Name        string    `form:"name"`
	SpecialtyID uuid.UUID `form:"specialty_id"`
}

func (f DoctorQueryParam) Validate() error {
	return nil
}
func (f DoctorQueryParam) Apply(sb *sqlbuilder.SelectBuilder) *sqlbuilder.SelectBuilder {
	trimmedName := strings.TrimSpace(f.Name)
	if trimmedName != "" {
		sb.Where(sb.Like("name", "%"+trimmedName+"%"))
	}
	if f.SpecialtyID != uuid.Nil {
		sb.Where(sb.Equal("specialty_id", f.SpecialtyID.String()))
	}
	return sb
}
