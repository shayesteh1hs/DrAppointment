package medical

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper functions

func newTestFilter(name string, specialtyID uuid.UUID) DoctorQueryParam {
	return DoctorQueryParam{
		Name:        name,
		SpecialtyID: specialtyID,
	}
}

func newTestSelectBuilder() *sqlbuilder.SelectBuilder {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("id", "name", "specialty_id", "phone_number", "avatar_url", "description", "created_at", "updated_at").
		From("doctors")
	return sb
}

func assertSQLContains(t *testing.T, sql, expected string) {
	t.Helper()
	assert.Contains(t, strings.ToLower(sql), strings.ToLower(expected))
}

func assertSQLNotContains(t *testing.T, sql, notExpected string) {
	t.Helper()
	assert.NotContains(t, strings.ToLower(sql), strings.ToLower(notExpected))
}

// assertParameterExists verifies an exact parameter exists in args with optional wrapping
func assertParameterExists(t *testing.T, args []interface{}, expectedValue string, wrapWith string) bool {
	t.Helper()
	for _, arg := range args {
		if str, ok := arg.(string); ok {
			if wrapWith != "" {
				// Check for wrapped value (e.g., %value%)
				wrapped := wrapWith + expectedValue + wrapWith
				if str == wrapped {
					return true
				}
			} else {
				// Check for exact match
				if str == expectedValue {
					return true
				}
			}
		}
	}
	return false
}

// Table-driven tests

func TestDoctorQueryParam_Apply(t *testing.T) {
	testSpecialtyID := uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479")

	tests := []struct {
		name            string
		filter          DoctorQueryParam
		expectWhere     bool
		expectName      bool
		expectSpecialty bool
		nameValue       string
		specialtyValue  string
	}{
		{
			name:            "empty filter adds no conditions",
			filter:          DoctorQueryParam{},
			expectWhere:     false,
			expectName:      false,
			expectSpecialty: false,
		},
		{
			name:            "name filter only adds LIKE condition",
			filter:          newTestFilter("Smith", uuid.Nil),
			expectWhere:     true,
			expectName:      true,
			expectSpecialty: false,
			nameValue:       "Smith",
		},
		{
			name:            "specialty filter only adds equality condition",
			filter:          newTestFilter("", testSpecialtyID),
			expectWhere:     true,
			expectName:      false,
			expectSpecialty: true,
			specialtyValue:  testSpecialtyID.String(),
		},
		{
			name:            "both filters add AND condition",
			filter:          newTestFilter("John", testSpecialtyID),
			expectWhere:     true,
			expectName:      true,
			expectSpecialty: true,
			nameValue:       "John",
			specialtyValue:  testSpecialtyID.String(),
		},
		{
			name:            "whitespace name is treated as empty",
			filter:          newTestFilter("   ", uuid.Nil),
			expectWhere:     false,
			expectName:      false,
			expectSpecialty: false,
		},
		{
			name:            "nil UUID is treated as empty",
			filter:          newTestFilter("", uuid.Nil),
			expectWhere:     false,
			expectName:      false,
			expectSpecialty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := newTestSelectBuilder()
			result := tt.filter.Apply(sb)
			require.NotNil(t, result)

			sql, args := result.Build()
			require.NotEmpty(t, sql)

			// Verify WHERE clause presence
			if tt.expectWhere {
				assertSQLContains(t, sql, "WHERE")
			} else {
				assertSQLNotContains(t, sql, "WHERE")
			}

			// Verify name LIKE condition
			if tt.expectName {
				assertSQLContains(t, sql, "name LIKE")
				assert.True(t, assertParameterExists(t, args, tt.nameValue, "%"),
					"Expected name parameter %%%s%% not found in args: %v", tt.nameValue, args)
			} else {
				assertSQLNotContains(t, sql, "name LIKE")
			}

			// Verify specialty_id equality condition
			if tt.expectSpecialty {
				assertSQLContains(t, sql, "specialty_id =")
				assert.True(t, assertParameterExists(t, args, tt.specialtyValue, ""),
					"Expected specialty_id parameter %s not found in args: %v", tt.specialtyValue, args)
			} else {
				assertSQLNotContains(t, sql, "specialty_id =")
			}

			// Verify AND is present when both filters exist
			if tt.expectName && tt.expectSpecialty {
				assertSQLContains(t, sql, "AND")
			}

			// Verify argument count
			expectedArgCount := 0
			if tt.expectName {
				expectedArgCount++
			}
			if tt.expectSpecialty {
				expectedArgCount++
			}
			assert.Equal(t, expectedArgCount, len(args),
				"Unexpected number of query arguments. SQL: %s", sql)
		})
	}
}
