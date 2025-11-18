package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.SugaredLogger
}

func NewHandler(logger *zap.SugaredLogger) *Handler {
	return &Handler{
		logger: logger,
	}
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	Stock       int     `json:"stock"`
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer("product-service").Start(r.Context(), "ListProducts")
	defer span.End()

	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))
	
	if limit == 0 {
		limit = 10 
	}

	span.SetAttributes(
		attribute.Int("limit", limit),
		attribute.Int("offset", offset),
	)

	// TODO: Implement repository call to get products
	// products, err := h.repo.ListProducts(ctx, limit, offset)
	
	// Mock data for now
	products := []Product{
		{
			ID:          "1",
			Name:        "Product 1",
			Description: "Description for product 1",
			Price:       19.99,
			ImageURL:    "https://example.com/product1.jpg",
			Stock:       100,
		},
		{
			ID:          "2",
			Name:        "Product 2",
			Description: "Description for product 2",
			Price:       29.99,
			ImageURL:    "https://example.com/product2.jpg",
			Stock:       50,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(products); err != nil {
		h.logger.Errorw("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer("product-service").Start(r.Context(), "GetProduct")
	defer span.End()
	productID := r.PathValue("id")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	span.SetAttributes(attribute.String("product.id", productID))

	// TODO: Implement repository call to get product
	// product, err := h.repo.GetProduct(ctx, productID)
	
	// Mock data for now
	product := Product{
		ID:          productID,
		Name:        "Product " + productID,
		Description: "Description for product " + productID,
		Price:       19.99,
		ImageURL:    "https://example.com/product" + productID + ".jpg",
		Stock:       100,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(product); err != nil {
		h.logger.Errorw("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer("product-service").Start(r.Context(), "CreateProduct")
	defer span.End()
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Implement repository call to create product
	// err := h.repo.CreateProduct(ctx, product)
	
	// Mock response for now
	product.ID = "new-product-id" // In a real implementation, this would be generated

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	
	if err := json.NewEncoder(w).Encode(product); err != nil {
		h.logger.Errorw("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer("product-service").Start(r.Context(), "UpdateProduct")
	defer span.End()
	productID := r.PathValue("id")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	span.SetAttributes(attribute.String("product.id", productID))

	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure the ID in the URL matches the ID in the body
	product.ID = productID

	// TODO: Implement repository call to update product
	// err := h.repo.UpdateProduct(ctx, product)
	
	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	// Write response
	if err := json.NewEncoder(w).Encode(product); err != nil {
		h.logger.Errorw("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer("product-service").Start(r.Context(), "DeleteProduct")
	defer span.End()

	productID := r.PathValue("id")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	span.SetAttributes(attribute.String("product.id", productID))

	// TODO: Implement repository call to delete product
	// err := h.repo.DeleteProduct(ctx, productID)
	
	// Set response headers
	w.WriteHeader(http.StatusNoContent)
}
