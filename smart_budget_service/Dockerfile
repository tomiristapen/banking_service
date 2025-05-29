# syntax=docker/dockerfile:1
FROM golang:1.23.4-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o smart_budget_service ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/smart_budget_service .
EXPOSE 8083
CMD ["./smart_budget_service"]
