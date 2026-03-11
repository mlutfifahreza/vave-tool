package grpc

import (
	"context"

	"github.com/vave-tool/backend/internal/domain"
	"github.com/vave-tool/backend/proto"
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
	products, err := s.service.ListProducts(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch products: %v", err)
	}

	pbProducts := make([]*proto.Product, 0, len(products))
	for _, p := range products {
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
	if p.Category != nil {
		pbProduct.Category = *p.Category
	}
	if p.SKU != nil {
		pbProduct.Sku = *p.SKU
	}

	return pbProduct
}
