package pagination

import (
	"drgo/internal/domain"

	"github.com/huandu/go-sqlbuilder"
)

type Result[T domain.ModelEntity] struct {
	Items []T         `json:"items"`
	Meta  interface{} `json:"meta"`
}

type Paginator[T domain.ModelEntity] interface {
	Paginate(sb *sqlbuilder.SelectBuilder) (string, []interface{})
	GetMeta(totalCount int) interface{}
	CreatePaginationResult(items []T, totalCount int) *Result[T]
}

type Params interface {
	Validate() error
}
