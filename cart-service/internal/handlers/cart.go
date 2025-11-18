package handlers

import (
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel"
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

type CartItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
}

type Cart struct {
	UserID string     `json:"user_id"`
	Items  []CartItem `json:"items"`
	Total  float64    `json:"total"`
}

type AddToCartRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer("cart-service").Start(r.Context(), "GetCart")
	defer span.End()

	userID := r.PathValue("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Add user ID to logs instead of span attributes
	h.logger.Infow("Getting cart", "user_id", userID)

	// TODO: Implement repository call to get cart
	// cart, err := h.repo.GetCart(ctx, userID)
	
	// Mock data for now
	cart := Cart{
		UserID: userID,
		Items: []CartItem{
			{
				ProductID:   "product-1",
				ProductName: "Sample Product 1",
				Quantity:    2,
				Price:       29.99,
				ImageURL:    "https://example.com/product1.jpg",
			},
			{
				ProductID:   "product-2",
				ProductName: "Sample Product 2",
				Quantity:    1,
				Price:       39.99,
				ImageURL:    "https://example.com/product2.jpg",
			},
		},
		Total: 99.97,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	// Write response
	if err := json.NewEncoder(w).Encode(cart); err != nil {
		h.logger.Errorw("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer("cart-service").Start(r.Context(), "AddToCart")
	defer span.End()

	userID := r.PathValue("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var req AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ProductID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}
	if req.Quantity <= 0 {
		http.Error(w, "Quantity must be greater than 0", http.StatusBadRequest)
		return
	}

	h.logger.Infow("Adding to cart", "user_id", userID, "product_id", req.ProductID, "quantity", req.Quantity)

	// TODO: Implement repository call to add item to cart
	// err := h.repo.AddToCart(ctx, userID, req.ProductID, req.Quantity)
	
	// Mock data for now
	cart := Cart{
		UserID: userID,
		Items: []CartItem{
			{
				ProductID:   "product-1",
				ProductName: "Sample Product 1",
				Quantity:    2,
				Price:       29.99,
				ImageURL:    "https://example.com/product1.jpg",
			},
			{
				ProductID:   req.ProductID,
				ProductName: "Sample Product " + req.ProductID,
				Quantity:    req.Quantity,
				Price:       39.99,
				ImageURL:    "https://example.com/product" + req.ProductID + ".jpg",
			},
		},
		Total: 99.97 + (39.99 * float64(req.Quantity)),
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	// Write response
	if err := json.NewEncoder(w).Encode(cart); err != nil {
		h.logger.Errorw("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer("cart-service").Start(r.Context(), "UpdateCartItem")
	defer span.End()

	userID := r.PathValue("user_id")
	productID := r.PathValue("product_id")
	
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	var req AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		http.Error(w, "Quantity must be greater than 0", http.StatusBadRequest)
		return
	}

	h.logger.Infow("Updating cart item", "user_id", userID, "product_id", productID, "quantity", req.Quantity)

	// TODO: Implement repository call to update item in cart
	// err := h.repo.UpdateCartItem(ctx, userID, productID, req.Quantity)
	
	// Mock data for now
	cart := Cart{
		UserID: userID,
		Items: []CartItem{
			{
				ProductID:   "product-1",
				ProductName: "Sample Product 1",
				Quantity:    2,
				Price:       29.99,
				ImageURL:    "https://example.com/product1.jpg",
			},
			{
				ProductID:   productID,
				ProductName: "Sample Product " + productID,
				Quantity:    req.Quantity,
				Price:       39.99,
				ImageURL:    "https://example.com/product" + productID + ".jpg",
			},
		},
		Total: 59.98 + (39.99 * float64(req.Quantity)),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(cart); err != nil {
		h.logger.Errorw("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer("cart-service").Start(r.Context(), "RemoveFromCart")
	defer span.End()

	userID := r.PathValue("user_id")
	productID := r.PathValue("product_id")
	
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	h.logger.Infow("Removing from cart", "user_id", userID, "product_id", productID)

	// TODO: Implement repository call to remove item from cart
	// err := h.repo.RemoveFromCart(ctx, userID, productID)
	
	// Mock data for now
	cart := Cart{
		UserID: userID,
		Items: []CartItem{
			{
				ProductID:   "product-1",
				ProductName: "Sample Product 1",
				Quantity:    2,
				Price:       29.99,
				ImageURL:    "https://example.com/product1.jpg",
			},
		},
		Total: 59.98,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(cart); err != nil {
		h.logger.Errorw("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer("cart-service").Start(r.Context(), "ClearCart")
	defer span.End()
	userID := r.PathValue("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	h.logger.Infow("Clearing cart", "user_id", userID)

	// TODO: Implement repository call to clear cart
	// err := h.repo.ClearCart(ctx, userID)
	
	// Mock data for now
	cart := Cart{
		UserID: userID,
		Items:  []CartItem{},
		Total:  0,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(cart); err != nil {
		h.logger.Errorw("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
