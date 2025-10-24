package medical

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/shayesteh1hs/DrAppointment/internal/domain/medical"
	filter "github.com/shayesteh1hs/DrAppointment/internal/filter/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
)

// Helper functions

func setupTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock
}

func newTestDoctor(name string) medical.Doctor {
	now := time.Now().Truncate(time.Second)
	return medical.Doctor{
		ID:          uuid.New(),
		Name:        name,
		SpecialtyID: uuid.New(),
		PhoneNumber: "1234567890",
		AvatarURL:   "avatar.jpg",
		Description: "Test description",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func newTestPaginator(t *testing.T, page, limit int) *pagination.LimitOffsetPaginator[medical.Doctor] {
	t.Helper()
	params := pagination.LimitOffsetParams{
		Page:    page,
		Limit:   limit,
		BaseURL: "http://localhost:8080/api/doctors",
	}
	require.NoError(t, params.Validate())
	return pagination.NewLimitOffsetPaginator[medical.Doctor](params)
}

func mockDoctorRows(doctors ...medical.Doctor) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"id", "name", "specialty_id", "phone_number",
		"avatar_url", "description", "created_at", "updated_at",
	})
	for _, d := range doctors {
		rows.AddRow(d.ID, d.Name, d.SpecialtyID, d.PhoneNumber,
			d.AvatarURL, d.Description, d.CreatedAt, d.UpdatedAt)
	}
	return rows
}

func assertDoctorEqual(t *testing.T, expected, actual medical.Doctor) {
	t.Helper()
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.SpecialtyID, actual.SpecialtyID)
	assert.Equal(t, expected.PhoneNumber, actual.PhoneNumber)
	assert.Equal(t, expected.AvatarURL, actual.AvatarURL)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.CreatedAt.Unix(), actual.CreatedAt.Unix())
	assert.Equal(t, expected.UpdatedAt.Unix(), actual.UpdatedAt.Unix())
}

// Table-driven tests

