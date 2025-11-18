# OpenTelemetry OTLP Migration Guide

## Overview

This document describes the migration from the deprecated Jaeger exporter to the recommended OpenTelemetry Protocol (OTLP) exporter in the shop-ecommerce application.

## Background

OpenTelemetry dropped support for the Jaeger exporter in July 2023. Jaeger now officially accepts and recommends using OTLP instead. This migration updates our codebase to use the recommended approach.

## Changes Made

1. **Updated Dependencies**:
   - Removed: `go.opentelemetry.io/otel/exporters/jaeger`
   - Added: 
     - `go.opentelemetry.io/otel/exporters/otlp/otlptrace`
     - `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc`

2. **Configuration Updates**:
   - Added new configuration fields for OTLP in the `Config` struct:
     - `OTLPEndpoint`: The host address of the OTLP collector
     - `OTLPPort`: The port of the OTLP collector (default: 4317)
   - Kept the Jaeger configuration fields for backward compatibility:
     - `JaegerAgentHost` (deprecated)
     - `JaegerAgentPort` (deprecated)

3. **Environment Variables**:
   - Added new environment variables:
     - `OTLP_ENDPOINT`: The host address of the OTLP collector
     - `OTLP_PORT`: The port of the OTLP collector
   - Kept the Jaeger environment variables for backward compatibility:
     - `JAEGER_AGENT_HOST` (deprecated)
     - `JAEGER_AGENT_PORT` (deprecated)

4. **Tracing Initialization**:
   - Updated the `InitTracer` function to use the OTLP exporter instead of Jaeger
   - Added context handling for the OTLP exporter
   - Added timeout configuration for the OTLP client

## How to Use

### Environment Configuration

Update your environment variables to include the OTLP configuration:

```env
# OpenTelemetry Protocol (OTLP) configuration
OTLP_ENDPOINT=jaeger  # or your collector hostname
OTLP_PORT=4317        # default OTLP gRPC port
```

### Docker Compose

The docker-compose.dev.yml file has been updated to include the OTLP configuration for the api-gateway service:

```yaml
api-gateway:
  environment:
    # OTLP configuration
    - OTLP_ENDPOINT=jaeger
    - OTLP_PORT=4317
```

### Kubernetes

If you're deploying to Kubernetes, make sure to update your ConfigMaps or Secrets to include the OTLP configuration.

## Compatibility

The OpenTelemetry API used in the handlers and middleware remains unchanged, so existing code that uses the OpenTelemetry API will continue to work without modification. The only changes are in the exporter configuration.

## Benefits of OTLP

- **Future-proof**: OTLP is the recommended protocol for OpenTelemetry
- **Vendor-neutral**: OTLP is a vendor-neutral protocol
- **Better performance**: OTLP is designed for high-performance tracing
- **More features**: OTLP supports more features than the Jaeger exporter

## Troubleshooting

If you encounter issues with the OTLP exporter:

1. Verify that the OTLP collector is running and accessible
2. Check that the OTLP endpoint and port are correctly configured
3. Ensure that the OTLP gRPC port (4317) is open in your firewall
4. Check the logs for any errors related to the OTLP exporter

## References

- [OpenTelemetry OTLP Documentation](https://opentelemetry.io/docs/specs/otlp/)
- [Jaeger OTLP Documentation](https://www.jaegertracing.io/docs/latest/apis/#opentelemetry-protocol-otlp)
