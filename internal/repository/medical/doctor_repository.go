package medical

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"

	domain "drgo/internal/domain/medical"
	filter "drgo/internal/filter/medical"
	"drgo/internal/pagination"
	querybuilder "drgo/internal/query_builder/medical"
)

type DoctorRepository interface {
	GetAllOffset(ctx context.Context, filters filter.DoctorQueryParam, paginator pagination.LimitOffsetPaginator[domain.Doctor]) (*pagination.Result[domain.Doctor], error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Doctor, error)
	Create(ctx context.Context, doctor *domain.Doctor) error
}

type doctorRepository struct {
	db *sql.DB
}

func (r *doctorRepository) GetAllOffset(ctx context.Context, filters filter.DoctorQueryParam, paginator pagination.LimitOffsetPaginator[domain.Doctor]) (*pagination.Result[domain.Doctor], error) {
	// Build query for count
	countQb := querybuilder.NewDoctorQueryBuilder()
	countQb.WithFilters(filters)
	countSb := countQb.CountBuilder()
	countQuery, countArgs := countSb.Build()

	var totalCount int
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count doctors: %w", err)
	}

	// Build query with filters and pagination
	qb := querybuilder.NewDoctorQueryBuilder().
		WithFilters(filters).
		WithOrderBy("doctors.id")

	sb := qb.Build()

	// Apply pagination
	paginatedQuery, paginatedArgs := paginator.Paginate(sb)

	// Handle cursor pagination result processing
	rows, err := r.db.QueryContext(ctx, paginatedQuery, paginatedArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to query doctors: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	doctors, err := r.scanDoctors(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan doctors: %w", err)
	}

	// Create pagination result using the paginator
	meta := paginator.GetMeta(totalCount)
	return &pagination.Result[domain.Doctor]{
		Items: doctors,
		Meta:  meta,
	}, nil
}

func (r *doctorRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Doctor, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("id", "name", "specialty", "phone_number", "avatar_url", "description", "created_at", "updated_at")
	sb.From("doctors")
	sb.Where(sb.Equal("id", id))

	query, args := sb.Build()
	row := r.db.QueryRowContext(ctx, query, args...)

	var doc domain.Doctor
	err := row.Scan(
		&doc.ID,
		&doc.Name,
		&doc.Specialty,
		&doc.PhoneNumber,
		&doc.AvatarURL,
		&doc.Description,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("doctor not found")
		}
		return nil, fmt.Errorf("failed to scan doctor: %w", err)
	}

	return &doc, nil
}

func (r *doctorRepository) scanDoctors(rows *sql.Rows) ([]domain.Doctor, error) {
	var doctors []domain.Doctor
	for rows.Next() {
		var doc domain.Doctor
		err := rows.Scan(
			&doc.ID,
			&doc.Name,
			&doc.Specialty,
			&doc.PhoneNumber,
			&doc.AvatarURL,
			&doc.Description,
			&doc.CreatedAt,
			&doc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		doctors = append(doctors, doc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return doctors, nil
}

func (r *doctorRepository) Create(ctx context.Context, doctor *domain.Doctor) error {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("doctors")
	ib.Cols(
		"id", "name", "specialty", "phone_number",
		"avatar_url", "description", "created_at", "updated_at",
	)
	ib.Values(
		doctor.ID,
		doctor.Name,
		doctor.Specialty,
		doctor.PhoneNumber,
		doctor.AvatarURL,
		doctor.Description,
		doctor.CreatedAt,
		doctor.UpdatedAt,
	)

	sql, args := ib.Build()

	_, err := r.db.ExecContext(ctx, sql, args...)
	return err
}
func NewDoctorRepository(db *sql.DB) DoctorRepository {
	return &doctorRepository{db: db}
}
