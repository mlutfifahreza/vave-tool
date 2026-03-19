package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/vave-tool/internal/domain"
	"github.com/vave-tool/internal/observability"
)

type categoryRepository struct {
	db      *sql.DB
	metrics *observability.Metrics
}

func NewCategoryRepository(db *sql.DB) domain.CategoryRepository {
	return &categoryRepository{
		db:      db,
		metrics: observability.GetMetrics(),
	}
}

func (r *categoryRepository) List(ctx context.Context, params domain.PaginationParams) ([]*domain.Category, error) {
	ctx, span := observability.StartSpan(ctx, "Repository.ListCategories")
	defer span.End()

	offset := (params.Page - 1) * params.Size

	start := time.Now()
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM categories
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, params.Size, offset)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "list_categories", time.Since(start), err)
	}
	if err != nil {
		observability.RecordError(span, err, "Failed to query categories")
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var c domain.Category
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
			&c.IsActive,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) Count(ctx context.Context) (int64, error) {
	ctx, span := observability.StartSpan(ctx, "Repository.CountCategories")
	defer span.End()

	start := time.Now()
	query := `SELECT COUNT(*) FROM categories WHERE is_active = true`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "count_categories", time.Since(start), err)
	}
	if err != nil {
		observability.RecordError(span, err, "Failed to count categories")
		return 0, err
	}

	return count, nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	ctx, span := observability.StartSpan(ctx, "Repository.GetCategoryByID")
	defer span.End()

	start := time.Now()
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	var c domain.Category
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID,
		&c.Name,
		&c.Description,
		&c.IsActive,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "get_category", time.Since(start), err)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			observability.RecordError(span, domain.ErrNotFound, "Category not found")
			return nil, domain.ErrNotFound
		}
		observability.RecordError(span, err, "Failed to query category")
		return nil, err
	}

	return &c, nil
}

func (r *categoryRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	ctx, span := observability.StartSpan(ctx, "Repository.CategoryExistsByID")
	defer span.End()

	start := time.Now()
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "exists_category", time.Since(start), err)
	}
	if err != nil {
		observability.RecordError(span, err, "Failed to check category existence")
		return false, err
	}

	return exists, nil
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	ctx, span := observability.StartSpan(ctx, "Repository.CreateCategory")
	defer span.End()

	start := time.Now()
	query := `
		INSERT INTO categories (id, name, description, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		category.ID,
		category.Name,
		category.Description,
		category.IsActive,
	).Scan(&category.CreatedAt, &category.UpdatedAt)

	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "create_category", time.Since(start), err)
	}

	if err != nil {
		observability.RecordError(span, err, "Failed to insert category")
	}

	return err
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	ctx, span := observability.StartSpan(ctx, "Repository.UpdateCategory")
	defer span.End()

	start := time.Now()
	query := `
		UPDATE categories
		SET name = $1, description = $2, is_active = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		category.Name,
		category.Description,
		category.IsActive,
		category.ID,
	).Scan(&category.UpdatedAt)

	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "update_category", time.Since(start), err)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			observability.RecordError(span, domain.ErrNotFound, "Category not found")
			return domain.ErrNotFound
		}
		observability.RecordError(span, err, "Failed to update category")
		return err
	}

	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	ctx, span := observability.StartSpan(ctx, "Repository.DeleteCategory")
	defer span.End()

	start := time.Now()
	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "delete_category", time.Since(start), err)
	}
	if err != nil {
		observability.RecordError(span, err, "Failed to delete category")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		observability.RecordError(span, err, "Failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		observability.RecordError(span, domain.ErrNotFound, "Category not found")
		return domain.ErrNotFound
	}

	return nil
}
