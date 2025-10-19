package pagination

import (
	"github.com/shayesteh1hs/DrAppointment/internal/domain"

	"github.com/huandu/go-sqlbuilder"
)

type Result[T domain.ModelEntity] struct {
	Items      []T     `json:"items"`
	TotalCount int     `json:"total_count"`
	Previous   *string `json:"previous,omitempty"`
	Next       *string `json:"next,omitempty"`
}

type Paginator[T domain.ModelEntity] interface {
	Paginate(sb *sqlbuilder.SelectBuilder) *sqlbuilder.SelectBuilder
	CreatePaginationResult(items []T, totalCount int) *Result[T]
}

type Params interface {
	Validate() error
}
