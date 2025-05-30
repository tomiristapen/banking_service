package main

import (
	"context"
	"log"
	"net"
	"os"

	grpcadapter "transaction_service/adapter/grpc"
	"transaction_service/infrastructure/db"
	"transaction_service/infrastructure/mq"
	transactionpb "transaction_service/proto"
	"transaction_service/usecase"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	grpcclient "transaction_service/infrastructure/grpcclient"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is not set")
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	dbConn := client.Database("transaction_service")
	txRepo := db.NewTransactionMongo(dbConn)

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	publisher := mq.NewNatsPublisher(nc, "transactions")

	
	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	if userServiceAddr == "" {
		userServiceAddr = ":50051" 
	}
	userConn, err := grpc.Dial(userServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("❌ failed to connect to user service: %v", err)
	}
	defer userConn.Close()
	userServiceClient := grpcclient.NewUserServiceClient(userConn)

	uc := usecase.NewTransactionUsecase(txRepo, publisher, userServiceClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "50053"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	transactionpb.RegisterTransactionServiceServer(grpcServer, grpcadapter.NewTransactionHandler(uc))
	log.Printf("TransactionService gRPC server started on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
