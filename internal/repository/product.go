package repository

import (
	"context"
	"database/sql"
	"errors"
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

func (r *productRepository) List(ctx context.Context) ([]*domain.Product, error) {
	start := time.Now()
	query := `
		SELECT id, name, description, price, stock_quantity, category, sku, is_active, created_at, updated_at
		FROM products
		WHERE is_active = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "list_products", time.Since(start), err)
	}
	if err != nil {
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
			&p.Category,
			&p.SKU,
			&p.IsActive,
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

func (r *productRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	start := time.Now()
	query := `
		SELECT id, name, description, price, stock_quantity, category, sku, is_active, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var p domain.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.StockQuantity,
		&p.Category,
		&p.SKU,
		&p.IsActive,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "get_product", time.Since(start), err)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	start := time.Now()
	query := `
		INSERT INTO products (name, description, price, stock_quantity, category, sku, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Description,
		product.Price,
		product.StockQuantity,
		product.Category,
		product.SKU,
		product.IsActive,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "create_product", time.Since(start), err)
	}

	return err
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	start := time.Now()
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, stock_quantity = $4, 
		    category = $5, sku = $6, is_active = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Description,
		product.Price,
		product.StockQuantity,
		product.Category,
		product.SKU,
		product.IsActive,
		product.ID,
	).Scan(&product.UpdatedAt)

	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "update_product", time.Since(start), err)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	start := time.Now()
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if r.metrics != nil {
		r.metrics.RecordDBCall(ctx, "delete_product", time.Since(start), err)
	}
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
