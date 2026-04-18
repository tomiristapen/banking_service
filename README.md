# Banking Microservices System

Backend system for banking operations built using a microservices architecture.
The system is designed with independent services communicating via gRPC and following clean architectural principles.

---

## Overview

This project demonstrates a modular backend system where each service is responsible for its own domain logic.
Services interact through well-defined interfaces and support both synchronous and asynchronous communication.

---

## Architecture

### Services

* api_gateway — handles incoming HTTP requests and routes them to services
* user_service — authentication and user management (JWT)
* transaction_service — transfer operations between accounts
* payment_service — payment processing
* smart_budget_service — budget management

### Communication

* gRPC with Protocol Buffers for service-to-service interaction
* Message queue (pub/sub) for asynchronous communication

---

## Tech Stack

* Go (Golang)
* gRPC, Protocol Buffers
* MongoDB
* JWT authentication
* grpc-gateway
* Clean Architecture (layered structure)

---

## Features

* Token-based authentication using JWT
* API Gateway for REST to gRPC translation
* Structured service-to-service communication
* Isolated user management service
* Transaction processing logic
* Payment handling service
* Budget tracking functionality
* Asynchronous messaging using pub/sub
* Separation of concerns with layered architecture

---

## How to Run

```bash
git clone https://github.com/tomiristapen/banking_service.git
cd banking_service
go run cmd/main.go
```

Each service can be started independently depending on configuration.

---

## Notes

This project focuses on backend architecture and service interaction.
It can be extended with containerization, monitoring, and service orchestration.

