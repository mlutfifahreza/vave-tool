package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/vave-tool/internal/domain"
	"github.com/vave-tool/internal/observability"
)

type subcategoryRepository struct {
	db      *sql.DB
	metrics *observability.Metrics
}

func NewSubcategoryRepository(db *sql.DB) domain.SubcategoryRepository {
	return &subcategoryRepository{
		db:      db,
		metrics: observability.GetMetrics(),
	}
}

func (r *subcategoryRepository) List(ctx context.Context, params domain.PaginationParams) ([]*domain.Subcategory, error) {
	ctx, span := observability.StartSpan(ctx, "Repository.ListSubcategories")
	defer span.End()

	offset := (params.Page - 1) * params.Size

	start := time.Now()
	query := `
		SELECT id, category_id, name, description, is_active, created_at, updated_at
		FROM subcategories
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, params.Size, offset)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "list_subcategories", time.Since(start), err)
	}
	if err != nil {
		observability.RecordError(span, err, "Failed to query subcategories")
		return nil, err
	}
	defer rows.Close()

	var subcategories []*domain.Subcategory
	for rows.Next() {
		var s domain.Subcategory
		err := rows.Scan(
			&s.ID,
			&s.CategoryID,
			&s.Name,
			&s.Description,
			&s.IsActive,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subcategories = append(subcategories, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subcategories, nil
}

func (r *subcategoryRepository) Count(ctx context.Context) (int64, error) {
	ctx, span := observability.StartSpan(ctx, "Repository.CountSubcategories")
	defer span.End()

	start := time.Now()
	query := `SELECT COUNT(*) FROM subcategories WHERE is_active = true`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "count_subcategories", time.Since(start), err)
	}
	if err != nil {
		observability.RecordError(span, err, "Failed to count subcategories")
		return 0, err
	}

	return count, nil
}

func (r *subcategoryRepository) GetByID(ctx context.Context, id string) (*domain.Subcategory, error) {
	ctx, span := observability.StartSpan(ctx, "Repository.GetSubcategoryByID")
	defer span.End()

	start := time.Now()
	query := `
		SELECT id, category_id, name, description, is_active, created_at, updated_at
		FROM subcategories
		WHERE id = $1
	`

	var s domain.Subcategory
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&s.CategoryID,
		&s.Name,
		&s.Description,
		&s.IsActive,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "get_subcategory", time.Since(start), err)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			observability.RecordError(span, domain.ErrNotFound, "Subcategory not found")
			return nil, domain.ErrNotFound
		}
		observability.RecordError(span, err, "Failed to query subcategory")
		return nil, err
	}

	return &s, nil
}

func (r *subcategoryRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	ctx, span := observability.StartSpan(ctx, "Repository.SubcategoryExistsByID")
	defer span.End()

	start := time.Now()
	query := `SELECT EXISTS(SELECT 1 FROM subcategories WHERE id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "exists_subcategory", time.Since(start), err)
	}
	if err != nil {
		observability.RecordError(span, err, "Failed to check subcategory existence")
		return false, err
	}

	return exists, nil
}

func (r *subcategoryRepository) GetByCategoryID(ctx context.Context, categoryID string, params domain.PaginationParams) ([]*domain.Subcategory, error) {
	ctx, span := observability.StartSpan(ctx, "Repository.GetSubcategoriesByCategoryID")
	defer span.End()

	offset := (params.Page - 1) * params.Size

	start := time.Now()
	query := `
		SELECT id, category_id, name, description, is_active, created_at, updated_at
		FROM subcategories
		WHERE category_id = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, categoryID, params.Size, offset)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "list_subcategories_by_category", time.Since(start), err)
	}
	if err != nil {
		observability.RecordError(span, err, "Failed to query subcategories by category")
		return nil, err
	}
	defer rows.Close()

	var subcategories []*domain.Subcategory
	for rows.Next() {
		var s domain.Subcategory
		err := rows.Scan(
			&s.ID,
			&s.CategoryID,
			&s.Name,
			&s.Description,
			&s.IsActive,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subcategories = append(subcategories, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subcategories, nil
}

func (r *subcategoryRepository) CountByCategoryID(ctx context.Context, categoryID string) (int64, error) {
	ctx, span := observability.StartSpan(ctx, "Repository.CountSubcategoriesByCategoryID")
	defer span.End()

	start := time.Now()
	query := `SELECT COUNT(*) FROM subcategories WHERE category_id = $1 AND is_active = true`

	var count int64
	err := r.db.QueryRowContext(ctx, query, categoryID).Scan(&count)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "count_subcategories_by_category", time.Since(start), err)
	}
	if err != nil {
		observability.RecordError(span, err, "Failed to count subcategories by category")
		return 0, err
	}

	return count, nil
}

func (r *subcategoryRepository) Create(ctx context.Context, subcategory *domain.Subcategory) error {
	ctx, span := observability.StartSpan(ctx, "Repository.CreateSubcategory")
	defer span.End()

	start := time.Now()
	query := `
		INSERT INTO subcategories (id, category_id, name, description, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		subcategory.ID,
		subcategory.CategoryID,
		subcategory.Name,
		subcategory.Description,
		subcategory.IsActive,
	).Scan(&subcategory.CreatedAt, &subcategory.UpdatedAt)

	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "create_subcategory", time.Since(start), err)
	}

	if err != nil {
		observability.RecordError(span, err, "Failed to insert subcategory")
	}

	return err
}

func (r *subcategoryRepository) Update(ctx context.Context, subcategory *domain.Subcategory) error {
	ctx, span := observability.StartSpan(ctx, "Repository.UpdateSubcategory")
	defer span.End()

	start := time.Now()
	query := `
		UPDATE subcategories
		SET category_id = $1, name = $2, description = $3, is_active = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		subcategory.CategoryID,
		subcategory.Name,
		subcategory.Description,
		subcategory.IsActive,
		subcategory.ID,
	).Scan(&subcategory.UpdatedAt)

	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "update_subcategory", time.Since(start), err)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			observability.RecordError(span, domain.ErrNotFound, "Subcategory not found")
			return domain.ErrNotFound
		}
		observability.RecordError(span, err, "Failed to update subcategory")
		return err
	}

	return nil
}

func (r *subcategoryRepository) Delete(ctx context.Context, id string) error {
	ctx, span := observability.StartSpan(ctx, "Repository.DeleteSubcategory")
	defer span.End()

	start := time.Now()
	query := `DELETE FROM subcategories WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "delete_subcategory", time.Since(start), err)
	}
	if err != nil {
		observability.RecordError(span, err, "Failed to delete subcategory")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		observability.RecordError(span, err, "Failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		observability.RecordError(span, domain.ErrNotFound, "Subcategory not found")
		return domain.ErrNotFound
	}

	return nil
}
