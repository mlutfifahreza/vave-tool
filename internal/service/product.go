package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vave-tool/backend/internal/domain"
)

type productService struct {
	repo        domain.ProductRepository
	redisClient *redis.Client
}

func NewProductService(repo domain.ProductRepository, redisClient *redis.Client) domain.ProductService {
	return &productService{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (s *productService) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	return s.repo.List(ctx)
}

func (s *productService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	cacheKey := fmt.Sprintf("product:%s", id)

	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var product domain.Product
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			return &product, nil
		}
	}

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if productJSON, err := json.Marshal(product); err == nil {
		s.redisClient.Set(ctx, cacheKey, productJSON, 15*time.Minute)
	}

	return product, nil
}

func (s *productService) CreateProduct(ctx context.Context, product *domain.Product) error {
	return s.repo.Create(ctx, product)
}

func (s *productService) UpdateProduct(ctx context.Context, product *domain.Product) error {
	if err := s.repo.Update(ctx, product); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("product:%s", product.ID)
	s.redisClient.Del(ctx, cacheKey)

	return nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("product:%s", id)
	s.redisClient.Del(ctx, cacheKey)

	return nil
}
