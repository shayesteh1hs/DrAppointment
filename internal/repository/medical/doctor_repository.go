package medical

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"

	domain "drgo/internal/domain/medical"
	filter "drgo/internal/filter/medical"
	"drgo/internal/pagination"
)

type DoctorRepository interface {
	GetAllPaginated(ctx context.Context, filters filter.DoctorQueryParam, paginator *pagination.LimitOffsetPaginator[domain.Doctor]) ([]domain.Doctor, error)
	Count(ctx context.Context, filters filter.DoctorQueryParam) (int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Doctor, error)
}

type doctorRepository struct {
	db *sql.DB
}

func (r *doctorRepository) GetAllPaginated(ctx context.Context, filters filter.DoctorQueryParam, paginator *pagination.LimitOffsetPaginator[domain.Doctor]) ([]domain.Doctor, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.From("doctors")
	sb = filters.Apply(sb)
	sb = paginator.Paginate(sb)

	query, args := sb.Build()
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}(rows)

	doctors, err := r.scanDoctors(rows)
	if err != nil {
		return nil, err
	}

	return doctors, nil
}

func (r *doctorRepository) Count(ctx context.Context, filters filter.DoctorQueryParam) (int, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("count(*)")
	sb.From("doctors")
	filters.Apply(sb)

	query, args := sb.Build()
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count doctors: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}(rows)

	var totalCount int
	err = rows.Scan(&totalCount)
	if err != nil {
		return 0, fmt.Errorf("failed to scan total count: %w", err)
	}
	return totalCount, nil
}

func (r *doctorRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Doctor, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("id", "name", "specialty_id", "phone_number", "avatar_url", "description", "created_at", "updated_at")
	sb.From("doctors")
	sb.Where(sb.Equal("id", id))

	query, args := sb.Build()
	row := r.db.QueryRowContext(ctx, query, args...)

	var doc domain.Doctor
	err := row.Scan(
		&doc.ID,
		&doc.Name,
		&doc.SpecialtyID,
		&doc.PhoneNumber,
		&doc.AvatarURL,
		&doc.Description,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("doctor not found: %s", id)
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
			&doc.SpecialtyID,
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

func NewDoctorRepository(db *sql.DB) DoctorRepository {
	return &doctorRepository{db: db}
}
