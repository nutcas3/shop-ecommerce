package handlers

import (
"bytes"
"io"
"net/http"
"strings"

"go.opentelemetry.io/otel"
"go.opentelemetry.io/otel/propagation"
)

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	Stock       int     `json:"stock"`
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("api-gateway").Start(r.Context(), "ListProducts")
	defer span.End()
	req, err := http.NewRequestWithContext(ctx, "GET", h.cfg.ProductServiceURL+"/api/products", nil)
	if err != nil {
		h.logger.Errorw("Failed to create request to product service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.URL.RawQuery = r.URL.RawQuery

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorw("Failed to send request to product service", "error", err)
		http.Error(w, "Failed to communicate with product service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorw("Failed to read response from product service", "error", err)
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
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("api-gateway").Start(r.Context(), "GetProduct")
	defer span.End()

	parts := strings.Split(r.URL.Path, "/")
	productID := parts[len(parts)-1]

	req, err := http.NewRequestWithContext(ctx, "GET", h.cfg.ProductServiceURL+"/api/products/"+productID, nil)
	if err != nil {
		h.logger.Errorw("Failed to create request to product service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorw("Failed to send request to product service", "error", err)
		http.Error(w, "Failed to communicate with product service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorw("Failed to read response from product service", "error", err)
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
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("api-gateway").Start(r.Context(), "CreateProduct")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", h.cfg.ProductServiceURL+"/api/products", bytes.NewBuffer(body))
	if err != nil {
		h.logger.Errorw("Failed to create request to product service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorw("Failed to send request to product service", "error", err)
		http.Error(w, "Failed to communicate with product service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorw("Failed to read response from product service", "error", err)
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
func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("api-gateway").Start(r.Context(), "UpdateProduct")
	defer span.End()

	parts := strings.Split(r.URL.Path, "/")
	productID := parts[len(parts)-1]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", h.cfg.ProductServiceURL+"/api/products/"+productID, bytes.NewBuffer(body))
	if err != nil {
		h.logger.Errorw("Failed to create request to product service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorw("Failed to send request to product service", "error", err)
		http.Error(w, "Failed to communicate with product service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorw("Failed to read response from product service", "error", err)
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
func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("api-gateway").Start(r.Context(), "DeleteProduct")
	defer span.End()

	parts := strings.Split(r.URL.Path, "/")
	productID := parts[len(parts)-1]

	req, err := http.NewRequestWithContext(ctx, "DELETE", h.cfg.ProductServiceURL+"/api/products/"+productID, nil)
	if err != nil {
		h.logger.Errorw("Failed to create request to product service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorw("Failed to send request to product service", "error", err)
		http.Error(w, "Failed to communicate with product service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorw("Failed to read response from product service", "error", err)
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
