package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/nutcase/shop-ecommerce/api-gateway/internal/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type CartItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type AddToCartRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("api-gateway").Start(r.Context(), "GetCart")
	defer span.End()

	userClaims, ok := r.Context().Value(middleware.UserKey).(*middleware.UserClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "GET", h.cfg.CartServiceURL+"/api/carts/"+userClaims.UserID, nil)
	if err != nil {
		h.logger.Errorw("Failed to create request to cart service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorw("Failed to send request to cart service", "error", err)
		http.Error(w, "Failed to communicate with cart service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorw("Failed to read response from cart service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("api-gateway").Start(r.Context(), "AddToCart")
	defer span.End()

	userClaims, ok := r.Context().Value(middleware.UserKey).(*middleware.UserClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var addToCartReq AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&addToCartReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	reqBody, err := json.Marshal(addToCartReq)
	if err != nil {
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", h.cfg.CartServiceURL+"/api/carts/"+userClaims.UserID+"/items", bytes.NewBuffer(reqBody))
	if err != nil {
		h.logger.Errorw("Failed to create request to cart service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorw("Failed to send request to cart service", "error", err)
		http.Error(w, "Failed to communicate with cart service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorw("Failed to read response from cart service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func (h *Handler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("api-gateway").Start(r.Context(), "UpdateCartItem")
	defer span.End()

	userClaims, ok := r.Context().Value(middleware.UserKey).(*middleware.UserClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var updateCartReq AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&updateCartReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	reqBody, err := json.Marshal(updateCartReq)
	if err != nil {
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", h.cfg.CartServiceURL+"/api/carts/"+userClaims.UserID+"/items/"+updateCartReq.ProductID, bytes.NewBuffer(reqBody))
	if err != nil {
		h.logger.Errorw("Failed to create request to cart service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorw("Failed to send request to cart service", "error", err)
		http.Error(w, "Failed to communicate with cart service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorw("Failed to read response from cart service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func (h *Handler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("api-gateway").Start(r.Context(), "RemoveFromCart")
	defer span.End()

	userClaims, ok := r.Context().Value(middleware.UserKey).(*middleware.UserClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", h.cfg.CartServiceURL+"/api/carts/"+userClaims.UserID+"/items/"+productID, nil)
	if err != nil {
		h.logger.Errorw("Failed to create request to cart service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorw("Failed to send request to cart service", "error", err)
		http.Error(w, "Failed to communicate with cart service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorw("Failed to read response from cart service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func (h *Handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("api-gateway").Start(r.Context(), "ClearCart")
	defer span.End()

	userClaims, ok := r.Context().Value(middleware.UserKey).(*middleware.UserClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", h.cfg.CartServiceURL+"/api/carts/"+userClaims.UserID, nil)
	if err != nil {
		h.logger.Errorw("Failed to create request to cart service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorw("Failed to send request to cart service", "error", err)
		http.Error(w, "Failed to communicate with cart service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorw("Failed to read response from cart service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}
