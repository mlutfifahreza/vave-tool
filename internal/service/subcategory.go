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

type subcategoryService struct {
	repo        domain.SubcategoryRepository
	redisClient *redis.Client
	logger      *observability.Logger
	metrics     *observability.Metrics
}

func NewSubcategoryService(repo domain.SubcategoryRepository, redisClient *redis.Client, logger *observability.Logger) domain.SubcategoryService {
	return &subcategoryService{
		repo:        repo,
		redisClient: redisClient,
		logger:      logger,
		metrics:     observability.GetMetrics(),
	}
}

func (s *subcategoryService) ListSubcategories(ctx context.Context, params domain.PaginationParams) ([]*domain.Subcategory, error) {
	ctx, span := observability.StartSpan(ctx, "SubcategoryService.ListSubcategories")
	defer span.End()

	cacheKey := fmt.Sprintf("subcategories:list:page:%d:size:%d", params.Page, params.Size)

	s.logger.Debug(ctx, "Checking cache for subcategory list", zap.Int("page", params.Page), zap.Int("size", params.Size))
	redisStart := time.Now()
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if s.metrics != nil {
		s.metrics.RecordRedisOp(ctx, "get", time.Since(redisStart), err)
	}
	if err == nil {
		var subcategories []*domain.Subcategory
		if err := json.Unmarshal([]byte(cached), &subcategories); err == nil {
			span.SetAttributes(attribute.Bool("cache_hit", true))
			if s.metrics != nil {
				s.metrics.RecordCacheAccess(ctx, "list_subcategories", true)
			}
			s.logger.Debug(ctx, "Subcategory list found in cache", zap.Int("count", len(subcategories)))
			return subcategories, nil
		}
	}

	span.SetAttributes(attribute.Bool("cache_hit", false))
	if s.metrics != nil {
		s.metrics.RecordCacheAccess(ctx, "list_subcategories", false)
	}
	s.logger.Debug(ctx, "Subcategory list not in cache, fetching from repository", zap.Int("page", params.Page), zap.Int("size", params.Size))

	subcategories, err := s.repo.List(ctx, params)
	if err != nil {
		observability.RecordError(span, err, "Failed to list subcategories")
		return nil, err
	}

	span.SetAttributes(
		attribute.Int("subcategory_count", len(subcategories)),
		attribute.Int("page", params.Page),
		attribute.Int("size", params.Size),
	)
	s.logger.Debug(ctx, "Subcategories fetched from repository",
		zap.Int("count", len(subcategories)),
		zap.Int("page", params.Page),
	)

	go func() {
		if subcategoriesJSON, err := json.Marshal(subcategories); err == nil {
			redisStart := time.Now()
			setErr := s.redisClient.Set(context.Background(), cacheKey, subcategoriesJSON, 30*time.Second).Err()
			if s.metrics != nil {
				s.metrics.RecordRedisOp(context.Background(), "set", time.Since(redisStart), setErr)
			}
			s.logger.Debug(context.Background(), "Subcategory list cached", zap.Int("count", len(subcategories)))
		}
	}()

	return subcategories, nil
}

