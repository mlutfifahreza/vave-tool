package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"

	"github.com/vave-tool/internal/domain"
	"github.com/vave-tool/internal/observability"
)

type categoryService struct {
	repo        domain.CategoryRepository
	redisClient *redis.Client
	logger      *observability.Logger
	metrics     *observability.Metrics
}

func NewCategoryService(repo domain.CategoryRepository, redisClient *redis.Client, logger *observability.Logger) domain.CategoryService {
	return &categoryService{
		repo:        repo,
		redisClient: redisClient,
		logger:      logger,
		metrics:     observability.GetMetrics(),
	}
}

func (s *categoryService) ListCategories(ctx context.Context, params domain.PaginationParams) ([]*domain.Category, error) {
	ctx, span := observability.StartSpan(ctx, "CategoryService.ListCategories")
	defer span.End()

	cacheKey := fmt.Sprintf("categories:list:page:%d:size:%d", params.Page, params.Size)

	s.logger.Debug(ctx, "Checking cache for category list", zap.Int("page", params.Page), zap.Int("size", params.Size))
	redisStart := time.Now()
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if s.metrics != nil {
		s.metrics.RecordRedisOp(ctx, "get", time.Since(redisStart), err)
	}
	if err == nil {
		var categories []*domain.Category
		if err := json.Unmarshal([]byte(cached), &categories); err == nil {
			span.SetAttributes(attribute.Bool("cache_hit", true))
			if s.metrics != nil {
				s.metrics.RecordCacheAccess(ctx, "list_categories", true)
			}
			s.logger.Debug(ctx, "Category list found in cache", zap.Int("count", len(categories)))
			return categories, nil
		}
	}

	span.SetAttributes(attribute.Bool("cache_hit", false))
	if s.metrics != nil {
		s.metrics.RecordCacheAccess(ctx, "list_categories", false)
	}
	s.logger.Debug(ctx, "Category list not in cache, fetching from repository", zap.Int("page", params.Page), zap.Int("size", params.Size))

	categories, err := s.repo.List(ctx, params)
	if err != nil {
		observability.RecordError(span, err, "Failed to list categories")
		return nil, err
	}

	span.SetAttributes(
		attribute.Int("category_count", len(categories)),
		attribute.Int("page", params.Page),
		attribute.Int("size", params.Size),
	)
	s.logger.Debug(ctx, "Categories fetched from repository",
		zap.Int("count", len(categories)),
		zap.Int("page", params.Page),
	)

	go func() {
		if categoriesJSON, err := json.Marshal(categories); err == nil {
			redisStart := time.Now()
			setErr := s.redisClient.Set(context.Background(), cacheKey, categoriesJSON, 30*time.Second).Err()
			if s.metrics != nil {
				s.metrics.RecordRedisOp(context.Background(), "set", time.Since(redisStart), setErr)
			}
			s.logger.Debug(context.Background(), "Category list cached", zap.Int("count", len(categories)))
		}
	}()

	return categories, nil
}

func (s *categoryService) GetCategory(ctx context.Context, id string) (*domain.Category, error) {
	ctx, span := observability.StartSpan(ctx, "CategoryService.GetCategory",
		attribute.String("category_id", id))
	defer span.End()

	cacheKey := fmt.Sprintf("categories:get:%s", id)

	s.logger.Debug(ctx, "Checking cache for category", zap.String("id", id))
	redisStart := time.Now()
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if s.metrics != nil {
		s.metrics.RecordRedisOp(ctx, "get", time.Since(redisStart), err)
	}
	if err == nil {
		var category domain.Category
		if err := json.Unmarshal([]byte(cached), &category); err == nil {
			span.SetAttributes(attribute.Bool("cache_hit", true))
			if s.metrics != nil {
				s.metrics.RecordCacheAccess(ctx, "get_category", true)
			}
			s.logger.Debug(ctx, "Category found in cache", zap.String("id", id))
			return &category, nil
		}
	}

	span.SetAttributes(attribute.Bool("cache_hit", false))
	if s.metrics != nil {
		s.metrics.RecordCacheAccess(ctx, "get_category", false)
	}
	s.logger.Debug(ctx, "Category not in cache, fetching from repository", zap.String("id", id))

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		observability.RecordError(span, err, "Failed to get category")
		return nil, err
	}

	span.SetAttributes(attribute.String("category_name", category.Name))
	s.logger.Debug(ctx, "Category fetched from repository", zap.String("id", id), zap.String("name", category.Name))

	go func() {
		if categoryJSON, err := json.Marshal(category); err == nil {
			redisStart := time.Now()
			setErr := s.redisClient.Set(context.Background(), cacheKey, categoryJSON, 30*time.Second).Err()
			if s.metrics != nil {
				s.metrics.RecordRedisOp(context.Background(), "set", time.Since(redisStart), setErr)
			}
			s.logger.Debug(context.Background(), "Category cached", zap.String("id", id))
		}
	}()

	return category, nil
}

