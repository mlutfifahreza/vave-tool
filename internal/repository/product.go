package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/vave-tool/internal/domain"
	"github.com/vave-tool/internal/observability"
)

type productRepository struct {
	db      *sql.DB
	metrics *observability.Metrics
}

func NewProductRepository(db *sql.DB) domain.ProductRepository {
	return &productRepository{
		db:      db,
		metrics: observability.GetMetrics(),
	}
}

func (r *productRepository) List(ctx context.Context, params domain.PaginationParams, filters domain.ProductFilterParams) ([]*domain.Product, error) {
	ctx, span := observability.StartSpan(ctx, "Repository.ListProducts")
	defer span.End()

	offset := (params.Page - 1) * params.Size

	query := `
		SELECT p.id, p.name, p.description, p.price, p.stock_quantity, 
		       p.category_id, c.name as category_name,
		       p.subcategory_id, s.name as subcategory_name,
		       p.sku, p.is_active, p.updated_by, p.created_at, p.updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN subcategories s ON p.subcategory_id = s.id
		WHERE p.is_active = true`

	args := []interface{}{}
	argIndex := 1

	if filters.CategoryID != nil && *filters.CategoryID != "" {
		query += ` AND p.category_id = $` + strconv.Itoa(argIndex)
		args = append(args, *filters.CategoryID)
		argIndex++
	}

	if filters.SubcategoryID != nil && *filters.SubcategoryID != "" {
		query += ` AND p.subcategory_id = $` + strconv.Itoa(argIndex)
		args = append(args, *filters.SubcategoryID)
		argIndex++
	}

	if filters.MinPrice != nil {
		query += ` AND p.price >= $` + strconv.Itoa(argIndex)
		args = append(args, *filters.MinPrice)
		argIndex++
	}

	if filters.MaxPrice != nil {
		query += ` AND p.price <= $` + strconv.Itoa(argIndex)
		args = append(args, *filters.MaxPrice)
		argIndex++
	}

	query += `
		ORDER BY p.created_at DESC
		LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)

	args = append(args, params.Size, offset)

	start := time.Now()
	rows, err := r.db.QueryContext(ctx, query, args...)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "list_products", time.Since(start), err)
	}
	if err != nil {
		observability.RecordError(span, err, "Failed to query products")
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.StockQuantity,
			&p.CategoryID,
			&p.CategoryName,
			&p.SubcategoryID,
			&p.SubcategoryName,
			&p.SKU,
			&p.IsActive,
			&p.UpdatedBy,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) Count(ctx context.Context, filters domain.ProductFilterParams) (int64, error) {
	ctx, span := observability.StartSpan(ctx, "Repository.CountProducts")
	defer span.End()

	query := `SELECT COUNT(*) FROM products p WHERE p.is_active = true`

	args := []interface{}{}
	argIndex := 1

	if filters.CategoryID != nil && *filters.CategoryID != "" {
		query += ` AND p.category_id = $` + strconv.Itoa(argIndex)
		args = append(args, *filters.CategoryID)
		argIndex++
	}

	if filters.SubcategoryID != nil && *filters.SubcategoryID != "" {
		query += ` AND p.subcategory_id = $` + strconv.Itoa(argIndex)
		args = append(args, *filters.SubcategoryID)
		argIndex++
	}

	if filters.MinPrice != nil {
		query += ` AND p.price >= $` + strconv.Itoa(argIndex)
		args = append(args, *filters.MinPrice)
		argIndex++
	}

	if filters.MaxPrice != nil {
		query += ` AND p.price <= $` + strconv.Itoa(argIndex)
		args = append(args, *filters.MaxPrice)
		argIndex++
	}

	start := time.Now()
	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "count_products", time.Since(start), err)
	}
	if err != nil {
		observability.RecordError(span, err, "Failed to count products")
		return 0, err
	}

	return count, nil
}

func (r *productRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	ctx, span := observability.StartSpan(ctx, "Repository.GetProductByID")
	defer span.End()

	start := time.Now()
	query := `
		SELECT p.id, p.name, p.description, p.price, p.stock_quantity, 
		       p.category_id, c.name as category_name,
		       p.subcategory_id, s.name as subcategory_name,
		       p.sku, p.is_active, p.updated_by, p.created_at, p.updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN subcategories s ON p.subcategory_id = s.id
		WHERE p.id = $1
	`

	var p domain.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.StockQuantity,
		&p.CategoryID,
		&p.CategoryName,
		&p.SubcategoryID,
		&p.SubcategoryName,
		&p.SKU,
		&p.IsActive,
		&p.UpdatedBy,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "get_product", time.Since(start), err)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			observability.RecordError(span, domain.ErrNotFound, "Product not found")
			return nil, domain.ErrNotFound
		}
		observability.RecordError(span, err, "Failed to query product")
		return nil, err
	}

	return &p, nil
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	ctx, span := observability.StartSpan(ctx, "Repository.CreateProduct")
	defer span.End()

	start := time.Now()
	query := `
		INSERT INTO products (name, description, price, stock_quantity, category_id, subcategory_id, sku, is_active, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Description,
		product.Price,
		product.StockQuantity,
		product.CategoryID,
		product.SubcategoryID,
		product.SKU,
		product.IsActive,
		product.UpdatedBy,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "create_product", time.Since(start), err)
	}

	if err != nil {
		observability.RecordError(span, err, "Failed to insert product")
	}

	return err
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	ctx, span := observability.StartSpan(ctx, "Repository.UpdateProduct")
	defer span.End()

	start := time.Now()
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, stock_quantity = $4, 
		    category_id = $5, subcategory_id = $6, sku = $7, is_active = $8, updated_by = $9, updated_at = CURRENT_TIMESTAMP
		WHERE id = $10
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Description,
		product.Price,
		product.StockQuantity,
		product.CategoryID,
		product.SubcategoryID,
		product.SKU,
		product.IsActive,
		product.UpdatedBy,
		product.ID,
	).Scan(&product.UpdatedAt)

	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "update_product", time.Since(start), err)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			observability.RecordError(span, domain.ErrNotFound, "Product not found")
			return domain.ErrNotFound
		}
		observability.RecordError(span, err, "Failed to update product")
		return err
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	ctx, span := observability.StartSpan(ctx, "Repository.DeleteProduct")
	defer span.End()

	start := time.Now()
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "delete_product", time.Since(start), err)
	}
	if err != nil {
		observability.RecordError(span, err, "Failed to delete product")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		observability.RecordError(span, err, "Failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		observability.RecordError(span, domain.ErrNotFound, "Product not found")
		return domain.ErrNotFound
	}

	return nil
}
