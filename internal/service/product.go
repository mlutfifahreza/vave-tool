package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"

	"github.com/vave-tool/internal/constants"
	"github.com/vave-tool/internal/domain"
	"github.com/vave-tool/internal/observability"
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

func (s *productService) ListProducts(ctx context.Context, params domain.PaginationParams, filters domain.ProductFilterParams) (*domain.PaginatedProducts, error) {
	ctx, span := observability.StartSpan(ctx, "ProductService.ListProducts")
	defer span.End()

	cacheKey := s.generateCacheKey(params, filters)

	s.logger.Debug(ctx, "Checking cache for product list",
		zap.Int("page", params.Page),
		zap.Int("size", params.Size),
		zap.String("cache_key", cacheKey))
	redisStart := time.Now()
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if s.metrics != nil {
		s.metrics.RecordRedisOp(ctx, "get", time.Since(redisStart), err)
	}
	if err == nil {
		var products []*domain.Product
		if err := json.Unmarshal([]byte(cached), &products); err == nil {
			span.SetAttributes(attribute.Bool("cache_hit", true))
			if s.metrics != nil {
				s.metrics.RecordCacheAccess(ctx, "list_products", true)
			}
			s.logger.Debug(ctx, "Product list found in cache", zap.Int("count", len(products)))
			return &domain.PaginatedProducts{
				Products: products,
			}, nil
		}
	}

	span.SetAttributes(attribute.Bool("cache_hit", false))
	if s.metrics != nil {
		s.metrics.RecordCacheAccess(ctx, "list_products", false)
	}
	s.logger.Debug(ctx, "Product list not in cache, fetching from repository",
		zap.Int("page", params.Page),
		zap.Int("size", params.Size))

	products, err := s.repo.List(ctx, params, filters)
	if err != nil {
		observability.RecordError(span, err, "Failed to list products")
		return nil, err
	}

	span.SetAttributes(
		attribute.Int("product_count", len(products)),
		attribute.Int("page", params.Page),
		attribute.Int("size", params.Size),
	)
	s.logger.Debug(ctx, "Products fetched from repository",
		zap.Int("count", len(products)),
		zap.Int("page", params.Page),
	)

	go func() {
		if productsJSON, err := json.Marshal(products); err == nil {
			redisStart := time.Now()
			setErr := s.redisClient.Set(context.Background(), cacheKey, productsJSON, 30*time.Second).Err()
			if s.metrics != nil {
				s.metrics.RecordRedisOp(context.Background(), "set", time.Since(redisStart), setErr)
			}
			s.logger.Debug(context.Background(), "Product list cached", zap.Int("count", len(products)))
		}
	}()

	return &domain.PaginatedProducts{
		Products: products,
	}, nil
}

func (s *productService) generateCacheKey(params domain.PaginationParams, filters domain.ProductFilterParams) string {
	key := fmt.Sprintf("products:list:page:%d:size:%d", params.Page, params.Size)

	if filters.CategoryID != nil && *filters.CategoryID != "" {
		key += fmt.Sprintf(":cat:%s", *filters.CategoryID)
	}

	if filters.SubcategoryID != nil && *filters.SubcategoryID != "" {
		key += fmt.Sprintf(":subcat:%s", *filters.SubcategoryID)
	}

	if filters.MinPrice != nil {
		key += fmt.Sprintf(":minp:%.2f", *filters.MinPrice)
	}

	if filters.MaxPrice != nil {
		key += fmt.Sprintf(":maxp:%.2f", *filters.MaxPrice)
	}

	return key
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

	go func() {
		if productJSON, err := json.Marshal(product); err == nil {
			redisStart := time.Now()
			setErr := s.redisClient.Set(context.Background(), cacheKey, productJSON, constants.ProductCacheTTL).Err()
			if s.metrics != nil {
				s.metrics.RecordRedisOp(context.Background(), "set", time.Since(redisStart), setErr)
			}
			s.logger.Debug(context.Background(), "Product cached", zap.String("product_id", id))
		}
	}()

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

	go func() {
		cacheKey := fmt.Sprintf(constants.ProductCacheKeyPrefix, product.ID)
		redisStart := time.Now()
		delErr := s.redisClient.Del(context.Background(), cacheKey).Err()
		if s.metrics != nil {
			s.metrics.RecordRedisOp(context.Background(), "delete", time.Since(redisStart), delErr)
		}
		s.logger.Debug(context.Background(), "Product cache invalidated", zap.String("product_id", product.ID))
	}()

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

	go func() {
		cacheKey := fmt.Sprintf(constants.ProductCacheKeyPrefix, id)
		redisStart := time.Now()
		delErr := s.redisClient.Del(context.Background(), cacheKey).Err()
		if s.metrics != nil {
			s.metrics.RecordRedisOp(context.Background(), "delete", time.Since(redisStart), delErr)
		}
		s.logger.Debug(context.Background(), "Product cache invalidated", zap.String("product_id", id))
	}()

	return nil
}
