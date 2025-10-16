package query_builder

import (
	"drgo/internal/filter"

	"github.com/huandu/go-sqlbuilder"
)

type QueryBuilder interface {
	WithFilters(filters filter.Filters) QueryBuilder
	WithOrderBy(orderBy string) QueryBuilder
	Build() *sqlbuilder.SelectBuilder
	CountBuilder() *sqlbuilder.SelectBuilder
}
