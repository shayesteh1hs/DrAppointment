package pagination

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"testing"

	"github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
)

func TestCursorParams_Validate_ValidParams(t *testing.T) {
	params := CursorParams{
		Cursor:   "",
		Ordering: "asc",
		Limit:    10,
		BaseURL:  "http://example.com/api",
	}

	err := params.Validate()

	assert.Nil(t, err)
	assert.True(t, params.IsValidated())
}

func TestCursorParams_Validate_ValidParamsWithCursor(t *testing.T) {
	validCursor := encodeCursor("123")
	params := CursorParams{
		Cursor:   validCursor,
		Ordering: "desc",
		Limit:    20,
		BaseURL:  "http://example.com/api",
	}

	err := params.Validate()

	assert.Nil(t, err)
	assert.True(t, params.IsValidated())
}

func TestCursorParams_Validate_InvalidOrdering(t *testing.T) {
	params := CursorParams{
		Cursor:   "",
		Ordering: "invalid",
		Limit:    10,
		BaseURL:  "http://example.com/api",
	}

	err := params.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ordering must be either 'asc' or 'desc'")
	assert.False(t, params.IsValidated())
}

func TestCursorParams_Validate_InvalidCursor(t *testing.T) {
	params := CursorParams{
		Cursor:   "invalid-cursor!!!",
		Ordering: "asc",
		Limit:    10,
		BaseURL:  "http://example.com/api",
	}

	err := params.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid cursor")
	assert.False(t, params.IsValidated())
}

func TestCursorParams_Validate_CaseInsensitiveOrdering(t *testing.T) {
	params := CursorParams{
		Cursor:   "",
		Ordering: "ASC",
		Limit:    10,
		BaseURL:  "http://example.com/api",
	}

	err := params.Validate()

	assert.Nil(t, err)
	assert.Equal(t, "asc", params.Ordering)
	assert.True(t, params.IsValidated())
}

func TestCursorParams_IsForward(t *testing.T) {
	params := CursorParams{Ordering: "asc"}
	assert.True(t, params.IsForward())

	params.Ordering = "desc"
	assert.False(t, params.IsForward())
}

func TestCursorParams_IsBackward(t *testing.T) {
	params := CursorParams{Ordering: "desc"}
	assert.True(t, params.IsBackward())

	params.Ordering = "asc"
	assert.False(t, params.IsBackward())
}

func TestCursorPaginator_Paginate_ForwardWithoutCursor(t *testing.T) {
	params := CursorParams{
		Cursor:   "",
		Ordering: "asc",
		Limit:    10,
		BaseURL:  "http://example.com/api",
	}

	// Validate the params
	_ = params.Validate()

	paginator := NewCursorPaginator[mockEntity](params)
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").From("test")

	err := paginator.Paginate(sb)

	assert.NoError(t, err)

	// Check that LIMIT and ORDER BY were applied correctly
	sql, args := sb.Build()
	assert.Contains(t, sql, "LIMIT $1")
	assert.Contains(t, sql, "ORDER BY id ASC")

	// Check that the correct limit value is in the args (limit + 1 for hasMore check)
	assert.Equal(t, params.Limit+1, args[0])
}

func TestCursorPaginator_Paginate_ForwardWithCursor(t *testing.T) {
	cursor := encodeCursor("123")
	params := CursorParams{
		Cursor:   cursor,
		Ordering: "asc",
		Limit:    10,
		BaseURL:  "http://example.com/api",
	}

	// Validate the params
	_ = params.Validate()

	paginator := NewCursorPaginator[mockEntity](params)
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").From("test")

	err := paginator.Paginate(sb)

	assert.NoError(t, err)

	// Check that WHERE, LIMIT and ORDER BY were applied correctly
	sql, args := sb.Build()
	assert.Contains(t, sql, "WHERE id > $1")
	assert.Contains(t, sql, "LIMIT $2")
	assert.Contains(t, sql, "ORDER BY id ASC")

	// Check that the correct values are in the args
	assert.Equal(t, "123", args[0])
	assert.Equal(t, params.Limit+1, args[1])
}

func TestCursorPaginator_Paginate_BackwardWithCursor(t *testing.T) {
	cursor := encodeCursor("123")
	params := CursorParams{
		Cursor:   cursor,
		Ordering: "desc",
		Limit:    10,
		BaseURL:  "http://example.com/api",
	}

	// Validate the params
	_ = params.Validate()

	paginator := NewCursorPaginator[mockEntity](params)
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").From("test")

	err := paginator.Paginate(sb)

	assert.NoError(t, err)

	// Check that WHERE, LIMIT and ORDER BY were applied correctly
	sql, args := sb.Build()
	assert.Contains(t, sql, "WHERE id < $1")
	assert.Contains(t, sql, "LIMIT $2")
	assert.Contains(t, sql, "ORDER BY id DESC")

	// Check that the correct values are in the args
	assert.Equal(t, "123", args[0])
	assert.Equal(t, params.Limit+1, args[1])
}

