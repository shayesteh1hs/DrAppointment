package medical

import (
	filter "drgo/internal/filter/medical"

	"github.com/huandu/go-sqlbuilder"
)

type DoctorQueryBuilder struct {
	sb      *sqlbuilder.SelectBuilder
	filters filter.DoctorQueryParam
}

func NewDoctorQueryBuilder() *DoctorQueryBuilder {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select(
		"doctors.id",
		"doctors.name",
		"doctors.specialty_id",
		"doctors.phone_number",
		"doctors.avatar_url",
		"doctors.description",
		"doctors.created_at",
		"doctors.updated_at",
	)
	sb.From("doctors")
	return &DoctorQueryBuilder{
		sb:      sb,
		filters: filter.DoctorQueryParam{},
	}
}

func (qb *DoctorQueryBuilder) WithFilters(filters filter.DoctorQueryParam) *DoctorQueryBuilder {
	qb.filters = filters
	return qb
}

func (qb *DoctorQueryBuilder) WithOrderBy(orderBy string) *DoctorQueryBuilder {
	qb.sb.OrderBy(orderBy)
	return qb
}

func (qb *DoctorQueryBuilder) Build() *sqlbuilder.SelectBuilder {
	qb.filters.Apply(qb.sb)
	return qb.sb
}

func (qb *DoctorQueryBuilder) CountBuilder() *sqlbuilder.SelectBuilder {
	countSb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	countSb.Select("COUNT(*)")
	countSb.From("doctors")
	qb.filters.Apply(countSb)
	return countSb
}
