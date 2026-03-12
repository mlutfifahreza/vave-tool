package domain

import (
	"context"
	"time"
)

type Product struct {
	ID            string    `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Description   *string   `json:"description,omitempty" db:"description"`
	Price         float64   `json:"price" db:"price"`
	StockQuantity int       `json:"stock_quantity" db:"stock_quantity"`
	Category      *string   `json:"category,omitempty" db:"category"`
	SKU           *string   `json:"sku,omitempty" db:"sku"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type PaginationParams struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type PaginatedProducts struct {
	Products []*Product `json:"products"`
}

type ProductRepository interface {
	List(ctx context.Context, params PaginationParams) ([]*Product, error)
	Count(ctx context.Context) (int64, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
}

type ProductService interface {
	ListProducts(ctx context.Context, params PaginationParams) (*PaginatedProducts, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	CreateProduct(ctx context.Context, product *Product) error
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, id string) error
}
