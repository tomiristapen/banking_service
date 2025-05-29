package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tomiristapen/banking_service/user_service/metrics"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	handler "github.com/tomiristapen/banking_service/user_service/adapter/grpc"
	infraMongo "github.com/tomiristapen/banking_service/user_service/infrastructure/mongo"
	"github.com/tomiristapen/banking_service/user_service/infrastructure/smtp"
	userpb "github.com/tomiristapen/banking_service/user_service/proto"
	"github.com/tomiristapen/banking_service/user_service/usecase"
)

func main() {
	metrics.Init()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	client, err := mongodriver.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}

	db := client.Database("user_service")

	userRepo := infraMongo.NewUserRepository(db)
	mail := smtp.NewMailer()
	userUC := usecase.NewUserUseCase(userRepo, mail)
	userHandler := handler.NewUserHandler(userUC)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, userHandler)

	reflection.Register(grpcServer)

	fmt.Printf("✅ UserService is running on port %s...\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// docker-compose logging config example:
// logging:
//   driver: loki
//   options:
//     loki-url: "http://localhost:3100/loki/api/v1/push"
