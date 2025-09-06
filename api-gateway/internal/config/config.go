package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port                int    `mapstructure:"API_GATEWAY_PORT"`
	JWTSecret           string `mapstructure:"JWT_SECRET"`
	IdentityServiceURL  string `mapstructure:"IDENTITY_SERVICE_URL"`
	ProductServiceURL   string `mapstructure:"PRODUCT_SERVICE_URL"`
	CartServiceURL      string `mapstructure:"CART_SERVICE_URL"`
	OrderServiceURL     string `mapstructure:"ORDER_SERVICE_URL"`
	// Deprecated: Use OTLPEndpoint and OTLPPort instead
	JaegerAgentHost     string `mapstructure:"JAEGER_AGENT_HOST"`
	// Deprecated: Use OTLPEndpoint and OTLPPort instead
	JaegerAgentPort     int    `mapstructure:"JAEGER_AGENT_PORT"`
	// OTLP exporter configuration
	OTLPEndpoint        string `mapstructure:"OTLP_ENDPOINT"`
	OTLPPort            int    `mapstructure:"OTLP_PORT"`
}

func Load() (*Config, error) {
	viper.SetDefault("API_GATEWAY_PORT", 8080)
	viper.SetDefault("JWT_SECRET", "default_secret_change_me")
	viper.SetDefault("IDENTITY_SERVICE_URL", "http://identity-service:8081")
	viper.SetDefault("PRODUCT_SERVICE_URL", "http://product-service:8082")
	viper.SetDefault("CART_SERVICE_URL", "http://cart-service:8083")
	viper.SetDefault("ORDER_SERVICE_URL", "http://order-service:8084")
	// For backward compatibility
	viper.SetDefault("JAEGER_AGENT_HOST", "jaeger")
	viper.SetDefault("JAEGER_AGENT_PORT", 6831)
	// OTLP defaults
	viper.SetDefault("OTLP_ENDPOINT", "jaeger")
	viper.SetDefault("OTLP_PORT", 4317)

	viper.AutomaticEnv()

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
