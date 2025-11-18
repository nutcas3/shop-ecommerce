package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nutcase/shop-ecommerce/identity-service/internal/handlers"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	handler := handlers.NewHandler(sugar)
	mux := http.NewServeMux()
	
	mux.HandleFunc("POST /api/auth/register", handler.Register)
	mux.HandleFunc("POST /api/auth/login", handler.Login)
	mux.HandleFunc("GET /api/users/{id}", handler.GetProfile)
	mux.HandleFunc("PUT /api/users/{id}", handler.UpdateProfile)
	mux.HandleFunc("POST /api/users/{id}/change-password", handler.ChangePassword)

	server := &http.Server{
		Addr:    ":8084",
		Handler: mux,
	}
	go func() {
		sugar.Infof("Starting identity service on port 8084")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	sugar.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		sugar.Fatalf("Server forced to shutdown: %v", err)
	}

	sugar.Info("Server exiting")
}
