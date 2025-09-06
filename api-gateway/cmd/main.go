package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nutcase/shop-ecommerce/api-gateway/internal/config"
	"github.com/nutcase/shop-ecommerce/api-gateway/internal/handlers"
	custommiddleware "github.com/nutcase/shop-ecommerce/api-gateway/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	cfg, err := config.Load()
	if err != nil {
		sugar.Fatalf("Failed to load configuration: %v", err)
	}

	tp, err := config.InitTracer(cfg)
	if err != nil {
		sugar.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			sugar.Errorf("Error shutting down tracer provider: %v", err)
		}
	}()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(custommiddleware.TracingMiddleware)
	h := handlers.NewHandler(cfg, sugar)

	r.Get("/health", h.HealthCheck)

	r.Route("/api/identity", func(r chi.Router) {
		r.Post("/register", h.RegisterUser)
		r.Post("/login", h.LoginUser)
		r.With(custommiddleware.AuthMiddleware(cfg.JWTSecret)).Get("/profile", h.GetUserProfile)
	})

	r.Route("/api/products", func(r chi.Router) {
		r.Get("/", h.ListProducts)
		r.Get("/{id}", h.GetProduct)
		r.With(custommiddleware.AuthMiddleware(cfg.JWTSecret)).Post("/", h.CreateProduct)
		r.With(custommiddleware.AuthMiddleware(cfg.JWTSecret)).Put("/{id}", h.UpdateProduct)
		r.With(custommiddleware.AuthMiddleware(cfg.JWTSecret)).Delete("/{id}", h.DeleteProduct)
	})

	r.Route("/api/cart", func(r chi.Router) {
		r.Use(custommiddleware.AuthMiddleware(cfg.JWTSecret))
		r.Get("/", h.GetCart)
		r.Post("/items", h.AddToCart)
		r.Put("/items/{id}", h.UpdateCartItem)
		r.Delete("/items/{id}", h.RemoveFromCart)
		r.Delete("/", h.ClearCart)
	})

	r.Route("/api/orders", func(r chi.Router) {
		r.Use(custommiddleware.AuthMiddleware(cfg.JWTSecret))
		r.Get("/", h.GetOrders)
		r.Get("/{id}", h.GetOrder)
		r.Post("/", h.CreateOrder)
		r.Post("/{id}/cancel", h.CancelOrder)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				sugar.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			sugar.Fatal(err)
		}
		serverStopCtx()
	}()

	sugar.Infof("Starting server on port %d", cfg.Port)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		sugar.Fatal(err)
	}

	<-serverCtx.Done()
}
