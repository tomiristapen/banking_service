package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	handler "github.com/tomiristapen/banking_service/payment_service/adapter/grpc"
	"github.com/tomiristapen/banking_service/payment_service/infrastructure/db"
	"github.com/tomiristapen/banking_service/payment_service/infrastructure/grpcclient"
	"github.com/tomiristapen/banking_service/payment_service/infrastructure/mq"
	pb "github.com/tomiristapen/banking_service/payment_service/proto"
	"github.com/tomiristapen/banking_service/payment_service/usecase"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found. Proceeding with system env variables")
	}
}

func main() {
	// Ensure logs directory exists and set log output to file
	if err := os.MkdirAll("/app/logs", 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}
	logFile, err := os.OpenFile("/app/logs/payment_service.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
	defer logFile.Close()

	// Подключение к MongoDB
	paymentCol, serviceCol := db.ConnectMongo()

	// Подключение к UserService по gRPC
	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	if userServiceAddr == "" {
		userServiceAddr = ":50051" // default port
	}
	userConn, err := grpc.Dial(userServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("❌ failed to connect to user service: %v", err)
	}
	defer userConn.Close()
	userClient := grpcclient.NewUserServiceClient(userConn)

	// Подключение к NATS
	mqPublisher := mq.NewMQPublisher()
	defer mqPublisher.Close()

	// Clean Architecture: repo → usecase → handler
	repo := db.NewPaymentMongoRepo(paymentCol, serviceCol)
	uc := usecase.NewPaymentUsecase(repo, userClient, mqPublisher)
	h := handler.NewPaymentHandler(uc)

	// gRPC server
	port := os.Getenv("PORT")
	if port == "" {
		port = "50052"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("❌ failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s, h)

	log.Printf("✅ PaymentService is running on port %s...\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("❌ failed to serve: %v", err)
	}
}
