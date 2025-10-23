package pagination

import (
	"errors"
	"net/url"
	"strconv"
	"testing"

	"github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
)

type mockEntity struct {
	ID string
}

func (m mockEntity) GetId() string {
	return m.ID
}

func TestLimitOffsetParams_Validate_ValidParams(t *testing.T) {
	params := LimitOffsetParams{
		Page:    1,
		Limit:   10,
		BaseURL: "http://example.com/api",
	}

	err := params.Validate()

	assert.Nil(t, err)
	assert.True(t, params.IsValidated())
}

func TestLimitOffsetParams_Validate_MissingBaseURL(t *testing.T) {
	params := LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}

	err := params.Validate()

	expectedErr := errors.New("base url is required")
	assert.Equal(t, expectedErr, err)
}

func TestLimitOffsetPaginator_Paginate_WithValidation(t *testing.T) {
	params := LimitOffsetParams{
		Page:    2,
		Limit:   10,
		BaseURL: "http://example.com/api",
	}

	// Validate the params
	_ = params.Validate()

	paginator := NewLimitOffsetPaginator[mockEntity](params)
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").From("test")

	err := paginator.Paginate(sb)

	assert.NoError(t, err)

	// Check that LIMIT and OFFSET were applied correctly
	sql, args := sb.Build()
	assert.Contains(t, sql, "LIMIT $1")
	assert.Contains(t, sql, "OFFSET $2")

	// Check that the correct values are in the args
	assert.Equal(t, params.Limit, args[0])
	assert.Equal(t, paginator.getOffset(), args[1])
}

func TestLimitOffsetPaginator_Paginate_WithoutValidation(t *testing.T) {
	params := LimitOffsetParams{
		Page:    2,
		Limit:   10,
		BaseURL: "http://example.com/api",
	}

	// Do not validate the params

	paginator := NewLimitOffsetPaginator[mockEntity](params)
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").From("test")

	err := paginator.Paginate(sb)

	assert.Error(t, err)
	assert.Equal(t, "params should be validated before paginating", err.Error())
}

func TestLimitOffsetPaginator_CreatePaginationResult_FirstPageWithMorePages(t *testing.T) {
	baseURL := "http://example.com/api"
	params := LimitOffsetParams{
		Page:    1,
		Limit:   10,
		BaseURL: baseURL,
	}

	items := generateMockItems(10)
	totalCount := 25

	validationErr := params.Validate()
	assert.NoError(t, validationErr, "Validation should succeed for valid params")

	paginator := NewLimitOffsetPaginator[mockEntity](params)
	result, resultErr := paginator.CreatePaginationResult(items, totalCount)

	// Check that CreatePaginationResult succeeds
	assert.NoError(t, resultErr, "CreatePaginationResult should not return an error for valid params")

	// Check result properties
	assert.Equal(t, items, result.Items)
	assert.Equal(t, totalCount, result.TotalCount)

	// Check previous link (should be nil for first page)
	assert.Nil(t, result.Previous)

	// Check next link (should exist for first page when there are more pages)
	assert.NotNil(t, result.Next)
	nextURL, err := url.Parse(*result.Next)
	assert.NoError(t, err)
	assert.Equal(t, strconv.Itoa(params.Page+1), nextURL.Query().Get("page"))
}

func TestLimitOffsetPaginator_CreatePaginationResult_MiddlePage(t *testing.T) {
	baseURL := "http://example.com/api"
	params := LimitOffsetParams{
		Page:    2,
		Limit:   10,
		BaseURL: baseURL,
	}

	items := generateMockItems(10)
	totalCount := 25

	validationErr := params.Validate()
	assert.NoError(t, validationErr, "Validation should succeed for valid params")

	paginator := NewLimitOffsetPaginator[mockEntity](params)
	result, resultErr := paginator.CreatePaginationResult(items, totalCount)

	// Check that CreatePaginationResult succeeds
	assert.NoError(t, resultErr, "CreatePaginationResult should not return an error for valid params")

	// Check result properties
	assert.Equal(t, items, result.Items)
	assert.Equal(t, totalCount, result.TotalCount)

	// Check previous link (should exist for middle page)
	assert.NotNil(t, result.Previous)
	prevURL, err := url.Parse(*result.Previous)
	assert.NoError(t, err)
	assert.Equal(t, strconv.Itoa(params.Page-1), prevURL.Query().Get("page"))

	// Check next link (should exist for middle page)
	assert.NotNil(t, result.Next)
	nextURL, err := url.Parse(*result.Next)
	assert.NoError(t, err)
	assert.Equal(t, strconv.Itoa(params.Page+1), nextURL.Query().Get("page"))
}