func (s *subcategoryService) GetSubcategory(ctx context.Context, id string) (*domain.Subcategory, error) {
	ctx, span := observability.StartSpan(ctx, "SubcategoryService.GetSubcategory",
		attribute.String("subcategory_id", id))
	defer span.End()

	cacheKey := fmt.Sprintf("subcategories:get:%s", id)

	s.logger.Debug(ctx, "Checking cache for subcategory", zap.String("id", id))
	redisStart := time.Now()
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if s.metrics != nil {
		s.metrics.RecordRedisOp(ctx, "get", time.Since(redisStart), err)
	}
	if err == nil {
		var subcategory domain.Subcategory
		if err := json.Unmarshal([]byte(cached), &subcategory); err == nil {
			span.SetAttributes(attribute.Bool("cache_hit", true))
			if s.metrics != nil {
				s.metrics.RecordCacheAccess(ctx, "get_subcategory", true)
			}
			s.logger.Debug(ctx, "Subcategory found in cache", zap.String("id", id))
			return &subcategory, nil
		}
	}

	span.SetAttributes(attribute.Bool("cache_hit", false))
	if s.metrics != nil {
		s.metrics.RecordCacheAccess(ctx, "get_subcategory", false)
	}
	s.logger.Debug(ctx, "Subcategory not in cache, fetching from repository", zap.String("id", id))

	subcategory, err := s.repo.GetByID(ctx, id)
	if err != nil {
		observability.RecordError(span, err, "Failed to get subcategory")
		return nil, err
	}

	span.SetAttributes(attribute.String("subcategory_name", subcategory.Name))
	s.logger.Debug(ctx, "Subcategory fetched from repository", zap.String("id", id), zap.String("name", subcategory.Name))

	go func() {
		if subcategoryJSON, err := json.Marshal(subcategory); err == nil {
			redisStart := time.Now()
			setErr := s.redisClient.Set(context.Background(), cacheKey, subcategoryJSON, 30*time.Second).Err()
			if s.metrics != nil {
				s.metrics.RecordRedisOp(context.Background(), "set", time.Since(redisStart), setErr)
			}
			s.logger.Debug(context.Background(), "Subcategory cached", zap.String("id", id))
		}
	}()

	return subcategory, nil
}

func (s *subcategoryService) GetSubcategoriesByCategory(ctx context.Context, categoryID string, params domain.PaginationParams) ([]*domain.Subcategory, error) {
	ctx, span := observability.StartSpan(ctx, "SubcategoryService.GetSubcategoriesByCategory",
		attribute.String("category_id", categoryID))
	defer span.End()

	cacheKey := fmt.Sprintf("subcategories:by_category:%s:page:%d:size:%d", categoryID, params.Page, params.Size)

	s.logger.Debug(ctx, "Checking cache for subcategories by category", zap.String("category_id", categoryID), zap.Int("page", params.Page), zap.Int("size", params.Size))
	redisStart := time.Now()
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if s.metrics != nil {
		s.metrics.RecordRedisOp(ctx, "get", time.Since(redisStart), err)
	}
	if err == nil {
		var subcategories []*domain.Subcategory
		if err := json.Unmarshal([]byte(cached), &subcategories); err == nil {
			span.SetAttributes(attribute.Bool("cache_hit", true))
			if s.metrics != nil {
				s.metrics.RecordCacheAccess(ctx, "list_subcategories_by_category", true)
			}
			s.logger.Debug(ctx, "Subcategories by category found in cache", zap.Int("count", len(subcategories)))
			return subcategories, nil
		}
	}

	span.SetAttributes(attribute.Bool("cache_hit", false))
	if s.metrics != nil {
		s.metrics.RecordCacheAccess(ctx, "list_subcategories_by_category", false)
	}
	s.logger.Debug(ctx, "Subcategories by category not in cache, fetching from repository", zap.String("category_id", categoryID), zap.Int("page", params.Page), zap.Int("size", params.Size))

	subcategories, err := s.repo.GetByCategoryID(ctx, categoryID, params)
	if err != nil {
		observability.RecordError(span, err, "Failed to get subcategories by category")
		return nil, err
	}

	span.SetAttributes(
		attribute.Int("subcategory_count", len(subcategories)),
		attribute.String("category_id", categoryID),
		attribute.Int("page", params.Page),
		attribute.Int("size", params.Size),
	)
	s.logger.Debug(ctx, "Subcategories by category fetched from repository",
		zap.Int("count", len(subcategories)),
		zap.String("category_id", categoryID),
		zap.Int("page", params.Page),
	)

	go func() {
		if subcategoriesJSON, err := json.Marshal(subcategories); err == nil {
			redisStart := time.Now()
			setErr := s.redisClient.Set(context.Background(), cacheKey, subcategoriesJSON, 30*time.Second).Err()
			if s.metrics != nil {
				s.metrics.RecordRedisOp(context.Background(), "set", time.Since(redisStart), setErr)
			}
			s.logger.Debug(context.Background(), "Subcategories by category cached", zap.Int("count", len(subcategories)))
		}
	}()

	return subcategories, nil
}

