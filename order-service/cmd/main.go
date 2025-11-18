package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nutcase/shop-ecommerce/order-service/internal/handlers"
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
	
	mux.HandleFunc("POST /api/orders", handler.CreateOrder)
	mux.HandleFunc("GET /api/orders", handler.GetOrders)
	mux.HandleFunc("GET /api/orders/{id}", handler.GetOrder)
	mux.HandleFunc("POST /api/orders/{id}/cancel", handler.CancelOrder)

	server := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	go func() {
		sugar.Infof("Starting order service on port 8082")
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
