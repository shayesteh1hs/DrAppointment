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
	Paginate(sb *sqlbuilder.SelectBuilder) error
	CreatePaginationResult(items []T, totalCount int) *Result[T]
	isValidated() bool // to check if the params are validated before paginating
}

type Params interface {
	Validate() error
}
