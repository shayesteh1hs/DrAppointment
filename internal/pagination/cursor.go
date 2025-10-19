package pagination

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"drgo/internal/domain"

	"github.com/huandu/go-sqlbuilder"
)

type CursorParams struct {
	Cursor   string `form:"cursor"`
	Ordering string `form:"ordering,default=asc"`
	Limit    int    `form:"limit,default=10" binding:"min=1,max=100"`
	BaseURL  string `form:"-"`
}

func (p *CursorParams) Validate() error {
	p.Ordering = strings.ToLower(p.Ordering)
	if p.Ordering != "asc" && p.Ordering != "desc" {
		return fmt.Errorf("ordering must be either 'asc' or 'desc'")
	}

	if p.Cursor != "" {
		if _, err := decodeCursor(p.Cursor); err != nil {
			return fmt.Errorf("invalid cursor: %w", err)
		}
	}

	return nil
}

func (p *CursorParams) IsForward() bool {
	return p.Ordering == "asc"
}

func (p *CursorParams) IsBackward() bool {
	return p.Ordering == "desc"
}

type CursorPaginator[T domain.ModelEntity] struct {
	params CursorParams
}

func NewCursorPaginator[T domain.ModelEntity](params CursorParams) *CursorPaginator[T] {
	return &CursorPaginator[T]{params: params}
}

func (p *CursorPaginator[T]) Paginate(sb *sqlbuilder.SelectBuilder) *sqlbuilder.SelectBuilder {
	// Fetch one extra item to determine if there's a next/previous page
	sb.Limit(p.params.Limit + 1)

	// Apply cursor condition if provided
	if p.params.Cursor != "" {
		cursorID, _ := decodeCursor(p.params.Cursor)

		if p.params.IsForward() {
			// Forward pagination: id > cursor
			sb.Where(sb.GreaterThan("id", cursorID))
		} else {
			// Backward pagination: id < cursor
			sb.Where(sb.LessThan("id", cursorID))
		}
	}

	if p.params.IsForward() {
		sb.OrderByAsc("id")
	} else {
		sb.OrderByDesc("id")
	}

	return sb
}

func (p *CursorPaginator[T]) CreatePaginationResult(items []T, totalCount int) *Result[T] {
	result := &Result[T]{
		Items:      items,
		TotalCount: totalCount,
	}

	hasMore := len(items) > p.params.Limit
	if hasMore {
		// Remove the extra item
		items = items[:p.params.Limit]
		result.Items = items
	}

	// If paginating backward, reverse the items to maintain correct order
	if p.params.IsBackward() {
		reverseSlice(result.Items)
	}

	if len(result.Items) == 0 {
		return result
	}

	firstItem := result.Items[0]
	lastItem := result.Items[len(result.Items)-1]

	firstID := firstItem.GetId()
	lastID := lastItem.GetId()

	// Generate previous link (backward pagination from first item)
	if p.params.Cursor != "" || p.params.IsBackward() {
		if prevURL := p.buildURL(firstID, "desc"); prevURL != "" {
			result.Previous = &prevURL
		}
	}

	// Generate next link (forward pagination from last item)
	if hasMore || p.params.IsBackward() {
		if nextURL := p.buildURL(lastID, "asc"); nextURL != "" {
			result.Next = &nextURL
		}
	}

	return result
}

func (p *CursorPaginator[T]) buildURL(id string, ordering string) string {
	if p.params.BaseURL == "" {
		return ""
	}

	cursor := encodeCursor(id)
	params := url.Values{}
	params.Set("cursor", cursor)
	params.Set("ordering", ordering)
	params.Set("limit", fmt.Sprintf("%d", p.params.Limit))

	return fmt.Sprintf("%s?%s", p.params.BaseURL, params.Encode())
}

func encodeCursor(value interface{}) string {
	str := fmt.Sprintf("%v", value)
	return base64.RawURLEncoding.EncodeToString([]byte(str))
}

func decodeCursor(cursor string) (string, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return "", fmt.Errorf("failed to decode cursor: %w", err)
	}
	return string(decoded), nil
}

func reverseSlice[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
