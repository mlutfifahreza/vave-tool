package handler

import (
	"net/http"

	"github.com/vave-tool/backend/internal/api/response"
	"github.com/vave-tool/backend/internal/domain"
)

type ProductHandler struct {
	service domain.ProductService
}

func NewProductHandler(service domain.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	products, err := h.service.ListProducts(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	response.Success(w, products)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("id")

	if id == "" {
		response.Error(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	product, err := h.service.GetProduct(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			response.Error(w, http.StatusNotFound, "Product not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to fetch product")
		return
	}

	response.Success(w, product)
}
