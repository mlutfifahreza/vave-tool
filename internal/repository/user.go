package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/vave-tool/internal/domain"
	"github.com/vave-tool/internal/observability"
)

type userRepository struct {
	db      *sql.DB
	metrics *observability.Metrics
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{
		db:      db,
		metrics: observability.GetMetrics(),
	}
}

func (r *userRepository) GetByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	ctx, span := observability.StartSpan(ctx, "Repository.GetUserByGoogleID")
	defer span.End()

	query := `
		SELECT id, google_id, email, name, picture, is_active, created_at, updated_at
		FROM users WHERE google_id = $1
	`

	start := time.Now()
	row := r.db.QueryRowContext(ctx, query, googleID)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "get_user_by_google_id", time.Since(start), nil)
	}
	return r.scanUser(row)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, span := observability.StartSpan(ctx, "Repository.GetUserByEmail")
	defer span.End()

	query := `
		SELECT id, google_id, email, name, picture, is_active, created_at, updated_at
		FROM users WHERE email = $1
	`

	start := time.Now()
	row := r.db.QueryRowContext(ctx, query, email)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "get_user_by_email", time.Since(start), nil)
	}
	return r.scanUser(row)
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	ctx, span := observability.StartSpan(ctx, "Repository.CreateUser")
	defer span.End()

	query := `
		INSERT INTO users (google_id, email, name, picture, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, google_id, email, name, picture, is_active, created_at, updated_at
	`

	now := time.Now()
	start := time.Now()
	row := r.db.QueryRowContext(ctx, query, user.GoogleID, user.Email, user.Name, user.Picture, true, now, now)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "create_user", time.Since(start), nil)
	}
	return r.scanUser(row)
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	ctx, span := observability.StartSpan(ctx, "Repository.UpdateUser")
	defer span.End()

	query := `
		UPDATE users SET name = $1, picture = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, google_id, email, name, picture, is_active, created_at, updated_at
	`

	now := time.Now()
	start := time.Now()
	row := r.db.QueryRowContext(ctx, query, user.Name, user.Picture, now, user.ID)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "update_user", time.Since(start), nil)
	}
	return r.scanUser(row)
}

func (r *userRepository) scanUser(row *sql.Row) (*domain.User, error) {
	var user domain.User
	var picture sql.NullString

	err := row.Scan(
		&user.ID,
		&user.GoogleID,
		&user.Email,
		&user.Name,
		&picture,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if picture.Valid {
		user.Picture = picture.String
	}

	return &user, nil
}
