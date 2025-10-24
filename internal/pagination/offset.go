package pagination

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/shayesteh1hs/DrAppointment/internal/domain"

	"github.com/huandu/go-sqlbuilder"
)

type LimitOffsetParams struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	Limit     int    `form:"limit,default=10" binding:"min=1,max=100"`
	BaseURL   string `form:"-"`
	validated bool
}

func (p *LimitOffsetParams) Validate() error {
	if p.BaseURL == "" {
		return errors.New("base url is required")
	}

	p.validated = true
	return nil
}

func (p *LimitOffsetParams) IsValidated() bool {
	return p.validated
}

type LimitOffsetPaginator[T domain.ModelEntity] struct {
	params LimitOffsetParams
}

func NewLimitOffsetPaginator[T domain.ModelEntity](params LimitOffsetParams) *LimitOffsetPaginator[T] {
	return &LimitOffsetPaginator[T]{params: params}
}

func (p *LimitOffsetPaginator[T]) getOffset() int {
	return p.params.Limit * (p.params.Page - 1)
}

func (p *LimitOffsetPaginator[T]) Paginate(sb *sqlbuilder.SelectBuilder) error {
	if !p.params.IsValidated() {
		return errors.New("params should be validated before paginating")
	}

	sb.Limit(p.params.Limit)
	sb.Offset(p.getOffset())

	return nil
}

func (p *LimitOffsetPaginator[T]) CreatePaginationResult(items []T, totalCount int) (*Result[T], error) {
	if !p.params.IsValidated() {
		return nil, errors.New("params should be validated before paginating")
	}
	result := &Result[T]{
		Items:      items,
		TotalCount: totalCount,
	}

	totalPages := (totalCount + p.params.Limit - 1) / p.params.Limit

	if p.params.Page > 1 {
		prevPage := p.params.Page - 1
		prevURL, err := p.buildURL(prevPage)
		if err != nil {
			return nil, err
		}
		result.Previous = &prevURL
	}

	if p.params.Page < totalPages {
		nextPage := p.params.Page + 1
		nextURL, err := p.buildURL(nextPage)
		if err != nil {
			return nil, err
		}
		result.Next = &nextURL
	}

	return result, nil
}

func (p *LimitOffsetPaginator[T]) buildURL(page int) (string, error) {
	if p.params.BaseURL == "" {
		return "", errors.New("base url is required")
	}

	// Parse existing URL to preserve query parameters
	u, err := url.Parse(p.params.BaseURL)
	if err != nil {
		log.Printf("failed to parse base URL: %v", err)
		return "", errors.New("failed to parse base url")
	}

	params := u.Query()
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("limit", fmt.Sprintf("%d", p.params.Limit))

	u.RawQuery = params.Encode()
	return u.String(), nil
}
