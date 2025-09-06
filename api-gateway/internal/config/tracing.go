package config

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// InitTracer initializes the OpenTelemetry tracer
func InitTracer(cfg *Config) (*tracesdk.TracerProvider, error) {
	// Create Jaeger exporter
	exp, err := jaeger.New(jaeger.WithAgentEndpoint(
		jaeger.WithAgentHost(cfg.JaegerAgentHost),
		jaeger.WithAgentPort(fmt.Sprintf("%d", cfg.JaegerAgentPort)),
	))
	if err != nil {
		return nil, err
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
