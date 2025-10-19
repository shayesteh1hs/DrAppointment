package pagination

import (
	"fmt"
	"net/url"

	"drgo/internal/domain"

	"github.com/huandu/go-sqlbuilder"
)

type LimitOffsetParams struct {
	Page    int    `form:"page,default=1" binding:"min=1"`
	Limit   int    `form:"limit,default=10" binding:"min=1,max=100"`
	BaseURL string `form:"-"`
}

func (p *LimitOffsetParams) Validate() error {
	return nil
}

type LimitOffsetPaginator[T domain.ModelEntity] struct {
	params LimitOffsetParams
}

func NewLimitOffsetPaginator[T domain.ModelEntity](params LimitOffsetParams) *LimitOffsetPaginator[T] {
	return &LimitOffsetPaginator[T]{params: params}
}

func (p *LimitOffsetPaginator[T]) Paginate(sb *sqlbuilder.SelectBuilder) *sqlbuilder.SelectBuilder {
	offset := (p.params.Page - 1) * p.params.Limit
	sb.Limit(p.params.Limit)
	sb.Offset(offset)
	return sb
}

func (p *LimitOffsetPaginator[T]) CreatePaginationResult(items []T, totalCount int) *Result[T] {
	result := &Result[T]{
		Items:      items,
		TotalCount: totalCount,
	}

	totalPages := (totalCount + p.params.Limit - 1) / p.params.Limit

	if p.params.Page > 1 {
		prevPage := p.params.Page - 1
		prevURL := p.buildURL(prevPage)
		result.Previous = &prevURL
	}

	if p.params.Page < totalPages {
		nextPage := p.params.Page + 1
		nextURL := p.buildURL(nextPage)
		result.Next = &nextURL
	}

	return result
}

func (p *LimitOffsetPaginator[T]) buildURL(page int) string {
	if p.params.BaseURL == "" {
		return ""
	}

	params := url.Values{}
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("limit", fmt.Sprintf("%d", p.params.Limit))

	return fmt.Sprintf("%s?%s", p.params.BaseURL, params.Encode())
}
