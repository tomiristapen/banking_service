package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	grpcadapter "github.com/tomiristapen/banking_service/smart_budget_service/adapter/grpc"
	"github.com/tomiristapen/banking_service/smart_budget_service/infrastructure/db"
	"github.com/tomiristapen/banking_service/smart_budget_service/infrastructure/mq"
	"github.com/tomiristapen/banking_service/smart_budget_service/metrics"
	smartbudgetpb "github.com/tomiristapen/banking_service/smart_budget_service/proto"
	"github.com/tomiristapen/banking_service/smart_budget_service/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracer() func() {
	exp, _ := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("smart_budget_service"),
		)),
	)
	otel.SetTracerProvider(tp)
	return func() { _ = tp.Shutdown(context.Background()) }
}

func main() {
	shutdown := initTracer()
	defer shutdown()

	_ = godotenv.Load()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	log.Printf("[SmartBudget] MONGO_URI: %s", mongoURI)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatalf("[SmartBudget] MongoDB ping failed: %v", err)
	}
	log.Println("[SmartBudget] Successfully connected and pinged MongoDB")
	col := client.Database("smart_budgeting_service").Collection("budget_expenses")
	repo := db.NewBudgetMongoRepo(col)
	uc := usecase.NewBudgetUsecase(repo)

	// NATS subscriber
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	subscriber := mq.NewMQSubscriber(natsURL, repo)
	go subscriber.SubscribePayments()

	// Prometheus metrics
	metrics.Init()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	// gRPC server
	port := os.Getenv("PORT")
	if port == "" {
		port = "50053"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	handler := grpcadapter.NewBudgetHandler(uc)
	server := grpc.NewServer()
	smartbudgetpb.RegisterSmartBudgetServiceServer(server, handler)
	log.Printf("✅ SmartBudgetService is running on port %s...", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	// Пример использования трейсера
	tracer := otel.Tracer("smart_budget_service")
	_, span := tracer.Start(context.Background(), "main")
	defer span.End()
}

// docker-compose logging config example:
// logging:
//   driver: loki
//   options:
//     loki-url: "http://localhost:3100/loki/api/v1/push"
