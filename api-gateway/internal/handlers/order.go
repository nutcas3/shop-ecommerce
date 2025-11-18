package handlers

import (
"bytes"
"encoding/json"
"io"
"net/http"
"strings"

"github.com/nutcase/shop-ecommerce/api-gateway/internal/middleware"
"go.opentelemetry.io/otel"
"go.opentelemetry.io/otel/propagation"
)

type CreateOrderRequest struct {
	ShippingAddress string `json:"shipping_address"`
	PaymentMethod   string `json:"payment_method"`
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("api-gateway").Start(r.Context(), "CreateOrder")
	defer span.End()

	userClaims, ok := r.Context().Value(middleware.UserKey).(*middleware.UserClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var createOrderReq CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&createOrderReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	reqBody, err := json.Marshal(map[string]interface{}{
		"user_id":          userClaims.UserID,
		"shipping_address": createOrderReq.ShippingAddress,
		"payment_method":   createOrderReq.PaymentMethod,
	})
	if err != nil {
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", h.cfg.OrderServiceURL+"/api/orders", bytes.NewBuffer(reqBody))
	if err != nil {
		h.logger.Errorw("Failed to create request to order service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorw("Failed to send request to order service", "error", err)
		http.Error(w, "Failed to communicate with order service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorw("Failed to read response from order service", "error", err)
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

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("api-gateway").Start(r.Context(), "GetOrders")
	defer span.End()
	userClaims, ok := r.Context().Value(middleware.UserKey).(*middleware.UserClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "GET", h.cfg.OrderServiceURL+"/api/orders?user_id="+userClaims.UserID, nil)
	if err != nil {
		h.logger.Errorw("Failed to create request to order service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorw("Failed to send request to order service", "error", err)
		http.Error(w, "Failed to communicate with order service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorw("Failed to read response from order service", "error", err)
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
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("api-gateway").Start(r.Context(), "GetOrder")
	defer span.End()

	_, ok := r.Context().Value(middleware.UserKey).(*middleware.UserClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	orderID := parts[len(parts)-1]

	req, err := http.NewRequestWithContext(ctx, "GET", h.cfg.OrderServiceURL+"/api/orders/"+orderID, nil)
	if err != nil {
		h.logger.Errorw("Failed to create request to order service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorw("Failed to send request to order service", "error", err)
		http.Error(w, "Failed to communicate with order service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorw("Failed to read response from order service", "error", err)
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
func (h *Handler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("api-gateway").Start(r.Context(), "CancelOrder")
	defer span.End()

	_, ok := r.Context().Value(middleware.UserKey).(*middleware.UserClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	orderID := parts[len(parts)-2]

	req, err := http.NewRequestWithContext(ctx, "POST", h.cfg.OrderServiceURL+"/api/orders/"+orderID+"/cancel", nil)
	if err != nil {
		h.logger.Errorw("Failed to create request to order service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorw("Failed to send request to order service", "error", err)
		http.Error(w, "Failed to communicate with order service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorw("Failed to read response from order service", "error", err)
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
