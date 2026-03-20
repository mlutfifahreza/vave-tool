package domain

import (
	"context"
	"time"
)

type Product struct {
	ID              string    `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Description     *string   `json:"description,omitempty" db:"description"`
	Price           float64   `json:"price" db:"price"`
	StockQuantity   int       `json:"stock_quantity" db:"stock_quantity"`
	CategoryID      *string   `json:"category_id,omitempty" db:"category_id"`
	CategoryName    *string   `json:"category_name,omitempty" db:"category_name"`
	SubcategoryID   *string   `json:"subcategory_id,omitempty" db:"subcategory_id"`
	SubcategoryName *string   `json:"subcategory_name,omitempty" db:"subcategory_name"`
	SKU             *string   `json:"sku,omitempty" db:"sku"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	UpdatedBy       *string   `json:"updated_by,omitempty" db:"updated_by"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type Category struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Subcategory struct {
	ID          string    `json:"id" db:"id"`
	CategoryID  string    `json:"category_id" db:"category_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type PaginationParams struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type ProductFilterParams struct {
	CategoryID    *string  `json:"category_id,omitempty"`
	SubcategoryID *string  `json:"subcategory_id,omitempty"`
	MinPrice      *float64 `json:"min_price,omitempty"`
	MaxPrice      *float64 `json:"max_price,omitempty"`
}

type PaginatedProducts struct {
	Products []*Product `json:"products"`
}

type ProductRepository interface {
	List(ctx context.Context, params PaginationParams, filters ProductFilterParams) ([]*Product, error)
	Count(ctx context.Context, filters ProductFilterParams) (int64, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
}

type CategoryRepository interface {
	List(ctx context.Context, params PaginationParams) ([]*Category, error)
	Count(ctx context.Context) (int64, error)
	GetByID(ctx context.Context, id string) (*Category, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	Create(ctx context.Context, category *Category) error
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id string) error
}

type SubcategoryRepository interface {
	List(ctx context.Context, params PaginationParams) ([]*Subcategory, error)
	Count(ctx context.Context) (int64, error)
	GetByID(ctx context.Context, id string) (*Subcategory, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	GetByCategoryID(ctx context.Context, categoryID string, params PaginationParams) ([]*Subcategory, error)
	CountByCategoryID(ctx context.Context, categoryID string) (int64, error)
	Create(ctx context.Context, subcategory *Subcategory) error
	Update(ctx context.Context, subcategory *Subcategory) error
	Delete(ctx context.Context, id string) error
}

type ProductService interface {
	ListProducts(ctx context.Context, params PaginationParams, filters ProductFilterParams) (*PaginatedProducts, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	CreateProduct(ctx context.Context, product *Product) error
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, id string) error
}

type CategoryService interface {
	ListCategories(ctx context.Context, params PaginationParams) ([]*Category, error)
	GetCategory(ctx context.Context, id string) (*Category, error)
	CreateCategory(ctx context.Context, category *Category) error
	UpdateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, id string) error
}

type SubcategoryService interface {
	ListSubcategories(ctx context.Context, params PaginationParams) ([]*Subcategory, error)
	GetSubcategory(ctx context.Context, id string) (*Subcategory, error)
	GetSubcategoriesByCategory(ctx context.Context, categoryID string, params PaginationParams) ([]*Subcategory, error)
	CreateSubcategory(ctx context.Context, subcategory *Subcategory) error
	UpdateSubcategory(ctx context.Context, subcategory *Subcategory) error
	DeleteSubcategory(ctx context.Context, id string) error
}
