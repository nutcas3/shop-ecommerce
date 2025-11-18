package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nutcase/shop-ecommerce/cart-service/internal/handlers"
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
	
	mux.HandleFunc("GET /api/carts/{user_id}", handler.GetCart)
	mux.HandleFunc("POST /api/carts/{user_id}/items", handler.AddToCart)
	mux.HandleFunc("PUT /api/carts/{user_id}/items/{product_id}", handler.UpdateCartItem)
	mux.HandleFunc("DELETE /api/carts/{user_id}/items/{product_id}", handler.RemoveFromCart)
	mux.HandleFunc("DELETE /api/carts/{user_id}", handler.ClearCart)
	server := &http.Server{
		Addr:    ":8083",
		Handler: mux,
	}
	go func() {
		sugar.Infof("Starting cart service on port 8083")
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
