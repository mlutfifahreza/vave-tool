package grpc

import (
	"context"

	"github.com/vave-tool/internal/domain"
	"github.com/vave-tool/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductServer struct {
	proto.UnimplementedProductServiceServer
	service domain.ProductService
}

func NewProductServer(service domain.ProductService) *ProductServer {
	return &ProductServer{
		service: service,
	}
}

func (s *ProductServer) ListProducts(ctx context.Context, req *proto.ListProductsRequest) (*proto.ListProductsResponse, error) {
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

	result, err := s.service.ListProducts(ctx, paginationParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch products: %v", err)
	}

	pbProducts := make([]*proto.Product, 0, len(result.Products))
	for _, p := range result.Products {
		pbProducts = append(pbProducts, toPBProduct(p))
	}

	return &proto.ListProductsResponse{
		Products: pbProducts,
	}, nil
}

func (s *ProductServer) GetProduct(ctx context.Context, req *proto.GetProductRequest) (*proto.GetProductResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "product id is required")
	}

	product, err := s.service.GetProduct(ctx, req.Id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch product: %v", err)
	}

	return &proto.GetProductResponse{
		Product: toPBProduct(product),
	}, nil
}

func (s *ProductServer) CreateProduct(ctx context.Context, req *proto.CreateProductRequest) (*proto.CreateProductResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "product name is required")
	}

	product := &domain.Product{
		Name:          req.Name,
		Price:         req.Price,
		StockQuantity: int(req.StockQuantity),
		IsActive:      req.IsActive,
	}

	if req.Description != "" {
		product.Description = &req.Description
	}
	if req.CategoryId != "" {
		product.CategoryID = &req.CategoryId
	}
	if req.SubcategoryId != "" {
		product.SubcategoryID = &req.SubcategoryId
	}
	if req.Sku != "" {
		product.SKU = &req.Sku
	}

	if err := s.service.CreateProduct(ctx, product); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return &proto.CreateProductResponse{
		Product: toPBProduct(product),
	}, nil
}

func (s *ProductServer) UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest) (*proto.UpdateProductResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "product id is required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "product name is required")
	}

	product := &domain.Product{
		ID:            req.Id,
		Name:          req.Name,
		Price:         req.Price,
		StockQuantity: int(req.StockQuantity),
		IsActive:      req.IsActive,
	}

	if req.Description != "" {
		product.Description = &req.Description
	}
	if req.CategoryId != "" {
		product.CategoryID = &req.CategoryId
	}
	if req.SubcategoryId != "" {
		product.SubcategoryID = &req.SubcategoryId
	}
	if req.Sku != "" {
		product.SKU = &req.Sku
	}

	if err := s.service.UpdateProduct(ctx, product); err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}

	return &proto.UpdateProductResponse{
		Product: toPBProduct(product),
	}, nil
}

func (s *ProductServer) DeleteProduct(ctx context.Context, req *proto.DeleteProductRequest) (*proto.DeleteProductResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "product id is required")
	}

	if err := s.service.DeleteProduct(ctx, req.Id); err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}

	return &proto.DeleteProductResponse{
		Success: true,
	}, nil
}

func toPBProduct(p *domain.Product) *proto.Product {
	pbProduct := &proto.Product{
		Id:            p.ID,
		Name:          p.Name,
		Price:         p.Price,
		StockQuantity: int32(p.StockQuantity),
		IsActive:      p.IsActive,
		CreatedAt:     timestamppb.New(p.CreatedAt),
		UpdatedAt:     timestamppb.New(p.UpdatedAt),
	}

	if p.Description != nil {
		pbProduct.Description = *p.Description
	}
	if p.CategoryID != nil {
		pbProduct.CategoryId = *p.CategoryID
	}
	if p.CategoryName != nil {
		pbProduct.CategoryName = *p.CategoryName
	}
	if p.SubcategoryID != nil {
		pbProduct.SubcategoryId = *p.SubcategoryID
	}
	if p.SubcategoryName != nil {
		pbProduct.SubcategoryName = *p.SubcategoryName
	}
	if p.SKU != nil {
		pbProduct.Sku = *p.SKU
	}

	return pbProduct
}
