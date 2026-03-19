package grpc

import (
	"context"

	"github.com/vave-tool/internal/domain"
	"github.com/vave-tool/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SubcategoryServer struct {
	proto.UnimplementedSubcategoryServiceServer
	service domain.SubcategoryService
}

func NewSubcategoryServer(service domain.SubcategoryService) *SubcategoryServer {
	return &SubcategoryServer{
		service: service,
	}
}

func (s *SubcategoryServer) ListSubcategories(ctx context.Context, req *proto.ListSubcategoriesRequest) (*proto.ListSubcategoriesResponse, error) {
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

	subcategories, err := s.service.ListSubcategories(ctx, paginationParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch subcategories: %v", err)
	}

	pbSubcategories := make([]*proto.Subcategory, 0, len(subcategories))
	for _, sc := range subcategories {
		pbSubcategories = append(pbSubcategories, toPBSubcategory(sc))
	}

	return &proto.ListSubcategoriesResponse{
		Subcategories: pbSubcategories,
	}, nil
}

func (s *SubcategoryServer) GetSubcategory(ctx context.Context, req *proto.GetSubcategoryRequest) (*proto.GetSubcategoryResponse, error) {
	subcategory, err := s.service.GetSubcategory(ctx, req.Id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "subcategory not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch subcategory: %v", err)
	}

	return &proto.GetSubcategoryResponse{
		Subcategory: toPBSubcategory(subcategory),
	}, nil
}

func (s *SubcategoryServer) GetSubcategoriesByCategory(ctx context.Context, req *proto.GetSubcategoriesByCategoryRequest) (*proto.GetSubcategoriesByCategoryResponse, error) {
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

	subcategories, err := s.service.GetSubcategoriesByCategory(ctx, req.CategoryId, paginationParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch subcategories by category: %v", err)
	}

	pbSubcategories := make([]*proto.Subcategory, 0, len(subcategories))
	for _, sc := range subcategories {
		pbSubcategories = append(pbSubcategories, toPBSubcategory(sc))
	}

	return &proto.GetSubcategoriesByCategoryResponse{
		Subcategories: pbSubcategories,
	}, nil
}

func (s *SubcategoryServer) CreateSubcategory(ctx context.Context, req *proto.CreateSubcategoryRequest) (*proto.CreateSubcategoryResponse, error) {
	subcategory := &domain.Subcategory{
		ID:         req.Id,
		CategoryID: req.CategoryId,
		Name:       req.Name,
		IsActive:   req.IsActive,
	}

	if req.Description != "" {
		subcategory.Description = &req.Description
	}

	err := s.service.CreateSubcategory(ctx, subcategory)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create subcategory: %v", err)
	}

	return &proto.CreateSubcategoryResponse{
		Subcategory: toPBSubcategory(subcategory),
	}, nil
}

func (s *SubcategoryServer) UpdateSubcategory(ctx context.Context, req *proto.UpdateSubcategoryRequest) (*proto.UpdateSubcategoryResponse, error) {
	subcategory := &domain.Subcategory{
		ID:         req.Id,
		CategoryID: req.CategoryId,
		Name:       req.Name,
		IsActive:   req.IsActive,
	}

	if req.Description != "" {
		subcategory.Description = &req.Description
	}

	err := s.service.UpdateSubcategory(ctx, subcategory)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "subcategory not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update subcategory: %v", err)
	}

	return &proto.UpdateSubcategoryResponse{
		Subcategory: toPBSubcategory(subcategory),
	}, nil
}

func (s *SubcategoryServer) DeleteSubcategory(ctx context.Context, req *proto.DeleteSubcategoryRequest) (*proto.DeleteSubcategoryResponse, error) {
	err := s.service.DeleteSubcategory(ctx, req.Id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "subcategory not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete subcategory: %v", err)
	}

	return &proto.DeleteSubcategoryResponse{
		Success: true,
	}, nil
}

func toPBSubcategory(sc *domain.Subcategory) *proto.Subcategory {
	pbSubcategory := &proto.Subcategory{
		Id:         sc.ID,
		CategoryId: sc.CategoryID,
		Name:       sc.Name,
		IsActive:   sc.IsActive,
		CreatedAt:  timestamppb.New(sc.CreatedAt),
		UpdatedAt:  timestamppb.New(sc.UpdatedAt),
	}

	if sc.Description != nil {
		pbSubcategory.Description = *sc.Description
	}

	return pbSubcategory
}
