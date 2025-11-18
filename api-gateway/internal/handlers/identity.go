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

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("api-gateway").Start(r.Context(), "RegisterUser")
	defer span.End()

	var registerReq RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	reqBody, err := json.Marshal(registerReq)
	if err != nil {
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", h.cfg.IdentityServiceURL+"/api/users/register", bytes.NewBuffer(reqBody))
	if err != nil {
		h.logger.Errorw("Failed to create request to identity service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorw("Failed to send request to identity service", "error", err)
		http.Error(w, "Failed to communicate with identity service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorw("Failed to read response from identity service", "error", err)
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
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("api-gateway").Start(r.Context(), "LoginUser")
	defer span.End()

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	reqBody, err := json.Marshal(loginReq)
	if err != nil {
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", h.cfg.IdentityServiceURL+"/api/users/login", bytes.NewBuffer(reqBody))
	if err != nil {
		h.logger.Errorw("Failed to create request to identity service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorw("Failed to send request to identity service", "error", err)
		http.Error(w, "Failed to communicate with identity service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorw("Failed to read response from identity service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set response status code and body
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

// GetUserProfile forwards the profile request to the identity service
func (h *Handler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("api-gateway").Start(r.Context(), "GetUserProfile")
	defer span.End()

	// Get user claims from context
	userClaims, ok := r.Context().Value(middleware.UserKey).(*middleware.UserClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create HTTP request to identity service
	req, err := http.NewRequestWithContext(ctx, "GET", h.cfg.IdentityServiceURL+"/api/users/"+userClaims.UserID, nil)
	if err != nil {
		h.logger.Errorw("Failed to create request to identity service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set headers
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	// Inject tracing context
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	// Send request to identity service
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorw("Failed to send request to identity service", "error", err)
		http.Error(w, "Failed to communicate with identity service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorw("Failed to read response from identity service", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set response status code and body
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}
