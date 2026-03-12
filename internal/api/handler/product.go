package handler

import (
	"encoding/json"
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

// List godoc
// @Summary List all products
// @Description Get a list of all products
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]domain.Product}
// @Failure 500 {object} response.Response
// @Router /api/products [get]
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	products, err := h.service.ListProducts(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	response.Success(w, products)
}

// GetByID godoc
// @Summary Get product by ID
// @Description Get a single product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} response.Response{data=domain.Product}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/products/{id} [get]
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

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

// Create godoc
// @Summary Create new product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body domain.Product true "Product object"
// @Success 201 {object} response.Response{data=domain.Product}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/products [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.CreateProduct(ctx, &product); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	response.Created(w, &product)
}

// Update godoc
// @Summary Update product
// @Description Update an existing product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body domain.Product true "Product object"
// @Success 200 {object} response.Response{data=domain.Product}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/products/{id} [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	if id == "" {
		response.Error(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product.ID = id

	if err := h.service.UpdateProduct(ctx, &product); err != nil {
		if err == domain.ErrNotFound {
			response.Error(w, http.StatusNotFound, "Product not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	response.Success(w, &product)
}

// Delete godoc
// @Summary Delete product
// @Description Delete a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/products/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	if id == "" {
		response.Error(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	if err := h.service.DeleteProduct(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			response.Error(w, http.StatusNotFound, "Product not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	response.Success(w, map[string]bool{"success": true})
}