func TestCursorPaginator_Paginate_WithoutValidation(t *testing.T) {
	params := CursorParams{
		Cursor:   "",
		Ordering: "asc",
		Limit:    10,
		BaseURL:  "http://example.com/api",
	}

	// Do not validate the params

	paginator := NewCursorPaginator[mockEntity](params)
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").From("test")

	err := paginator.Paginate(sb)

	assert.Error(t, err)
	assert.Equal(t, "params should be validated before paginating", err.Error())
}

func TestCursorPaginator_CreatePaginationResult_ForwardWithMoreItems(t *testing.T) {
	baseURL := "http://example.com/api"
	params := CursorParams{
		Cursor:   "",
		Ordering: "asc",
		Limit:    10,
		BaseURL:  baseURL,
	}

	// Generate 11 items (limit + 1 to simulate hasMore)
	items := generateMockItems(11)
	totalCount := 25

	validationErr := params.Validate()
	assert.NoError(t, validationErr, "Validation should succeed for valid params")

	paginator := NewCursorPaginator[mockEntity](params)
	result, resultErr := paginator.CreatePaginationResult(items, totalCount)

	// Check that CreatePaginationResult succeeds
	assert.NoError(t, resultErr, "CreatePaginationResult should not return an error for valid params")

	// Check result properties
	assert.Equal(t, 10, len(result.Items)) // Should be trimmed to limit
	assert.Equal(t, totalCount, result.TotalCount)

	// Check previous link (should be nil for first page without cursor)
	assert.Nil(t, result.Previous)

	// Check next link (should exist because we have more items)
	assert.NotNil(t, result.Next)
	nextURL, err := url.Parse(*result.Next)
	assert.NoError(t, err)

	// Decode the cursor to verify it's the last item's ID
	cursor := nextURL.Query().Get("cursor")
	decodedCursor, err := decodeCursor(cursor)
	assert.NoError(t, err)
	assert.Equal(t, "9", decodedCursor) // Last item ID (0-indexed)
	assert.Equal(t, "asc", nextURL.Query().Get("ordering"))
}

func TestCursorPaginator_CreatePaginationResult_BackwardWithMoreItems(t *testing.T) {
	baseURL := "http://example.com/api"
	cursor := encodeCursor("20")
	params := CursorParams{
		Cursor:   cursor,
		Ordering: "desc",
		Limit:    10,
		BaseURL:  baseURL,
	}

	// Generate 11 items (limit + 1 to simulate hasMore)
	items := generateMockItemsReverse(11, 19) // Items 19, 18, 17, ..., 9
	totalCount := 25

	validationErr := params.Validate()
	assert.NoError(t, validationErr, "Validation should succeed for valid params")

	paginator := NewCursorPaginator[mockEntity](params)
	result, resultErr := paginator.CreatePaginationResult(items, totalCount)

	// Check that CreatePaginationResult succeeds
	assert.NoError(t, resultErr, "CreatePaginationResult should not return an error for valid params")

	// Check result properties
	assert.Equal(t, 10, len(result.Items)) // Should be trimmed to limit
	assert.Equal(t, totalCount, result.TotalCount)

	// Items should be reversed for backward pagination
	assert.Equal(t, "10", result.Items[0].GetId()) // First item after reversal
	assert.Equal(t, "19", result.Items[9].GetId()) // Last item after reversal

	// Check previous link (should exist for backward pagination)
	assert.NotNil(t, result.Previous)
	prevURL, err := url.Parse(*result.Previous)
	assert.NoError(t, err)
	assert.Equal(t, "desc", prevURL.Query().Get("ordering"))

	// Check next link (should exist for backward pagination)
	assert.NotNil(t, result.Next)
	nextURL, err := url.Parse(*result.Next)
	assert.NoError(t, err)
	assert.Equal(t, "asc", nextURL.Query().Get("ordering"))
}

func TestCursorPaginator_CreatePaginationResult_NoMoreItems(t *testing.T) {
	baseURL := "http://example.com/api"
	params := CursorParams{
		Cursor:   "",
		Ordering: "asc",
		Limit:    10,
		BaseURL:  baseURL,
	}

	// Generate exactly limit items (no more items)
	items := generateMockItems(10)
	totalCount := 10

	validationErr := params.Validate()
	assert.NoError(t, validationErr, "Validation should succeed for valid params")

	paginator := NewCursorPaginator[mockEntity](params)
	result, resultErr := paginator.CreatePaginationResult(items, totalCount)

	// Check that CreatePaginationResult succeeds
	assert.NoError(t, resultErr, "CreatePaginationResult should not return an error for valid params")

	// Check result properties
	assert.Equal(t, items, result.Items)
	assert.Equal(t, totalCount, result.TotalCount)

	// Check previous link (should be nil for first page without cursor)
	assert.Nil(t, result.Previous)

	// Check next link (should be nil because no more items)
	assert.Nil(t, result.Next)
}

