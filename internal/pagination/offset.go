package pagination

import (
	"fmt"

	"drgo/internal/domain"

	"github.com/huandu/go-sqlbuilder"
)

type LimitOffsetParams struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=10" binding:"min=1,max=100"`
}

func (p *LimitOffsetParams) Validate() error {
	if p.Page < 1 {
		return fmt.Errorf("page must be >= 1")
	}
	if p.Limit < 1 || p.Limit > 100 {
		return fmt.Errorf("limit must be between 1 and 100")
	}
	return nil
}

type LimitOffsetMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

type LimitOffsetPaginator[T domain.ModelEntity] struct {
	params LimitOffsetParams
}

func NewLimitOffsetPaginator[T domain.ModelEntity](params LimitOffsetParams) LimitOffsetPaginator[T] {
	return LimitOffsetPaginator[T]{params: params}
}

func (p *LimitOffsetPaginator[T]) Paginate(sb *sqlbuilder.SelectBuilder) (string, []interface{}) {
	offset := (p.params.Page - 1) * p.params.Limit
	sb.Limit(p.params.Limit)
	sb.Offset(offset)

	query, args := sb.Build()

	return query, args
}

func (p *LimitOffsetPaginator[T]) GetMeta(totalCount int) LimitOffsetMeta {
	totalPages := (totalCount + p.params.Limit - 1) / p.params.Limit
	return LimitOffsetMeta{
		Page:       p.params.Page,
		Limit:      p.params.Limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}
}

func (p *LimitOffsetPaginator[T]) CreatePaginationResult(items []T, totalCount int) *Result[T] {
	meta := p.GetMeta(totalCount)

	return &Result[T]{
		Items: items,
		Meta:  meta,
	}
}
