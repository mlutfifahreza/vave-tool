package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vave-tool/internal/api/response"
	"github.com/vave-tool/internal/domain"
	"github.com/vave-tool/internal/observability"
	"go.uber.org/zap"
)

type ProductHandler struct {
	service domain.ProductService
	logger  *observability.Logger
}

func NewProductHandler(service domain.ProductService, logger *observability.Logger) *ProductHandler {
	return &ProductHandler{
		service: service,
		logger:  logger,
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

	h.logger.Info(ctx, "Listing all products")

	products, err := h.service.ListProducts(ctx)
	if err != nil {
		h.logger.Error(ctx, "Failed to fetch products", zap.Error(err))
		response.Error(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	h.logger.Info(ctx, "Products fetched successfully", zap.Int("count", len(products)))
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
		h.logger.Warn(ctx, "Product ID is missing")
		response.Error(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	h.logger.Info(ctx, "Fetching product by ID", zap.String("product_id", id))

	product, err := h.service.GetProduct(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			h.logger.Warn(ctx, "Product not found", zap.String("product_id", id))
			response.Error(w, http.StatusNotFound, "Product not found")
			return
		}
		h.logger.Error(ctx, "Failed to fetch product", zap.String("product_id", id), zap.Error(err))
		response.Error(w, http.StatusInternalServerError, "Failed to fetch product")
		return
	}

	h.logger.Info(ctx, "Product fetched successfully", zap.String("product_id", id))
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
		h.logger.Warn(ctx, "Invalid request body", zap.Error(err))
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	h.logger.Info(ctx, "Creating new product", zap.String("product_name", product.Name))

	if err := h.service.CreateProduct(ctx, &product); err != nil {
		h.logger.Error(ctx, "Failed to create product", zap.Error(err))
		response.Error(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	h.logger.Info(ctx, "Product created successfully", zap.String("product_id", product.ID))
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
		h.logger.Warn(ctx, "Product ID is missing")
		response.Error(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		h.logger.Warn(ctx, "Invalid request body", zap.Error(err))
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product.ID = id

	h.logger.Info(ctx, "Updating product", zap.String("product_id", id))

	if err := h.service.UpdateProduct(ctx, &product); err != nil {
		if err == domain.ErrNotFound {
			h.logger.Warn(ctx, "Product not found", zap.String("product_id", id))
			response.Error(w, http.StatusNotFound, "Product not found")
			return
		}
		h.logger.Error(ctx, "Failed to update product", zap.String("product_id", id), zap.Error(err))
		response.Error(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	h.logger.Info(ctx, "Product updated successfully", zap.String("product_id", id))
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
		h.logger.Warn(ctx, "Product ID is missing")
		response.Error(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	h.logger.Info(ctx, "Deleting product", zap.String("product_id", id))

	if err := h.service.DeleteProduct(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			h.logger.Warn(ctx, "Product not found", zap.String("product_id", id))
			response.Error(w, http.StatusNotFound, "Product not found")
			return
		}
		h.logger.Error(ctx, "Failed to delete product", zap.String("product_id", id), zap.Error(err))
		response.Error(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	h.logger.Info(ctx, "Product deleted successfully", zap.String("product_id", id))
	response.Success(w, map[string]bool{"success": true})
}
