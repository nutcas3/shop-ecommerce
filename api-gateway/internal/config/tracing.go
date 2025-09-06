package config

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// InitTracer initializes the OpenTelemetry tracer
func InitTracer(cfg *Config) (*tracesdk.TracerProvider, error) {
	ctx := context.Background()

	// Configure OTLP exporter client options
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%d", cfg.OTLPEndpoint, cfg.OTLPPort)),
		otlptracegrpc.WithInsecure(), // For development; use secure connection in production
		otlptracegrpc.WithTimeout(5 * time.Second),
	}

	// Create OTLP exporter
	exp, err := otlptrace.New(ctx, otlptracegrpc.NewClient(opts...))
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create trace provider with the exporter
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("api-gateway"),
		)),
	)

	// Set global trace provider
	otel.SetTracerProvider(tp)

	return tp, nil
}
