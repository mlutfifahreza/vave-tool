package grpc

import (
	"context"

	"github.com/vave-tool/internal/domain"
	"github.com/vave-tool/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CategoryServer struct {
	proto.UnimplementedCategoryServiceServer
	service domain.CategoryService
}

func NewCategoryServer(service domain.CategoryService) *CategoryServer {
	return &CategoryServer{
		service: service,
	}
}

func (s *CategoryServer) ListCategories(ctx context.Context, req *proto.ListCategoriesRequest) (*proto.ListCategoriesResponse, error) {
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}

	size := int(req.Size)
	if size <= 0 {
		size = 100
	}

	paginationParams := domain.PaginationParams{
		Page: page,
		Size: size,
	}

	categories, err := s.service.ListCategories(ctx, paginationParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch categories: %v", err)
	}

	pbCategories := make([]*proto.Category, 0, len(categories))
	for _, c := range categories {
		pbCategories = append(pbCategories, toPBCategory(c))
	}

	return &proto.ListCategoriesResponse{
		Categories: pbCategories,
	}, nil
}

func (s *CategoryServer) GetCategory(ctx context.Context, req *proto.GetCategoryRequest) (*proto.GetCategoryResponse, error) {
	category, err := s.service.GetCategory(ctx, req.Id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "category not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch category: %v", err)
	}

	return &proto.GetCategoryResponse{
		Category: toPBCategory(category),
	}, nil
}

func (s *CategoryServer) CreateCategory(ctx context.Context, req *proto.CreateCategoryRequest) (*proto.CreateCategoryResponse, error) {
	category := &domain.Category{
		ID:       req.Id,
		Name:     req.Name,
		IsActive: req.IsActive,
	}

	if req.Description != "" {
		category.Description = &req.Description
	}

	err := s.service.CreateCategory(ctx, category)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create category: %v", err)
	}

	return &proto.CreateCategoryResponse{
		Category: toPBCategory(category),
	}, nil
}

func (s *CategoryServer) UpdateCategory(ctx context.Context, req *proto.UpdateCategoryRequest) (*proto.UpdateCategoryResponse, error) {
	category := &domain.Category{
		ID:       req.Id,
		Name:     req.Name,
		IsActive: req.IsActive,
	}

	if req.Description != "" {
		category.Description = &req.Description
	}

	err := s.service.UpdateCategory(ctx, category)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "category not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update category: %v", err)
	}

	return &proto.UpdateCategoryResponse{
		Category: toPBCategory(category),
	}, nil
}

func (s *CategoryServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*proto.DeleteCategoryResponse, error) {
	err := s.service.DeleteCategory(ctx, req.Id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "category not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete category: %v", err)
	}

	return &proto.DeleteCategoryResponse{
		Success: true,
	}, nil
}

func toPBCategory(c *domain.Category) *proto.Category {
	pbCategory := &proto.Category{
		Id:        c.ID,
		Name:      c.Name,
		IsActive:  c.IsActive,
		CreatedAt: timestamppb.New(c.CreatedAt),
		UpdatedAt: timestamppb.New(c.UpdatedAt),
	}

	if c.Description != nil {
		pbCategory.Description = *c.Description
	}

	return pbCategory
}
