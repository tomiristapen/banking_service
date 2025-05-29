# API Gateway for Banking Microservices

This API Gateway exposes HTTP REST endpoints for all gRPC endpoints of the following microservices:
- user_service
- transaction_service
- payment_service
- smart_budget_service

## Features
- REST → gRPC translation using grpc-gateway
- Clean architecture (handlers, router, config)
- .env config for service addresses
- OpenAPI/Swagger support (planned)

## Setup
1. Copy proto files from each microservice into `api_gateway/proto/`.
2. Generate grpc-gateway code for each proto.
3. Run `go mod tidy` in `api_gateway`.
4. Start the gateway:
   ```sh
   go run main.go
   ```

## .env example
```
USER_SERVICE_ADDR=localhost:50051
TRANSACTION_SERVICE_ADDR=localhost:50054
PAYMENT_SERVICE_ADDR=localhost:50052
SMART_BUDGET_SERVICE_ADDR=localhost:50053
PORT=8080
```

## TODO
- Add OpenAPI/Swagger docs
- Add authentication/authorization if needed