func TestDoctorRepository_GetAllPaginated(t *testing.T) {
	ctx := context.Background()
	doctor1 := newTestDoctor("Dr. John Smith")
	doctor2 := newTestDoctor("Dr. Jane Johnson")

	// Base query for selecting doctor fields
	baseSelectQuery := `SELECT id, name, specialty_id, phone_number, avatar_url, description, created_at, updated_at FROM doctors`

	tests := []struct {
		name      string
		filter    filter.DoctorQueryParam
		page      int
		limit     int
		mockSetup func(sqlmock.Sqlmock)
		want      []medical.Doctor
		wantErr   string
	}{
		{
			name:   "no filters returns all doctors",
			filter: filter.DoctorQueryParam{},
			page:   1,
			limit:  10,
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(
					baseSelectQuery+` LIMIT \$1 OFFSET \$2`,
				).WithArgs(10, 0).WillReturnRows(mockDoctorRows(doctor1, doctor2))
			},
			want: []medical.Doctor{doctor1, doctor2},
		},
		{
			name:   "name filter",
			filter: filter.DoctorQueryParam{Name: "John"},
			page:   1,
			limit:  10,
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(
					baseSelectQuery+` WHERE name LIKE \$1 LIMIT \$2 OFFSET \$3`,
				).WithArgs("%John%", 10, 0).WillReturnRows(mockDoctorRows(doctor1))
			},
			want: []medical.Doctor{doctor1},
		},
		{
			name:   "specialty filter",
			filter: filter.DoctorQueryParam{SpecialtyID: doctor1.SpecialtyID},
			page:   1,
			limit:  10,
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(
					baseSelectQuery+` WHERE specialty_id = \$1 LIMIT \$2 OFFSET \$3`,
				).WithArgs(doctor1.SpecialtyID.String(), 10, 0).WillReturnRows(mockDoctorRows(doctor1))
			},
			want: []medical.Doctor{doctor1},
		},
		{
			name:   "multiple filters",
			filter: filter.DoctorQueryParam{Name: "John", SpecialtyID: doctor1.SpecialtyID},
			page:   1,
			limit:  10,
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(
					baseSelectQuery+` WHERE name LIKE \$1 AND specialty_id = \$2 LIMIT \$3 OFFSET \$4`,
				).WithArgs("%John%", doctor1.SpecialtyID.String(), 10, 0).WillReturnRows(mockDoctorRows(doctor1))
			},
			want: []medical.Doctor{doctor1},
		},
		{
			name:   "empty result",
			filter: filter.DoctorQueryParam{Name: "NonExistent"},
			page:   1,
			limit:  10,
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(
					baseSelectQuery+` WHERE name LIKE \$1 LIMIT \$2 OFFSET \$3`,
				).WithArgs("%NonExistent%", 10, 0).WillReturnRows(mockDoctorRows())
			},
			want: []medical.Doctor{},
		},
		{
			name:   "database error",
			filter: filter.DoctorQueryParam{},
			page:   1,
			limit:  10,
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(
					baseSelectQuery+` LIMIT \$1 OFFSET \$2`,
				).WithArgs(10, 0).WillReturnError(errors.New("connection lost"))
			},
			wantErr: "connection lost",
		},
		{
			name:   "context cancelled",
			filter: filter.DoctorQueryParam{},
			page:   1,
			limit:  10,
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(
					baseSelectQuery+` LIMIT \$1 OFFSET \$2`,
				).WithArgs(10, 0).WillReturnError(context.Canceled)
			},
			wantErr: "context canceled",
		},
		{
			name:   "pagination page 2",
			filter: filter.DoctorQueryParam{},
			page:   2,
			limit:  5,
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(
					baseSelectQuery+` LIMIT \$1 OFFSET \$2`,
				).WithArgs(5, 5).WillReturnRows(mockDoctorRows(doctor2))
			},
			want: []medical.Doctor{doctor2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupTestDB(t)
			defer db.Close()

			repo := NewDoctorRepository(db)
			paginator := newTestPaginator(t, tt.page, tt.limit)

			tt.mockSetup(mock)

			got, err := repo.GetAllPaginated(ctx, tt.filter, paginator)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, len(tt.want), len(got))
				for i := range tt.want {
					assertDoctorEqual(t, tt.want[i], got[i])
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDoctorRepository_Count(t *testing.T) {
	ctx := context.Background()
	specialtyID := uuid.New()

	tests := []struct {
		name      string
		filter    filter.DoctorQueryParam
		mockSetup func(sqlmock.Sqlmock)
		want      int
		wantErr   string
	}{
		{
			name:   "no filters",
			filter: filter.DoctorQueryParam{},
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT count\(\*\) FROM doctors`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(15))
			},
			want: 15,
		},
		{
			name:   "name filter",
			filter: filter.DoctorQueryParam{Name: "John"},
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT count\(\*\) FROM doctors WHERE name LIKE \$1`).
					WithArgs("%John%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
			},
			want: 3,
		},
		{
			name:   "specialty filter",
			filter: filter.DoctorQueryParam{SpecialtyID: specialtyID},
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT count\(\*\) FROM doctors WHERE specialty_id = \$1`).
					WithArgs(specialtyID.String()).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(7))
			},
			want: 7,
		},
		{
			name:   "zero result",
			filter: filter.DoctorQueryParam{Name: "NonExistent"},
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT count\(\*\) FROM doctors WHERE name LIKE \$1`).
					WithArgs("%NonExistent%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			want: 0,
		},
		{
			name:   "database error",
			filter: filter.DoctorQueryParam{},
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT count\(\*\) FROM doctors`).
					WillReturnError(errors.New("connection timeout"))
			},
			wantErr: "connection timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupTestDB(t)
			defer db.Close()

			repo := NewDoctorRepository(db)
			tt.mockSetup(mock)

			got, err := repo.Count(ctx, tt.filter)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Equal(t, 0, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDoctorRepository_GetByID(t *testing.T) {
	ctx := context.Background()
	doctor := newTestDoctor("Dr. John Smith")

	// Base query for selecting doctor fields
	baseSelectQuery := `SELECT id, name, specialty_id, phone_number, avatar_url, description, created_at, updated_at FROM doctors`

	tests := []struct {
		name      string
		id        uuid.UUID
		mockSetup func(sqlmock.Sqlmock)
		want      *medical.Doctor
		wantErr   string
	}{
		{
			name: "found",
			id:   doctor.ID,
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(
					baseSelectQuery + ` WHERE id = \$1`,
				).WithArgs(doctor.ID).WillReturnRows(mockDoctorRows(doctor))
			},
			want: &doctor,
		},
		{
			name: "not found",
			id:   uuid.New(),
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(
					baseSelectQuery + ` WHERE id = \$1`,
				).WithArgs(sqlmock.AnyArg()).WillReturnError(sql.ErrNoRows)
			},
			wantErr: "doctor not found",
		},
		{
			name: "database error",
			id:   doctor.ID,
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(
					baseSelectQuery + ` WHERE id = \$1`,
				).WithArgs(doctor.ID).WillReturnError(errors.New("connection lost"))
			},
			wantErr: "connection lost",
		},
		{
			name: "nil UUID",
			id:   uuid.Nil,
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(
					baseSelectQuery + ` WHERE id = \$1`,
				).WithArgs(uuid.Nil).WillReturnError(sql.ErrNoRows)
			},
			wantErr: "doctor not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupTestDB(t)
			defer db.Close()

			repo := NewDoctorRepository(db)
			tt.mockSetup(mock)

			got, err := repo.GetByID(ctx, tt.id)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				assertDoctorEqual(t, *tt.want, *got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
