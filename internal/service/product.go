package service

import (
	"context"

	"github.com/vave-tool/backend/internal/domain"
)

type productService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) domain.ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	return s.repo.List(ctx)
}

func (s *productService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *productService) CreateProduct(ctx context.Context, product *domain.Product) error {
	return s.repo.Create(ctx, product)
}

func (s *productService) UpdateProduct(ctx context.Context, product *domain.Product) error {
	return s.repo.Update(ctx, product)
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