func (s *subcategoryService) CreateSubcategory(ctx context.Context, subcategory *domain.Subcategory) error {
	ctx, span := observability.StartSpan(ctx, "SubcategoryService.CreateSubcategory",
		attribute.String("subcategory_name", subcategory.Name))
	defer span.End()

	// Check if ID already exists
	exists, err := s.repo.ExistsByID(ctx, subcategory.ID)
	if err != nil {
		observability.RecordError(span, err, "Failed to check subcategory existence")
		return err
	}
	if exists {
		err := fmt.Errorf("subcategory with id %s already exists", subcategory.ID)
		observability.RecordError(span, err, "Subcategory ID already exists")
		return err
	}

	s.logger.Debug(ctx, "Creating subcategory", zap.String("name", subcategory.Name), zap.String("category_id", subcategory.CategoryID), zap.String("id", subcategory.ID))

	err = s.repo.Create(ctx, subcategory)
	if err != nil {
		observability.RecordError(span, err, "Failed to create subcategory")
		return err
	}

	span.SetAttributes(attribute.String("subcategory_id", subcategory.ID))
	s.logger.Info(ctx, "Subcategory created", zap.String("id", subcategory.ID), zap.String("name", subcategory.Name))

	// Invalidate cache
	go func() {
		redisStart := time.Now()
		keys := []string{
			"subcategories:list:*",
			fmt.Sprintf("subcategories:by_category:%s:*", subcategory.CategoryID),
		}
		for _, key := range keys {
			setErr := s.redisClient.Del(context.Background(), key).Err()
			if s.metrics != nil {
				s.metrics.RecordRedisOp(context.Background(), "del", time.Since(redisStart), setErr)
			}
		}
		s.logger.Debug(context.Background(), "Subcategory list cache invalidated")
	}()

	return nil
}

func (s *subcategoryService) UpdateSubcategory(ctx context.Context, subcategory *domain.Subcategory) error {
	ctx, span := observability.StartSpan(ctx, "SubcategoryService.UpdateSubcategory",
		attribute.String("subcategory_id", subcategory.ID))
	defer span.End()

	s.logger.Debug(ctx, "Updating subcategory", zap.String("id", subcategory.ID), zap.String("name", subcategory.Name))

	err := s.repo.Update(ctx, subcategory)
	if err != nil {
		observability.RecordError(span, err, "Failed to update subcategory")
		return err
	}

	span.SetAttributes(attribute.String("subcategory_name", subcategory.Name))
	s.logger.Info(ctx, "Subcategory updated", zap.String("id", subcategory.ID), zap.String("name", subcategory.Name))

	// Invalidate cache
	go func() {
		redisStart := time.Now()
		keys := []string{
			fmt.Sprintf("subcategories:get:%s", subcategory.ID),
			"subcategories:list:*",
			fmt.Sprintf("subcategories:by_category:%s:*", subcategory.CategoryID),
		}
		for _, key := range keys {
			setErr := s.redisClient.Del(context.Background(), key).Err()
			if s.metrics != nil {
				s.metrics.RecordRedisOp(context.Background(), "del", time.Since(redisStart), setErr)
			}
		}
		s.logger.Debug(context.Background(), "Subcategory cache invalidated")
	}()

	return nil
}

func (s *subcategoryService) DeleteSubcategory(ctx context.Context, id string) error {
	ctx, span := observability.StartSpan(ctx, "SubcategoryService.DeleteSubcategory",
		attribute.String("subcategory_id", id))
	defer span.End()

	s.logger.Debug(ctx, "Deleting subcategory", zap.String("id", id))

	err := s.repo.Delete(ctx, id)
	if err != nil {
		observability.RecordError(span, err, "Failed to delete subcategory")
		return err
	}

	s.logger.Info(ctx, "Subcategory deleted", zap.String("id", id))

	// Invalidate cache
	go func() {
		redisStart := time.Now()
		keys := []string{
			fmt.Sprintf("subcategories:get:%s", id),
			"subcategories:list:*",
		}
		for _, key := range keys {
			setErr := s.redisClient.Del(context.Background(), key).Err()
			if s.metrics != nil {
				s.metrics.RecordRedisOp(context.Background(), "del", time.Since(redisStart), setErr)
			}
		}
		s.logger.Debug(context.Background(), "Subcategory cache invalidated")
	}()

	return nil
}
