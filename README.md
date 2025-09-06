# eShopOnSteroids - Go & Rust Implementation

A cloud-native online shop powered by Golang, Rust, containers, and Kubernetes.

## Architecture

This project implements a microservices architecture where each service is responsible for a single business capability. The services are strategically implemented in either Go or Rust based on their specific requirements:

### Golang Services

- **API Gateway**: Routes requests, validates authentication, transforms requests/responses
- **Identity Service**: Manages authentication, authorization, and user profiles
- **Product Service**: Manages product catalog and basic inventory information
- **Cart Service**: Manages shopping cart operations with short-lived data

### Rust Services

- **Order Service**: Handles critical order processing and fulfillment
- **Payment Service**: Processes payments with high security requirements
- **Inventory Service**: Manages real-time inventory tracking and stock updates

### Communication Patterns

- **Synchronous**: gRPC for direct service-to-service communication
- **Asynchronous**: Event-driven messaging via NATS/RabbitMQ

### Database Technologies

- Identity Service: PostgreSQL
- Product Service: MongoDB
- Cart Service: Redis
- Order Service: PostgreSQL
- Payment Service: Stateless with transaction logs in PostgreSQL
- Inventory Service: PostgreSQL

### Observability

- Distributed Tracing: OpenTelemetry
- Metrics: Prometheus and Grafana
- Logging: Structured logging with ELK stack

## Setup

### Prerequisites

- Docker and Docker Compose
- Go 1.21+
- Rust 1.70+
- kubectl (for Kubernetes deployment)

### Development

1. Clone the repository
```
git clone https://github.com/nutcas3/shop-ecommerce.git

cd shop-ecommerce
```

2. Create the environment file
```
cp .env.example .env
# Edit .env with your configuration
```

3. Start the development environment
```
docker-compose -f docker-compose.dev.yml up
```

4. Access the application at http://localhost:8080

### Production

For production deployment, Kubernetes manifests are provided in the `deployment` directory.

```
kubectl apply -f deployment/
```

## Project Structure

```
shop-ecommerce/
├── api-gateway/         # Go implementation of API Gateway
├── identity-service/    # Go implementation of Identity Service
├── product-service/     # Go implementation of Product Service
├── cart-service/        # Go implementation of Cart Service
├── order-service/       # Rust implementation of Order Service
├── payment-service/     # Rust implementation of Payment Service
├── inventory-service/   # Rust implementation of Inventory Service
├── proto/              # Shared Protocol Buffer definitions
├── observability/      # Observability stack configuration
└── deployment/         # Kubernetes deployment manifests
```

## License

MIT