func TestCursorPaginator_CreatePaginationResult_EmptyItems(t *testing.T) {
	baseURL := "http://example.com/api"
	params := CursorParams{
		Cursor:   "",
		Ordering: "asc",
		Limit:    10,
		BaseURL:  baseURL,
	}

	items := []mockEntity{}
	totalCount := 0

	validationErr := params.Validate()
	assert.NoError(t, validationErr, "Validation should succeed for valid params")

	paginator := NewCursorPaginator[mockEntity](params)
	result, resultErr := paginator.CreatePaginationResult(items, totalCount)

	// Check that CreatePaginationResult succeeds
	assert.NoError(t, resultErr, "CreatePaginationResult should not return an error for valid params")

	// Check result properties
	assert.Equal(t, items, result.Items)
	assert.Equal(t, totalCount, result.TotalCount)

	// Check links (should be nil for empty results)
	assert.Nil(t, result.Previous)
	assert.Nil(t, result.Next)
}

func TestCursorPaginator_CreatePaginationResult_WithoutValidation(t *testing.T) {
	params := CursorParams{
		Cursor:   "",
		Ordering: "asc",
		Limit:    10,
		BaseURL:  "http://example.com/api",
	}

	// Do not validate the params
	items := generateMockItems(10)
	totalCount := 25

	paginator := NewCursorPaginator[mockEntity](params)
	result, resultErr := paginator.CreatePaginationResult(items, totalCount)

	assert.Error(t, resultErr)
	assert.Equal(t, "params should be validated before paginating", resultErr.Error())
	assert.Nil(t, result)
}

func TestCursorPaginator_BuildURL_WithExistingQueryParams(t *testing.T) {
	params := CursorParams{
		Cursor:   "",
		Ordering: "asc",
		Limit:    10,
		BaseURL:  "http://example.com/api?filter=test",
	}

	_ = params.Validate()
	paginator := NewCursorPaginator[mockEntity](params)

	result, err := paginator.buildURL("123", "asc")
	assert.NoError(t, err)

	assert.Contains(t, result, "http://example.com/api?")
	assert.Contains(t, result, "ordering=asc")
	assert.Contains(t, result, "limit=10")
	assert.Contains(t, result, "filter=test")

	// Verify cursor is properly encoded
	u, err := url.Parse(result)
	assert.NoError(t, err)
	cursor := u.Query().Get("cursor")
	decodedCursor, err := decodeCursor(cursor)
	assert.NoError(t, err)
	assert.Equal(t, "123", decodedCursor)
}

func TestCursorPaginator_BuildURL_EmptyBaseURL(t *testing.T) {
	params := CursorParams{
		Cursor:   "",
		Ordering: "asc",
		Limit:    10,
		BaseURL:  "",
	}

	paginator := NewCursorPaginator[mockEntity](params)

	result, err := paginator.buildURL("123", "asc")
	assert.Error(t, err)
	assert.Equal(t, "base url is required", err.Error())
	assert.Empty(t, result)
}

func TestEncodeCursor(t *testing.T) {
	testCases := []struct {
		input    interface{}
		expected string
	}{
		{"123", base64.RawURLEncoding.EncodeToString([]byte("123"))},
		{456, base64.RawURLEncoding.EncodeToString([]byte("456"))},
		{"", base64.RawURLEncoding.EncodeToString([]byte(""))},
	}

	for _, tc := range testCases {
		result := encodeCursor(tc.input)
		assert.Equal(t, tc.expected, result)
	}
}

func TestDecodeCursor(t *testing.T) {
	testCases := []struct {
		input       string
		expected    string
		shouldError bool
	}{
		{base64.RawURLEncoding.EncodeToString([]byte("123")), "123", false},
		{base64.RawURLEncoding.EncodeToString([]byte("456")), "456", false},
		{base64.RawURLEncoding.EncodeToString([]byte("")), "", false},
		{"invalid-base64!!!", "", true},
	}

	for _, tc := range testCases {
		result, err := decodeCursor(tc.input)
		if tc.shouldError {
			assert.Error(t, err)
			assert.Empty(t, result)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		}
	}
}

func TestReverseSlice(t *testing.T) {
	testCases := []struct {
		input    []mockEntity
		expected []mockEntity
	}{
		{
			[]mockEntity{{ID: "1"}, {ID: "2"}, {ID: "3"}},
			[]mockEntity{{ID: "3"}, {ID: "2"}, {ID: "1"}},
		},
		{
			[]mockEntity{{ID: "1"}, {ID: "2"}},
			[]mockEntity{{ID: "2"}, {ID: "1"}},
		},
		{
			[]mockEntity{{ID: "1"}},
			[]mockEntity{{ID: "1"}},
		},
		{
			[]mockEntity{},
			[]mockEntity{},
		},
	}

	for _, tc := range testCases {
		// Make a copy to avoid modifying the original
		input := make([]mockEntity, len(tc.input))
		copy(input, tc.input)

		reverseSlice(input)
		assert.Equal(t, tc.expected, input)
	}
}

// Helper function to generate mock items with reverse order
func generateMockItemsReverse(count int, startID int) []mockEntity {
	items := make([]mockEntity, count)
	for i := 0; i < count; i++ {
		items[i] = mockEntity{ID: fmt.Sprintf("%d", startID-i)}
	}
	return items
}
