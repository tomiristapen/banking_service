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
	paymentCol, serviceCol := db.ConnectMongo()

	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	if userServiceAddr == "" {
		userServiceAddr = ":50051" 
	}
	userConn, err := grpc.Dial(userServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("❌ failed to connect to user service: %v", err)
	}
	defer userConn.Close()
	userClient := grpcclient.NewUserServiceClient(userConn)

	mqPublisher := mq.NewMQPublisher()
	defer mqPublisher.Close()

	// repo → usecase → handler
	repo := db.NewPaymentMongoRepo(paymentCol, serviceCol)
	uc := usecase.NewPaymentUsecase(repo, userClient, mqPublisher)
	h := handler.NewPaymentHandler(uc)

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
