package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vave-tool/internal/constants"
	"github.com/vave-tool/internal/domain"
	"github.com/vave-tool/internal/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type productService struct {
	repo        domain.ProductRepository
	redisClient *redis.Client
	logger      *observability.Logger
	metrics     *observability.Metrics
}

func NewProductService(repo domain.ProductRepository, redisClient *redis.Client, logger *observability.Logger) domain.ProductService {
	return &productService{
		repo:        repo,
		redisClient: redisClient,
		logger:      logger,
		metrics:     observability.GetMetrics(),
	}
}

func (s *productService) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	ctx, span := observability.StartSpan(ctx, "ProductService.ListProducts")
	defer span.End()

	s.logger.Debug(ctx, "Fetching products from repository")

	products, err := s.repo.List(ctx)
	if err != nil {
		observability.RecordError(span, err, "Failed to list products")
		return nil, err
	}

	span.SetAttributes(attribute.Int("product_count", len(products)))
	s.logger.Debug(ctx, "Products fetched from repository", zap.Int("count", len(products)))

	return products, nil
}

func (s *productService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	ctx, span := observability.StartSpan(ctx, "ProductService.GetProduct",
		attribute.String("product_id", id),
	)
	defer span.End()

	cacheKey := fmt.Sprintf(constants.ProductCacheKeyPrefix, id)

	s.logger.Debug(ctx, "Checking cache for product", zap.String("product_id", id))
	redisStart := time.Now()
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if s.metrics != nil {
		s.metrics.RecordRedisOp(ctx, "get", time.Since(redisStart), err)
	}
	if err == nil {
		var product domain.Product
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			span.SetAttributes(attribute.Bool("cache_hit", true))
			if s.metrics != nil {
				s.metrics.RecordCacheAccess(ctx, "get_product", true)
			}
			s.logger.Debug(ctx, "Product found in cache", zap.String("product_id", id))
			return &product, nil
		}
	}

	span.SetAttributes(attribute.Bool("cache_hit", false))
	if s.metrics != nil {
		s.metrics.RecordCacheAccess(ctx, "get_product", false)
	}
	s.logger.Debug(ctx, "Product not in cache, fetching from repository", zap.String("product_id", id))

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		observability.RecordError(span, err, "Failed to get product")
		return nil, err
	}

	if productJSON, err := json.Marshal(product); err == nil {
		redisStart := time.Now()
		setErr := s.redisClient.Set(ctx, cacheKey, productJSON, constants.ProductCacheTTL).Err()
		if s.metrics != nil {
			s.metrics.RecordRedisOp(ctx, "set", time.Since(redisStart), setErr)
		}
		s.logger.Debug(ctx, "Product cached", zap.String("product_id", id))
	}

	return product, nil
}

func (s *productService) CreateProduct(ctx context.Context, product *domain.Product) error {
	ctx, span := observability.StartSpan(ctx, "ProductService.CreateProduct",
		attribute.String("product_name", product.Name),
	)
	defer span.End()

	s.logger.Debug(ctx, "Creating product in repository", zap.String("product_name", product.Name))

	err := s.repo.Create(ctx, product)
	if err != nil {
		observability.RecordError(span, err, "Failed to create product")
		return err
	}

	span.SetAttributes(attribute.String("product_id", product.ID))
	s.logger.Debug(ctx, "Product created in repository", zap.String("product_id", product.ID))

	return nil
}

func (s *productService) UpdateProduct(ctx context.Context, product *domain.Product) error {
	ctx, span := observability.StartSpan(ctx, "ProductService.UpdateProduct",
		attribute.String("product_id", product.ID),
	)
	defer span.End()

	s.logger.Debug(ctx, "Updating product in repository", zap.String("product_id", product.ID))

	if err := s.repo.Update(ctx, product); err != nil {
		observability.RecordError(span, err, "Failed to update product")
		return err
	}

	cacheKey := fmt.Sprintf(constants.ProductCacheKeyPrefix, product.ID)
	redisStart := time.Now()
	delErr := s.redisClient.Del(ctx, cacheKey).Err()
	if s.metrics != nil {
		s.metrics.RecordRedisOp(ctx, "delete", time.Since(redisStart), delErr)
	}
	s.logger.Debug(ctx, "Product cache invalidated", zap.String("product_id", product.ID))

	return nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	ctx, span := observability.StartSpan(ctx, "ProductService.DeleteProduct",
		attribute.String("product_id", id),
	)
	defer span.End()

	s.logger.Debug(ctx, "Deleting product from repository", zap.String("product_id", id))

	if err := s.repo.Delete(ctx, id); err != nil {
		observability.RecordError(span, err, "Failed to delete product")
		return err
	}

	cacheKey := fmt.Sprintf(constants.ProductCacheKeyPrefix, id)
	redisStart := time.Now()
	delErr := s.redisClient.Del(ctx, cacheKey).Err()
	if s.metrics != nil {
		s.metrics.RecordRedisOp(ctx, "delete", time.Since(redisStart), delErr)
	}
	s.logger.Debug(ctx, "Product cache invalidated", zap.String("product_id", id))

	return nil
}