func (s *categoryService) CreateCategory(ctx context.Context, category *domain.Category) error {
	ctx, span := observability.StartSpan(ctx, "CategoryService.CreateCategory",
		attribute.String("category_name", category.Name))
	defer span.End()

	// Check if ID already exists
	exists, err := s.repo.ExistsByID(ctx, category.ID)
	if err != nil {
		observability.RecordError(span, err, "Failed to check category existence")
		return err
	}
	if exists {
		err := fmt.Errorf("category with id %s already exists", category.ID)
		observability.RecordError(span, err, "Category ID already exists")
		return err
	}

	s.logger.Debug(ctx, "Creating category", zap.String("name", category.Name), zap.String("id", category.ID))

	err = s.repo.Create(ctx, category)
	if err != nil {
		observability.RecordError(span, err, "Failed to create category")
		return err
	}

	span.SetAttributes(attribute.String("category_id", category.ID))
	s.logger.Info(ctx, "Category created", zap.String("id", category.ID), zap.String("name", category.Name))

	// Invalidate cache
	go func() {
		redisStart := time.Now()
		keys := []string{"categories:list:*"}
		for _, key := range keys {
			setErr := s.redisClient.Del(context.Background(), key).Err()
			if s.metrics != nil {
				s.metrics.RecordRedisOp(context.Background(), "del", time.Since(redisStart), setErr)
			}
		}
		s.logger.Debug(context.Background(), "Category list cache invalidated")
	}()

	return nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, category *domain.Category) error {
	ctx, span := observability.StartSpan(ctx, "CategoryService.UpdateCategory",
		attribute.String("category_id", category.ID))
	defer span.End()

	s.logger.Debug(ctx, "Updating category", zap.String("id", category.ID), zap.String("name", category.Name))

	err := s.repo.Update(ctx, category)
	if err != nil {
		observability.RecordError(span, err, "Failed to update category")
		return err
	}

	span.SetAttributes(attribute.String("category_name", category.Name))
	s.logger.Info(ctx, "Category updated", zap.String("id", category.ID), zap.String("name", category.Name))

	// Invalidate cache
	go func() {
		redisStart := time.Now()
		keys := []string{
			fmt.Sprintf("categories:get:%s", category.ID),
			"categories:list:*",
		}
		for _, key := range keys {
			setErr := s.redisClient.Del(context.Background(), key).Err()
			if s.metrics != nil {
				s.metrics.RecordRedisOp(context.Background(), "del", time.Since(redisStart), setErr)
			}
		}
		s.logger.Debug(context.Background(), "Category cache invalidated")
	}()

	return nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id string) error {
	ctx, span := observability.StartSpan(ctx, "CategoryService.DeleteCategory",
		attribute.String("category_id", id))
	defer span.End()

	s.logger.Debug(ctx, "Deleting category", zap.String("id", id))

	err := s.repo.Delete(ctx, id)
	if err != nil {
		observability.RecordError(span, err, "Failed to delete category")
		return err
	}

	s.logger.Info(ctx, "Category deleted", zap.String("id", id))

	// Invalidate cache
	go func() {
		redisStart := time.Now()
		keys := []string{
			fmt.Sprintf("categories:get:%s", id),
			"categories:list:*",
		}
		for _, key := range keys {
			setErr := s.redisClient.Del(context.Background(), key).Err()
			if s.metrics != nil {
				s.metrics.RecordRedisOp(context.Background(), "del", time.Since(redisStart), setErr)
			}
		}
		s.logger.Debug(context.Background(), "Category cache invalidated")
	}()

	return nil
}