func TestLimitOffsetPaginator_CreatePaginationResult_LastPage(t *testing.T) {
	baseURL := "http://example.com/api"
	params := LimitOffsetParams{
		Page:    3,
		Limit:   10,
		BaseURL: baseURL,
	}

	items := generateMockItems(5)
	totalCount := 25

	validationErr := params.Validate()
	assert.NoError(t, validationErr, "Validation should succeed for valid params")

	paginator := NewLimitOffsetPaginator[mockEntity](params)
	result, resultErr := paginator.CreatePaginationResult(items, totalCount)

	// Check that CreatePaginationResult succeeds
	assert.NoError(t, resultErr, "CreatePaginationResult should not return an error for valid params")

	// Check result properties
	assert.Equal(t, items, result.Items)
	assert.Equal(t, totalCount, result.TotalCount)

	// Check previous link (should exist for last page)
	assert.NotNil(t, result.Previous)
	prevURL, err := url.Parse(*result.Previous)
	assert.NoError(t, err)
	assert.Equal(t, strconv.Itoa(params.Page-1), prevURL.Query().Get("page"))

	// Check next link (should not exist for last page)
	assert.Nil(t, result.Next)
}

func TestLimitOffsetPaginator_CreatePaginationResult_SinglePage(t *testing.T) {
	baseURL := "http://example.com/api"
	params := LimitOffsetParams{
		Page:    1,
		Limit:   10,
		BaseURL: baseURL,
	}

	items := generateMockItems(5)
	totalCount := 5

	validationErr := params.Validate()
	assert.NoError(t, validationErr, "Validation should succeed for valid params")

	paginator := NewLimitOffsetPaginator[mockEntity](params)
	result, resultErr := paginator.CreatePaginationResult(items, totalCount)

	// Check that CreatePaginationResult succeeds
	assert.NoError(t, resultErr, "CreatePaginationResult should not return an error for valid params")

	// Check result properties
	assert.Equal(t, items, result.Items)
	assert.Equal(t, totalCount, result.TotalCount)

	// Check previous link (should not exist for single page)
	assert.Nil(t, result.Previous)

	// Check next link (should not exist for single page)
	assert.Nil(t, result.Next)
}

func TestBuildURL_WithExistingQueryParams(t *testing.T) {
	params := LimitOffsetParams{
		Page:    2,
		Limit:   10,
		BaseURL: "http://example.com/api?filter=test",
	}

	_ = params.Validate()
	paginator := NewLimitOffsetPaginator[mockEntity](params)

	result, err := paginator.buildURL(2)

	assert.NoError(t, err)
	assert.Contains(t, result, "http://example.com/api?")
	assert.Contains(t, result, "page=2")
	assert.Contains(t, result, "limit=10")
	assert.Contains(t, result, "filter=test")
}

func TestBuildURL_EmptyURL(t *testing.T) {
	params := LimitOffsetParams{
		Page:    2,
		Limit:   10,
		BaseURL: "",
	}

	paginator := NewLimitOffsetPaginator[mockEntity](params)

	result, err := paginator.buildURL(2)

	assert.Error(t, err)
	assert.Empty(t, result)
}

func generateMockItems(count int) []mockEntity {
	items := make([]mockEntity, count)
	for i := 0; i < count; i++ {
		items[i] = mockEntity{ID: strconv.Itoa(i)}
	}
	return items
}
