package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nutcase/shop-ecommerce/product-service/internal/handlers"
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
	
	mux.HandleFunc("GET /api/products", handler.ListProducts)
	mux.HandleFunc("GET /api/products/{id}", handler.GetProduct)
	mux.HandleFunc("POST /api/products", handler.CreateProduct)
	mux.HandleFunc("PUT /api/products/{id}", handler.UpdateProduct)
	mux.HandleFunc("DELETE /api/products/{id}", handler.DeleteProduct)

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	go func() {
		sugar.Infof("Starting product service on port 8081")
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
