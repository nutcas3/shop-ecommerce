package handlers

import (
	"github.com/nutcase/shop-ecommerce/api-gateway/internal/config"
	"go.uber.org/zap"
)

type Handler struct {
	cfg    *config.Config
	logger *zap.SugaredLogger
}

func NewHandler(cfg *config.Config, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		cfg:    cfg,
		logger: logger,
	}
}
